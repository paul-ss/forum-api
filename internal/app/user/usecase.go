package user

import "github.com/paul-ss/forum-api/internal/domain"

type IUsecase interface {
	CreateUser(username string, req *domain.UserCreate) ([]domain.User, error)
	GetUser(username string) (*domain.User, error)
	UpdateUser(username string, req *domain.UserUpdate) (*domain.User, error)
	ClearAll() error
	GetStats() (*domain.Status, error)
}
