package fyne

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (a *App) NewDisplayZone() *fyne.Container {

	return container.NewVBox(widget.NewLabel("未完成"), a.undoDisplay, widget.NewLabel("已完成"), a.doneDisplay)
}
