package usecase

import (
	"github.com/paul-ss/forum-api/internal/app/thread"
	"github.com/paul-ss/forum-api/internal/domain"
	"github.com/paul-ss/forum-api/internal/domain/query"
)

type Usecase struct {
	db thread.IRepository
}

func New(db thread.IRepository) *Usecase {
	return &Usecase{
		db: db,
	}
}

func (uc *Usecase) CreatePosts(threadId interface{}, req []domain.PostCreate) ([]domain.Post, error) {
	return uc.db.CreatePosts(threadId, req)
}

func (uc *Usecase) GetThread(threadId interface{}) (*domain.Thread, error) {
	return uc.db.GetThread(threadId)
}

func (uc *Usecase) UpdateThread(threadId interface{}, req *domain.ThreadUpdate) (*domain.Thread, error) {
	return uc.db.UpdateThread(threadId, req)
}

func (uc *Usecase) GetPosts(threadId interface{}, q *query.GetThreadPosts) ([]domain.Post, error) {
	return uc.db.GetPosts(threadId, q)
}

func (uc *Usecase) VoteThread(threadId interface{}, req *domain.Vote) (*domain.Thread, error) {
	return uc.db.VoteThread(threadId, req)
}
