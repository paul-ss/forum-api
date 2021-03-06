package domain

import "time"

type Thread struct {
	Id int32 `json:"id"`
	Title string `json:"title"`
	Author string `json:"author"`
	Forum string `json:"forum"`
	Message string `json:"message"`
	Votes int32 `json:"votes"`
	Slug string `json:"slug"`
	Created time.Time `json:"created"`
}

type ThreadCreate struct {
	Title string `json:"title"`
	Author string `json:"author"`
	Message string `json:"message"`
	Created time.Time `json:"created"`
	Slug string `json:"slug"`
}

type ThreadUpdate struct {
	Title *string `json:"title"`
	Message *string `json:"message"`
}