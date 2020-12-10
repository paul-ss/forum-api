package domain

type Status struct {
	User int32 `json:"user"`
	Forum int32 `json:"forum"`
	Thread int32 `json:"thread"`
	Post int64 `json:"post"`
}

type Error struct {
	Message string `json:"message"`
}
