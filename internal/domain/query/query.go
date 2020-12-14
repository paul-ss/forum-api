package query

type GetForumUsers struct{
	Limit int32
	Since string
	Desc bool
}