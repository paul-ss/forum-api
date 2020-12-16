package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	config "github.com/paul-ss/forum-api/configs/go"
	"github.com/paul-ss/forum-api/internal/domain"
	domainErr "github.com/paul-ss/forum-api/internal/domain/errors"
	"github.com/paul-ss/forum-api/internal/domain/query"
	"github.com/paul-ss/forum-api/internal/utils"
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

func valuesPosts(i int, req []domain.PostCreate, args *[]interface{}) string {
	query := []string{}
	pIdxAdd := 1
	for _, t := range req {
		query = append(query,
			"( " +
			fmt.Sprintf("(WITH par AS (SELECT path FROM posts WHERE id = $%d) ", i) +
			"SELECT CASE WHEN ((SELECT path FROM par) IS NOT NULL) THEN " +
			fmt.Sprintf("(SELECT path FROM par) || (last_value + %d) ", pIdxAdd) +
			fmt.Sprintf("WHEN ($%d < 1) THEN ", i) +
			fmt.Sprintf("ARRAY[last_value + %d]  ", pIdxAdd) +
			"ELSE null END FROM seq), " +
			fmt.Sprintf("$%d, $%d, ", i + 1, i + 2) +
			"(SELECT slug FROM t), " +
			"(SELECT forum_id FROM t), " +
			"(SELECT id FROM t) " +
			")", ",")

		*args = append(*args, t.Parent, t.Author, t.Message)
		i += 3
		pIdxAdd += 1
	}

	query = query[:len(query) - 1]
	return strings.Join(query, "")
}

func createPostSelectThread(id interface{}) string {
	 if _, ok := id.(string); ok {
		 return "(SELECT forum_id, slug, id FROM threads WHERE slug = $1)"
	 } else {
		 return "(SELECT forum_id, slug, id FROM threads WHERE id = $1)"
	 }
}

func createPostError(err error) error {
	pErr := err.(*pgconn.PgError)

	switch pErr.ColumnName {
	case "path":
		return domainErr.PostNotExists
	default:
		return domainErr.ThreadNotExists
	}
}

func (r *Repository) CreatePosts(threadId interface{}, req []domain.PostCreate) ([]domain.Post, error) {
	args := []interface{}{}
	args = append(args, threadId)
	rows, err := r.db.Query(context.Background(),
		"WITH t AS " + createPostSelectThread(threadId) + "," +
				"seq AS (SELECT last_value FROM posts_id_seq) " +
			"INSERT INTO posts (path, author, message, forum_slug, forum_id, thread_id) " +
			"VALUES  " + valuesPosts(2, req, &args) +
			"RETURNING id, path[(array_length(path, 1) - 1)], author, message, isEdited, forum_slug, thread_id, created ",
			args...)

	if err != nil {
		config.Lg("thread_repo", "CreatePosts").Error("Query: ", err.Error())
		return nil, createPostError(err)
	}

	defer rows.Close()

	posts := []domain.Post{}
	for rows.Next() {
		p := domain.Post{}
		slug := sql.NullString{}
		parent := sql.NullInt64{}

		err := rows.Scan(&p.Id, &parent, &p.Author, &p.Message, &p.IsEdited, &slug, &p.ThreadId, &p.Created)
		if err != nil {
			config.Lg("thread_repo", "CreatePosts").Error("Scan: ", err.Error())
			return nil, err
		}

		p.ForumSlug = slug.String
		p.Parent = parent.Int64
		posts = append(posts, p)
	}

	if err := rows.Err(); err != nil {
		config.Lg("thread_repo", "CreatePosts").Error("Rows: ", err.Error())
		return nil, createPostError(err)
	}

	return posts, nil
}

func getThreadCond(id interface{}, param int) string {
	if _, ok := id.(string); ok {
		return fmt.Sprintf("WHERE slug = $%d", param)
	} else {
		return fmt.Sprintf("WHERE id = $%d", param)
	}
}

func (r *Repository) GetThread(threadId interface{}) (*domain.Thread, error) {
	t := domain.Thread{}
	err := r.db.QueryRow(context.Background(),
		"SELECT id, title, author, forum_title, message, votes, slug, created " +
			"FROM threads " +
			getThreadCond(threadId, 1),
			threadId).
		Scan(&t.Id, &t.Title, &t.Author, &t.Forum, &t.Message, &t.Votes, &t.Slug, &t.Created)

	if err != nil {
		config.Lg("thread_repo", "GetThread").Error("Query: ", err.Error())
		return nil, err
	}

	return &t, nil
}

func (r *Repository) UpdateThread(threadId interface{}, req *domain.ThreadUpdate) (*domain.Thread, error) {
	t := domain.Thread{}
	err := r.db.QueryRow(context.Background(),
		"UPDATE threads " +
		"SET title = $1, message = $2 " +
		getThreadCond(threadId, 3) +
		"RETURNING id, title, author, forum_title, message, votes, slug, created ",
		req.Title, req.Message, threadId).
		Scan(&t.Id, &t.Title, &t.Author, &t.Forum, &t.Message, &t.Votes, &t.Slug, &t.Created)

	if err != nil {
		config.Lg("thread_repo", "UpdateThread").Error("Query: ", err.Error())
		return nil, err
	}

	return &t, nil
}


func getPostsCond(id interface{}, param int) string {
	if _, ok := id.(string); ok {
		return fmt.Sprintf("(SELECT id FROM thread WHERE slug = $%d) ", param)
	} else {
		return fmt.Sprintf("(SELECT $%d AS id) ", param)
	}
}


func getPostsSort(q *query.GetThreadPosts, ) string {
	switch q.Sort {
	case "flat":
		return "AND id > $2 " +
			"ORDER BY created " + utils.DESC(q.Desc) +
			"LIMIT $3 "
	case "tree":
		return "AND id > $2 " +
			"ORDER BY path[1] " + utils.DESC(q.Desc) +  ", path " +
			"LIMIT $3 "
	default:
		return "AND id > $2 " +
			"ORDER BY path[1] " + utils.DESC(q.Desc) +  ", path " +
			"LIMIT $3 (SELECT COUNT(1) FROM posts " +
				"WHERE array_length(path, 1) = 1 AND " +
				"thread_id = (SELECT id FROM t) AND " +
				"id > $2 " +
				"ORDER BY id DESC " +
				"LIMIT 1"
	}
}

func (r *Repository) GetPosts(threadId interface{}, q *query.GetThreadPosts) ([]domain.Post, error) {
	rows, err := r.db.Query(context.Background(),
		"WITH t AS " + getPostsCond(threadId, 1) +
			"SELECT id, path[(array_length(path, 1) - 1)], author, message, isEdited, forum_slug, thread_id, created " +
			"FROM posts " +
			"WHERE thread_id = (SELECT id FROM t) " +
			getPostsSort(q),)


}