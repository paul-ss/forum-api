package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	config "github.com/paul-ss/forum-api/configs/go"
	"github.com/paul-ss/forum-api/internal/domain"
	domainErr "github.com/paul-ss/forum-api/internal/domain/errors"
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

func updateUserFields(args *[]interface{}, req *domain.UserUpdate, startIdx int) string {
	query := []string{}
	if req.FullName != nil {
		query = append(query, fmt.Sprintf(" fullname = $%d", startIdx), ", ")
		startIdx += 1
		*args = append(*args, *req.FullName)
	}

	if req.About != nil {
		query = append(query, fmt.Sprintf(" about = $%d", startIdx), ", ")
		startIdx += 1
		*args = append(*args, *req.About)
	}

	if req.Email != nil {
		query = append(query, fmt.Sprintf(" email = $%d", startIdx), ", ")
		startIdx += 1
		*args = append(*args, *req.Email)
	}

	if len(query) == 0 {
		query = append(query, " email = email ")
	} else {
		query = query[:len(query)-1]
	}

	return strings.Join(query, "")
}

func (r *Repository) CreateUser(username string, req *domain.UserCreate) ([]domain.User, error) {
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		config.Lg("user_repo", "CreateUser").Error("Begin: ", err.Error())
		return nil, err
	}
	defer tx.Rollback(context.Background())

	rows, err := r.db.Query(context.Background(),
		"SELECT nickname, fullname, about, email " +
		"FROM users " +
		"WHERE LOWER(nickname) = LOWER($1) OR LOWER(email) = LOWER($2) ",
		username, req.Email)

	if err != nil {
		config.Lg("user_repo", "CreateUser").Error("Query 1: ", err.Error())
		return nil, err
	}

	users := []domain.User{}
	for rows.Next() {

		u := domain.User{}
		err := rows.Scan(&u.Nickname, &u.FullName, &u.About, &u.Email)
		if err != nil {
			config.Lg("user_repo", "CreateUser").Error("Scan: ", err.Error())
			return nil, err
		}

		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		config.Lg("user_repo", "CreateUser").Error("Rows: ", err.Error())
		return nil, err
	}

	if len(users) > 0 {
		config.Lg("user_repo", "CreateUser").Error("len = 0")
		return users, domainErr.AlreadyExists
	}

	if _, err := tx.Exec(context.Background(),
		"INSERT INTO users (nickname, fullname, about, email) " +
			"VALUES ($1, $2, $3, $4) ",
			username, req.FullName, req.About, req.Email); err != nil {
		config.Lg("user_repo", "CreateUser").Error("Query 2: ", err.Error())
		return nil, err
	}


	if err := tx.Commit(context.Background()); err != nil {
		config.Lg("user_repo", "CreateUser").Error("Commit : ", err.Error())
		return nil, err
	}

	users = append(users, domain.User{
		Nickname: username,
		FullName: req.FullName,
		Email: req.Email,
		About: req.About,
	})

	return users, nil
}

func (r *Repository) GetUser(username string) (*domain.User, error) {
	u := domain.User{}
	err := r.db.QueryRow(context.Background(),
		"SELECT nickname, fullname, about, email " +
			"FROM users " +
			"WHERE LOWER(nickname) = LOWER($1)",
			username).Scan(&u.Nickname, &u.FullName, &u.About, &u.Email)

	if err != nil {
		config.Lg("user_repo", "GetUser").Error("Query : ", err.Error())
		return nil, err
	}

	return &u, nil
}

func (r *Repository) UpdateUser(username string, req *domain.UserUpdate) (*domain.User, error) {
	u := domain.User{}
	args := []interface{}{}
	args = append(args, username)

	err := r.db.QueryRow(context.Background(),
		"UPDATE users " +
		"SET " + updateUserFields(&args, req, 2) +
		"WHERE LOWER(nickname) = LOWER($1) " +
		"RETURNING nickname, fullname, about, email ",
		args...).Scan(&u.Nickname, &u.FullName, &u.About, &u.Email)

	if err != nil {
		switch e := err.(type) {
		case *pgconn.PgError:
			if e.SQLState() == "23505" {
				return nil, domainErr.AlreadyExists
			}
		}

		config.Lg("user_repo", "GetUser").Error("Query : ", err.Error())
		return nil, err
	}

	return &u, nil
}


// STATS ====

func (r *Repository) ClearAll() error {
	if _, err := r.db.Exec(context.Background(),
		"TRUNCATE TABLE users, forums, threads, posts, votes, stats "); err != nil {
		config.Lg("user_repo", "ClearAll").Error("Exec : ", err.Error())
		return err
	}

	return nil
}


func (r *Repository) GetStats() (*domain.Status, error) {
	s := domain.Status{}
	err := r.db.QueryRow(context.Background(),
		"SELECT " +
				"CASE WHEN u.is_called THEN u.last_value ELSE 0 END, " +
				"CASE WHEN f.is_called THEN f.last_value ELSE 0 END, " +
				"CASE WHEN t.is_called THEN t.last_value ELSE 0 END, " +
				"CASE WHEN p.is_called THEN p.last_value ELSE 0 END " +
			"FROM users_id_seq u, forums_id_seq f, threads_id_seq t, pidseq p ").
		Scan(&s.User, &s.Forum, &s.Thread, &s.Post)

	if err != nil {
		config.Lg("user_repo", "GetStats").Error(err.Error())
		return nil, err
	}

	return &s, nil
}