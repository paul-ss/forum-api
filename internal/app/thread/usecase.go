package thread

import (
	"github.com/paul-ss/forum-api/internal/domain"
	"github.com/paul-ss/forum-api/internal/domain/query"
)

type IUsecase interface {
	CreatePosts(threadId interface{}, req []domain.PostCreate) ([]domain.Post, error)
	GetThread(threadId interface{}) (*domain.Thread, error)
	UpdateThread(threadId interface{}, req *domain.ThreadUpdate) (*domain.Thread, error)
	GetPosts(threadId interface{}, q *query.GetThreadPosts) ([]domain.Post, error)
	VoteThread(threadId interface{}, req *domain.Vote) (*domain.Thread, error)
}
