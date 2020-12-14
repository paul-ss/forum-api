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



//func TestGetPost(t *testing.T) {
//	r := New(db)
//
//	p, err := r.GetPost(6)
//
//	assert.Nil(t, err)
//	fmt.Println(p)
//}


//func TestGetAuthor(t *testing.T) {
//	r := New(db)
//
//	p, err := r.GetAuthor(1)
//
//	assert.Nil(t, err)
//	fmt.Println(p)
//}

//func TestGetThread(t *testing.T) {
//	r := New(db)
//
//	p, err := r.GetThread(3)
//
//	assert.Nil(t, err)
//	fmt.Println(p)
//}

//func TestGetForum(t *testing.T) {
//	r := New(db)
//
//	p, err := r.GetForum(1)
//
//	assert.Nil(t, err)
//	fmt.Println(p)
//}

func TestUpdPost(t *testing.T) {
	r := New(db)

	p, err := r.UpdatePost(4, &domain.PostUpdate{Message: "upd"})

	assert.Nil(t, err)
	fmt.Println(p)
}