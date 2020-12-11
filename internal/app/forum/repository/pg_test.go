package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	config "github.com/paul-ss/forum-api/configs/go"
	"github.com/paul-ss/forum-api/internal/domain"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var db *pgxpool.Pool

func TestMain(m *testing.M) {
	config.Conf = config.NewConfig()
	conf, err := pgxpool.ParseConfig(fmt.Sprintf(
		"user=%s password=%s host=%s dbname=%s sslmode=%s pool_max_conns=%s",
		config.Conf.Db.Postgres.Username,
		config.Conf.Db.Postgres.Password,
		config.Conf.Db.Postgres.Host,
		config.Conf.Db.Postgres.DbName,
		config.Conf.Db.Postgres.SslMode,
		config.Conf.Db.Postgres.MaxConn,
	))
	if err != nil {
		fmt.Println("Error ", err.Error())
	}

	db, err = pgxpool.ConnectConfig(context.Background(), conf)
	code := m.Run()

	db.Close()
	os.Exit(code)
}





func TestStore(t *testing.T) {
	r := New(db)

	f := domain.Forum{
		User: "username",
		Title: "f_title",
		Slug: "slug",
	}

	fres, err := r.StoreForum(&f)

	assert.Nil(t, err)
	fmt.Println(fres)

}


