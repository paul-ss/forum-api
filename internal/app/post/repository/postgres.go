package repository

import (
	"context"
	"database/sql"
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

func (r *Repository) GetPost(postId int64) (*domain.Post, error) {
	p := domain.Post{}
	slug := sql.NullString{}
	err := r.db.QueryRow(context.Background(),
		"SELECT id, path[(array_length(path, 1) - 1)], author, message, isEdited, forum_slug, thread_id, created " +
		"FROM posts " +
		"WHERE id = $1 ",
		postId).
		Scan(&p.Id, &p.Parent, &p.Author, &p.Message, &p.IsEdited, &slug, &p.ThreadId, &p.Created)

	if err != nil {
		config.Lg("post_repo", "GetPost").Error("Query: ", err.Error())
		return nil, err
	}

	p.ForumSlug = slug.String
	return &p, nil
}

func (r *Repository) GetAuthor(postId int64) (*domain.User, error) {
	u := domain.User{}
	err := r.db.QueryRow(context.Background(),
		"WITH p AS (SELECT author FROM posts WHERE id = $1) " +
			"SELECT nickname, fullname, about, email " +
			"FROM users " +
			"WHERE nickname = (SELECT author FROM p) ",
		postId).
		Scan(&u.Nickname, &u.FullName, &u.About, &u.Email)

	if err != nil {
		config.Lg("post_repo", "GetAuthor").Error("Query: ", err.Error())
		return nil, err
	}

	return &u, nil
}

func (r *Repository) GetThread(postId int64) (*domain.Thread, error) {
	t := domain.Thread{}
	slug := sql.NullString{}
	err := r.db.QueryRow(context.Background(),
		"WITH t AS (SELECT thread_id FROM posts WHERE id = $1) " +
			"SELECT id, title, author, forum_title, message, votes, slug, created " +
			"FROM threads " +
			"WHERE id = (SELECT thread_id FROM t) ",
		postId).
		Scan(&t.Id, &t.Title, &t.Author, &t.Forum, &t.Message, &t.Votes, &slug, &t.Created)

	if err != nil {
		config.Lg("post_repo", "GetThread").Error("Query: ", err.Error())
		return nil, err
	}

	t.Slug = slug.String
	return &t, nil
}

func (r *Repository) GetForum(postId int64) (*domain.Forum, error) {
	f := domain.Forum{}
	err := r.db.QueryRow(context.Background(),
		"WITH t AS (SELECT forum_id FROM posts WHERE id = $1) " +
			"SELECT id, title, nickname, slug, posts, threads " +
			"FROM forums " +
			"WHERE id = (SELECT forum_id FROM t) ",
		postId).
		Scan(&f.Id, &f.Title, &f.User, &f.Slug, &f.Posts, &f.Threads)

	if err != nil {
		config.Lg("post_repo", "GetThread").Error("Query: ", err.Error())
		return nil, err
	}

	return &f, nil
}

func (r *Repository) UpdatePost(postId int64, rq *domain.PostUpdate) (*domain.Post, error) {
	p := domain.Post{}
	slug := sql.NullString{}
	err := r.db.QueryRow(context.Background(),
		"UPDATE posts " +
			"SET message = $1, isEdited = true " +
			"WHERE id = $2 " +
			"RETURNING id, path[(array_length(path, 1) - 1)], author, message, isEdited, forum_slug, thread_id, created ",
		rq.Message, postId).
		Scan(&p.Id, &p.Parent, &p.Author, &p.Message, &p.IsEdited, &slug, &p.ThreadId, &p.Created)

	if err != nil {
		config.Lg("post_repo", "UpdatePost").Error("Query: ", err.Error())
		return nil, err
	}

	p.ForumSlug = slug.String
	return &p, nil
}