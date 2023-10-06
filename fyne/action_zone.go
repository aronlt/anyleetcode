package fyne

import (
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"anyleetcode/leetcode"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/aronlt/toolkit/ds"
)

func (a *App) clearDisplay() {
	a.doneDisplay.Objects = a.doneDisplay.Objects[:0]
	a.undoDisplay.Objects = a.undoDisplay.Objects[:0]
	a.doneProblemMap = make(map[string]*widget.Hyperlink)
	a.undoProblemMap = make(map[string]*widget.Hyperlink)
}
func (a *App) refreshDisplay() {
	a.doneDisplay.Refresh()
	a.undoDisplay.Refresh()
}

func (a *App) newOnTapFunc(url string, link *widget.Hyperlink) func() {
	return func() {
		if _, ok := a.doneProblemMap[url]; ok {
			delete(a.doneProblemMap, url)
			a.undoProblemMap[url] = link
		} else {
			delete(a.undoProblemMap, url)
			a.doneProblemMap[url] = link
		}
		addDisplay := func(problemMap map[string]*widget.Hyperlink, display *fyne.Container) {
			display.Objects = display.Objects[:0]
			links := make([]*widget.Hyperlink, 0)
			for _, v := range problemMap {
				links = append(links, v)
			}
			getIndex := func(text string) int {
				index := text[:strings.Index(text, ".")]
				num, err := strconv.Atoi(index)
				if err != nil {
					panic(err)
				}
				return num
			}
			sort.Slice(links, func(i, j int) bool {
				indexa := getIndex(links[i].Text)
				indexb := getIndex(links[j].Text)
				return indexa < indexb
			})
			ds.SliceIter(links, func(links []*widget.Hyperlink, i int) {
				display.Objects = append(display.Objects, links[i])
			})
		}
		addDisplay(a.undoProblemMap, a.undoDisplay)
		addDisplay(a.doneProblemMap, a.doneDisplay)

		err := fyne.CurrentApp().OpenURL(link.URL)
		if err != nil {
			fyne.LogError("Failed to open url", err)
		}
	}
}

func (a *App) NewGenButton() *widget.Button {
	button := widget.NewButton("生成题目列表", func() {
		result, err := a.lcApi.Query(&leetcode.SearchCond{
			Difficulty:          a.sDiff.Keys(),
			TopicTags:           a.sTags.Keys(),
			AcRate:              a.rate,
			SubmissionCountRank: a.submitCountRank,
			Count:               10,
		})
		if err != nil {
			panic(err)
		}

		a.clearDisplay()
		a.saveFileSelect.Selected = ""
		if len(result) != 0 {
			sort.Slice(result, func(i, j int) bool {
				if result[i].AcRate == result[j].AcRate {
					return result[i].SubmissionCount > result[j].SubmissionCount
				}
				return result[i].AcRate > result[j].AcRate
			})

			ds.SliceIter(result, func(r []*leetcode.Problem, i int) {
				u, _ := url.Parse(r[i].Url)
				title := fmt.Sprintf("%d. title:%s;  ac_rate:%d; total_submit:%d; topics:%s", i+1, r[i].Title, r[i].AcRate, r[i].SubmissionCount, strings.Join(r[i].TopicTags, ","))
				link := widget.NewHyperlink(title, u)
				link.OnTapped = a.newOnTapFunc(r[i].Url, link)
				a.undoDisplay.Add(link)
			})
		}
		a.refreshDisplay()
	})
	return button
}

func (a *App) NewLoadButton() *widget.Button {
	button := widget.NewButton("加载文件", func() {
		if a.saveFileSelect.Selected == "" {
			return
		}
		result, err := a.lcApi.LoadResult(a.saveFileSelect.Selected)
		if err != nil {
			panic(err)
		}
		a.clearDisplay()
		ds.SliceIter(result, func(b []*leetcode.HyperLink, i int) {
			u, _ := url.Parse(b[i].Link)
			link := widget.NewHyperlink(b[i].Text, u)
			link.OnTapped = a.newOnTapFunc(b[i].Link, link)
			if b[i].Type == leetcode.Done {
				a.doneDisplay.Objects = append(a.doneDisplay.Objects, link)
				a.doneProblemMap[b[i].Link] = link
			} else {
				a.undoDisplay.Objects = append(a.undoDisplay.Objects, link)
				a.undoProblemMap[b[i].Link] = link
			}
		})
		a.refreshDisplay()
		a.dialog.Show()
	})
	return button
}

func (a *App) NewDeleteButton() *widget.Button {
	button := widget.NewButton("删除文件", func() {
		if a.saveFileSelect.Selected == "" {
			return
		}
		err := a.lcApi.RemoveResult(a.saveFileSelect.Selected)
		if err != nil {
			panic(err)
		}
		a.saveFileSelect.Selected = ""
		a.ReloadSetSaveFileSelect()
		a.clearDisplay()
		a.refreshDisplay()
		a.dialog.Show()
	})
	return button
}

func (a *App) NewSaveButton() *widget.Button {
	button := widget.NewButton("保存文件", func() {
		if len(a.undoDisplay.Objects) == 0 && len(a.doneDisplay.Objects) == 0 {
			return
		}
		links := make([]*leetcode.HyperLink, 0, len(a.doneDisplay.Objects))
		for _, obj := range a.doneDisplay.Objects {
			link := obj.(*widget.Hyperlink)
			text := link.Text
			u := link.URL.String()
			hyperLink := &leetcode.HyperLink{
				Text: text,
				Link: u,
				Type: leetcode.Done,
			}
			links = append(links, hyperLink)
		}
		for _, obj := range a.undoDisplay.Objects {
			link := obj.(*widget.Hyperlink)
			text := link.Text
			u := link.URL.String()
			hyperLink := &leetcode.HyperLink{
				Text: text,
				Link: u,
				Type: leetcode.Undo,
			}
			links = append(links, hyperLink)
		}
		var err error
		err = a.lcApi.StoreResult(links, a.saveFileSelect.Selected)
		if err != nil {
			panic(err)
		}
		a.ReloadSetSaveFileSelect()
		a.dialog.Show()
	})
	return button
}

func (a *App) NewActionZone() *fyne.Container {
	genButton := a.NewGenButton()
	loadButton := a.NewLoadButton()
	saveButton := a.NewSaveButton()
	deleteButton := a.NewDeleteButton()
	return container.NewHBox(genButton, loadButton, saveButton, deleteButton)
}
