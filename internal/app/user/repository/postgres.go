package repository

import (
	"context"
	"github.com/jackc/pgconn"
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
		"WHERE nickname = $1 OR email = $2 ",
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
			"WHERE nickname = $1",
			username).Scan(&u.Nickname, &u.FullName, &u.About, &u.Email)

	if err != nil {
		config.Lg("user_repo", "GetUser").Error("Query : ", err.Error())
		return nil, err
	}

	return &u, nil
}

func (r *Repository) UpdateUser(username string, req *domain.UserCreate) (*domain.User, error) {
	u := domain.User{}

	err := r.db.QueryRow(context.Background(),
		"UPDATE users " +
		"SET fullname = $2, about = $3, email = $4 " +
		"WHERE nickname = $1 " +
		"RETURNING nickname, fullname, about, email ",
		username, req.FullName, req.About, req.Email).Scan(&u.Nickname, &u.FullName, &u.About, &u.Email)

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