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
	"encoding/json"
	"sort"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/deblasis/edgex-foundry-datamonitor/config"
	"github.com/deblasis/edgex-foundry-datamonitor/services"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos"
)

type homePageHandler struct {
	appState *services.AppManager

	Key widget.TreeNodeID

	totalNumberEventsBinding   binding.ExternalInt
	totalNumberReadingsBinding binding.ExternalInt

	eventsPerSecondLastMinute   binding.ExternalFloat
	readingsPerSecondLastMinute binding.ExternalFloat

	//eventsTable    fyne.CanvasObject
	dashboardTable               *widget.Table
	dashboardTableDataMapBinding *[]binding.DataMap
	dashboardTableLock           sync.Mutex

	tableContainer *fyne.Container
	dashboardStats *fyne.Container
}

func NewHomePageHandler(appState *services.AppManager) *homePageHandler {
	p := &homePageHandler{
		Key:      HomePageKey,
		appState: appState,
	}

	p.updateTable()

	return p
}
func (p *homePageHandler) SetInitialState() {
	p.dashboardTableDataMapBinding = &[]binding.DataMap{}
	p.dashboardTable = p.renderDashboardTable()
	p.tableContainer = container.NewMax(p.dashboardTable)

}
func (p *homePageHandler) RehydrateSession() {}
func (p *homePageHandler) SetupBindings() {
	eventProcessor := p.appState.GetEventProcessor()

	p.totalNumberEventsBinding = binding.BindInt(config.Int(eventProcessor.TotalNumberEvents))
	p.totalNumberReadingsBinding = binding.BindInt(config.Int(eventProcessor.TotalNumberReadings))

	p.eventsPerSecondLastMinute = binding.BindFloat(config.Float(eventProcessor.EventsPerSecondLastMinute))
	p.readingsPerSecondLastMinute = binding.BindFloat(config.Float(eventProcessor.ReadingsPerSecondLastMinute))

	p.dashboardStats = container.NewCenter(container.NewGridWithRows(3,
		container.NewGridWithColumns(4,
			widget.NewLabelWithStyle("Total Number of Events", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithData(binding.IntToString(p.totalNumberEventsBinding)),
			widget.NewLabelWithStyle("Total Number of Readings", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithData(binding.IntToString(p.totalNumberReadingsBinding)),
		),
		container.NewGridWithColumns(4,
			widget.NewLabelWithStyle("Events per second", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithData(binding.FloatToString(p.eventsPerSecondLastMinute)),
			widget.NewLabelWithStyle("Readings per second", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabelWithData(binding.FloatToString(p.readingsPerSecondLastMinute)),
		),
		layout.NewSpacer(),
	))
}

func (p *homePageHandler) OnEventReceived(event dtos.Event) {

	if p.appState.GetConnectionState() != services.ClientConnected {
		return
	}
	eventProcessor := p.appState.GetEventProcessor()
	p.totalNumberEventsBinding.Set(eventProcessor.TotalNumberEvents)
	p.totalNumberReadingsBinding.Set(eventProcessor.TotalNumberReadings)
	p.eventsPerSecondLastMinute.Set(eventProcessor.EventsPerSecondLastMinute)
	p.readingsPerSecondLastMinute.Set(eventProcessor.ReadingsPerSecondLastMinute)

	p.updateTable()
	if p.dashboardTable != nil {
		p.dashboardTable.Refresh()
	}

}

func (p *homePageHandler) updateTable() {
	sortAsc := fyne.CurrentApp().Preferences().BoolWithFallback(config.PrefEventsTableSortOrderAscending, config.DefaultEventsTableSortOrderAscending)

	p.dashboardTableLock.Lock()
	defer p.dashboardTableLock.Unlock()
	p.dashboardTableDataMapBinding = &[]binding.DataMap{}

	events := p.appState.GetEventProcessor().LastEvents.Get()

	evts := make([]*dtos.Event, len(events))
	copy(evts, events)

	if !sortAsc {
		sort.Slice(evts, func(i, j int) bool {
			return evts[i].Origin > evts[j].Origin
		})
	}

	for _, row := range evts {
		tags, _ := json.MarshalIndent(row.Tags, "", "    ")
		eventJson, _ := json.MarshalIndent(row, "", "    ")
		r := eventRow{
			Id:            row.Id,
			DeviceName:    row.DeviceName,
			ProfileName:   row.ProfileName,
			Created:       row.Created,
			Origin:        row.Origin,
			ReadingsCount: int64(len(row.Readings)),
			Tags:          string(tags),
			Json:          string(eventJson),
		}
		*p.dashboardTableDataMapBinding = append(*p.dashboardTableDataMapBinding, binding.BindStruct(&r))
	}

}

func (p *homePageHandler) renderDashboardTable() *widget.Table {

	table := widget.NewTable(
		func() (int, int) { return 6, 2 },
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			label := o.(*widget.Label)
			switch i.Row {
			case 0:
				switch i.Col {
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
				label.TextStyle = fyne.TextStyle{Bold: false}

				dm := *p.dashboardTableDataMapBinding
				if dm == nil || i.Row > len(dm) {
					label.SetText("")
					break
				}

				row := dm[i.Row-1]

				switch i.Col {
				case 0:
					id, _ := row.GetItem("DeviceName")
					o.(*widget.Label).Bind(id.(binding.String))
				case 1:
					origin, _ := row.GetItem("Origin")
					v, _ := origin.(binding.Int).Get()
					o.(*widget.Label).SetText(time.Unix(0, int64(v)).String())
				default:
					label.SetText("")
				}
			}

		})
	table.SetColumnWidth(0, 350)
	table.SetColumnWidth(1, 350)

	return table
}
