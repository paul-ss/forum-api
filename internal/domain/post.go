package domain

import (
	"time"
)

type Post struct {
	Description string `json:"description"`
	Id int64 `json:"id"`
	Parent int64 `json:"parent"`
	Author string `json:"author"`
	Message string `json:"message"`
	IsEdited bool `json:"isEdited"`
	Forum string `json:"forum"`
	ThreadId int32 `json:"thread"`
	Created time.Time `json:"created"`
}

type PostUpdate struct {
	Message string `json:"message"`
}

type PostFull struct {
	Post Post `json:"post"`
	Author User `json:"author"`
	Thread Thread `json:"thread"`
	Forum Forum `json:"forum"`
}



