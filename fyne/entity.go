package fyne

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type ProblemUnit struct {
	Index int
	Link  *widget.Hyperlink
	Title string
	Done  string
	Tags  string
}

type Display struct {
	Box        *fyne.Container
	ProblemMap map[string]*ProblemUnit
}

func (d *Display) Clean() {
	d.ProblemMap = make(map[string]*ProblemUnit)
	d.Box.RemoveAll()
}

func (d *Display) Refresh() {
	d.Box.Refresh()
}

func NewDisplay() *Display {
	return &Display{
		Box:        container.NewVBox(),
		ProblemMap: make(map[string]*ProblemUnit),
	}
}
