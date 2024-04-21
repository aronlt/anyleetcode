package fyne

import (
	"fmt"
	"net/url"
	"sort"
	"strings"

	"anyleetcode/leetcode"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/aronlt/toolkit/ds"
)

func (a *App) clearDisplay() {
	a.UndoDisplay.Clean()
	a.DoneDisplay.Clean()
}

func (a *App) refreshDisplay() {
	a.DoneDisplay.Refresh()
	a.UndoDisplay.Refresh()
}

func (a *App) newOnTapFunc(url string, link *widget.Hyperlink) func() {
	return func() {
		if v, ok := a.DoneDisplay.ProblemMap[url]; ok {
			delete(a.DoneDisplay.ProblemMap, url)
			a.UndoDisplay.ProblemMap[url] = v
		} else if v, ok = a.UndoDisplay.ProblemMap[url]; ok {
			delete(a.UndoDisplay.ProblemMap, url)
			a.DoneDisplay.ProblemMap[url] = v
		}
		addDisplay := func(problemMap map[string]*ProblemUnit, display *fyne.Container) {
			display.RemoveAll()
			units := make([]*ProblemUnit, 0)
			for _, v := range problemMap {
				units = append(units, v)
			}

			sort.Slice(units, func(i, j int) bool {
				return units[i].Index < units[j].Index
			})
			ds.SliceIter(units, func(links []*ProblemUnit, i int) {
				a.displayProblemUnit(display, links[i])
			})
			display.Refresh()
		}
		addDisplay(a.UndoDisplay.ProblemMap, a.UndoDisplay.Box)
		addDisplay(a.DoneDisplay.ProblemMap, a.DoneDisplay.Box)
		a.refreshDisplay()

		err := fyne.CurrentApp().OpenURL(link.URL)
		if err != nil {
			fyne.LogError("Failed to open url", err)
		}
	}
}

func (a *App) NewCleanButton() *widget.Button {
	button := widget.NewButton("重置", func() {
		a.sTags.Clear()
		a.sExcludeTags.Clear()
		a.sDiff.Clear()
		a.rate = 0
		a.problemStatus = leetcode.All
		a.submitCountRank = 0
		a.clearDisplay()
		a.Init()
	})

	return button

}

func (a *App) NewGenButton() *widget.Button {
	button := widget.NewButton("生成题目", func() {
		result, err := a.lcApi.Query(&leetcode.SearchCond{
			Difficulty:          a.sDiff.Keys(),
			TopicTags:           a.sTags.Keys(),
			ExcludeTopicTags:    a.sExcludeTags.Keys(),
			ProblemStatus:       a.problemStatus,
			AcRate:              a.rate,
			Cookie:              a.cookie,
			SubmissionCountRank: a.submitCountRank,
			Count:               10,
		})
		if err != nil {
			panic(err)
		}

		a.clearDisplay()
		if len(result) != 0 {
			sort.Slice(result, func(i, j int) bool {
				if result[i].AcRate == result[j].AcRate {
					return result[i].SubmissionCount > result[j].SubmissionCount
				}
				return result[i].AcRate > result[j].AcRate
			})

			ds.SliceIter(result, func(r []*leetcode.Problem, i int) {
				u, _ := url.Parse(r[i].Url)
				ac := "否"
				if r[i].HasAC {
					ac = "是"
				}
				unit := &ProblemUnit{
					Index: i + 1,
					Title: fmt.Sprintf("%d. %s", i+1, r[i].Title),
					Done:  fmt.Sprintf("是否做过: %s", ac),
					Tags:  fmt.Sprintf("标签: %s", strings.Join(r[i].TopicTags, ",")),
				}
				link := widget.NewHyperlink(unit.Title, u)
				link.OnTapped = a.newOnTapFunc(r[i].Url, link)
				unit.Link = link

				a.displayProblemUnit(a.UndoDisplay.Box, unit)
				a.UndoDisplay.ProblemMap[r[i].Url] = unit
				a.UndoDisplay.Refresh()
			})
		}
		a.refreshDisplay()
	})
	return button
}

func (a *App) displayProblemUnit(display *fyne.Container, unit *ProblemUnit) {
	display.Add(container.NewHBox(unit.Link, widget.NewLabel(unit.Done), widget.NewLabel(unit.Tags)))
}

func (a *App) NewRefreshDoneSlugButton() *widget.Button {
	button := widget.NewButton("刷新完成题目", func() {
		if a.cookie != "" {
			err := a.lcApi.LoadDone(a.cookie)
			if err != nil {
				panic(err)
			}
		}
	})
	return button
}

func (a *App) NewActionZone() *fyne.Container {
	genButton := a.NewGenButton()
	refreshDoneSlugButton := a.NewRefreshDoneSlugButton()
	cleanButton := a.NewCleanButton()
	return container.NewHBox(genButton, cleanButton, refreshDoneSlugButton)
}
