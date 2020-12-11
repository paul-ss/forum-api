package domain

type Forum struct {
	Id int32 `json:"-"`
	Title string `json:"tittle"`
	User string `json:"user"`
	Slug string `json:"slug"`
	Posts int64 `json:"posts"`
	Threads	int32 `json:"threads"`
}