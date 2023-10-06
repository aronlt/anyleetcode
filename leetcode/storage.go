package leetcode

import (
	"embed"
	"encoding/json"
	"path/filepath"
	"time"

	"anyleetcode/common"

	"github.com/aronlt/toolkit/terror"
	"github.com/aronlt/toolkit/tio"
	"github.com/aronlt/toolkit/tutils"
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

func (s *Storage) LoadResultFiles() ([]string, error) {
	_, files, err := tio.ReadDir(common.GetResultPath())
	return files, err
}

func (s *Storage) LoadResultContent(path string) ([]*HyperLink, error) {
	content, err := tio.ReadFile(filepath.Join(common.GetResultPath(), path))
	if err != nil {
		err = terror.Wrap(err, "call ReadFile fail, filepath:"+path)
		return nil, err
	}
	result := make([]*HyperLink, 0)
	err = json.Unmarshal(content, &result)
	if err != nil {
		err = terror.Wrap(err, "call Unmarshal fail, filepath:"+path)
		return nil, err
	}
	return result, nil
}

func (s *Storage) StoreResult(links []*HyperLink, path string) error {
	content, err := json.MarshalIndent(links, "", " ")
	if err != nil {
		panic(err)
	}
	if path == "" {
		t := tutils.TimeToString(time.Now(), "2006-01-02_15:04:05")
		path = t + ".json"
	}
	_, err = tio.WriteFile(filepath.Join(common.GetResultPath(), path), content, false)
	if err != nil {
		err = terror.Wrap(err, "call WriteFile fail")
		return err
	}
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
		tio.Mkdirs(common.GetResultPath())
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
