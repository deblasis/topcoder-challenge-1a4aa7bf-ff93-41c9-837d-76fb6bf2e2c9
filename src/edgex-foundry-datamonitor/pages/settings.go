package pages

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/widget"
	"github.com/deblasis/edgex-foundry-datamonitor/config"
	"github.com/deblasis/edgex-foundry-datamonitor/data"
	"github.com/deblasis/edgex-foundry-datamonitor/state"
)

// dialogScreen loads demos of the dialogs we support
func settingsScreen(win fyne.Window, appState *state.AppManager) fyne.CanvasObject {
	a := fyne.CurrentApp()
	preferences := a.Preferences()

	hostname := widget.NewEntry()
	hostname.SetPlaceHolder(fmt.Sprintf("Insert Redis host (default: %v)", config.RedisDefaultHost))
	hostname.Validator = data.StringNotEmptyValidator

	port := widget.NewEntry()
	port.SetPlaceHolder(fmt.Sprintf("Insert Redis port (default: %v)", config.RedisDefaultPort))
	port.Validator = validation.NewRegexp(`\d`, "Must contain a number")

	shouldConnectAutomatically := widget.NewCheckWithData("Connect at startup", binding.NewBool())
	eventsSortedAscendingly := widget.NewCheckWithData("Sort events ascendingly", binding.NewBool())

	//read from settings
	hostname.SetText(preferences.StringWithFallback(config.PrefRedisHost, config.RedisDefaultHost))

	port.SetText(fmt.Sprintf("%d", preferences.IntWithFallback(config.PrefRedisPort, config.RedisDefaultPort)))
	shouldConnectAutomatically.SetChecked(preferences.BoolWithFallback(config.PrefShouldConnectAtStartup, config.DefaultShouldConnectAtStartup))
	eventsSortedAscendingly.SetChecked(preferences.BoolWithFallback(config.PrefEventsTableSortOrderAscending, config.DefaultEventsTableSortOrderAscending))

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Hostname", Widget: hostname, HintText: "EdgeX Redis Pub/Sub hostname"},
			{Text: "Port", Widget: port, HintText: "EdgeX Redis Pub/Sub port"},
			{
				Text:     "",
				Widget:   shouldConnectAutomatically,
				HintText: "",
			},
			{
				Text:     "",
				Widget:   eventsSortedAscendingly,
				HintText: "",
			},
		},
		OnCancel: func() {
			hostname.Text = config.RedisDefaultHost
			port.Text = fmt.Sprintf("%d", config.RedisDefaultPort)

			shouldConnectAutomatically.SetChecked(config.DefaultShouldConnectAtStartup)
			eventsSortedAscendingly.SetChecked(config.DefaultEventsTableSortOrderAscending)

			fmt.Println("Cancelled")
		},
		OnSubmit: func() {
			log.Println("Settings form submitted")

			preferences.SetString(config.PrefRedisHost, strings.TrimSpace(hostname.Text))

			p, _ := strconv.Atoi(port.Text)
			preferences.SetInt(config.PrefRedisPort, p)

			preferences.SetBool(config.PrefShouldConnectAtStartup, shouldConnectAutomatically.Checked)
			preferences.SetBool(config.PrefEventsTableSortOrderAscending, eventsSortedAscendingly.Checked)

			//preferences.SetString(preferenceCurrentTutorial, uid)
			a.SendNotification(&fyne.Notification{
				Title:   "EdgeX Redis Pub/Sub Connection Settings",
				Content: fmt.Sprintf("%v:%v", hostname.Text, port.Text),
			})
		},

		CancelText: "Reset defaults",
	}

	return container.NewCenter(
		container.NewVBox(
			widget.NewLabelWithStyle("Please enter EdgeX Redis Pub/Sub Connection Settings", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			container.NewMax(
				container.NewHBox(
					form,
				)),
		))

}
