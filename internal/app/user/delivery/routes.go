package delivery

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/paul-ss/forum-api/internal/app/user/repository"
	"github.com/paul-ss/forum-api/internal/app/user/usecase"
)

func CreateUserDelivery(db *pgxpool.Pool, handler *gin.RouterGroup) {
	r := repository.New(db)
	uc := usecase.New(r)
	d := New(uc)

	handler.POST(fmt.Sprintf("/user/:%s/create", NicknameParam), d.CreateUser)
	handler.GET(fmt.Sprintf("/user/:%s/profile", NicknameParam), d.GetUser)
	handler.POST(fmt.Sprintf("/user/:%s/profile", NicknameParam), d.UpdateUser)
	handler.POST("/service/clear", d.ClearAll)
	handler.GET("/service/status", d.GetStats)
}
