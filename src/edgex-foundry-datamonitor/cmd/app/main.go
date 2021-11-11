// Copyright 2021 Alessandro De Blasis <alex@deblasis.net>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/deblasis/edgex-foundry-datamonitor/bundled"
	"github.com/deblasis/edgex-foundry-datamonitor/config"
	"github.com/deblasis/edgex-foundry-datamonitor/messaging"
	"github.com/deblasis/edgex-foundry-datamonitor/pages"
	"github.com/deblasis/edgex-foundry-datamonitor/services"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var topWindow fyne.Window

func main() {
	a := app.NewWithID("edgex-datamonitor")
	a.SetIcon(bundled.ResourceBgxPng)
	logLifecycle(a)
	w := a.NewWindow("EdgeX Data Monitor")
	topWindow = w
	w.SetMaster()

	cfg := config.GetConfig()
	client, err := messaging.NewClient(cfg)
	if err != nil {
		uerr := errors.New("Error while initializing client")
		dialog.ShowError(uerr, topWindow)
		log.Println(err)
	}

	events := make(chan *dtos.Event)
	ep := services.NewEventProcessor(events)
	go ep.Run()

	client.OnConnect = func() bool {
		messages, errs := client.Subscribe(config.DefaultEventsTopic)

		ok := make(chan bool, 1)

		go func() {
		LOOP:
			for {
				select {
				case err := <-errs:
					if client.IsConnecting {
						if strings.Contains(err.Error(), "redis: client is closed") {
							//handling "redis: client is closed on connect" which is ok because it's then set by go-mod-messaging and the error is ignored
							continue
						}
						uerr := errors.New("Error while subscribing to Redis")
						dialog.ShowError(uerr, topWindow)
						log.Println(err)
						client.IsConnecting = false
						ok <- false
						break LOOP
					}
				case msgEnvelope := <-messages:
					event, _ := messaging.ParseEvent(msgEnvelope.Payload)
					events <- event
					select {
					case ok <- true:
					default:
					}
				}
			}
		}()

		return <-ok
	}

	AppManager := services.NewAppManager(client, cfg, ep)

	shouldConnect := a.Preferences().BoolWithFallback(config.PrefShouldConnectAtStartup, false)

	if shouldConnect {
		a.SendNotification(&fyne.Notification{
			Title:   "Connecting...",
			Content: fmt.Sprintf("Connecting to %v:%v", cfg.GetRedisHost(), cfg.GetRedisPort()),
		})
		if err = client.Connect(); err != nil {
			uerr := errors.New(fmt.Sprintf("Cannot connect\n%s", err))
			dialog.ShowError(uerr, topWindow)
			log.Println(err)
		}

	}

	content := container.NewMax()
	title := widget.NewLabel("Component name")
	intro := widget.NewLabel("An introduction would probably go\nhere, as well as a")
	intro.Wrapping = fyne.TextWrapWord
	setPage := func(uid string, t pages.Page, appMgr *services.AppManager) {
		if fyne.CurrentDevice().IsMobile() {
			child := a.NewWindow(t.Title)
			topWindow = child
			child.SetContent(t.View(topWindow, appMgr))
			child.Show()
			child.SetOnClosed(func() {
				topWindow = w
			})
			return
		}

		title.SetText(t.Title)
		intro.SetText(t.Intro)
		draw := func(cnt *fyne.Container) {
			cnt.Objects = []fyne.CanvasObject{t.View(w, appMgr)}
			cnt.Refresh()
		}
		draw(content)

		appMgr.SetCurrentContainer(content, draw)
	}

	page := container.NewBorder(container.NewVBox(title, widget.NewSeparator(), intro), nil, nil, nil, content)

	navBar := makeNav(setPage, AppManager)
	AppManager.SetNav(navBar.(*fyne.Container))

	if fyne.CurrentDevice().IsMobile() {
		w.SetContent(navBar)
	} else {
		split := container.NewHSplit(navBar, page)
		split.Offset = 0.2
		w.SetContent(split)
	}
	w.Resize(fyne.NewSize(1024, 768))
	w.ShowAndRun()
}

func logLifecycle(a fyne.App) {
	a.Lifecycle().SetOnStarted(func() {
		log.Println("Lifecycle: Started")
	})
	a.Lifecycle().SetOnStopped(func() {
		log.Println("Lifecycle: Stopped")
	})
	a.Lifecycle().SetOnEnteredForeground(func() {
		log.Println("Lifecycle: Entered Foreground")
	})
	a.Lifecycle().SetOnExitedForeground(func() {
		log.Println("Lifecycle: Exited Foreground")
	})
}

func makeNav(setPage func(_ string, page pages.Page, _ *services.AppManager), appMgr *services.AppManager) fyne.CanvasObject {
	a := fyne.CurrentApp()

	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return pages.PageIndex[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := pages.PageIndex[uid]

			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Nav widgets")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			t, ok := pages.Pages[uid]
			if !ok {
				fyne.LogError("Missing panel: "+uid, nil)
				return
			}
			obj.(*widget.Label).SetText(t.Title)
		},
		OnSelected: func(uid string) {
			if t, ok := pages.Pages[uid]; ok {
				setPage(uid, t, appMgr)
			}
		},
	}

	tree.Select("home")

	themes := container.New(layout.NewGridLayout(2),
		widget.NewButton("Dark", func() {
			a.Settings().SetTheme(theme.DarkTheme())
		}),
		widget.NewButton("Light", func() {
			a.Settings().SetTheme(theme.LightTheme())
		}),
	)

	disconnectBtn := widget.NewButtonWithIcon("Disconnect", theme.LogoutIcon(), func() {
		appMgr.Disconnect()
		appMgr.Refresh()
	})

	switch appMgr.GetConnectionState() {
	case services.ClientConnected:
		disconnectBtn.Show()
	default:
		disconnectBtn.Hide()
	}

	buttons := container.NewVBox(
		disconnectBtn,
		themes,
	)

	return container.NewBorder(nil, buttons, nil, nil, tree)
}
