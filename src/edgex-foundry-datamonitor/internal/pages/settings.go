package pages

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/deblasis/edgex-foundry-datamonitor/internal/config"
	"github.com/deblasis/edgex-foundry-datamonitor/internal/state"
)

// dialogScreen loads demos of the dialogs we support
func settingsScreen(win fyne.Window, appState *state.AppManager) fyne.CanvasObject {
	a := fyne.CurrentApp()
	hostname := widget.NewEntry()
	hostname.SetPlaceHolder(config.RedisDefaultHost)

	port := widget.NewEntry()
	port.SetPlaceHolder(fmt.Sprintf("%d", config.RedisDefaultPort))

	//hostname.Validator = validation.NewRegexp(`\w{1,}@\w{1,}\.\w{1,4}`, "not a valid email")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Hostname", Widget: hostname, HintText: "EdgeX Redis Pub/Sub hostname"},
			{Text: "Port", Widget: port, HintText: "EdgeX Redis Pub/Sub port"},
		},
		OnCancel: func() {
			hostname.Text = config.RedisDefaultHost
			port.Text = fmt.Sprintf("%d", config.RedisDefaultPort)
			fmt.Println("Cancelled")
		},
		OnSubmit: func() {
			fmt.Println("Form submitted")
			//a.Preferences().SetString(preferenceCurrentTutorial, uid)
			a.SendNotification(&fyne.Notification{
				Title:   "EdgeX Redis Pub/Sub Connection Settings",
				Content: fmt.Sprintf("%v:%v", hostname.Text, port.Text),
			})
		},

		CancelText: "Reset defaults",
	}

	return container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle("Please enter EdgeX Redis Pub/Sub Connection Settings", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		container.NewHBox(
			form,
		),
	))

}
