package usecase

import (
	config "github.com/paul-ss/forum-api/configs/go"
	"github.com/paul-ss/forum-api/internal/app/post"
	"github.com/paul-ss/forum-api/internal/domain"
)

type Usecase struct {
	db post.IRepository
}

func New(db post.IRepository) *Usecase {
	return &Usecase{
		db: db,
	}
}

func (uc *Usecase) GetPostFull(postId int64, related []string) (*domain.PostFull, error) {
	resp := domain.PostFull{}
	for _, elem := range related {
		switch elem {
		case "user":
			u, err := uc.db.GetAuthor(postId)
			if err != nil {
				config.Lg("post_uc", "GetPostFull").Error("GetAuthor : " + err.Error())
				return nil, err
			}
			resp.Author = u
		case "forum":
			f, err := uc.db.GetForum(postId)
			if err != nil {
				config.Lg("post_uc", "GetPostFull").Error("GetForum : " + err.Error())
				return nil, err
			}
			resp.Forum = f
		case "thread":
			t, err := uc.db.GetThread(postId)
			if err != nil {
				config.Lg("post_uc", "GetPostFull").Error("GetThread : " + err.Error())
				return nil, err
			}
			resp.Thread = t
		}
	}

	p, err := uc.db.GetPost(postId)
	if err != nil {
		config.Lg("post_uc", "GetPostFull").Error("GetPost : " + err.Error())
		return nil, err
	}
	resp.Post = p

	return &resp, nil
}


func (uc *Usecase) UpdatePost(postId int64, rq *domain.PostUpdate) (*domain.Post, error) {
	return uc.db.UpdatePost(postId, rq)
}