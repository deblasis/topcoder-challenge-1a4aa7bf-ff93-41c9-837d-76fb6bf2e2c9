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
	"strconv"
	"strings"
	"time"

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

func dataScreen(win fyne.Window, appManager *services.AppManager) fyne.CanvasObject {

	connectionState := appManager.GetConnectionState()
	redisHost, redisPort := appManager.GetRedisHostPort()

	disconnectedContent := container.NewCenter(container.NewVBox(
		widget.NewCard("You are currently disconnected from EdgeX Foundry",
			fmt.Sprintf("Would you like to connect to %v:%d?", redisHost, redisPort),
			container.NewCenter(
				widget.NewButtonWithIcon("Connect", theme.LoginIcon(), func() {
					if err := appManager.Connect(); err != nil {
						uerr := errors.New(fmt.Sprintf("Cannot connect\n%s", err))
						dialog.ShowError(uerr, win)
						log.Printf("cannot connect: %v", err)
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

	renderCell := func(row, col int, label *widget.Label) {

		label.SetText(fmt.Sprintf("(%v,%v)", row, col))

	}

	eventsTable := widget.NewTable(
		func() (int, int) { return 100001, 7 },
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
					label.SetText("Id")
					label.TextStyle = fyne.TextStyle{Bold: true}
				case 1:
					label.SetText("Device Name")
					label.TextStyle = fyne.TextStyle{Bold: true}
				case 2:
					label.SetText("Profile Name")
					label.TextStyle = fyne.TextStyle{Bold: true}
				case 3:
					label.SetText("Created")
					label.TextStyle = fyne.TextStyle{Bold: true}
				case 4:
					label.SetText("Origin")
					label.TextStyle = fyne.TextStyle{Bold: true}
				case 5:
					label.SetText("Tags")
					label.TextStyle = fyne.TextStyle{Bold: true}
				case 6:
					label.SetText("Readings Count")
					label.TextStyle = fyne.TextStyle{Bold: true}
				default:
					label.SetText("")
				}

			default:
				renderCell(id.Row-1, id.Col, label)
			}

		})

	eventsTable.SetColumnWidth(0, 350)
	eventsTable.SetColumnWidth(1, 350)
	eventsTable.SetColumnWidth(2, 350)
	eventsTable.SetColumnWidth(3, 350)
	eventsTable.SetColumnWidth(4, 350)
	eventsTable.SetColumnWidth(5, 350)
	eventsTable.SetColumnWidth(6, 350)

	//TODO bind this
	detail := widget.NewEntry()
	detail.Text = `{"property": "value", "propertyInt": 12 }
	{"property": "value", "propertyInt": 12 }{"property": "value", "propertyInt": 12 }{"property": "value", "propertyInt": 12 }
	{"property": "value", "propertyInt": 12 }{"property": "value", "propertyInt": 12 }{"property": "value", "propertyInt": 12 }
	{"property": "value", "propertyInt": 12 }{"property": "value", "propertyInt": 12 }
	{"property": "value", "propertyInt": 12 }{"property": "value", "propertyInt": 12 }
	{"property": "value", "propertyInt": 12 }{"property": "value", "propertyInt": 12 }
	{"property": "value", "propertyInt": 12 }{"property": "value", "propertyInt": 12 }
	{"property": "value", "propertyInt": 12 }{"property": "value", "propertyInt": 12 }
	{"property": "value", "propertyInt": 12 }{"property": "value", "propertyInt": 12 }
	{"property": "value", "propertyInt": 12 }{"property": "value", "propertyInt": 12 }
	{"property": "value", "propertyInt": 12 }{"property": "value", "propertyInt": 12 }`

	detail.Disabled()

	copyToClipboardBtn := widget.NewButtonWithIcon("Copy to clipboard", theme.ContentCopyIcon(), func() {
		fyne.Clipboard.SetContent(win.Clipboard(), detail.Text)
	})

	detailBox := container.NewBorder(
		container.NewBorder(nil, nil, nil, copyToClipboardBtn, widget.NewLabelWithStyle("Selected item", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})),
		nil, nil, nil,
		container.NewVScroll(container.NewMax(detail)),
	)
	dlg := dialog.NewCustom("Detail", "Close", detailBox, win)
	dlg.Resize(fyne.NewSize(800, 1000))

	eventsTable.OnSelected = func(id widget.TableCellID) {
		if id.Row == 0 {
			return
		}
		dlg.Show()
	}

	content := container.NewBorder(
		container.NewVBox(
			container.NewGridWithColumns(2,
				radioGroup,
				searchBox,
			),
			widget.NewSeparator(),
			bufferSizeContainer,
		),
		nil, nil, nil, eventsTable,
	)

	return content

}

func renderDataEventsTable(events []*dtos.Event, sortAsc bool) fyne.CanvasObject {

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

type dataPageHandler struct {
	appState *services.AppManager
	Key      widget.TreeNodeID

	dataType           *widget.RadioGroup
	search             *widget.Entry
	searchBtn          *widget.Button
	resetSearchBtn     *widget.Button
	applyBufferSizeBtn *widget.Button
	bufferSize         *widget.Entry
	bufferSizeBinding  binding.Int
	bufferUsageBinding binding.Float
	bufferProgress     *widget.ProgressBar
}

func (p *dataPageHandler) SetInitialState() {
	preferences := fyne.CurrentApp().Preferences()

	if p.dataType.Selected == "" {
		p.dataType.SetSelected(config.DataTypeEvents)
	}
	p.searchBtn.Disable()
	p.resetSearchBtn.Disable()
	p.applyBufferSizeBtn.Disable()

	currentBuffsize := preferences.IntWithFallback(config.PrefBufferSizeInDataPage, config.DefaultBufferSizeInDataPage)

	b := p.bufferSizeBinding
	b.Set(currentBuffsize)
	p.bufferSize.SetText(fmt.Sprintf("%d", currentBuffsize))

	p.bufferProgress.Max = float64(currentBuffsize)
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
		log.Printf("bufferSize is %v", *bufferSize)
		//p.bufferSize.SetText(fmt.Sprintf("%d", *bufferSize))
		//p.bufferSize.SetText("100")
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
		log.Printf("DO SEARCH %v \n", p.search.Text)
		p.resetSearchBtn.Enable()
	}

	b := p.bufferSizeBinding
	b.AddListener(binding.NewDataListener(func() {
		v, _ := b.Get()

		p.appState.SetDataPageBufferSize(v)

		//TODO handle underflow
		p.bufferProgress.Max = float64(v)
		p.bufferProgress.Refresh()

		log.Printf("bufferSizeBinding CHANGED to %v", v)
		//retriggering validation, updating the binding alone doesn't do it
	}))

	p.bufferSize.OnChanged = func(s string) {
		if p.bufferSize.Validate() != nil {
			p.applyBufferSizeBtn.Disable()
			return
		}
		log.Printf("bufferSize changed to %v", s)
		boundValue, _ := b.Get()

		if fmt.Sprintf("%d", boundValue) != s && s != "" {
			p.applyBufferSizeBtn.Enable()
		} else {
			p.applyBufferSizeBtn.Disable()
		}
	}

	p.applyBufferSizeBtn.OnTapped = func() {
		v, err := strconv.Atoi(p.bufferSize.Text)
		if err != nil {
			return
		}

		b.Set(v)
		p.bufferProgress.Max = float64(v)
		p.bufferProgress.Refresh()

		log.Println("applyBufferSizeBtn tapped")
	}

	p.resetSearchBtn.OnTapped = func() {
		p.appState.SetDataPageSearch("")
		p.search.Text = ""
		p.search.Refresh()
		p.searchBtn.Disable()
		log.Println("filter reset")
	}

	currentDataType := p.dataType.Selected
	p.setBufferUsageBindingByDataType(currentDataType)

	p.bufferUsageBinding.AddListener(binding.NewDataListener(func() {
		log.Println("updated bufferUsageBinding")
		p.bufferProgress.Refresh()
	}))
	p.bufferProgress.Bind(p.bufferUsageBinding)

	p.dataType.OnChanged = func(currentDataType string) {
		//change table
		//change bindings
		log.Printf("changed dataType to %v", currentDataType)
		p.setBufferUsageBindingByDataType(currentDataType)
	}

}

func (p *dataPageHandler) setBufferUsageBindingByDataType(dataType string) {
	switch dataType {
	case config.DataTypeEvents:
		p.bufferUsageBinding = binding.BindFloat(config.Float(float64(p.appState.GetDB().GetEventsCount())))
	case config.DataTypeReadings:
		p.bufferUsageBinding = binding.BindFloat(config.Float(float64(p.appState.GetDB().GetReadingsCount())))
	}
}

func (p *dataPageHandler) updateBufferUsageBindingByDataType(currentDataType string) {
	switch currentDataType {
	case config.DataTypeEvents:
		p.bufferUsageBinding.Set(float64(p.appState.GetDB().GetEventsCount()))
	case config.DataTypeReadings:
		p.bufferUsageBinding.Set(float64(p.appState.GetDB().GetReadingsCount()))
	}
}

func (p *dataPageHandler) OnEventReceived(event dtos.Event) {

	if p.appState.GetConnectionState() != services.ClientConnected {
		return
	}

	p.updateBufferUsageBindingByDataType(p.dataType.Selected)

}

func NewDataPageHandler(appState *services.AppManager) *dataPageHandler {

	p := &dataPageHandler{
		appState: appState,
		Key:      DataPageKey,
	}

	p.dataType = widget.NewRadioGroup([]string{config.DataTypeEvents, config.DataTypeReadings}, func(dataType string) {
		log.Printf("Selected %s", dataType)
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
