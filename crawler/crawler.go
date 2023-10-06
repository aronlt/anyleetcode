package main

import (
	"encoding/json"
	"time"

	"anyleetcode/common"
	"anyleetcode/leetcode"
	"anyleetcode/utils"

	"github.com/aronlt/toolkit/ds"
	"github.com/aronlt/toolkit/terror"
	"github.com/dustyRAIN/leetcode-api-go/leetcodeapi"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

func main() {
	problems, err := QueryAllProblems()
	if err != nil {
		panic(err)
	}
	err = utils.WriteToFile(problems, common.DataFilePath)
	if err != nil {
		panic(err)
	}
	err = StoreDifficulty(problems)
	if err != nil {
		panic(err)
	}

	err = StoreTags(problems)
	if err != nil {
		panic(err)
	}
}

func QueryAllProblems() ([]*leetcode.Problem, error) {
	result := make([]*leetcode.Problem, 0)
	problems := make([]leetcodeapi.Problem, 0)
	for {
		allProblemList, err := leetcodeapi.GetAllProblems(len(problems), 50)
		if err != nil {
			err = terror.Wrap(err, "call leetcodeapi.GetAllProblems fail")
			return result, err
		}
		if len(allProblemList.Problems) == 0 {
			logrus.Infof("get empty leetcode problem list")
			break
		}
		problems = append(problems, allProblemList.Problems...)
		logrus.Infof("get leetcode problem success, nums:%d, total:%d, now:%d", len(allProblemList.Problems), allProblemList.Total, len(result))
		time.Sleep(5 * time.Second)
	}
	ds.SliceIter(problems, func(a []leetcodeapi.Problem, i int) {
		topics := make([]string, 0)
		ds.SliceIter(a[i].TopicTags, func(b []leetcodeapi.TopicTag, j int) {
			topics = append(topics, b[j].Slug)
		})
		status := make(map[string]interface{})
		err := json.Unmarshal([]byte(a[i].Stats), &status)
		if err != nil {
			logrus.WithError(err).Errorf("call Unmarshal fail")
			return
		}
		acceptedRaw := cast.ToInt64(status["totalAcceptedRaw"])
		submissionRaw := cast.ToInt64(status["totalSubmissionRaw"])

		lcProblem := &leetcode.Problem{
			AcRate:          int(a[i].AcRate),
			Difficulty:      a[i].Difficulty,
			SubmissionCount: submissionRaw,
			AcceptCount:     acceptedRaw,
			Title:           a[i].Title,
			Url:             "https://leetcode.cn/problems/" + a[i].TitleSlug + "/",
			TopicTags:       topics,
		}
		result = append(result, lcProblem)
	})
	return result, nil
}

func StoreDifficulty(problems []*leetcode.Problem) error {
	set := ds.NewSet[string]()
	for _, problem := range problems {
		set.Insert(problem.Difficulty)
	}
	d := set.Keys()
	err := utils.WriteToFile(d, common.DifficultyFilePath)
	return err
}

func StoreTags(problems []*leetcode.Problem) error {
	tagSet := ds.NewSet[string]()
	for _, problem := range problems {
		tagSet.InsertN(problem.TopicTags...)
	}
	tags := tagSet.Keys()
	err := utils.WriteToFile(tags, common.TagsFilePath)
	return err
}
