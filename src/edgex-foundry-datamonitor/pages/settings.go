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
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/widget"
	"github.com/deblasis/edgex-foundry-datamonitor/config"
	"github.com/deblasis/edgex-foundry-datamonitor/data"
	"github.com/deblasis/edgex-foundry-datamonitor/services"
)

// dialogScreen loads demos of the dialogs we support
func settingsScreen(win fyne.Window, appState *services.AppManager) fyne.CanvasObject {
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

	dataPageBufferSize := widget.NewEntry()
	dataPageBufferSize.SetPlaceHolder("* required")
	dataPageBufferSize.Validator = data.MinMaxValidator(config.MinBufferSize, config.MaxBufferSize, data.ErrInvalidBufferSize)

	//read from settings
	hostname.SetText(preferences.StringWithFallback(config.PrefRedisHost, config.RedisDefaultHost))

	port.SetText(fmt.Sprintf("%d", preferences.IntWithFallback(config.PrefRedisPort, config.RedisDefaultPort)))
	shouldConnectAutomatically.SetChecked(preferences.BoolWithFallback(config.PrefShouldConnectAtStartup, config.DefaultShouldConnectAtStartup))
	eventsSortedAscendingly.SetChecked(preferences.BoolWithFallback(config.PrefEventsTableSortOrderAscending, config.DefaultEventsTableSortOrderAscending))
	dataPageBufferSize.SetText(fmt.Sprintf("%d", preferences.IntWithFallback(config.PrefBufferSizeInDataPage, config.DefaultBufferSizeInDataPage)))

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
			{Text: "Initial buffer size in Data page", Widget: dataPageBufferSize},
		},
		OnSubmit: func() {
			log.Info("Settings form submitted")

			preferences.SetString(config.PrefRedisHost, strings.TrimSpace(hostname.Text))

			p, _ := strconv.Atoi(port.Text)
			preferences.SetInt(config.PrefRedisPort, p)

			preferences.SetBool(config.PrefShouldConnectAtStartup, shouldConnectAutomatically.Checked)
			preferences.SetBool(config.PrefEventsTableSortOrderAscending, eventsSortedAscendingly.Checked)

			a.SendNotification(&fyne.Notification{
				Title:   "EdgeX Redis Pub/Sub Connection Settings",
				Content: fmt.Sprintf("%v:%v", hostname.Text, port.Text),
			})
		},

		CancelText: "Reset defaults",
	}

	form.OnCancel = func() {
		hostname.Text = config.RedisDefaultHost
		port.Text = fmt.Sprintf("%d", config.RedisDefaultPort)

		shouldConnectAutomatically.SetChecked(config.DefaultShouldConnectAtStartup)
		eventsSortedAscendingly.SetChecked(config.DefaultEventsTableSortOrderAscending)

		hostname.Validate()
		port.Validate()
		form.Refresh()
		log.Info("Settings reset to default")
	}

	return container.NewMax(
		container.NewVBox(
			widget.NewLabelWithStyle("Please enter EdgeX Redis Pub/Sub Connection Settings", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			container.NewCenter(
				container.NewHBox(
					form,
				)),
		))

}
