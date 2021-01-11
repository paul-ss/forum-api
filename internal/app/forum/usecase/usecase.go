package usecase

import (
	"github.com/paul-ss/forum-api/internal/app/forum"
	"github.com/paul-ss/forum-api/internal/domain"
	"github.com/paul-ss/forum-api/internal/domain/query"
)

type Usecase struct {
	db forum.IRepository
}

func New(db forum.IRepository) *Usecase {
	return &Usecase{
		db: db,
	}
}

func (uc *Usecase) StoreForum(f *domain.Forum) (*domain.Forum, error) {
	return uc.db.StoreForum(f)
}

func (uc *Usecase) GetForumBySlug(slug string) (*domain.Forum, error) {
	return uc.db.GetForumBySlug(slug)
}

func (uc *Usecase) StoreThread(slug string, tc domain.ThreadCreate) (*domain.Thread, error) {
	return uc.db.StoreThread(slug, tc)
}

func (uc *Usecase) GetUsers(slug string, q *query.GetForumUsers) ([]domain.User, error) {
	return uc.db.GetUsers(slug, q)
}

func (uc *Usecase) GetThreads(slug string, q *query.GetForumThreads) ([]domain.Thread, error) {
	return uc.db.GetThreads(slug, q)
}
