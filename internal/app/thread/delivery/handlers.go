package delivery

import (
	"github.com/gin-gonic/gin"
	config "github.com/paul-ss/forum-api/configs/go"
	"github.com/paul-ss/forum-api/internal/app/thread"
	"github.com/paul-ss/forum-api/internal/domain"
	domainErr "github.com/paul-ss/forum-api/internal/domain/errors"
	"github.com/paul-ss/forum-api/internal/domain/query"
	"github.com/paul-ss/forum-api/internal/utils"
)

const (
	SlugOrIdPar = "slug-or-id"
)

type Delivery struct {
	uc thread.IUsecase
}

func New(uc thread.IUsecase) *Delivery {
	return &Delivery{
		uc: uc,
	}
}

func (d *Delivery) CreatePosts(c *gin.Context) {
	id, err := utils.GetIntOrStringParam(c, SlugOrIdPar)
	if err != nil {
		config.Lg("thread_http", "CreatePosts").Error(err.Error())
		c.JSON(400, domain.Error{err.Error()})
		return
	}

	req := []domain.PostCreate{}
	if err := c.BindJSON(&req); err != nil {
		config.Lg("thread_http", "CreatePosts").Error(err.Error())
		c.JSON(400, domain.Error{err.Error()})
		return
	}

	resp, err := d.uc.CreatePosts(id, req)
	if err != nil {
		switch err {
		case domainErr.PostNotExists:
			c.JSON(409, domain.Error{err.Error()})
		case domainErr.ThreadNotExists:
			c.JSON(404, domain.Error{err.Error()})
		default:
			c.JSON(500, domain.Error{err.Error()})
		}

		config.Lg("thread_http", "CreatePosts").Error(err.Error())
		return
	}

	c.JSON(201, resp)
}

func (d *Delivery) GetThread(c *gin.Context) {
	id, err := utils.GetIntOrStringParam(c, SlugOrIdPar)
	if err != nil {
		config.Lg("thread_http", "GetThread").Error(err.Error())
		c.JSON(400, domain.Error{err.Error()})
		return
	}

	resp, err := d.uc.GetThread(id)
	if err != nil {
		c.JSON(404, domain.Error{err.Error()})
		config.Lg("thread_http", "GetThread").Error(err.Error())
		return
	}

	c.JSON(200, resp)
}

func (d *Delivery) UpdateThread(c *gin.Context) {
	id, err := utils.GetIntOrStringParam(c, SlugOrIdPar)
	if err != nil {
		config.Lg("thread_http", "UpdateThread").Error(err.Error())
		c.JSON(400, domain.Error{err.Error()})
		return
	}

	req := domain.ThreadUpdate{}
	if err := c.BindJSON(&req); err != nil {
		config.Lg("thread_http", "UpdateThread").Error(err.Error())
		c.JSON(400, domain.Error{err.Error()})
		return
	}

	resp, err := d.uc.UpdateThread(id, &req)
	if err != nil {
		c.JSON(404, domain.Error{err.Error()})
		config.Lg("thread_http", "UpdateThread").Error(err.Error())
		return
	}

	c.JSON(200, resp)
}

func (d *Delivery) GetPosts(c *gin.Context) {
	id, err := utils.GetIntOrStringParam(c, SlugOrIdPar)
	if err != nil {
		config.Lg("thread_http", "GetPosts").Error(err.Error())
		c.JSON(400, domain.Error{err.Error()})
		return
	}

	q := query.GetThreadPosts{}
	if err := q.Init(c); err != nil {
		config.Lg("thread_http", "GetPosts").Error(err.Error())
		c.JSON(400, domain.Error{err.Error()})
		return
	}


	resp, err := d.uc.GetPosts(id, &q)
	if err != nil {
		c.JSON(404, domain.Error{err.Error()})
		config.Lg("thread_http", "GetPosts").Error(err.Error())
		return
	}

	c.JSON(200, resp)
}

func (d *Delivery) VoteThread(c *gin.Context) {
	id, err := utils.GetIntOrStringParam(c, SlugOrIdPar)
	if err != nil {
		config.Lg("thread_http", "VoteThread").Error(err.Error())
		c.JSON(400, domain.Error{err.Error()})
		return
	}

	req := domain.Vote{}
	if err := c.BindJSON(&req); err != nil {
		config.Lg("thread_http", "VoteThread").Error(err.Error())
		c.JSON(400, domain.Error{err.Error()})
		return
	}

	resp, err := d.uc.VoteThread(id, &req)
	if err != nil {
		c.JSON(404, domain.Error{err.Error()})
		config.Lg("thread_http", "VoteThread").Error(err.Error())
		return
	}

	c.JSON(200, resp)
}

