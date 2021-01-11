package query

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)



type GetForumUsers struct{
	Limit int32
	Since string
	Desc bool
}

func (q *GetForumUsers) Init(c *gin.Context) error {
	p, ok := c.GetQuery("limit")
	if !ok {
		return fmt.Errorf("limit param not exists")
	}

	limit, err := strconv.Atoi(p)
	if err != nil {
		return fmt.Errorf("can't convert limit to int")
	}

	q.Limit = int32(limit)


	p, ok = c.GetQuery("since")
	if !ok {
		return fmt.Errorf("since param not exists")
	}

	q.Since = p


	p, ok = c.GetQuery("desc")
	if !ok {
		return fmt.Errorf("desc param not exists")
	}

	desc, err := strconv.ParseBool(p)
	if err != nil {
		return fmt.Errorf("can't convert desc to bool")
	}

	q.Desc = desc
	return nil
}


type GetForumThreads struct{
	Limit int32
	Since time.Time
	Desc bool
}

func (q *GetForumThreads) Init(c *gin.Context) error {
	p, ok := c.GetQuery("limit")
	if !ok {
		return fmt.Errorf("limit param not exists")
	}

	limit, err := strconv.Atoi(p)
	if err != nil {
		return fmt.Errorf("can't convert limit to int")
	}

	q.Limit = int32(limit)


	p, ok = c.GetQuery("since")
	if !ok {
		return fmt.Errorf("since param not exists")
	}

	tme, err := time.Parse(time.RFC3339, p)
	if err != nil {
		return fmt.Errorf("error parce time")
	}

	q.Since = tme


	p, ok = c.GetQuery("desc")
	if !ok {
		return fmt.Errorf("desc param not exists")
	}

	desc, err := strconv.ParseBool(p)
	if err != nil {
		return fmt.Errorf("can't convert desc to bool")
	}

	q.Desc = desc
	return nil
}

type GetThreadPosts struct{
	Limit int32
	Since int64
	Sort string
	Desc bool
}

func (q *GetThreadPosts) Init(c *gin.Context) error {
	p, ok := c.GetQuery("limit")
	if !ok {
		return fmt.Errorf("limit param not exists")
	}

	limit, err := strconv.Atoi(p)
	if err != nil {
		return fmt.Errorf("can't convert limit to int")
	}

	q.Limit = int32(limit)


	p, ok = c.GetQuery("since")
	if !ok {
		return fmt.Errorf("since param not exists")
	}

	since, err := strconv.ParseInt(p, 10, 64)
	if err != nil {
		return fmt.Errorf("error parce since")
	}

	q.Since = since


	p, ok = c.GetQuery("sort")
	if !ok {
		return fmt.Errorf("sort param not exists")
	}

	q.Sort = p


	p, ok = c.GetQuery("desc")
	if !ok {
		return fmt.Errorf("desc param not exists")
	}

	desc, err := strconv.ParseBool(p)
	if err != nil {
		return fmt.Errorf("can't convert desc to bool")
	}

	q.Desc = desc
	return nil
}