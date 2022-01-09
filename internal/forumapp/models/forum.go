package models

type Forum struct {
	Id      int64  `json:"id,omitempty"`
	Title   string `json:"title"`
	User    string `json:"user"`
	Slug    string `json:"slug"`
	Posts   int64  `json:"posts,omitempty"`
	Threads int32  `json:"threads,omitempty"`
}
