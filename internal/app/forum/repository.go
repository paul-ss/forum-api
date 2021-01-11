package forum

import (
	"github.com/paul-ss/forum-api/internal/domain"
	"github.com/paul-ss/forum-api/internal/domain/query"
)

type IRepository interface {
	StoreForum(f *domain.Forum) (*domain.Forum, error)
	GetForumBySlug(slug string) (*domain.Forum, error)
	StoreThread(slug string, tc domain.ThreadCreate) (*domain.Thread, error)
	GetUsers(slug string, q *query.GetForumUsers) ([]domain.User, error)
	GetThreads(slug string, q *query.GetForumThreads) ([]domain.Thread, error)
}
