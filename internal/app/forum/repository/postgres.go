package repository

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	config "github.com/paul-ss/forum-api/configs/go"
	"github.com/paul-ss/forum-api/internal/domain"
	domainErr "github.com/paul-ss/forum-api/internal/domain/errors"
	query "github.com/paul-ss/forum-api/internal/domain/query"
	"github.com/paul-ss/forum-api/internal/utils"
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


func (r *Repository) StoreThread(slug string, tc domain.ThreadCreate) (*domain.Thread, error) {
	t := domain.Thread{}
	forumId := 0
	err := r.db.QueryRow(context.Background(),
		"WITH f AS (SELECT title, id FROM forums WHERE slug = $1) " +
			"INSERT INTO threads (title, author, forum_tittle, forum_id, message, slug, created) " +
			"VALUES ($2, $3, (SELECT title FROM f), (SELECT id FROM f), $4, $5, $6) " +
			"RETURNING id, title, author, forum_tittle, message, slug, created, votes, (SELECT id FROM f) ",
	slug, tc.Title, tc.Author, tc.Message, utils.RandomSlug(), utils.GetCurrentTime(tc.Created)).
		Scan(&t.Id, &t.Title, &t.Author, &t.Forum, &t.Message, &t.Slug, &t.Created, &t.Votes, &forumId)

	if err != nil {
		config.Lg("forum_repo", "StoreThread").Info(err.Error())
		er := r.db.QueryRow(context.Background(),
			"SELECT id, title, author, forum_tittle, message, slug, created, votes " +
				"FROM threads " +
				"WHERE forum_id = $1 ",
			forumId).Scan(&t.Id, &t.Title, &t.Author, &t.Forum, &t.Message, &t.Slug, &t.Created, &t.Votes)

		if er != nil {
			config.Lg("forum_repo", "StoreThread").Error("Select: ", err.Error())
			return nil, er
		}

		return &t, domainErr.DuplicateKeyError
	}

	return &t, nil
}

func (r *Repository) GetUsers(slug string, q *query.GetForumUsers) ([]domain.User, error) {
	count := 0
	if err := r.db.QueryRow(context.Background(),
		"SELECT COUNT(1) " +
			"FROM forums " +
			"WHERE slug = $1", slug).Scan(&count); err != nil {

		config.Lg("forum_repo", "GetUsers").Error("Query1: ", err.Error())
		return nil, err
	}

	if count == 0 {
		config.Lg("forum_repo", "GetUsers").Error("Query1: count == 0")
		return nil, domainErr.NotExists
	}

	rows, err := r.db.Query(context.Background(),
		"WITH f AS (SELECT title, id FROM forums WHERE slug = $1) " +
		"SELECT u.nickname, u.fullname, u.about, u.email " +
		"FROM " +
		"(SELECT DISTINCT author FROM threads WHERE forum_id = (SELECT id from f) " +
		"UNION " +
		"SELECT DISTINCT author FROM posts WHERE forum_id = (SELECT id from f)) AS a " +
		"JOIN users u ON a.author = u.nickname " +
		"WHERE lower(nickname) > lower ($2)" +
		"ORDER BY lower(nickname) " + utils.DESC(q.Desc) +
		"LIMIT ($3) ",
		slug, q.Since, q.Limit)

	// NOTE: maybe if it's DESC, you should change > to < ?

	if err != nil {
		config.Lg("forum_repo", "GetUsers").Error("Query: ", err.Error())
		return nil, err
	}

	defer rows.Close()

	users := []domain.User{}
	for rows.Next() {
		u := domain.User{}

		err := rows.Scan(&u.Nickname, &u.FullName, &u.About, &u.Email)
		if err != nil {
			config.Lg("forum_repo", "GetUsers").Error("Scan: ", err.Error())
			return nil, err
		}

		users = append(users, u)
	}

	return users, nil
}


func (r *Repository) GetThreads(slug string, q *query.GetForumThreads) ([]domain.Thread, error) {
	count := 0
	if err := r.db.QueryRow(context.Background(),
		"SELECT COUNT(1) " +
		"FROM forums " +
		"WHERE slug = $1", slug).Scan(&count); err != nil {

		config.Lg("forum_repo", "GetThreads").Error("Query1: ", err.Error())
		return nil, err
	}

	if count == 0 {
		config.Lg("forum_repo", "GetUsers").Error("Query1: count == 0")
		return nil, domainErr.NotExists
	}

	rows, err := r.db.Query(context.Background(),
		"WITH f AS (SELECT id FROM forums WHERE slug = $1) " +
			"SELECT id, title, author, forum_title, message, votes, slug, created " +
			"FROM threads " +
			"WHERE forum_id = (SELECT id FROM f) AND created > $2 " +
			"ORDER BY created " + utils.DESC(q.Desc) +
			"LIMIT $3 ",
		slug, q.Since, q.Limit)

	// NOTE: maybe if it's DESC, you should change > to < ?

	if err != nil {
		config.Lg("forum_repo", "GetThreads").Error("Query2: ", err.Error())
		return nil, err
	}

	defer rows.Close()

	threads := []domain.Thread{}
	for rows.Next() {
		t := domain.Thread{}

		err := rows.Scan(&t.Id, &t.Title, &t.Author, &t.Forum, &t.Message, &t.Votes, &t.Slug, &t.Created)
		if err != nil {
			config.Lg("forum_repo", "GetThreads").Error("Scan: ", err.Error())
			return nil, err
		}

		threads = append(threads, t)
	}

	return threads, nil
}