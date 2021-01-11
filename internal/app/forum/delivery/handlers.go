package delivery

import (
	"github.com/gin-gonic/gin"
	config "github.com/paul-ss/forum-api/configs/go"
	"github.com/paul-ss/forum-api/internal/app/forum"
	"github.com/paul-ss/forum-api/internal/domain"
	domainErr "github.com/paul-ss/forum-api/internal/domain/errors"
	"github.com/paul-ss/forum-api/internal/domain/query"
)

const (
	SlugParam = "slug"
)

type Delivery struct {
	uc forum.IUsecase
}

func New(uc forum.IUsecase) *Delivery {
	return &Delivery{
		uc: uc,
	}
}

func (d *Delivery) CreateForum(c *gin.Context) {
	req := domain.Forum{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, domain.Error{err.Error()})
		config.Lg("forum_http", "CreateForum").Error(err.Error())
		return
	}

	resp, err := d.uc.StoreForum(&req)
	if err != nil {
		switch err {
		case domainErr.DuplicateKeyError:
			c.JSON(409, resp)
		default:
			c.JSON(404, domain.Error{err.Error()})
		}

		config.Lg("forum_http", "CreateForum").Error(err.Error())
		return
	}

	c.JSON(201, resp)
}


func (d *Delivery) GetForumBySlug(c *gin.Context) {
	slug := c.Param(SlugParam)

	resp, err := d.uc.GetForumBySlug(slug)
	if err != nil {
		c.JSON(404, domain.Error{err.Error()})
		config.Lg("forum_http", "GetForumBySlug").Error(err.Error())
		return
	}

	c.JSON(200, resp)
}

func (d *Delivery) CreateThread(c *gin.Context) {
	slug := c.Param(SlugParam)

	req := domain.ThreadCreate{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, domain.Error{err.Error()})
		config.Lg("forum_http", "CreateThread").Error(err.Error())
		return
	}

	resp, err := d.uc.StoreThread(slug, req)
	if err != nil {
		switch err {
		case domainErr.DuplicateKeyError:
			c.JSON(409, resp)
		default:
			c.JSON(404, domain.Error{err.Error()})
		}

		config.Lg("forum_http", "CreateThread").Error(err.Error())
		return
	}

	c.JSON(201, resp)
}

func (d *Delivery) GetUsers(c *gin.Context) {
	slug := c.Param(SlugParam)

	q := query.GetForumUsers{}
	if err := q.Init(c); err != nil {
		c.JSON(400, domain.Error{err.Error()})
		config.Lg("forum_http", "GetUsers").Error(err.Error())
		return
	}

	resp, err := d.uc.GetUsers(slug, &q)
	if err != nil {
		c.JSON(404, domain.Error{err.Error()})
		config.Lg("forum_http", "GetUsers").Error(err.Error())
		return
	}

	c.JSON(200, resp)
}

func (d *Delivery) GetThreads(c *gin.Context) {
	slug := c.Param(SlugParam)

	q := query.GetForumThreads{}
	if err := q.Init(c); err != nil {
		c.JSON(400, domain.Error{err.Error()})
		config.Lg("forum_http", "GetThreads").Error(err.Error())
		return
	}

	resp, err := d.uc.GetThreads(slug, &q)
	if err != nil {
		c.JSON(404, domain.Error{err.Error()})
		config.Lg("forum_http", "GetThreads").Error(err.Error())
		return
	}

	c.JSON(200, resp)
}

