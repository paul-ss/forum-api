package domain

import "time"

type Thread struct {
	Description string `json:"description"`
	Id int32 `json:"id"`
	Title string `json:"title"`
	Author string `json:"author"`
	Forum string `json:"forum"`
	Message string `json:"message"`
	Votes int32 `json:"votes"`
	Slug string `json:"slug"`
	Created time.Time `json:"created"`
}

type ThreadUpdate struct {
	Description string `json:"description"`
	Title string `json:"title"`
	Message string `json:"message"`
}