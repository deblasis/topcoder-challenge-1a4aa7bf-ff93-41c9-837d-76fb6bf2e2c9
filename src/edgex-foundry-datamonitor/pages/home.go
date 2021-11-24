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
package pages

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/deblasis/edgex-foundry-datamonitor/bundled"
	"github.com/deblasis/edgex-foundry-datamonitor/services"
)

func homeScreen(w fyne.Window, appManager *services.AppManager) fyne.CanvasObject {

	h := appManager.GetPageHandler(HomePageKey).(*homePageHandler)

	var contentContainer *fyne.Container

	logo := canvas.NewImageFromResource(bundled.ResourceCompanyLogoLgPng)
	logo.FillMode = canvas.ImageFillStretch
	if fyne.CurrentDevice().IsMobile() {
		logo.SetMinSize(fyne.NewSize(171, 57))
	} else {
		logo.SetMinSize(fyne.NewSize(500, 165))
	}

	redisHost, redisPort := appManager.GetRedisHostPort()
	connectionState := appManager.GetConnectionState()

	connectingContent := container.NewCenter(container.NewVBox(
		container.NewHBox(widget.NewProgressBarInfinite()),
	))

	disconnectedContent := container.NewCenter(container.NewVBox(
		logo,
		widget.NewCard("You are currently disconnected from EdgeX Foundry",
			fmt.Sprintf("Would you like to connect to %v:%d?", redisHost, redisPort),
			container.NewCenter(
				widget.NewButtonWithIcon("Connect", theme.LoginIcon(), func() {
					if err := appManager.Connect(); err != nil {
						uerr := fmt.Errorf("Cannot connect\n%s", err)
						dialog.ShowError(uerr, w)
						log.Errorf("cannot connect: %v", err)
					}
					appManager.Refresh()
				}),
			),
		),
	))

	h.SetInitialState()
	h.RehydrateSession()
	h.SetupBindings()

	connectedContent := h.dashboardStats
	h.dashboardStats.Hide()
	h.tableContainer.Hide()

	switch connectionState {
	case services.ClientConnected:
		contentContainer = connectedContent
		h.dashboardStats.Show()
		h.tableContainer = container.NewMax(h.dashboardTable)
	case services.ClientConnecting:
		contentContainer = connectingContent
		h.dashboardStats.Hide()
		h.tableContainer.Hide()
	case services.ClientDisconnected:
		contentContainer = disconnectedContent
		h.dashboardStats.Hide()
		h.tableContainer = container.NewMax()
	}

	home := container.NewBorder(
		contentContainer,
		nil, nil, nil, h.tableContainer,
	)

	return home

}
