package leetcode

import (
	"anyleetcode/utils"
	"encoding/json"
	"fmt"
	"github.com/aronlt/toolkit/tjson"
	"github.com/sirupsen/logrus"
	"sort"

	"github.com/aronlt/toolkit/ds"
	"github.com/aronlt/toolkit/terror"
)

type Manager struct {
	problems   []*Problem
	tags       []string
	difficulty []string
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
	return &Manager{problems: problems, tags: tags, difficulty: difficulty}
}

func (m *Manager) checkAC(cookie string, slug string) bool {
	args := fmt.Sprintf(`{"query":"\n    query submissionList($offset: Int!, $limit: Int!, $lastKey: String, $questionSlug: String!, $lang: String, $status: SubmissionStatusEnum) {\n  submissionList(\n    offset: $offset\n    limit: $limit\n    lastKey: $lastKey\n    questionSlug: $questionSlug\n    lang: $lang\n    status: $status\n  ) {\n    lastKey\n    hasNext\n    submissions {\n      id\n      title\n      status\n      statusDisplay\n      lang\n      langName: langVerboseName\n      runtime\n      timestamp\n      url\n      isPending\n      memory\n      submissionComment {\n        comment\n        flagType\n      }\n    }\n  }\n}\n    ","variables":{"questionSlug":"%s","offset":0,"limit":20,"lastKey":null,"status":null},"operationName":"submissionList"}`, slug)
	resp, err := utils.DoQuery([]byte(args), cookie)
	if err != nil {
		logrus.WithError(err).Errorf("call DoQuery fail")
		return false
	}
	type Submission struct {
		ID            string `json:"id"`
		Title         string `json:"title"`
		Status        string `json:"status"`
		StatusDisplay string `json:"statusDisplay"`
		Lang          string `json:"lang"`
		LangName      string `json:"langName"`
		Runtime       string `json:"runtime"`
		Timestamp     string `json:"timestamp"`
		URL           string `json:"url"`
		IsPending     string `json:"isPending"`
		Memory        string `json:"memory"`
	}
	raw, err := tjson.GetRawMessage(resp, "data.submissionList.submissions")
	if err != nil {
		logrus.WithError(err).Errorf("call GetRawMessage fail, resp: %s", string(resp))
		return false
	}
	values := make([]*Submission, 0)
	err = json.Unmarshal(raw, &values)
	if err != nil {
		logrus.WithError(err).Errorf("call GetRawMessage fail, resp: %s", string(resp))
		return false
	}
	for _, v := range values {
		if v.Status == "AC" {
			return true
		}
	}
	return false
}

func (m *Manager) LoadTags() []string {
	return m.tags
}
func (m *Manager) LoadDifficulty() []string {
	return m.difficulty
}

func (m *Manager) Query(cond *SearchCond) ([]*Problem, error) {
	problems := m.problems

	result := ds.SliceGetFilter(problems, func(i int) bool {
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
		if len(cond.TopicTags) != 0 {
			setA := ds.SetFromSlice(cond.TopicTags)
			setB := ds.SetFromSlice(problems[i].TopicTags)
			if setA.Intersection(setB).IsEmpty() {
				return false
			}
		}
		return true
	})
	if cond.SubmissionCountRank != 0 {
		sort.Slice(problems, func(i, j int) bool {
			return problems[i].SubmissionCount > problems[j].SubmissionCount
		})
		size := int(float32(cond.SubmissionCountRank) / float32(100) * float32(len(problems)))
		problems = problems[:size]
	}
	count := ds.SliceMinUnpack(len(result), cond.Count)
	ds.SliceOpShuffle(result)
	result = result[:count]
	if cond.Cookie != "" {
		ds.SliceIterV2(result, func(i int) {
			if result[i].CheckAC == false {
				result[i].HasAC = m.checkAC(cond.Cookie, result[i].Slug)
				result[i].CheckAC = true
			}
		})
	}
	return result, nil
}
