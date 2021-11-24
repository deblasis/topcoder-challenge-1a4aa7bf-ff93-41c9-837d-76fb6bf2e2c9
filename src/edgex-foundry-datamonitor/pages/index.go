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
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/deblasis/edgex-foundry-datamonitor/services"
)

type Page struct {
	Title, Intro string
	View         func(w fyne.Window, appMgr *services.AppManager) fyne.CanvasObject
}

var (
	Pages = map[widget.TreeNodeID]Page{
		HomePageKey:     {Title: "Home", Intro: "", View: homeScreen},
		DataPageKey:     {Title: "Data", Intro: "", View: dataScreen},
		SettingsPageKey: {Title: "Settings", Intro: "", View: settingsScreen},
	}

	//PageIndex  defines how our pages should be laid out in the index tree
	PageIndex = map[widget.TreeNodeID][]widget.TreeNodeID{
		"": {HomePageKey, DataPageKey, SettingsPageKey},
	}
)

const (
	HomePageKey     widget.TreeNodeID = "home"
	DataPageKey     widget.TreeNodeID = "data"
	SettingsPageKey widget.TreeNodeID = "settings"
)
