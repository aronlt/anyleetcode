package leetcode

type Api interface {
	Query(cond *SearchCond) ([]*Problem, error)
	LoadTags() ([]string, error)
	LoadDifficulty() ([]string, error)
	LoadResultFiles() ([]string, error)
	LoadResult(path string) ([]*HyperLink, error)
	StoreResult(links []*HyperLink, path string) error
	RemoveResult(path string) error
}

func NewApi() Api {
	return NewManager()
}
