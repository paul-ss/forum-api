package post

import "github.com/paul-ss/forum-api/internal/domain"

type IUsecase interface {
	GetPostFull(postId int64, related []string) (*domain.PostFull, error)
	UpdatePost(postId int64, rq *domain.PostUpdate) (*domain.Post, error)
}
