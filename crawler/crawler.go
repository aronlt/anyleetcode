package main

import (
	"anyleetcode/common"
	"anyleetcode/leetcode"
	"anyleetcode/utils"

	"encoding/json"
	"fmt"
	"time"

	"github.com/aronlt/toolkit/ds"
	"github.com/aronlt/toolkit/terror"
	"github.com/aronlt/toolkit/tjson"
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

func QueryTopics(titleSlug string) ([]string, error) {
	args := fmt.Sprintf(`{"query":"\n    query singleQuestionTopicTags($titleSlug: String!) {\n  question(titleSlug: $titleSlug) {\n    topicTags {\n      name\n      slug\n      translatedName\n    }\n  }\n}\n    ","variables":{"titleSlug":"%s"},"operationName":"singleQuestionTopicTags"}`, titleSlug)
	resp, err := utils.DoQuery([]byte(args))
	if err != nil {
		return nil, err
	}

	/*
		{
		    "data": {
		        "question": {
		            "topicTags": [
		                {
		                    "name": "Array",
		                    "slug": "array",
		                    "translatedName": "数组"
		                },
		                {
		                    "name": "Hash Table",
		                    "slug": "hash-table",
		                    "translatedName": "哈希表"
		                }
		            ]
		        }
		    }
		}
	*/
	v, err := tjson.GetRawMessage(resp, "data.question.topicTags")
	if err != nil {
		logrus.WithError(err).Errorf("call tjson.GetString fail, resp:%s", string(resp))
		return nil, err
	}
	type M struct {
		Name           string `json:"name"`
		Slug           string `json:"slug"`
		TranslatedName string `json:"translatedName"`
	}
	result := make([]*M, 0)

	err = json.Unmarshal(v, &result)
	if err != nil {
		logrus.WithError(err).Errorf("call Unmarshal fail, v:%s", string(v))
		return nil, err
	}

	names := make([]string, 0, len(result))
	ds.SliceIterV2(result, func(i int) {
		names = append(names, result[i].TranslatedName)
	})

	return names, nil

}

func QueryTitle(titleSlug string) (string, error) {
	args := fmt.Sprintf(`{"query":"\n    query questionTranslations($titleSlug: String!) {\n  question(titleSlug: $titleSlug) {\n    translatedTitle\n  }\n}\n    ","variables":{"titleSlug":"%s"},"operationName":"questionTranslations"}`, titleSlug)
	resp, err := utils.DoQuery([]byte(args))
	if err != nil {
		return "", err
	}

	/*
		{"data":{"question":{"translatedTitle":"\u4e24\u6570\u4e4b\u548c"}}}
	*/
	v, err := tjson.GetString(resp, "data.question.translatedTitle")
	if err != nil {
		logrus.WithError(err).Errorf("call  tjson.GetString fail, resp:%s", string(resp))
		return "", err
	}
	return v, err
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
		logrus.Infof("get leetcode problem success, nums:%d, total:%d, now:%d", len(allProblemList.Problems), allProblemList.Total, len(problems))
		time.Sleep(5 * time.Second)
	}
	count := 0
	ds.SliceIter(problems, func(a []leetcodeapi.Problem, i int) {
		status := make(map[string]interface{})
		err := json.Unmarshal([]byte(a[i].Stats), &status)
		if err != nil {
			logrus.WithError(err).Errorf("call Unmarshal fail")
			return
		}
		acceptedRaw := cast.ToInt64(status["totalAcceptedRaw"])
		submissionRaw := cast.ToInt64(status["totalSubmissionRaw"])

		title, err := QueryTitle(a[i].TitleSlug)
		if err != nil {
			logrus.WithError(err).Errorf("call QueryTitle fail, title slug:%s", a[i].TitleSlug)
			return
		}
		logrus.Infof("query title success, index:%d, result:%s", i, title)
		topics, err := QueryTopics(a[i].TitleSlug)
		if err != nil {
			logrus.WithError(err).Errorf("call QueryTopics fail, title slug:%s", a[i].TitleSlug)
			return
		}
		logrus.Infof("query topics success, index:%d, result:%v", i, topics)
		count++
		if count%50 == 0 {
			time.Sleep(1 * time.Second)
		}
		lcProblem := &leetcode.Problem{
			AcRate:          int(a[i].AcRate),
			Difficulty:      a[i].Difficulty,
			SubmissionCount: submissionRaw,
			AcceptCount:     acceptedRaw,
			Title:           title,
			Slug:            a[i].TitleSlug,
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
