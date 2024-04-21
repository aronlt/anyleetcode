package fyne

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (a *App) NewDisplayZone() *fyne.Container {
	return container.NewVBox(widget.NewLabel("未访问"), a.UndoDisplay.Box, widget.NewLabel("已访问"), a.DoneDisplay.Box)
}
