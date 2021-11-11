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
	"github.com/deblasis/edgex-foundry-datamonitor/services"
)

// Page defines the data structure for a tutorial
type Page struct {
	Title, Intro string
	View         func(w fyne.Window, appMgr *services.AppManager) fyne.CanvasObject
}

var (
	// Tutorials defines the metadata for each tutorial
	Pages = map[string]Page{
		"home":     {"Home", "", homeScreen},
		"data":     {"Data", "", dataScreen},
		"settings": {"Settings", "", settingsScreen},
		// "canvas": {"Canvas",
		// 	"See the canvas capabilities.",
		// 	canvasScreen,
		// },
		// "animations": {"Animations",
		// 	"See how to animate components.",
		// 	makeAnimationScreen,
		// },
		// "icons": {"Theme Icons",
		// 	"Browse the embedded icons.",
		// 	iconScreen,
		// },
		// "containers": {"Containers",
		// 	"Containers group other widgets and canvas objects, organising according to their layout.\n" +
		// 		"Standard containers are illustrated in this section, but developers can also provide custom " +
		// 		"layouts using the fyne.NewContainerWithLayout() constructor.",
		// 	containerScreen,
		// },
		// "apptabs": {"AppTabs",
		// 	"A container to help divide up an application into functional areas.",
		// 	makeAppTabsTab,
		// },
		// "border": {"Border",
		// 	"A container that positions items around a central content.",
		// 	makeBorderLayout,
		// },
		// "box": {"Box",
		// 	"A container arranges items in horizontal or vertical list.",
		// 	makeBoxLayout,
		// },
		// "center": {"Center",
		// 	"A container to that centers child elements.",
		// 	makeCenterLayout,
		// },
		// "doctabs": {"DocTabs",
		// 	"A container to display a single document from a set of many.",
		// 	makeDocTabsTab,
		// },
		// "grid": {"Grid",
		// 	"A container that arranges all items in a grid.",
		// 	makeGridLayout,
		// },
		// "split": {"Split",
		// 	"A split container divides the container in two pieces that the user can resize.",
		// 	makeSplitTab,
		// },
		// "scroll": {"Scroll",
		// 	"A container that provides scrolling for it's content.",
		// 	makeScrollTab,
		// },
		// "widgets": {"Widgets",
		// 	"In this section you can see the features available in the toolkit widget set.\n" +
		// 		"Expand the tree on the left to browse the individual tutorial elements.",
		// 	widgetScreen,
		// },
		// "accordion": {"Accordion",
		// 	"Expand or collapse content panels.",
		// 	makeAccordionTab,
		// },
		// "button": {"Button",
		// 	"Simple widget for user tap handling.",
		// 	makeButtonTab,
		// },
		// "card": {"Card",
		// 	"Group content and widgets.",
		// 	makeCardTab,
		// },
		// "entry": {"Entry",
		// 	"Different ways to use the entry widget.",
		// 	makeEntryTab,
		// },
		// "form": {"Form",
		// 	"Gathering input widgets for data submission.",
		// 	makeFormTab,
		// },
		// "input": {"Input",
		// 	"A collection of widgets for user input.",
		// 	makeInputTab,
		// },
		// "text": {"Text",
		// 	"Text handling widgets.",
		// 	makeTextTab,
		// },
		// "toolbar": {"Toolbar",
		// 	"A row of shortcut icons for common tasks.",
		// 	makeToolbarTab,
		// },
		// "progress": {"Progress",
		// 	"Show duration or the need to wait for a task.",
		// 	makeProgressTab,
		// },
		// "collections": {"Collections",
		// 	"Collection widgets provide an efficient way to present lots of content.\n" +
		// 		"The List, Table, and Tree provide a cache and re-use mechanism that make it possible to scroll through thousands of elements.\n" +
		// 		"Use this for large data sets or for collections that can expand as users scroll.",
		// 	collectionScreen,
		// },
		// "list": {"List",
		// 	"A vertical arrangement of cached elements with the same styling.",
		// 	makeListTab,
		// },
		// "table": {"Table",
		// 	"A two dimensional cached collection of cells.",
		// 	makeTableTab,
		// },
		// "tree": {"Tree",
		// 	"A tree based arrangement of cached elements with the same styling.",
		// 	makeTreeTab,
		// },
		// "windows": {"Windows",
		// 	"Window function demo.",
		// 	windowScreen,
		// },
		// "binding": {"Data Binding",
		// 	"Connecting widgets to a data source.",
		// 	bindingScreen},
		// "advanced": {"Advanced",
		// 	"Debug and advanced information.",
		// 	advancedScreen,
		// },
	}

	//PageIndex  defines how our tutorials should be laid out in the index tree
	PageIndex = map[string][]string{
		"": {"home", "data", "settings"}, //"canvas", "animations", "icons", "widgets", "collections", "containers", "windows", "binding", "advanced",

		// "collections": {"list", "table", "tree"},
		// "containers":  {"apptabs", "border", "box", "center", "doctabs", "grid", "scroll", "split"},
		// "widgets":     {"accordion", "button", "card", "entry", "form", "input", "progress", "text", "toolbar"},
	}
)
