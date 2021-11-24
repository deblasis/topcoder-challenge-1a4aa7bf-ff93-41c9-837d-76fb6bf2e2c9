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
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/deblasis/edgex-foundry-datamonitor/config"
	"github.com/deblasis/edgex-foundry-datamonitor/data"
	"github.com/deblasis/edgex-foundry-datamonitor/services"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos"
)

type dataPageHandler struct {
	appState *services.AppManager
	Key      widget.TreeNodeID

	sortAsc bool

	dataType           *widget.RadioGroup
	search             *widget.Entry
	searchBtn          *widget.Button
	resetSearchBtn     *widget.Button
	applyBufferSizeBtn *widget.Button
	bufferSize         *widget.Entry
	bufferSizeBinding  binding.Int
	bufferUsageBinding binding.Float

	statusText *widget.Label

	bufferProgress *widget.ProgressBar
	tableHeading   *fyne.Container

	table         *widget.Table
	eventsTable   *widget.Table
	readingsTable *widget.Table
	tableViewLock sync.RWMutex

	tableContainer              *fyne.Container
	eventsTableDataMapBinding   *[]binding.DataMap
	readingsTableDataMapBinding *[]binding.DataMap

	tableDataLock sync.RWMutex

	jsonDetail *widget.Entry
}

func NewDataPageHandler(appState *services.AppManager) *dataPageHandler {

	p := &dataPageHandler{
		appState: appState,
		Key:      DataPageKey,
	}

	p.dataType = widget.NewRadioGroup([]string{config.DataTypeEvents, config.DataTypeReadings}, func(dataType string) {
		log.Debugf("Selected %s", dataType)
		p.appState.SetDataPageSelectedDataType(dataType)
	})
	p.dataType.Horizontal = true
	p.dataType.Required = true

	p.search = widget.NewEntry()
	p.search.SetPlaceHolder("Type here to loosely search")

	p.searchBtn = widget.NewButtonWithIcon("Search", theme.SearchIcon(), func() {})
	p.resetSearchBtn = widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {})

	p.bufferSizeBinding = binding.NewInt()
	p.bufferSize = widget.NewEntryWithData(binding.IntToString(p.bufferSizeBinding))

	// p.bufferUsageBinding = binding.BindFloat(config.Float(float64(p.appState.GetDB().Ge)))

	p.bufferSize.SetPlaceHolder("Buffer size")
	p.bufferSize.Validator = data.MinMaxValidator(config.MinBufferSize, config.MaxBufferSize, data.ErrInvalidBufferSize)
	p.applyBufferSizeBtn = widget.NewButtonWithIcon("", theme.DocumentSaveIcon(), func() {})

	p.bufferProgress = widget.NewProgressBar()

	p.bufferProgress.TextFormatter = func() string {
		return fmt.Sprintf("Buffered %d/%d", int(p.bufferProgress.Value), int(p.bufferProgress.Max))
	}

	return p
}

func (p *dataPageHandler) SetInitialState() {

	if p.dataType.Selected == "" {
		p.dataType.SetSelected(config.DataTypeEvents)
	}
	p.searchBtn.Disable()
	p.resetSearchBtn.Disable()
	p.applyBufferSizeBtn.Disable()

	//preferences := fyne.CurrentApp().Preferences()
	// defaultBufSize := preferences.IntWithFallback(config.PrefBufferSizeInDataPage, config.DefaultBufferSizeInDataPage)
	// b := p.bufferSizeBinding
	// b.Set(defaultBufSize)
	// p.bufferSize.SetText(fmt.Sprintf("%d", defaultBufSize))
	// p.bufferProgress.Max = float64(defaultBufSize)

	p.eventsTableDataMapBinding = &[]binding.DataMap{}
	p.readingsTableDataMapBinding = &[]binding.DataMap{}

	p.eventsTable = p.renderEventsTable()
	p.readingsTable = p.renderReadingsTable()

	p.setTableByDataType(p.dataType.Selected, false)

	//Json detail dialog
	p.jsonDetail = widget.NewEntry()
	p.jsonDetail.Disabled()

	win := fyne.CurrentApp().Driver().AllWindows()[0]

	copyToClipboardBtn := widget.NewButtonWithIcon("Copy to clipboard", theme.ContentCopyIcon(), func() {
		fyne.Clipboard.SetContent(win.Clipboard(), p.jsonDetail.Text)
	})

	detailBox := container.NewBorder(
		container.NewBorder(nil, nil, nil, copyToClipboardBtn, widget.NewLabelWithStyle("Selected item", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})),
		nil, nil, nil,
		container.NewVScroll(container.NewMax(p.jsonDetail)),
	)
	dlg := dialog.NewCustom("Detail", "Close", detailBox, win)
	dlg.Resize(fyne.NewSize(800, 1000))

	p.eventsTable.OnSelected = func(id widget.TableCellID) {
		if id.Row == 0 {
			return
		}

		p.tableViewLock.RLock()
		defer p.tableViewLock.RUnlock()
		dm := *p.eventsTableDataMapBinding
		js, _ := dm[id.Row-1].GetItem("Json")

		v, _ := js.(binding.String).Get()

		p.jsonDetail.SetText(v)
		dlg.Show()
		p.eventsTable.UnselectAll()
	}

	p.readingsTable.OnSelected = func(id widget.TableCellID) {
		if id.Row == 0 {
			return
		}

		p.tableViewLock.RLock()
		defer p.tableViewLock.RUnlock()
		dm := *p.readingsTableDataMapBinding
		js, _ := dm[id.Row-1].GetItem("Json")

		v, _ := js.(binding.String).Get()

		p.jsonDetail.SetText(v)
		dlg.Show()
		p.readingsTable.UnselectAll()
	}

}

func (p *dataPageHandler) RehydrateSession() {
	sessionSearch := p.appState.GetDataPageSearch()
	if sessionSearch != nil {
		p.search.Text = config.StringVal(sessionSearch)
		p.resetSearchBtn.Enable()
	}
	bufferSize := p.appState.GetDataPageBufferSize()
	if bufferSize != nil {
		b := p.bufferSizeBinding
		b.Set(*bufferSize)
		log.Debugf("bufferSize is %v", *bufferSize)
		//p.bufferSize.SetText(fmt.Sprintf("%d", *bufferSize))
		//p.bufferSize.SetText("100")
	} else {
		preferences := fyne.CurrentApp().Preferences()
		defaultBufSize := preferences.IntWithFallback(config.PrefBufferSizeInDataPage, config.DefaultBufferSizeInDataPage)
		b := p.bufferSizeBinding
		b.Set(defaultBufSize)
		p.bufferSize.SetText(fmt.Sprintf("%d", defaultBufSize))
		p.bufferProgress.Max = float64(defaultBufSize)
	}
}

func (p *dataPageHandler) SetupBindings() {

	p.search.OnChanged = func(s string) {
		if strings.Trim(s, " ") != "" {
			p.searchBtn.Enable()
		} else {
			p.searchBtn.Disable()
		}
	}

	p.searchBtn.OnTapped = func() {
		p.appState.SetDataPageSearch(p.search.Text)
		p.resetSearchBtn.Enable()

		p.updateTableByDataType(p.dataType.Selected)
		p.updateStatusByDataType(p.dataType.Selected)
	}
	p.search.OnSubmitted = func(s string) {
		if p.search.Validate() != nil {
			return
		}
		p.searchBtn.OnTapped()
	}

	b := p.bufferSizeBinding
	b.AddListener(binding.NewDataListener(func() {
		v, _ := b.Get()

		p.appState.SetDataPageBufferSize(v)

		p.bufferProgress.Max = float64(v)
		p.bufferProgress.Refresh()

		log.Debugf("bufferSizeBinding CHANGED to %v", v)
		//retriggering validation, updating the binding alone doesn't do it
	}))

	p.applyBufferSizeBtn.OnTapped = func() {
		v, err := strconv.Atoi(p.bufferSize.Text)
		if err != nil {
			return
		}

		b.Set(v)
		p.bufferProgress.Max = float64(v)
		p.bufferProgress.Refresh()

		log.Debug("applyBufferSizeBtn tapped")
	}

	p.bufferSize.OnChanged = func(s string) {
		if p.bufferSize.Validate() != nil {
			p.applyBufferSizeBtn.Disable()
			return
		}
		log.Debugf("bufferSize changed to %v", s)
		boundValue, _ := b.Get()

		if fmt.Sprintf("%d", boundValue) != s && s != "" {
			p.applyBufferSizeBtn.Enable()
		} else {
			p.applyBufferSizeBtn.Disable()
		}
	}

	p.bufferSize.OnSubmitted = func(s string) {
		if p.bufferSize.Validate() != nil {
			return
		}
		p.applyBufferSizeBtn.OnTapped()
	}

	p.resetSearchBtn.OnTapped = func() {
		log.Debug("filter reset")
		p.appState.SetDataPageSearch("")
		p.search.Text = ""
		p.search.Refresh()
		p.searchBtn.Disable()

		// updating both because the user could switch datatype with a running filter
		// and see the one in the background inconsistent with the UI
		p.updateTableByDataType(config.DataTypeEvents)
		p.updateTableByDataType(config.DataTypeReadings)

		p.updateStatusByDataType(p.dataType.Selected)
	}

	currentDataType := p.dataType.Selected
	p.setBufferUsageBindingByDataType(currentDataType)

	sortorder := "descendingly"
	if fyne.CurrentApp().Preferences().BoolWithFallback(config.PrefEventsTableSortOrderAscending, config.DefaultEventsTableSortOrderAscending) {
		sortorder = "ascendingly"
	}

	p.statusText = widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	p.tableHeading = container.NewHBox(
		p.statusText,
		layout.NewSpacer(),
		widget.NewLabelWithStyle(fmt.Sprintf("sorted %v by timestamp", sortorder), fyne.TextAlignTrailing, fyne.TextStyle{Italic: true}),
	)

	p.bufferUsageBinding.AddListener(binding.NewDataListener(func() {
		log.Debug("updated bufferUsageBinding")
		p.bufferProgress.Refresh()

		p.tableDataLock.RLock()
		defer p.tableDataLock.RUnlock()
	}))

	p.bufferProgress.Bind(p.bufferUsageBinding)

	p.dataType.OnChanged = func(currentDataType string) {

		//change bindings
		log.Debugf("changed dataType to %v", currentDataType)
		p.setBufferUsageBindingByDataType(currentDataType)

		p.updateStatusByDataType(p.dataType.Selected)

		//change table
		p.setTableByDataType(currentDataType, true)

	}

}

func (p *dataPageHandler) setTableByDataType(currentDataType string, refresh bool) {
	p.tableViewLock.Lock()
	defer p.tableViewLock.Unlock()

	switch currentDataType {
	case config.DataTypeEvents:
		p.table = p.eventsTable
		p.eventsTable.Show()
		p.readingsTable.Hide()

	case config.DataTypeReadings:
		p.table = p.readingsTable
		p.eventsTable.Hide()
		p.readingsTable.Show()
	default:
		log.Fatalf("unhandled type %v", currentDataType)
	}

	if refresh && p.tableContainer != nil {
		p.updateTableByDataType(p.dataType.Selected)
		p.tableContainer.Refresh()
	}
}

func (p *dataPageHandler) renderReadingsTable() *widget.Table {
	t := widget.NewTable(
		func() (int, int) {
			p.tableDataLock.RLock()
			defer p.tableDataLock.RUnlock()
			return len(*p.readingsTableDataMapBinding) + 1, 10
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("---fdaec17c-c0fc-4a04-982e-31a08a0bb776---")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			p.tableDataLock.RLock()
			defer p.tableDataLock.RUnlock()

			label := o.(*widget.Label)
			switch i.Row {
			case 0:
				label.TextStyle = fyne.TextStyle{Bold: true}
				switch i.Col {

				case 0:
					label.SetText("Id")
				case 1:
					label.SetText("Device Name")
				case 2:
					label.SetText("Resource Name")
				case 3:
					label.SetText("Profile Name")
				case 4:
					label.SetText("Value Type")
				case 5:
					label.SetText("Value")
				case 6:
					label.SetText("Binary Value")
				case 7:
					label.SetText("Media Type")
				case 8:
					label.SetText("Origin")
				case 9:
					label.SetText("Created")
					label.TextStyle = fyne.TextStyle{Bold: true}
				default:
					label.SetText("")
				}
			default:
				label.TextStyle = fyne.TextStyle{Bold: false}

				dm := *p.readingsTableDataMapBinding
				if dm == nil || i.Row > len(dm) {
					label.SetText("")
					break
				}

				row := dm[i.Row-1]

				switch i.Col {
				case 0:
					id, _ := row.GetItem("Id")
					o.(*widget.Label).Bind(id.(binding.String))
				case 1:
					eventName, _ := row.GetItem("DeviceName")
					o.(*widget.Label).Bind(eventName.(binding.String))
				case 2:
					resourceName, _ := row.GetItem("ResourceName")
					o.(*widget.Label).Bind(resourceName.(binding.String))
				case 3:
					profileName, _ := row.GetItem("ProfileName")
					o.(*widget.Label).Bind(profileName.(binding.String))
				case 4:
					valueType, _ := row.GetItem("ValueType")
					o.(*widget.Label).Bind(valueType.(binding.String))
				case 5:
					value, _ := row.GetItem("Value")
					o.(*widget.Label).Bind(value.(binding.String))
				case 6:
					binaryValue, _ := row.GetItem("BinaryValue")
					o.(*widget.Label).Bind(binaryValue.(binding.String))

					//o.(*widget.Label).Bind(binaryValue.(binding.String))
				case 7:
					mediaType, _ := row.GetItem("MediaType")
					o.(*widget.Label).Bind(mediaType.(binding.String))
				case 8:
					origin, _ := row.GetItem("Origin")
					v, _ := origin.(binding.Int).Get()
					o.(*widget.Label).SetText(time.Unix(0, int64(v)).String())
				case 9:
					created, _ := row.GetItem("Created")
					v, _ := created.(binding.Int).Get()
					txt := time.Unix(0, int64(v)).String()
					if v == 0 {
						txt = ""
					}
					o.(*widget.Label).SetText(txt)
				default:
					label.SetText("")
				}

			}

		},
	)

	//BinaryValue can be smaller
	t.SetColumnWidth(6, 110)
	//MediaType can be smaller
	t.SetColumnWidth(7, 100)

	return t
}

func (p *dataPageHandler) renderEventsTable() *widget.Table {
	t := widget.NewTable(
		func() (int, int) {
			p.tableDataLock.RLock()
			defer p.tableDataLock.RUnlock()
			return len(*p.eventsTableDataMapBinding) + 1, 7
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("---fdaec17c-c0fc-4a04-982e-31a08a0bb776---")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			p.tableDataLock.RLock()
			defer p.tableDataLock.RUnlock()

			label := o.(*widget.Label)
			switch i.Row {
			case 0:
				label.TextStyle = fyne.TextStyle{Bold: true}
				switch i.Col {

				case 0:
					label.SetText("Id")
				case 1:
					label.SetText("Device Name")
				case 2:
					label.SetText("Profile Name")
				case 3:
					label.SetText("Origin")
				case 4:
					label.SetText("Readings")
				case 5:
					label.SetText("Tags")
				case 6:
					label.SetText("Created")
				default:
					label.SetText("")
				}
			default:
				label.TextStyle = fyne.TextStyle{Bold: false}

				dm := *p.eventsTableDataMapBinding
				if dm == nil || i.Row > len(dm) {
					label.SetText("")
					break
				}

				row := dm[i.Row-1]

				switch i.Col {
				case 0:
					id, _ := row.GetItem("Id")
					o.(*widget.Label).Bind(id.(binding.String))
				case 1:
					eventName, _ := row.GetItem("DeviceName")
					o.(*widget.Label).Bind(eventName.(binding.String))
				case 2:
					profileName, _ := row.GetItem("ProfileName")
					o.(*widget.Label).Bind(profileName.(binding.String))
				// case 3:
				// 	created, _ := row.GetItem("Created")
				// 	v, _ := created.(binding.Int).Get()
				// 	o.(*widget.Label).SetText(time.Unix(0, int64(v)).String())
				case 3:
					origin, _ := row.GetItem("Origin")
					v, _ := origin.(binding.Int).Get()
					o.(*widget.Label).SetText(time.Unix(0, int64(v)).String())
				case 4:
					readings, _ := row.GetItem("ReadingsCount")
					v, _ := readings.(binding.Int).Get()
					o.(*widget.Label).SetText(fmt.Sprintf("%d", v))
				case 5:
					tags, _ := row.GetItem("Tags")
					v, _ := tags.(binding.String).Get()
					o.(*widget.Label).SetText(v)
				case 6:
					created, _ := row.GetItem("Created")
					v, _ := created.(binding.Int).Get()
					txt := time.Unix(0, int64(v)).String()
					if v == 0 {
						txt = ""
					}
					o.(*widget.Label).SetText(txt)
				default:
					label.SetText("")
				}

			}

		},
	)

	//Readings can be smaller
	t.SetColumnWidth(4, 95)

	return t
}

func (p *dataPageHandler) setBufferUsageBindingByDataType(dataType string) {
	switch dataType {
	case config.DataTypeEvents:
		p.bufferUsageBinding = binding.BindFloat(config.Float(float64(p.appState.GetDB().GetTotalEventsCount())))
	case config.DataTypeReadings:
		p.bufferUsageBinding = binding.BindFloat(config.Float(float64(p.appState.GetDB().GetTotalReadingsCount())))
	}
}

func (p *dataPageHandler) updateBufferUsageBindingByDataType(currentDataType string) {

	if currentDataType == "" || p.bufferUsageBinding == nil || p.bufferProgress == nil {
		return
	}
	p.appState.RLock()
	defer p.appState.RUnlock()
	log.Debugf("updateBufferUsageBindingByDataType for %v", currentDataType)

	switch currentDataType {
	case config.DataTypeEvents:
		p.bufferUsageBinding.Set(float64(p.appState.GetDB().GetTotalEventsCount()))
	case config.DataTypeReadings:
		p.bufferUsageBinding.Set(float64(p.appState.GetDB().GetTotalReadingsCount()))
	}
	p.bufferProgress.Bind(p.bufferUsageBinding)

}

func (p *dataPageHandler) updateStatusByDataType(currentDataType string) {

	if currentDataType == "" || p.statusText == nil {
		return
	}
	p.appState.RLock()
	defer p.appState.RUnlock()
	log.Debugf("updateStatusByDataType for %v", currentDataType)

	rowCount := 0
	recordType := ""

	switch currentDataType {
	case config.DataTypeEvents:
		rowCount = int(p.appState.GetDB().GetEventsCount())
		recordType = "events"
	case config.DataTypeReadings:
		rowCount = int(p.appState.GetDB().GetReadingsCount())
		recordType = "readings"
	}

	filter := p.appState.GetDataPageSearch()
	txt := fmt.Sprintf("Last %v %v", rowCount, recordType)
	if filter != nil && *filter != "" {
		txt = txt + fmt.Sprintf(" matching \"%v\" (case-insensitive)", *filter)
	}
	p.statusText.SetText(txt)

}

func (p *dataPageHandler) updateTableByDataType(currentDataType string) {

	if currentDataType == "" {
		return
	}

	db := p.appState.GetDB()
	log.Debugf("updating datatable for %v", currentDataType)

	sortAsc := fyne.CurrentApp().Preferences().BoolWithFallback(config.PrefEventsTableSortOrderAscending, config.DefaultEventsTableSortOrderAscending)

	if currentDataType == config.DataTypeEvents {
		p.tableDataLock.Lock()
		defer p.tableDataLock.Unlock()
		p.eventsTableDataMapBinding = &[]binding.DataMap{}

		events := db.GetEvents()

		evts := make([]dtos.Event, len(events))
		copy(evts, events)

		if !sortAsc {
			log.Debugf("sorting events desc")
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
			*p.eventsTableDataMapBinding = append(*p.eventsTableDataMapBinding, binding.BindStruct(&r))
		}
	} else if currentDataType == config.DataTypeReadings {
		p.tableDataLock.Lock()
		defer p.tableDataLock.Unlock()
		p.readingsTableDataMapBinding = &[]binding.DataMap{}

		readings := db.GetReadings()

		rdngs := make([]dtos.BaseReading, len(readings))
		copy(rdngs, readings)

		if !p.sortAsc {
			log.Debugf("sorting readings desc")
			sort.Slice(rdngs, func(i, j int) bool {
				return rdngs[i].Origin > rdngs[j].Origin
			})
		}

		for _, row := range rdngs {
			readingJson, _ := json.MarshalIndent(row, "", "    ")
			r := readingRow{
				Id:           row.Id,
				Created:      row.Created,
				Origin:       row.Origin,
				DeviceName:   row.DeviceName,
				ProfileName:  row.ProfileName,
				ResourceName: row.ResourceName,
				ValueType:    row.ValueType,
				BinaryValue:  string(row.BinaryValue),
				MediaType:    row.MediaType,
				Value:        row.Value,
				Json:         string(readingJson),
			}
			*p.readingsTableDataMapBinding = append(*p.readingsTableDataMapBinding, binding.BindStruct(&r))
		}
	}

}

func (p *dataPageHandler) OnEventReceived(event dtos.Event) {

	if p.appState.GetConnectionState() != services.ClientConnected {
		return
	}

	p.updateTableByDataType(p.dataType.Selected)
	p.updateBufferUsageBindingByDataType(p.dataType.Selected)
	p.updateStatusByDataType(p.dataType.Selected)
	if p.table != nil {
		p.table.Refresh()
	}

}
