package leetcode

type Api interface {
	Query(cond *SearchCond) ([]*Problem, error)
	LoadTags() []string
	LoadDifficulty() []string
}

func NewApi() Api {
	return NewManager()
}
