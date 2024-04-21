package fyne

import (
	"anyleetcode/common"
	"embed"
	"github.com/aronlt/toolkit/tio"
	"github.com/golang/freetype/truetype"
	"os"
	"sort"
	"strings"

	"anyleetcode/leetcode"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"github.com/aronlt/toolkit/ds"
)

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

//go:embed font
var f embed.FS

type App struct {
	lcApi           leetcode.Api
	UndoDisplay     *Display
	DoneDisplay     *Display
	window          fyne.Window
	sTags           ds.BuiltinSet[string]
	sExcludeTags    ds.BuiltinSet[string]
	sDiff           ds.BuiltinSet[string]
	problemStatus   leetcode.ProblemStatus
	cookie          string
	rate            int64
	submitCountRank int64
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
}

func NewApp() *App {
	api := leetcode.NewApi()
	window := app.New().NewWindow("Leetcode Helper")
	window.Resize(fyne.Size{
		Width:  1000,
		Height: 1000,
	})

	return &App{window: window, lcApi: api,
		sTags:         ds.NewSet[string](),
		sExcludeTags:  ds.NewSet[string](),
		sDiff:         ds.NewSet[string](),
		problemStatus: leetcode.All,
		DoneDisplay:   NewDisplay(),
		UndoDisplay:   NewDisplay(),
	}
}

func (a *App) Start() {
	a.window.ShowAndRun()
}
