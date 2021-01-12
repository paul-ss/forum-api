package delivery

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/paul-ss/forum-api/internal/app/forum/usecase"
	"github.com/paul-ss/forum-api/internal/app/forum/repository"
)

func CreateForumDelivery(db *pgxpool.Pool, handler *gin.RouterGroup) {
	r := repository.New(db)
	uc := usecase.New(r)
	d := New(uc)

	handler.POST(fmt.Sprintf("/forum/:%s", SlugParam), d.CreateForum) // /forum/create ????

	handler.GET(fmt.Sprintf("/forum/:%s/details", SlugParam), d.GetForumBySlug)
	handler.POST(fmt.Sprintf("/forum/:%s/create", SlugParam), d.CreateThread)
	handler.GET(fmt.Sprintf("/forum/:%s/users", SlugParam), d.GetUsers)
	handler.GET(fmt.Sprintf("/forum/:%s/threads", SlugParam), d.GetThreads)
}