package leetcode

import (
	"anyleetcode/common"
	"embed"
	"encoding/json"

	"github.com/aronlt/toolkit/terror"
	"github.com/sirupsen/logrus"
)

//go:embed data
var f embed.FS

type Storage struct {
	problems   []*Problem
	tags       []string
	difficulty []string
}

func NewStorage() *Storage {
	return &Storage{problems: make([]*Problem, 0)}
}

func LoadFromFile[T any](filepath string) ([]T, error) {
	var data []T
	content, err := f.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(content, &data)
	if err != nil {
		err = terror.Wrap(err, "call Unmarshal fail")
		return nil, err
	}
	return data, err
}

func (s *Storage) LoadTags() error {
	tags, err := LoadFromFile[string](common.TagsFilePath)
	if err != nil {
		return terror.Wrap(err, "call LoadFromFile fail")
	}
	s.tags = tags
	return nil
}

func (s *Storage) LoadDifficulty() error {
	d, err := LoadFromFile[string](common.DifficultyFilePath)
	if err != nil {
		return terror.Wrap(err, "call LoadFromFile fail")
	}
	s.difficulty = d
	return nil
}

func (s *Storage) Load() ([]*Problem, []string, []string, error) {
	if len(s.problems) != 0 {
		return s.problems, s.tags, s.difficulty, nil
	}

	var problems []*Problem
	problems, err := LoadFromFile[*Problem](common.DataFilePath)
	if err != nil {
		logrus.WithError(err).Errorf("call ReadFile fail")
		err = terror.Wrap(err, "call LoadFromFile fail")
		return nil, nil, nil, err
	} else {
		s.problems = problems
		err = s.LoadTags()
		if err != nil {
			return nil, nil, nil, err
		}
		err = s.LoadDifficulty()
		if err != nil {
			return nil, nil, nil, err
		}
	}
	return s.problems, s.tags, s.difficulty, nil
}
