package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	config "github.com/paul-ss/forum-api/configs/go"
	"github.com/paul-ss/forum-api/internal/domain"
	"strings"
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
	parent := sql.NullInt64{}
	err := r.db.QueryRow(context.Background(),
		"SELECT id, path[(array_length(path, 1) - 1)], author, message, isEdited, forum_slug, thread_id, created " +
		"FROM posts " +
		"WHERE id = $1 ",
		postId).
		Scan(&p.Id, &parent, &p.Author, &p.Message, &p.IsEdited, &slug, &p.ThreadId, &p.Created)

	if err != nil {
		config.Lg("post_repo", "GetPost").Error("Query: ", err.Error())
		return nil, err
	}

	p.Parent = parent.Int64
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
			"SELECT id, title, author, forum_slug, message, votes, slug, created " +
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


func updatePostFields(args *[]interface{}, req *domain.PostUpdate, startIdx int) string {
	query := []string{}
	if req.Message != nil {
		query = append(query, fmt.Sprintf(" isEdited = message <> $%d, message = $%d ", startIdx, startIdx))
		startIdx += 1
		*args = append(*args, *req.Message)
	}

	if len(query) == 0 {
		query = append(query, " message = message ")
	}

	return strings.Join(query, "")
}

func (r *Repository) UpdatePost(postId int64, rq *domain.PostUpdate) (*domain.Post, error) {
	p := domain.Post{}
	slug := sql.NullString{}
	parent := sql.NullInt64{}
	args := []interface{}{}
	args = append(args, postId)
	err := r.db.QueryRow(context.Background(),
		"UPDATE posts " +
			"SET " + updatePostFields(&args, rq, 2) +
			"WHERE id = $1 " +
			"RETURNING id, path[(array_length(path, 1) - 1)], author, message, isEdited, forum_slug, thread_id, created ",
		args...).
		Scan(&p.Id, &parent, &p.Author, &p.Message, &p.IsEdited, &slug, &p.ThreadId, &p.Created)

	if err != nil {
		config.Lg("post_repo", "UpdatePost").Error("Query: ", err.Error())
		return nil, err
	}

	p.Parent = parent.Int64
	p.ForumSlug = slug.String
	return &p, nil
}