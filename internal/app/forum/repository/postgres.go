package repository

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	config "github.com/paul-ss/forum-api/configs/go"
	"github.com/paul-ss/forum-api/internal/domain"
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
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		config.Lg("forum_repo", "StoreForum").Error("Begin: ", err.Error())
		return nil, err
	}

	err = tx.QueryRow(context.Background(),
		"WITH ins AS (SELECT true as was) " +
			"INSERT INTO forums (title, nickname, slug) " +
			"VALUES ($1, $2, $3) " +
			"ON CONFLICT (slug) DO UPDATE " +
				"SET ins.was = false " +
			"RETURNING id, posts, threads",
		f.Title, f.User, f.Slug).Scan(&f.Id, &f.Posts, &f.Threads)


	if err != nil {
		config.Lg("forum_repo", "StoreForum").Info(err.Error())
		er := tx.QueryRow(context.Background(),
			"SELECT id, title, nickname, slug, posts, threads " +
				"FROM forums " +
				"WHERE slug = $1 ",
			f.Slug).Scan(&f.Id, &f.Title, &f.User, &f.Slug, &f.Posts, &f.Threads)

		if er != nil {
			config.Lg("forum_repo", "StoreForum").Error("Select: ", err.Error())
			tx.Rollback(context.Background())
			return nil, er
		}
	}

	if err := tx.Commit(context.Background()); err != nil {
		config.Lg("forum_repo", "StoreForum").Error("Commit: ", err.Error())
		tx.Rollback(context.Background())
		return nil, err
	}

	return f, nil
}
