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
	"errors"
	"fmt"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/deblasis/edgex-foundry-datamonitor/bundled"
	"github.com/deblasis/edgex-foundry-datamonitor/config"
	"github.com/deblasis/edgex-foundry-datamonitor/services"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos"
)

func homeScreen(w fyne.Window, appManager *services.AppManager) fyne.CanvasObject {
	a := fyne.CurrentApp()
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
	eventProcessor := appManager.GetEventProcessor()

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
						uerr := errors.New(fmt.Sprintf("Cannot connect\n%s", err))
						dialog.ShowError(uerr, w)
						log.Printf("cannot connect: %v", err)
					}
					appManager.Refresh()
				}),
			),
		),
	))

	eventsTable := renderEventsTable(eventProcessor.LastEvents.Get(), false)
	tableContainer := container.NewMax(eventsTable)

	totalNumberEventsBinding := binding.BindInt(config.Int(eventProcessor.TotalNumberEvents))
	totalNumberReadingsBinding := binding.BindInt(config.Int(eventProcessor.TotalNumberReadings))

	eventsPerSecondLastMinute := binding.BindFloat(config.Float(eventProcessor.EventsPerSecondLastMinute))
	readingsPerSecondLastMinute := binding.BindFloat(config.Float(eventProcessor.ReadingsPerSecondLastMinute))

	dashboardStats := container.NewCenter(container.NewGridWithRows(2,
		container.NewGridWithColumns(4,
			widget.NewLabelWithStyle("Total Number of Events", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithData(binding.IntToString(totalNumberEventsBinding)),
			widget.NewLabelWithStyle("Total Number of Readings", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithData(binding.IntToString(totalNumberReadingsBinding)),
		),
		container.NewGridWithColumns(4,
			widget.NewLabelWithStyle("Events per second", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithData(binding.FloatToString(eventsPerSecondLastMinute)),
			widget.NewLabelWithStyle("Readings per second", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithData(binding.FloatToString(readingsPerSecondLastMinute)),
		),
	))
	connectedContent := dashboardStats

	dashboardStats.Hide()
	tableContainer.Hide()

	switch connectionState {
	case services.ClientConnected:
		contentContainer = connectedContent
		dashboardStats.Show()
		tableContainer = container.NewMax(eventsTable)
	case services.ClientConnecting:
		contentContainer = connectingContent
		dashboardStats.Hide()
		tableContainer.Hide()
	case services.ClientDisconnected:
		contentContainer = disconnectedContent
		dashboardStats.Hide()
		tableContainer = container.NewMax()
	}

	home := container.NewGridWithRows(2,
		contentContainer,
		tableContainer,
	)

	go func() {
		for {
			time.Sleep(100 * time.Millisecond)
			if appManager.GetConnectionState() != services.ClientConnected {
				continue
			}
			//log.Printf("refreshing UI: events %v\n", ep.TotalNumberEvents)
			totalNumberEventsBinding.Set(eventProcessor.TotalNumberEvents)
			totalNumberReadingsBinding.Set(eventProcessor.TotalNumberReadings)
			eventsPerSecondLastMinute.Set(eventProcessor.EventsPerSecondLastMinute)
			readingsPerSecondLastMinute.Set(eventProcessor.ReadingsPerSecondLastMinute)
			eventsTable = renderEventsTable(eventProcessor.LastEvents.Get(), a.Preferences().BoolWithFallback(config.PrefEventsTableSortOrderAscending, config.DefaultEventsTableSortOrderAscending))

			if len(tableContainer.Objects) > 0 {
				tableContainer.Objects[0] = eventsTable
			}

		}
	}()

	return home

}

func renderEventsTable(events []*dtos.Event, sortAsc bool) fyne.CanvasObject {

	// the slice is fifo, we reverse it so that the first element is the most recent
	evts := make([]*dtos.Event, len(events))
	copy(evts, events)

	if !sortAsc {
		for i, j := 0, len(evts)-1; i < j; i, j = i+1, j-1 {
			evts[i], evts[j] = evts[j], evts[i]
		}
	}

	renderCell := func(row, col int, label *widget.Label) {

		if len(evts) == 0 || row >= len(evts) {
			label.SetText("")
			return
		}

		event := evts[row]
		switch col {
		case 0:
			label.SetText(event.DeviceName)
		case 1:
			label.SetText(time.Unix(0, event.Origin).String())
		default:
			label.SetText("")
		}

	}

	table := widget.NewTable(
		func() (int, int) { return 6, 2 },
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)
			switch id.Row {
			case 0:
				switch id.Col {
				// case 0:
				// 	label.SetText("ID")
				case 0:
					label.SetText("Device Name")
					label.TextStyle = fyne.TextStyle{Bold: true}
				case 1:
					label.SetText("Origin Timestamp")
					label.TextStyle = fyne.TextStyle{Bold: true}
				default:
					label.SetText("")
				}

			default:
				renderCell(id.Row-1, id.Col, label)
			}

		})
	// t.SetColumnWidth(0, 34)
	table.SetColumnWidth(0, 350)
	table.SetColumnWidth(1, 350)

	sortorder := "ascendingly"
	if !sortAsc {
		sortorder = "descendingly"
	}
	return container.NewBorder(
		container.NewVBox(layout.NewSpacer(), container.NewHBox(
			widget.NewLabelWithStyle("Last 5 events", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			layout.NewSpacer(),
			widget.NewLabelWithStyle(fmt.Sprintf("sorted %v by timestamp", sortorder), fyne.TextAlignTrailing, fyne.TextStyle{Italic: true}),
		)),
		nil,
		nil,
		nil,
		table,
	)

}
