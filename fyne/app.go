package fyne

import (
	"embed"
	"os"
	"sort"
	"strings"

	"anyleetcode/common"
	"anyleetcode/leetcode"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/aronlt/toolkit/ds"
	"github.com/aronlt/toolkit/tio"
	"github.com/golang/freetype/truetype"
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
	rate            int
	submitCountRank int
}

func init() {
	fontData, err := f.ReadFile("font/simkai.ttf")
	if err != nil {
		panic(err)
	}
	_, err = truetype.Parse(fontData)
	if err != nil {
		panic(err)
	}
	tio.WriteFile(common.GetFontPath(), fontData, false)

	os.Setenv("FYNE_FONT", common.GetFontPath())
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
	window := app.New().NewWindow("力扣选题")
	window.Resize(fyne.Size{
		Width:  1000,
		Height: 1000,
	})
	successDialog := dialog.NewInformation("结果", "成功", window)
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
