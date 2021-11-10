package pages

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/deblasis/edgex-foundry-datamonitor/internal/state"
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
	case binding.Bool:
		return widget.NewCheckWithData("", val)
	case binding.Float:
		s := widget.NewSliderWithData(0, 1, val)
		s.Step = 0.01
		return s
	case binding.Int:
		return widget.NewEntryWithData(binding.IntToString(val))
	case binding.String:
		return widget.NewEntryWithData(val)
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

	connectedContent := func() *fyne.Container {

		// stateStruct := struct {
		// 	TotalNumberEvents   int
		// 	TotalNumberReadings int
		// }{}

		// boundState := binding.BindStruct(&stateStruct)

		// go func() {
		// 	stateStruct.TotalNumberEvents = appManager.GetEventProcessor().TotalNumberEvents
		// 	stateStruct.TotalNumberReadings = appManager.GetEventProcessor().TotalNumberReadings
		// 	time.Sleep(500 * time.Millisecond)
		// }()
		// form := stateWithData(boundState)

		// totalEvents := binding.BindInt(&appManager.GetEventProcessor().TotalNumberEvents)
		// totalReadings := binding.BindInt(&appManager.GetEventProcessor().TotalNumberReadings)

		return container.NewVBox(
			container.NewCenter(
				//form,
				widget.NewSeparator(),
				widget.NewButtonWithIcon("Disconnect", theme.LogoutIcon(), func() {
					if err := appManager.Disconnect(); err != nil {
						uerr := errors.New(fmt.Sprintf("Cannot disconnect\n%s", err))
						dialog.ShowError(uerr, w)
						//TODO: log this
					}
					cnt, drawFn := appManager.GetCurrentContainer()
					drawFn(cnt)
				}),
			),
		)

		// left := widget.NewMultiLineEntry()
		// left.Wrapping = fyne.TextWrapWord
		// left.SetText("Long text is looooooooooooooong")
		// right := container.NewVSplit(
		// 	widget.NewLabel("Label"),
		// 	widget.NewButton("Button", func() { fmt.Println("button tapped!") }),
		// )

		// content := container.NewVBox(
		// 	makeTableTab(w, appManager),
		// 	container.NewHSplit(container.NewVScroll(left), right),
		// )
		// return container.NewCenter(
		// 	container.NewVBox(content,
		// 		container.NewCenter(

		// 			widget.NewButtonWithIcon("Disconnect", theme.LogoutIcon(), func() {
		// 				if err := appManager.Disconnect(); err != nil {
		// 					uerr := errors.New(fmt.Sprintf("Cannot disconnect\n%s", err))
		// 					dialog.ShowError(uerr, w)
		// 					//TODO: log this
		// 				}
		// 				cnt, drawFn := appManager.GetCurrentContainer()
		// 				drawFn(cnt)
		// 			}),
		// 		)))

	}()

	connectingContent := container.NewCenter(container.NewVBox(
		container.NewHBox(widget.NewProgressBarInfinite()),
	))

	disconnectedContent := container.NewVBox(
		widget.NewCard("You are currently disconnected from EdgeX Foundry",
			fmt.Sprintf("Would you like to connect to %v:%d?", redisHost, redisPort),
			container.NewCenter(
				widget.NewButtonWithIcon("Connect", theme.LoginIcon(), func() {
					if err := appManager.Connect(); err != nil {
						uerr := errors.New(fmt.Sprintf("Cannot connect\n%s", err))
						dialog.ShowError(uerr, w)
						//TODO: log this
					}
					cnt, drawFn := appManager.GetCurrentContainer()
					drawFn(cnt)
				}),
			),
		),
	)

	stateStruct := struct {
		TotalNumberEvents           int
		TotalNumberReadings         int
		EventsPerSecondLastMinute   float64
		ReadingsPerSecondLastMinute float64

		LastDeviceNames string
	}{}

	boundState := binding.BindStruct(&stateStruct)

	ep := appManager.GetEventProcessor()
	form := stateWithData(boundState)
	go func() {
		for {
			stateStruct.TotalNumberEvents = ep.TotalNumberEvents
			stateStruct.TotalNumberReadings = ep.TotalNumberReadings

			stateStruct.EventsPerSecondLastMinute = ep.EventsPerSecondLastMinute
			stateStruct.ReadingsPerSecondLastMinute = ep.ReadingsPerSecondLastMinute

			stateStruct.LastDeviceNames = strings.Join(ep.LastDeviceNames.Get(), "|")

			time.Sleep(500 * time.Millisecond)
			boundState.Reload()
			form.Refresh()
		}
	}()

	form.Hide()
	switch appManager.GetConnectionState() {
	case state.Connected:
		contentContainer = connectedContent
		form.Show()
	case state.Connecting:
		contentContainer = connectingContent
		form.Hide()
	case state.Disconnected:
		contentContainer = disconnectedContent
		form.Hide()
	}

	return container.NewVBox(container.NewCenter(
		container.NewHBox(
			contentContainer,
			form,
		),
	))
}
