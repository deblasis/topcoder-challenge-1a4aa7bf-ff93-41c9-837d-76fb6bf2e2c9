package pages

import (
	"errors"
	"fmt"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
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
		return container.NewVBox(
			container.NewCenter(
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

	return container.NewVBox(container.NewCenter(
		container.NewHBox(
			contentContainer,
		),
	))
}
