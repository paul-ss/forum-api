package delivery

import (
	"github.com/gin-gonic/gin"
	config "github.com/paul-ss/forum-api/configs/go"
	"github.com/paul-ss/forum-api/internal/app/post"
	"github.com/paul-ss/forum-api/internal/domain"
	"github.com/paul-ss/forum-api/internal/utils"
	"strings"
)

const (
	IdParam = "id"
	RelatedQuery = "related"
)

type Delivery struct {
	uc post.IUsecase
}

func New(uc post.IUsecase) *Delivery {
	return &Delivery{
		uc: uc,
	}
}


func findQueryRelated(url string) []string {
	res := []string{}
	if strings.Contains(url, RelatedQuery + "=") {
		if strings.Contains(url, "user") {
			res = append(res, "user")
		}
		if strings.Contains(url, "forum") {
			res = append(res, "forum")
		}
		if strings.Contains(url, "thread") {
			res = append(res, "thread")
		}
	}

	return res
}

func (d *Delivery) GetPostFull(c *gin.Context) {
	config.Lg("post_http", "GetPostFull").Info("Query: " + c.Request.URL.String())

	//qArr, ok := c.GetQueryArray(RelatedQuery)
	//if !ok {
	//	config.Lg("post_http", "GetPostFull").Info("Can't get query array")
	//	//c.JSON(400, domain.Error{"Can't get query array"})
	//	//return
	//}

	qArr := findQueryRelated(c.Request.URL.String())

	id, err := utils.GetInt64Param(c, IdParam)
	if err != nil {
		config.Lg("post_http", "GetPostFull").Error(err.Error())
		c.JSON(400, domain.Error{err.Error()})
		return
	}

	resp, err := d.uc.GetPostFull(id, qArr)
	if err != nil {
		c.JSON(404, domain.Error{err.Error()})
		return
	}

	c.JSON(200, resp)
}

func (d *Delivery) UpdatePost(c *gin.Context) {
	id, err := utils.GetInt64Param(c, IdParam)
	if err != nil {
		config.Lg("post_http", "UpdPost").Error(err.Error())
		c.JSON(400, domain.Error{err.Error()})
		return
	}

	req := domain.PostUpdate{}
	if err := c.BindJSON(&req); err != nil {
		config.Lg("post_http", "UpdPost").Error(err.Error())
		c.JSON(400, domain.Error{err.Error()})
		return
	}

	resp, err := d.uc.UpdatePost(id, &req)
	if err != nil {
		c.JSON(404, domain.Error{err.Error()})
		return
	}

	c.JSON(200, resp)
}
