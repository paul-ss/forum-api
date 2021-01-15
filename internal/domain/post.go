package domain

import (
	"time"
)

type Post struct {
	Id        int64     `json:"id"`
	Parent    int64     `json:"parent"`
	Author    string    `json:"author"`
	Message   string    `json:"message"`
	IsEdited  bool      `json:"isEdited"`
	ForumSlug string    `json:"forum"`
	ThreadId  int32     `json:"thread"`
	Created   time.Time `json:"created"`
}

type PostCreate struct {
	Parent    int64     `json:"parent"`
	Author    string    `json:"author"`
	Message   string    `json:"message"`
}

type PostUpdate struct {
	Message string `json:"message"`
}

type PostFull struct {
	Post *Post `json:"post,omitempty"`
	Author *User `json:"author,omitempty"`
	Thread *Thread `json:"thread,omitempty"`
	Forum *Forum `json:"forum,omitempty"`
}



