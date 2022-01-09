package models

type User struct {
	Id       int64  `json:"id,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	Fullname string `json:"fullname"`
	About    string `json:"about,omitempty"`
	Email    string `json:"email"`
}
