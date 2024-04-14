package fyne

import (
	"embed"
	"sort"
	"strings"

	"anyleetcode/leetcode"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/aronlt/toolkit/ds"
)

//go:embed font
var f embed.FS

type App struct {
	lcApi           leetcode.Api
	doneDisplay     *fyne.Container
	undoDisplay     *fyne.Container
	doneProblemMap  map[string]*widget.Hyperlink
	undoProblemMap  map[string]*widget.Hyperlink
	dialog          dialog.Dialog
	window          fyne.Window
	sTags           ds.BuiltinSet[string]
	sDiff           ds.BuiltinSet[string]
	cookie          string
	rate            int
	submitCountRank int
}

func (a *App) SetToString(set ds.BuiltinSet[string]) string {
	builder := strings.Builder{}
	keys := set.Keys()
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	for i, key := range keys {
		builder.WriteString(key)
		if i != len(keys)-1 {
			builder.WriteString(",")
		}
	}
	return builder.String()
}

func (a *App) Init() {
	editZone := a.NewEditZone()
	displayZone := a.NewDisplayZone()
	main := container.NewVBox(
		editZone,
		displayZone,
	)
	a.window.SetContent(main)
	a.clearDisplay()
}

func NewApp() *App {
	api := leetcode.NewApi()
	window := app.New().NewWindow("Leetcode Helper")
	window.Resize(fyne.Size{
		Width:  1000,
		Height: 1000,
	})
	successDialog := dialog.NewInformation("Result", "Success", window)
	return &App{window: window, lcApi: api,
		sTags:          ds.NewSet[string](),
		sDiff:          ds.NewSet[string](),
		doneProblemMap: make(map[string]*widget.Hyperlink),
		undoProblemMap: make(map[string]*widget.Hyperlink),
		dialog:         successDialog,
		doneDisplay:    container.NewVBox(),
		undoDisplay:    container.NewVBox(),
	}
}

func (a *App) Start() {
	a.window.ShowAndRun()
}
