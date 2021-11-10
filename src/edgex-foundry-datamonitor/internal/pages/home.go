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
		items[i] = widget.NewFormItem(k, createBoundItem(data))
	}

	return widget.NewForm(items...)
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

	switch appManager.GetConnectionState() {
	case state.Connected:
		contentContainer = connectedContent
	case state.Connecting:
		contentContainer = connectingContent
	case state.Disconnected:
		contentContainer = disconnectedContent
	}

	a := binding.NewInt()
	lbl := widget.NewLabelWithData(binding.IntToString(a))
	go func() {
		a.Set(time.Now().Nanosecond())
		time.Sleep(500 * time.Millisecond)
		lbl.Refresh()
	}()

	stateStruct := struct {
		TotalNumberEvents           int
		TotalNumberReadings         int
		EventsPerSecondLastMinute   float64
		ReadingsPerSecondLastMinute float64
	}{}

	boundState := binding.BindStruct(&stateStruct)

	ep := appManager.GetEventProcessor()
	form := stateWithData(boundState)
	go func() {
		for {
			stateStruct.TotalNumberEvents = ep.TotalNumberEvents
			stateStruct.TotalNumberReadings = ep.TotalNumberReadings
			time.Sleep(500 * time.Millisecond)
			boundState.Reload()
			form.Refresh()
		}
	}()

	return container.NewVBox(container.NewCenter(
		container.NewHBox(
			lbl,
			contentContainer,
			form,
		),
	))
}
