package query

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"strconv"
	"time"
)


type GetForumUsers struct{
	Limit int32
	Since string
	Desc bool
}

func (q *GetForumUsers) Init(c *gin.Context) error {
	p := c.DefaultQuery("limit", strconv.Itoa(math.MaxInt32))

	limit, err := strconv.Atoi(p)
	if err != nil {
		return fmt.Errorf("can't convert limit to int")
	}

	q.Limit = int32(limit)


	p = c.DefaultQuery("since", "")

	q.Since = p


	p = c.DefaultQuery("desc", "false")


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
	p := c.DefaultQuery("limit", "0")

	limit, err := strconv.Atoi(p)
	if err != nil {
		return fmt.Errorf("can't convert limit to int")
	}

	q.Limit = int32(limit)

	// "2006-01-02T15:04:05Z"

	p, ok := c.GetQuery("since")
	if ok {
		tme, err := time.Parse(time.RFC3339, p)
		if err != nil {
			return fmt.Errorf("error parce time")
		}

		q.Since = tme
	} else {
		q.Since = time.Time{}
	}


	p = c.DefaultQuery("desc", "false")

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
	p := c.DefaultQuery("limit", "0")

	limit, err := strconv.Atoi(p)
	if err != nil {
		return fmt.Errorf("can't convert limit to int")
	}

	q.Limit = int32(limit)


	p = c.DefaultQuery("since", "0")

	since, err := strconv.ParseInt(p, 10, 64)
	if err != nil {
		return fmt.Errorf("error parce since")
	}

	q.Since = since


	p = c.DefaultQuery("sort", "flat")

	q.Sort = p


	p = c.DefaultQuery("desc", "false")

	desc, err := strconv.ParseBool(p)
	if err != nil {
		return fmt.Errorf("can't convert desc to bool")
	}

	q.Desc = desc
	return nil
}