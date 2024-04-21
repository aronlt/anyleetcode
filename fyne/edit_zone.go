package fyne

import (
	"anyleetcode/leetcode"
	"fmt"
	"github.com/aronlt/toolkit/ds"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (a *App) tagContainer(name string, tags []string, tagSet ds.BuiltinSet[string]) *fyne.Container {
	tagLabel := widget.NewLabel(a.SetToString(tagSet))
	tagZone := container.NewHBox(
		widget.NewLabel(name),
		widget.NewSelect(tags, func(tag string) {
			if tagSet.Has(tag) {
				tagSet.Delete(tag)
			} else {
				tagSet.Insert(tag)
			}
			tagLabel.SetText(a.SetToString(tagSet))
		}))

	selectedTagsZone := container.NewHBox(
		widget.NewLabel(fmt.Sprintf("选择%s", name)),
		tagLabel,
	)

	return container.NewHBox(tagZone, selectedTagsZone)
}

func (a *App) NewDiffZone(diffs []string) *fyne.Container {
	return a.tagContainer("选择难度:", diffs, a.sDiff)
}

func (a *App) NewExcludeTagZone(tags []string) *fyne.Container {
	return a.tagContainer("排除标签:", tags, a.sExcludeTags)
}

func (a *App) NewTagZone(tags []string) *fyne.Container {
	return a.tagContainer("标签:", tags, a.sTags)
}

func (a *App) NewProblemStatus() *fyne.Container {
	status := []string{"全部", "仅未做", "仅完成"}
	statusLabel := widget.NewLabel("")
	statusZone := container.NewHBox(
		widget.NewLabel("题目状态:"),
		widget.NewSelect(status, func(s string) {
			index := ds.SliceIncludeIndex(status, s)
			statusLabel.SetText(s)
			a.problemStatus = leetcode.ProblemStatus(index + 1)
		}))
	selectedStatusZone := container.NewHBox(
		widget.NewLabel("选择状态:"),
		statusLabel,
	)
	return container.NewHBox(statusZone, selectedStatusZone)
}

func (a *App) rateContainer(name string, value *int64) *fyne.Container {
	rates := []string{"0", "10", "20", "30", "40", "50", "60", "70", "80", "90"}
	rateLabel := widget.NewLabel("")
	rateZone := container.NewHBox(
		widget.NewLabel(name),
		widget.NewSelect(rates, func(s string) {
			*value, _ = strconv.ParseInt(s, 10, 64)
			rateLabel.SetText(fmt.Sprintf("%s%%", s))
		}))
	selectedRateLabel := container.NewHBox(
		widget.NewLabel(fmt.Sprintf("选择%s", name)),
		rateLabel)
	return container.NewHBox(rateZone, selectedRateLabel)
}

func (a *App) NewRateZone() *fyne.Container {
	return a.rateContainer("通过率:", &a.rate)
}

func (a *App) NewRankZone() *fyne.Container {
	return a.rateContainer("提交率:", &a.submitCountRank)
}

func (a *App) NewCookieZone() *fyne.Container {
	input := widget.NewEntry()
	input.OnCursorChanged = func() {
		a.cookie = input.Text
	}
	return container.NewVBox(
		widget.NewLabel("输入LeetCode Cookie:"),
		input,
	)
}

func (a *App) NewEditZone() *fyne.Container {
	tags := a.lcApi.LoadTags()
	difficulty := a.lcApi.LoadDifficulty()
	diffZone := a.NewDiffZone(difficulty)
	tagZone := a.NewTagZone(tags)
	excludeTagZone := a.NewExcludeTagZone(tags)
	problemStatusZone := a.NewProblemStatus()
	rateZone := a.NewRateZone()
	submitCountZone := a.NewRankZone()
	actionZone := a.NewActionZone()
	cookieZone := a.NewCookieZone()
	zone := container.NewVBox(
		tagZone,
		excludeTagZone,
		diffZone,
		rateZone,
		problemStatusZone,
		submitCountZone,
		cookieZone,
		actionZone,
	)
	return zone
}
