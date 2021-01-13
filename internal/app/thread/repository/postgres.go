package repository

import (
	"context"
	"database/sql"
	"errors"
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
		"SELECT id, title, author, forum_slug, message, votes, slug, created " +
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
		"RETURNING id, title, author, forum_slug, message, votes, slug, created ",
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
		return fmt.Sprintf(" (SELECT id FROM threads WHERE slug = $%d) ", param)
	} else {
		return fmt.Sprintf(" (SELECT id FROM threads WHERE id = $%d) ", param)
	}
}


func getPostsSort(q *query.GetThreadPosts) string {
	switch q.Sort {
	case "flat":
		// desc + limit returns strange result
		return "AND id > $2 " +
			"ORDER BY created " + utils.DESC(q.Desc) +
			"LIMIT $3 "
	case "tree":
		return "AND id > $2 " +
			"ORDER BY p.path[1] " + utils.DESC(q.Desc) +  ", p.path " +
			"LIMIT $3 "
	case "parent_tree" :
		return "AND path[1] >= (SELECT min FROM ls)  AND path[1] <= (SELECT max FROM ls)" +
			"ORDER BY p.path[1] " + utils.DESC(q.Desc) +  ", p.path "
	default:
		return "error"
	}
}

func getPostsWith(q *query.GetThreadPosts) string {
	if q.Sort == "parent_tree" {
		return " , ls AS " +
			"(SELECT min(id) AS min, max(id) AS max " +
			"FROM (SELECT id FROM posts " +
			"WHERE thread_id = (SELECT id FROM t) AND array_length(path, 1) = 1 AND id > $2 " +
			"LIMIT $3) as unused_alias) "
	}

	return " "
}

func (r *Repository) GetPosts(threadId interface{}, q *query.GetThreadPosts) ([]domain.Post, error) {
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())

	counter := 0
	if err := tx.QueryRow(context.Background(),
			"WITH t AS " + getPostsCond(threadId, 1) +
			"SELECT COUNT(*) FROM t ",
			threadId).Scan(&counter); err != nil {
		return nil, errors.New("Query 1: " + err.Error())
	}

	if counter == 0 {
		return nil, domainErr.NotExists
	}


	rows, err := tx.Query(context.Background(),
		"WITH t AS " + getPostsCond(threadId, 1) +
			getPostsWith(q) +
			"SELECT id, path[(array_length(path, 1) - 1)], author, message, isEdited, forum_slug, thread_id, created, " +
			"CASE WHEN (SELECT id FROM t) IS NOT NULL THEN 1 ELSE 0 END as exists " +
			"FROM posts p " +
			"WHERE thread_id = (SELECT id FROM t) " +
			getPostsSort(q),
			threadId, q.Since, q.Limit)

	if err != nil {
		config.Lg("thread_repo", "GetPosts").Error("Query: ", err.Error())
		return nil, err
	}



	posts := []domain.Post{}
	for rows.Next() {
		p := domain.Post{}
		slug := sql.NullString{}
		parent := sql.NullInt64{}

		exists := 0

		err := rows.Scan(&p.Id, &parent, &p.Author, &p.Message, &p.IsEdited, &slug, &p.ThreadId, &p.Created, &exists)
		if err != nil {
			config.Lg("thread_repo", "GetPosts").Error("Scan: ", err.Error())
			return nil, err
		}

		if exists != 1 {
			return nil, domainErr.NotExists
		}

		p.ForumSlug = slug.String
		p.Parent = parent.Int64
		posts = append(posts, p)
	}

	if err := rows.Err(); err != nil {
		config.Lg("thread_repo", "GetPosts").Error("Rows: ", err.Error())
		return nil, createPostError(err)
	}

	if err := tx.Commit(context.Background()); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *Repository) VoteThread(threadId interface{}, req *domain.Vote) (*domain.Thread, error) {
	tx, er := r.db.Begin(context.Background())
	if er != nil {
		config.Lg("thread_repo", "VoteThread").Error("Begin: " + er.Error())
		return nil, er
	}
	defer tx.Rollback(context.Background())

	thrIdInt := 0
	err := tx.QueryRow(context.Background(),
		"WITH t AS " + getPostsCond(threadId, 1) +
			"INSERT into votes (nickname, thread_id, voice) " +
			"VALUES ($2, (SELECT id FROM t), $3) " +
			"ON CONFLICT(thread_id, nickname) DO UPDATE " +
			"SET voice = $3 " +
			"RETURNING thread_id ",
			threadId, req.Nickname, req.Voice).Scan(&thrIdInt)

	if err != nil {
		config.Lg("thread_repo", "VoteThread").Error("Query 1: " + err.Error())
		return nil, err
	}

	t := domain.Thread{}
	if err := tx.QueryRow(context.Background(),
		"SELECT id, title, author, forum_slug, message, votes, slug, created " +
			"FROM threads " +
			"WHERE id = $1 ", thrIdInt).
			Scan(&t.Id, &t.Title, &t.Author, &t.Forum, &t.Message, &t.Votes, &t.Slug, &t.Created); err != nil {
		config.Lg("thread_repo", "VoteThread").Error("Query 2: " + err.Error())
		return nil, err
	}

	if err := tx.Commit(context.Background()); err != nil {
		config.Lg("thread_repo", "VoteThread").Error("Commit: " + err.Error())
		return nil, err
	}


	return &t, nil
}