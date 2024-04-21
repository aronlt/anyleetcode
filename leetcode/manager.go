package leetcode

import (
	"anyleetcode/common"
	"anyleetcode/utils"
	"encoding/json"
	"fmt"
	"github.com/aronlt/toolkit/tio"
	"github.com/aronlt/toolkit/tjson"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/aronlt/toolkit/ds"
	"github.com/aronlt/toolkit/terror"
)

type Manager struct {
	problems   []*Problem
	tags       []string
	difficulty []string
	doneSlugs  ds.BuiltinSet[string]
}

func NewManager() *Manager {
	storage := NewStorage()
	problems, tags, difficulty, err := storage.Load()
	sort.Slice(tags, func(i, j int) bool {
		return tags[i] < tags[j]
	})
	sort.Slice(difficulty, func(i, j int) bool {
		return difficulty[i] < difficulty[j]
	})
	if err != nil {
		err = terror.Wrap(err, "call storage.Load fail")
		panic(err)
	}
	done := LoadDoneSlugs()
	return &Manager{problems: problems, tags: tags, difficulty: difficulty, doneSlugs: done}
}

func (m *Manager) LoadDone(cookie string) error {
	skip := 0
	limit := 100
	for {
		args := fmt.Sprintf(`{"query":"\n    query problemsetQuestionList($categorySlug: String, $limit: Int, $skip: Int, $filters: QuestionListFilterInput) {\n  problemsetQuestionList(\n    categorySlug: $categorySlug\n    limit: $limit\n    skip: $skip\n    filters: $filters\n  ) {\n    hasMore\n    total\n    questions {\n      acRate\n      difficulty\n      freqBar\n      frontendQuestionId\n      isFavor\n      paidOnly\n      solutionNum\n      status\n      title\n      titleCn\n      titleSlug\n      topicTags {\n        name\n        nameTranslated\n        id\n        slug\n      }\n      extra {\n        hasVideoSolution\n        topCompanyTags {\n          imgUrl\n          slug\n          numSubscribed\n        }\n      }\n    }\n  }\n}\n    ","variables":{"categorySlug":"all-code-essentials","skip":%d,"limit":%d,"filters":{"status":"AC"}},"operationName":"problemsetQuestionList"}`, skip, limit)
		resp, err := utils.DoQuery([]byte(args), cookie)
		if err != nil {
			logrus.WithError(err).Errorf("call DoQuery fail")
			return err
		}
		type Slug struct {
			Slug string `json:"titleSlug"`
		}

		raw, err := tjson.GetRawMessage(resp, "data.problemsetQuestionList.questions")
		if err != nil {
			logrus.WithError(err).Errorf("call GetRawMessage fail, resp: %s", string(resp))
			return err
		}
		values := make([]*Slug, 0)
		err = json.Unmarshal(raw, &values)
		if err != nil {
			logrus.WithError(err).Errorf("call GetRawMessage fail, resp: %s", string(resp))
			return err
		}
		for _, v := range values {
			m.doneSlugs.Insert(v.Slug)
		}
		hasMore, err := tjson.GetBool(resp, "data.problemsetQuestionList.hasMore")
		if err != nil {
			logrus.WithError(err).Errorf("call GetRawMessage fail, resp: %s", string(resp))
			return err
		}
		if hasMore {
			skip += limit
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}
	WriteDoneSlugs(m.doneSlugs)
	return nil
}

func (m *Manager) LoadTags() []string {
	return m.tags
}
func (m *Manager) LoadDifficulty() []string {
	return m.difficulty
}

func WriteDoneSlugs(set ds.BuiltinSet[string]) {
	slugs := ds.SetToSlice(set)
	fp := filepath.Join(common.GetCachePath(), "done.txt")
	content := strings.Join(slugs, "\n")
	_, err := tio.WriteFile(fp, []byte(content), false)
	if err != nil {
		logrus.WithError(err).Errorf("call WriteFile fail, fp:%s", fp)
	}
	logrus.Infof("call WriteFile success, fp:%s", fp)
}
func LoadDoneSlugs() ds.BuiltinSet[string] {
	result := ds.BuiltinSet[string]{}
	fp := filepath.Join(common.GetCachePath(), "done.txt")
	lines, err := tio.ReadLines(fp)
	if err != nil {
		logrus.Errorf("call ReadLines fail, fp:%s", fp)
		return result
	}
	for _, line := range lines {
		result.Insert(string(line))
	}
	logrus.Infof("load done slugs success, fp:%s", fp)
	return result
}

func (m *Manager) Query(cond *SearchCond) ([]*Problem, error) {
	problems := m.problems
	for i := 0; i < 3; i++ {
		ds.SliceOpShuffle(problems)
	}
	result := ds.SliceGetFilter(problems, func(i int) bool {
		problems[i].HasAC = m.doneSlugs.Has(problems[i].Slug)
		if len(cond.Difficulty) != 0 {
			if !ds.SliceInclude(cond.Difficulty, problems[i].Difficulty) {
				return false
			}
		}
		if cond.AcRate != 0 {
			if problems[i].AcRate <= cond.AcRate {
				return false
			}
		}
		if cond.SubmissionCountRank != 0 {
			if problems[i].SubmissionCount <= cond.SubmissionCountRank {
				return false
			}
		}
		if len(cond.TopicTags) != 0 {
			setA := ds.SetFromSlice(cond.TopicTags)
			setB := ds.SetFromSlice(problems[i].TopicTags)
			if setA.Intersection(setB).IsEmpty() {
				return false
			}
		}
		if len(cond.ExcludeTopicTags) != 0 {
			setA := ds.SetFromSlice(cond.ExcludeTopicTags)
			setB := ds.SetFromSlice(problems[i].TopicTags)
			if setA.Intersection(setB).IsEmpty() == false {
				return false
			}
		}
		if cond.ProblemStatus != All {
			if cond.ProblemStatus == OnlyDone && problems[i].HasAC == false {
				return false
			}
			if cond.ProblemStatus == OnlyUndo && problems[i].HasAC == true {
				return false
			}
		}
		return true
	})

	count := ds.SliceMinUnpack(len(result), cond.Count)
	result = result[:count]
	return result, nil
}
