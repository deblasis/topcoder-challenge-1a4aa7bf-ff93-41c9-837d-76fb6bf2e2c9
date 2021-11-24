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
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/deblasis/edgex-foundry-datamonitor/services"
)

func dataScreen(win fyne.Window, appManager *services.AppManager) fyne.CanvasObject {

	connectionState := appManager.GetConnectionState()
	redisHost, redisPort := appManager.GetRedisHostPort()

	disconnectedContent := container.NewCenter(container.NewVBox(
		widget.NewCard("You are currently disconnected from EdgeX Foundry",
			fmt.Sprintf("Would you like to connect to %v:%d?", redisHost, redisPort),
			container.NewCenter(
				widget.NewButtonWithIcon("Connect", theme.LoginIcon(), func() {
					if err := appManager.Connect(); err != nil {
						uerr := fmt.Errorf("Cannot connect\n%s", err)
						dialog.ShowError(uerr, win)
						log.Error("cannot connect: %v", err)
					}
					appManager.Refresh()
				}),
			),
		),
	))
	if connectionState == services.ClientDisconnected {
		return disconnectedContent
	}

	h := appManager.GetPageHandler(DataPageKey).(*dataPageHandler)

	// It will have a radio button to select between events or readings
	radioGroup := container.NewVBox(
		widget.NewLabelWithStyle("Show", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		h.dataType,
	)
	searchBox := container.NewVBox(
		widget.NewLabelWithStyle("Filter", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		container.NewBorder(nil, nil, nil, container.NewHBox(h.searchBtn, h.resetSearchBtn), container.NewMax(h.search)),
	)

	bufferSizeContainer := container.NewGridWithColumns(2,
		container.NewBorder(nil, nil, widget.NewLabel("Buffer size"), h.applyBufferSizeBtn, h.bufferSize),
		h.bufferProgress,
	)

	h.SetInitialState()
	h.RehydrateSession()
	h.SetupBindings()

	h.updateTableByDataType(h.dataType.Selected)

	h.tableContainer = container.NewBorder(
		container.NewVBox(
			widget.NewSeparator(),
			h.tableHeading,
			widget.NewSeparator(),
		),
		nil,
		nil,
		nil,
		container.NewMax(h.eventsTable, h.readingsTable),
	)

	content := container.NewBorder(
		container.NewVBox(
			container.NewGridWithColumns(2,
				radioGroup,
				searchBox,
			),
			widget.NewSeparator(),
			bufferSizeContainer,
		),
		nil, nil, nil, h.tableContainer,
	)

	return content

}
