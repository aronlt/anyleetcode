package leetcode

type Api interface {
	Query(cond *SearchCond) ([]*Problem, error)
	LoadDone(cookie string) error
	LoadTags() []string
	LoadDifficulty() []string
}

func NewApi() Api {
	return NewManager()
}
