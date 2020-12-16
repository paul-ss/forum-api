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

func valuesPosts(i int, req []domain.PostCreate, args *[]interface{}) string {
	query := []string{}
	pIdxAdd := 0
	for _, t := range req {
		query = append(query,
			"( " +
			fmt.Sprintf("(WITH par AS (SELECT path FROM posts WHERE id = $%d) ", i) +
			"SELECT CASE WHEN ((SELECT path FROM par) IS NOT NULL) THEN " +
			fmt.Sprintf("(SELECT path FROM par) || (last_value + %d) ", pIdxAdd) +
			fmt.Sprintf("WHEN ($%d < 1) THEN ", i) +
			fmt.Sprintf("ARRAY[last_value + %d]  ", pIdxAdd) +
			"ELSE null END FROM posts_id_seq), " +
			fmt.Sprintf("$%d, $%d, ", i + 1, i + 2) +
			"(SELECT slug FROM t), " +
			"(SELECT forum_id FROM t), " +
			"$1 " +
			")", ",")

		*args = append(*args, t.Parent, t.Author, t.Message)
		i += 3
		pIdxAdd += 1
	}

	query = query[:len(query) - 1]
	return strings.Join(query, "")
}

func (r *Repository) CreatePostsById(threadId int32, req []domain.PostCreate) ([]domain.Post, error) {
	args := []interface{}{}
	args = append(args, threadId)
	rows, err := r.db.Query(context.Background(),
		"WITH t AS (SELECT forum_id, slug FROM threads WHERE id = $1)," +
				"seq AS (SELECT last_value FROM posts_id_seq) " +
			"INSERT INTO posts (path, author, message, forum_slug, forum_id, thread_id) " +
			"VALUES  " + valuesPosts(2, req, &args) +
			"RETURNING id, path[(array_length(path, 1) - 1)], author, message, isEdited, forum_slug, thread_id, created ",
			args...)

	if err != nil {
		config.Lg("thread_repo", "CreatePostsById").Error("Query: ", err.Error())
		return nil, err
	}

	defer rows.Close()

	posts := []domain.Post{}
	for rows.Next() {
		p := domain.Post{}
		slug := sql.NullString{}
		parent := sql.NullInt64{}

		err := rows.Scan(&p.Id, &parent, &p.Author, &p.Message, &p.IsEdited, &slug, &p.ThreadId, &p.Created)
		if err != nil {
			config.Lg("thread_repo", "CreatePostsById").Error("Scan: ", err.Error())
			return nil, err
		}

		p.ForumSlug = slug.String
		p.Parent = parent.Int64
		posts = append(posts, p)
	}

	if err := rows.Err(); err != nil {
		config.Lg("thread_repo", "CreatePostsById").Error("Rows: ", err.Error())
		return nil, err
	}

	return posts, nil
}
