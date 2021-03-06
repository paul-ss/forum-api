package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/paul-ss/forum-api/internal/domain/query"
	"strconv"
	"strings"
	"time"
)

func RandomSlug() string {
	return strings.ReplaceAll(time.Now().String(), " ", "")
}

func GetCurrentTime(t time.Time) time.Time {
	nt := time.Time{}
	if t == nt {
		return time.Now()
	}

	return t
}

func DESC(d bool) string {
	if d {
		return " DESC "
	}
	return " "
}

func LessIfDESC(q *query.GetThreadPosts) string {
	if q.Since > 0 {
		if q.Desc {
			return " < "
		}
		return " > "
	} else {
		return " > "
	}
}

func LessIfDESCU(q *query.GetForumUsers) string {
	if q.Since != "" {
		if q.Desc {
			return " < "
		}
		return " > "
	} else {
		return " > "
	}
}

func GetInt64Param(c *gin.Context, param string) (int64, error) {
	strParam, ok := c.Params.Get(param)
	if !ok {
		return 0, fmt.Errorf("param '%s' not found", param)
	}

	return strconv.ParseInt(strParam, 10, 64)
}

func GetIntOrStringParam(c *gin.Context, param string) (interface{}, error) {
	strParam, ok := c.Params.Get(param)
	if !ok {
		return nil, fmt.Errorf("param '%s' not found", param)
	}

	intPar, err := strconv.Atoi(strParam)
	if err != nil {
		return strParam, nil
	}

	return intPar, nil
}