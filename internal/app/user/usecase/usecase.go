package usecase

import (
	"github.com/paul-ss/forum-api/internal/app/user"
	"github.com/paul-ss/forum-api/internal/domain"
)

type Usecase struct {
	db user.IRepository
}

func New(db user.IRepository) *Usecase {
	return &Usecase{
		db: db,
	}
}

func (uc *Usecase) CreateUser(username string, req *domain.UserCreate) ([]domain.User, error) {
	return uc.db.CreateUser(username, req)
}

func (uc *Usecase) GetUser(username string) (*domain.User, error) {
	return uc.db.GetUser(username)
}

func (uc *Usecase) UpdateUser(username string, req *domain.UserCreate) (*domain.User, error) {
	return uc.db.UpdateUser(username, req)
}

func (uc *Usecase) ClearAll() error {
	return uc.db.ClearAll()
}

func (uc *Usecase) GetStats() (*domain.Status, error) {
	return uc.db.GetStats()
}