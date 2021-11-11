package pages

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/deblasis/edgex-foundry-datamonitor/config"
	"github.com/deblasis/edgex-foundry-datamonitor/state"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos"
)

func parseURL(urlStr string) *url.URL {
	link, err := url.Parse(urlStr)
	if err != nil {
		fyne.LogError("Could not parse URL", err)
	}

	return link
}

func stateWithData(data binding.DataMap) *widget.Form {
	keys := data.Keys()
	items := make([]*widget.FormItem, len(keys))
	for i, k := range keys {
		data, err := data.GetItem(k)
		if err != nil {
			items[i] = widget.NewFormItem(k, widget.NewLabel(err.Error()))
		}
		items[i] = widget.NewFormItem(k, createItem(data))
	}

	return widget.NewForm(items...)
}

func createItem(v binding.DataItem) fyne.CanvasObject {
	switch val := v.(type) {
	case binding.Float:
		return widget.NewLabelWithData(binding.FloatToString(val))
	case binding.Int:
		return widget.NewLabelWithData(binding.IntToString(val))
	case binding.String:
		return widget.NewLabelWithData(val)
	default:
		return widget.NewLabel(fmt.Sprintf("%T", val))
	}
}

func homeScreen(w fyne.Window, appManager *state.AppManager) fyne.CanvasObject {
	// logo := canvas.NewImageFromResource(data.FyneScene)
	// logo.FillMode = canvas.ImageFillContain
	// if fyne.CurrentDevice().IsMobile() {
	// 	logo.SetMinSize(fyne.NewSize(171, 125))
	// } else {
	// 	logo.SetMinSize(fyne.NewSize(228, 167))
	// }

	redisHost, redisPort := appManager.GetRedisHostPort()

	var contentContainer *fyne.Container

	connectedContent := container.NewVBox()

	connectingContent := container.NewCenter(container.NewVBox(
		container.NewHBox(widget.NewProgressBarInfinite()),
	))

	disconnectedContent := container.NewCenter(container.NewVBox(
		widget.NewCard("You are currently disconnected from EdgeX Foundry",
			fmt.Sprintf("Would you like to connect to %v:%d?", redisHost, redisPort),
			container.NewCenter(
				widget.NewButtonWithIcon("Connect", theme.LoginIcon(), func() {
					if err := appManager.Connect(); err != nil {
						uerr := errors.New(fmt.Sprintf("Cannot connect\n%s", err))
						dialog.ShowError(uerr, w)
						//TODO: log this
					}
					appManager.Refresh()
				}),
			),
		),
	))

	// stateStruct := struct {
	// 	TotalNumberEvents           int
	// 	TotalNumberReadings         int
	// 	EventsPerSecondLastMinute   float64
	// 	ReadingsPerSecondLastMinute float64

	// 	// currently fyne doesn't support automatic comparison for slices
	// 	// I might push that feature upstream, meanwhile I am using a string
	// 	// to hold the bound data and then process it while rendering
	// 	LastEvents string
	// }{}

	//boundState := binding.BindStruct(&stateStruct)

	ep := appManager.GetEventProcessor()

	table := renderEventsTable(ep.LastEvents.Get(), false)
	tableContainer := container.NewMax(table)

	totalNumberEventsBinding := binding.BindInt(config.Int(ep.TotalNumberEvents))
	totalNumberReadingsBinding := binding.BindInt(config.Int(ep.TotalNumberReadings))

	eventsPerSecondLastMinute := binding.BindFloat(config.Float(ep.EventsPerSecondLastMinute))
	readingsPerSecondLastMinute := binding.BindFloat(config.Float(ep.ReadingsPerSecondLastMinute))

	dashboardStats := container.NewGridWithRows(2,
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
	)

	//form := stateWithData(boundState)

	//form.Hide()

	dashboardStats.Hide()
	tableContainer.Hide()

	switch appManager.GetConnectionState() {
	case state.Connected:
		contentContainer = connectedContent
		dashboardStats.Show()
		tableContainer = container.NewMax(table)
	case state.Connecting:
		contentContainer = connectingContent
		dashboardStats.Hide()
		tableContainer.Hide()
	case state.Disconnected:
		contentContainer = disconnectedContent
		dashboardStats.Hide()
		tableContainer = container.NewMax()
	}

	home := container.NewMax(
		container.NewGridWithRows(2,
			container.NewHBox(
				contentContainer,
				dashboardStats,
			),
			tableContainer,
		),
	)

	go func() {
		for {
			time.Sleep(100 * time.Millisecond)
			if appManager.GetConnectionState() != state.Connected {
				continue
			}
			//log.Printf("refreshing UI: events %v\n", ep.TotalNumberEvents)
			totalNumberEventsBinding.Set(ep.TotalNumberEvents)
			totalNumberReadingsBinding.Set(ep.TotalNumberReadings)
			eventsPerSecondLastMinute.Set(ep.EventsPerSecondLastMinute)
			readingsPerSecondLastMinute.Set(ep.ReadingsPerSecondLastMinute)
			table = renderEventsTable(ep.LastEvents.Get(), false)

			//dashboardStats.Refresh()
			// home.Objects[0].(*fyne.Container).Objects[1] = table
			if len(tableContainer.Objects) > 0 {
				tableContainer.Objects[0] = table
			}
			//home.Refresh()

		}
	}()

	return home

}

func renderEventsTable(events []*dtos.Event, sortAsc bool) *widget.Table {

	// the slice is fifo, we reverse it so that the first element is the most recent
	evts := events
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
		// case 0:

		// 	id := row + 1
		// 	if !sortAsc {
		// 		id = 5 - id + 1
		// 	}

		// 	label.SetText(fmt.Sprintf("%d", id))
		case 0:
			label.SetText(event.DeviceName)
		case 1:
			label.SetText(time.Unix(0, event.Origin).String())
		default:
			label.SetText("")
		}

	}

	t := widget.NewTable(
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
				case 1:
					label.SetText("Origin Timestamp")
				default:
					label.SetText("")
				}

			default:
				renderCell(id.Row-1, id.Col, label)
			}

		})
	// t.SetColumnWidth(0, 34)
	t.SetColumnWidth(0, 350)
	t.SetColumnWidth(1, 350)

	return t

}
