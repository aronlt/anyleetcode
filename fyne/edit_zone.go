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
		widget.NewLabel("Difficulty:"),
		widget.NewSelect(diffs, func(diff string) {
			if a.sDiff.Has(diff) {
				a.sDiff.Delete(diff)
			} else {
				a.sDiff.Insert(diff)
			}
			diffLabel.SetText(a.SetToString(a.sDiff))
		}))
	selectedDiffZone := container.NewHBox(
		widget.NewLabel("Select Difficulty:"),
		diffLabel,
	)

	return container.NewVBox(diffZone, selectedDiffZone)
}

func (a *App) NewTagZone(tags []string) *fyne.Container {
	tagLabel := widget.NewLabel(a.SetToString(a.sTags))
	tagZone := container.NewHBox(
		widget.NewLabel("Tags:"),
		widget.NewSelect(tags, func(tag string) {
			if a.sTags.Has(tag) {
				a.sTags.Delete(tag)
			} else {
				a.sTags.Insert(tag)
			}
			tagLabel.SetText(a.SetToString(a.sTags))
		}))

	selectedTagsZone := container.NewHBox(
		widget.NewLabel("Select Tags:"),
		tagLabel,
	)

	return container.NewVBox(tagZone, selectedTagsZone)
}

func (a *App) NewRateZone() *fyne.Container {
	rates := []string{"0", "10", "20", "30", "40", "50", "60", "70", "80", "90"}
	return container.NewHBox(
		widget.NewLabel("AC Rate:"),
		widget.NewSelect(rates, func(s string) {
			a.rate, _ = strconv.Atoi(s)
		}))
}

func (a *App) NewRankZone() *fyne.Container {
	rates := []string{"0", "10", "20", "30", "40", "50", "60", "70", "80", "90"}
	return container.NewHBox(
		widget.NewLabel("Submission Rate:"),
		widget.NewSelect(rates, func(s string) {
			a.submitCountRank, _ = strconv.Atoi(s)
		}))
}

func (a *App) NewCookieZone() *fyne.Container {
	input := widget.NewEntry()
	input.OnCursorChanged = func() {
		a.cookie = input.Text
	}
	return container.NewVBox(
		widget.NewLabel("LeetCode Cookie:"),
		input,
	)
}

func (a *App) NewEditZone() *fyne.Container {
	tags := a.lcApi.LoadTags()
	difficulty := a.lcApi.LoadDifficulty()
	diffZone := a.NewDiffZone(difficulty)
	tagZone := a.NewTagZone(tags)
	rateZone := a.NewRateZone()
	submitCountZone := a.NewRankZone()
	actionZone := a.NewActionZone()
	cookieZone := a.NewCookieZone()
	zone := container.NewVBox(
		tagZone,
		diffZone,
		rateZone,
		submitCountZone,
		cookieZone,
		actionZone,
	)
	return zone
}
