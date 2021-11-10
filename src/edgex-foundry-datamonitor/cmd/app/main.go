package main

import (
	"errors"
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/deblasis/edgex-foundry-datamonitor/eventsprocessor"
	"github.com/deblasis/edgex-foundry-datamonitor/internal/config"
	"github.com/deblasis/edgex-foundry-datamonitor/internal/pages"
	"github.com/deblasis/edgex-foundry-datamonitor/internal/state"
	"github.com/deblasis/edgex-foundry-datamonitor/messaging"
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
	a.SetIcon(theme.FyneLogo()) //TODO change
	logLifecycle(a)
	w := a.NewWindow("EdgeX Data Monitor")
	topWindow = w
	w.SetMaster()

	cfg := config.GetConfig()
	client, err := messaging.NewClient(cfg)
	if err != nil {
		uerr := errors.New("Error while initializing client")
		dialog.ShowError(uerr, topWindow)
		//TODO: log this
	}

	events := make(chan *dtos.Event)
	ep := eventsprocessor.New(events)
	go ep.Run()

	messages, _ := client.Subscribe(config.DefaultEventsTopic)
	go func() {
		for {
			select {
			// case e := <-errors:
			//TODO errors
			case msgEnvelope := <-messages:
				event, _ := messaging.ParseEvent(msgEnvelope.Payload)
				//TODO errors
				events <- event
			}
		}
	}()

	AppManager := state.NewAppManager(client, cfg, ep)

	shouldConnect := a.Preferences().BoolWithFallback(config.PrefShouldConnectAtStartup, false)

	if shouldConnect && cfg.RedisHost != nil && cfg.RedisPort != nil {
		a.SendNotification(&fyne.Notification{
			Title:   "Connecting...",
			Content: fmt.Sprintf("Connecting to %v:%v", cfg.RedisHost, cfg.RedisPort),
		})
		if err = client.Connect(); err != nil {
			uerr := errors.New(fmt.Sprintf("Cannot connect\n%s", err))
			dialog.ShowError(uerr, topWindow)
			//TODO: log this
		}

	}

	content := container.NewMax()
	title := widget.NewLabel("Component name")
	intro := widget.NewLabel("An introduction would probably go\nhere, as well as a")
	intro.Wrapping = fyne.TextWrapWord
	setPage := func(uid string, t pages.Page, appMgr *state.AppManager) {
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

	page := container.NewBorder(
		container.NewVBox(title, widget.NewSeparator(), intro), nil, nil, nil, content)
	if fyne.CurrentDevice().IsMobile() {
		w.SetContent(makeNav(setPage, AppManager))
	} else {
		split := container.NewHSplit(makeNav(setPage, AppManager), page)
		split.Offset = 0.2
		w.SetContent(split)
	}
	w.Resize(fyne.NewSize(800, 600))
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

func makeNav(setPage func(_ string, page pages.Page, _ *state.AppManager), state *state.AppManager) fyne.CanvasObject {
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
			return widget.NewLabel("Collection Widgets")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			t, ok := pages.Pages[uid]
			if !ok {
				fyne.LogError("Missing tutorial panel: "+uid, nil)
				return
			}
			obj.(*widget.Label).SetText(t.Title)
		},
		OnSelected: func(uid string) {
			if t, ok := pages.Pages[uid]; ok {
				//a.Preferences().SetString(preferenceCurrentTutorial, uid)
				setPage(uid, t, state)
			}
		},
	}

	//TODO refactor
	tree.Select("home")

	themes := container.New(layout.NewGridLayout(2),
		widget.NewButton("Dark", func() {
			a.Settings().SetTheme(theme.DarkTheme())
		}),
		widget.NewButton("Light", func() {
			a.Settings().SetTheme(theme.LightTheme())
		}),
	)

	return container.NewBorder(nil, themes, nil, nil, tree)
}

func shortcutFocused(s fyne.Shortcut, w fyne.Window) {
	if focused, ok := w.Canvas().Focused().(fyne.Shortcutable); ok {
		focused.TypedShortcut(s)
	}
}
