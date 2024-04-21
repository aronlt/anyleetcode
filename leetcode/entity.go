package leetcode

type HyperLinkType int

const Done HyperLinkType = 1
const Undo HyperLinkType = 2

type ProblemStatus int

const All ProblemStatus = 1
const OnlyUndo ProblemStatus = 2
const OnlyDone ProblemStatus = 3

type Problem struct {
	AcRate          int64
	AcceptCount     int64
	SubmissionCount int64
	Difficulty      string
	Title           string
	Url             string
	TopicTags       []string
	HasAC           bool
	Slug            string
}

type HyperLink struct {
	Text string        `json:"text"`
	Link string        `json:"link"`
	Type HyperLinkType `json:"type"`
}

type SearchCond struct {
	AcRate              int64
	SubmissionCountRank int64
	Difficulty          []string
	TopicTags           []string
	ExcludeTopicTags    []string
	ProblemStatus       ProblemStatus
	Count               int
	Cookie              string
}
