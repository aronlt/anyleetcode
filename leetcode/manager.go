package leetcode

import (
	"os"
	"path/filepath"
	"sort"

	"anyleetcode/common"

	"github.com/aronlt/toolkit/ds"
	"github.com/aronlt/toolkit/terror"
)

type Manager struct {
	storage *Storage
}

func NewManager() *Manager {
	return &Manager{storage: NewStorage()}
}

func (m *Manager) RemoveResult(path string) error {
	return os.Remove(filepath.Join(common.GetResultPath(), path))
}

func (m *Manager) LoadTags() ([]string, error) {
	_, tags, _, err := m.storage.Load()
	if err != nil {
		return nil, terror.Wrap(err, "call storage.Load fail")
	}
	sort.Slice(tags, func(i, j int) bool {
		return tags[i] < tags[j]
	})
	return tags, nil
}
func (m *Manager) LoadDifficulty() ([]string, error) {
	_, _, difficulty, err := m.storage.Load()
	if err != nil {
		return nil, terror.Wrap(err, "call storage.Load fail")
	}
	sort.Slice(difficulty, func(i, j int) bool {
		return difficulty[i] < difficulty[j]
	})
	return difficulty, nil
}

func (m *Manager) LoadResultFiles() ([]string, error) {
	files, err := m.storage.LoadResultFiles()
	if err != nil {
		err = terror.Wrap(err, "call LoadResultFiles fail")
		return nil, err
	}
	filenames := make([]string, 0, len(files))
	for _, file := range files {
		name := filepath.Base(file)
		filenames = append(filenames, name)
	}
	return filenames, nil
}

func (m *Manager) LoadResult(path string) ([]*HyperLink, error) {
	return m.storage.LoadResultContent(path)
}

func (m *Manager) StoreResult(links []*HyperLink, path string) error {
	return m.storage.StoreResult(links, path)
}

func (m *Manager) Query(cond *SearchCond) ([]*Problem, error) {
	problems, _, _, err := m.storage.Load()
	if err != nil {
		err = terror.Wrap(err, "call storage.Load fail")
		return nil, err
	}
	if cond.SubmissionCountRank != 0 {
		sort.Slice(problems, func(i, j int) bool {
			return problems[i].SubmissionCount > problems[j].SubmissionCount
		})
		size := int(float32(cond.SubmissionCountRank) / float32(100) * float32(len(problems)))
		problems = problems[:size]
	}
	result := ds.SliceGetFilter(problems, func(i int) bool {
		if cond.AcRate != 0 {
			if problems[i].AcRate <= cond.AcRate {
				return false
			}
		}
		if len(cond.Difficulty) != 0 {
			if !ds.SliceInclude(cond.Difficulty, problems[i].Difficulty) {
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
	count := ds.SliceMinUnpack(len(result), cond.Count)
	ds.SliceOpShuffle(result)
	return result[:count], nil
}
