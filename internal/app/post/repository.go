package post

import "github.com/paul-ss/forum-api/internal/domain"

type IRepository interface {
	GetPost(postId int64) (*domain.Post, error)
	GetAuthor(postId int64) (*domain.User, error)
	GetThread(postId int64) (*domain.Thread, error)
	GetForum(postId int64) (*domain.Forum, error)
	UpdatePost(postId int64, rq *domain.PostUpdate) (*domain.Post, error)
}
