package user

import "github.com/paul-ss/forum-api/internal/domain"

type IRepository interface {
	CreateUser(username string, req *domain.UserCreate) ([]domain.User, error)
	GetUser(username string) (*domain.User, error)
	UpdateUser(username string, req *domain.UserCreate) (*domain.User, error)
}
