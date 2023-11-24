package fyne

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (a *App) NewDiffZone(diffs []string) *fyne.Container {
	diffLabel := widget.NewLabel(a.SetToString(a.sDiff))
	diffZone := container.NewHBox(
		widget.NewLabel("难度:"),
		widget.NewSelect(diffs, func(diff string) {
			if a.sDiff.Has(diff) {
				a.sDiff.Delete(diff)
			} else {
				a.sDiff.Insert(diff)
			}
			diffLabel.SetText(a.SetToString(a.sDiff))
		}))
	selectedDiffZone := container.NewHBox(
		widget.NewLabel("选择难度:"),
		diffLabel,
	)

	return container.NewVBox(diffZone, selectedDiffZone)
}

func (a *App) NewTagZone(tags []string) *fyne.Container {
	tagLabel := widget.NewLabel(a.SetToString(a.sTags))
	tagZone := container.NewHBox(
		widget.NewLabel("标签:"),
		widget.NewSelect(tags, func(tag string) {
			if a.sTags.Has(tag) {
				a.sTags.Delete(tag)
			} else {
				a.sTags.Insert(tag)
			}
			tagLabel.SetText(a.SetToString(a.sTags))
		}))

	selectedTagsZone := container.NewHBox(
		widget.NewLabel("选择标签:"),
		tagLabel,
	)

	return container.NewVBox(tagZone, selectedTagsZone)
}

func (a *App) NewRateZone() *fyne.Container {
	rates := []string{"0", "10", "20", "30", "40", "50", "60", "70", "80", "90"}
	return container.NewHBox(
		widget.NewLabel("通过率:"),
		widget.NewSelect(rates, func(s string) {
			a.rate, _ = strconv.Atoi(s)
		}))
}

func (a *App) NewRankZone() *fyne.Container {
	rates := []string{"0", "10", "20", "30", "40", "50", "60", "70", "80", "90"}
	return container.NewHBox(
		widget.NewLabel("提交量排名:"),
		widget.NewSelect(rates, func(s string) {
			a.submitCountRank, _ = strconv.Atoi(s)
		}))
}

func (a *App) NewEditZone() *fyne.Container {
	tags, err := a.lcApi.LoadTags()
	if err != nil {
		panic(err)
	}

	difficulty, err := a.lcApi.LoadDifficulty()
	if err != nil {
		panic(err)
	}

	diffZone := a.NewDiffZone(difficulty)
	tagZone := a.NewTagZone(tags)
	rateZone := a.NewRateZone()
	submitCountZone := a.NewRankZone()
	actionZone := a.NewActionZone()
	zone := container.NewVBox(
		tagZone,
		diffZone,
		rateZone,
		submitCountZone,
		actionZone,
	)
	return zone
}
