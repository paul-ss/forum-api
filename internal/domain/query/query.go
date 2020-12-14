package query

import "time"

type GetForumUsers struct{
	Limit int32
	Since string
	Desc bool
}

type GetForumThreads struct{
	Limit int32
	Since time.Time
	Desc bool
}