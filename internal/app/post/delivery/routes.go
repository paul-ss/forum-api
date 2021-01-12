package delivery

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/paul-ss/forum-api/internal/app/post/repository"
	"github.com/paul-ss/forum-api/internal/app/post/usecase"
)

func CreatePostDelivery(db *pgxpool.Pool, handler *gin.RouterGroup) {
	r := repository.New(db)
	uc := usecase.New(r)
	d := New(uc)

	handler.GET(fmt.Sprintf("/post/:%s/details", IdParam), d.GetPostFull)
	handler.POST(fmt.Sprintf("/post/:%s/details", IdParam), d.UpdatePost)
}