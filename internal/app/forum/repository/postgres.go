package repository

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	config "github.com/paul-ss/forum-api/configs/go"
	"github.com/paul-ss/forum-api/internal/domain"
	domainErr "github.com/paul-ss/forum-api/internal/domain/errors"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) StoreForum(f *domain.Forum) (*domain.Forum, error) {
	err := r.db.QueryRow(context.Background(),
			"INSERT INTO forums (title, nickname, slug) " +
			"VALUES ($1, $2, $3) " +
			"RETURNING id, posts, threads ",
		f.Title, f.User, f.Slug).Scan(&f.Id, &f.Posts, &f.Threads)

	if err != nil {
		config.Lg("forum_repo", "StoreForum").Info(err.Error())
		er := r.db.QueryRow(context.Background(),
			"SELECT id, title, nickname, slug, posts, threads " +
				"FROM forums " +
				"WHERE slug = $1 ",
			f.Slug).Scan(&f.Id, &f.Title, &f.User, &f.Slug, &f.Posts, &f.Threads)

		if er != nil {
			config.Lg("forum_repo", "StoreForum").Error("Select: ", err.Error())
			return f, er
		}

		return f, domainErr.DuplicateKeyError
	}

	return f, nil
}

func (r *Repository) GetForumBySlug(slug string) (*domain.Forum, error) {
	f := domain.Forum{}
	err := r.db.QueryRow(context.Background(),
		"SELECT id, title, nickname, slug, posts, threads " +
			"FROM forums " +
			"WHERE slug = $1 ",
		slug).Scan(&f.Id, &f.Title, &f.User, &f.Slug, &f.Posts, &f.Threads)

	if err != nil {
		config.Lg("forum_repo", "StoreForum").Error("Select: ", err.Error())
		return nil, err
	}

	return &f, nil
}