package leetcode

type HyperLinkType int

const Done HyperLinkType = 1
const Undo HyperLinkType = 2

type Problem struct {
	AcRate          int
	AcceptCount     int64
	SubmissionCount int64
	Difficulty      string
	Title           string
	Url             string
	TopicTags       []string
	HasAC           bool
	CheckAC         bool
	Slug            string
}

type HyperLink struct {
	Text string        `json:"text"`
	Link string        `json:"link"`
	Type HyperLinkType `json:"type"`
}

type SearchCond struct {
	AcRate              int
	SubmissionCountRank int
	Difficulty          []string
	TopicTags           []string
	Count               int
	Cookie              string
}
