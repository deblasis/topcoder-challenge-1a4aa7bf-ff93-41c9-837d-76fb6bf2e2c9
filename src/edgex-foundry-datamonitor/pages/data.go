package pages

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/deblasis/edgex-foundry-datamonitor/services"
)

func dataScreen(win fyne.Window, appState *services.AppManager) fyne.CanvasObject {
	return container.NewMax(
		container.NewVBox(
			container.NewCenter(
				container.NewHBox(
					widget.NewLabelWithStyle("TODO", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
				)),
		))
}
