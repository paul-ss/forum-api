package delivery

import (
	"github.com/gin-gonic/gin"
	config "github.com/paul-ss/forum-api/configs/go"
	"github.com/paul-ss/forum-api/internal/app/user"
	"github.com/paul-ss/forum-api/internal/domain"
	domainErr "github.com/paul-ss/forum-api/internal/domain/errors"
)

const (
	NicknameParam = "nickname"
)

type Delivery struct {
	uc user.IUsecase
}

func New(uc user.IUsecase) *Delivery {
	return &Delivery{
		uc: uc,
	}
}

func (d *Delivery) CreateUser(c *gin.Context) {
	username := c.Param(NicknameParam)

	req := domain.UserCreate{}
	if err := c.BindJSON(&req); err != nil {
		config.Lg("user_http", "CreateUser").Error(err.Error())
		c.JSON(400, domain.Error{err.Error()})
		return
	}

	resp, err := d.uc.CreateUser(username, &req)
	if err != nil {
		switch (err) {
		case domainErr.AlreadyExists:
			c.JSON(409, resp)
		default:
			c.JSON(500, domain.Error{err.Error()})
		}

		config.Lg("user_http", "CreateUser").Error("Usecase: " + err.Error())
		return
	}

	c.JSON(201, resp[0])
}

func (d *Delivery) GetUser(c *gin.Context) {
	username := c.Param(NicknameParam)

	resp, err := d.uc.GetUser(username)
	if err != nil {
		config.Lg("user_http", "GetUser").Error("Usecase: " + err.Error())
		c.JSON(404, domain.Error{err.Error()})
		return
	}

	c.JSON(200, resp)
}

func (d *Delivery) UpdateUser(c *gin.Context) {
	username := c.Param(NicknameParam)

	req := domain.UserUpdate{}
	if err := c.BindJSON(&req); err != nil {
		config.Lg("user_http", "UpdUser").Error(err.Error())
		c.JSON(400, domain.Error{err.Error()})
		return
	}

	resp, err := d.uc.UpdateUser(username, &req)
	if err != nil {
		switch err {
		case domainErr.AlreadyExists:
			c.JSON(409, domain.Error{err.Error()})
		default:
			c.JSON(404, domain.Error{err.Error()})
		}

		config.Lg("user_http", "UpdUser").Error("Usecase: " + err.Error())
		return
	}

	c.JSON(200, resp)
}


func (d *Delivery) ClearAll(c *gin.Context) {
	if err := d.uc.ClearAll(); err != nil {
		config.Lg("user_http", "ClearAll").Error(err.Error())
		c.JSON(500, domain.Error{err.Error()})
		return
	}

	c.Status(200)
}

func (d *Delivery) GetStats(c *gin.Context) {
	resp, err := d.uc.GetStats()
	if err != nil {
		config.Lg("user_http", "GetStats").Error(err.Error())
		c.JSON(500, domain.Error{err.Error()})
		return
	}

	c.JSON(200, resp)
}

