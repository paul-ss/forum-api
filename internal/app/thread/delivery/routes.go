package delivery

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/paul-ss/forum-api/internal/app/thread/repository"
	"github.com/paul-ss/forum-api/internal/app/thread/usecase"
)

func CreateThreadDelivery(db *pgxpool.Pool, handler *gin.Engine) {
	r := repository.New(db)
	uc := usecase.New(r)
	d := New(uc)

	handler.POST(fmt.Sprintf("/thread/:%s/create", SlugOrIdPar), d.CreatePosts)
	handler.GET(fmt.Sprintf("/thread/:%s/details", SlugOrIdPar), d.GetThread)
	handler.POST(fmt.Sprintf("/thread/:%s/details", SlugOrIdPar), d.UpdateThread)
	handler.GET(fmt.Sprintf("/thread/:%s/posts", SlugOrIdPar), d.GetPosts)
	handler.POST(fmt.Sprintf("/thread/:%s/vote", SlugOrIdPar), d.VoteThread)
}
