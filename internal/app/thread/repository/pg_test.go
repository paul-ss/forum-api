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



//func TestCreatePosts(t *testing.T) {
//	r := New(db)
//
//	req := []domain.PostCreate{
//		{
//			Parent: 0,
//			Author: "username",
//			Message: "hello",
//		},
//		{
//			Parent: 0,
//			Author: "username2",
//			Message: "hi",
//		},
//
//	}
//
//	p, err := r.CreatePosts(1986, req)
//
//	assert.Nil(t, err)
//	fmt.Println(p)
//}


//func TestGetThread(t *testing.T) {
//	r := New(db)
//
//	p, err := r.GetThread("slugurhouweiur")
//
//	assert.Nil(t, err)
//	fmt.Println(p)
//}


//func TestUpdateThread(t *testing.T) {
//	r := New(db)
//
//	p, err := r.UpdateThread("23", &domain.ThreadUpdate{
//		Message: "upd_message",
//		Title: "upd_title",
//	})
//
//	assert.Nil(t, err)
//	fmt.Println(p)
//}

//
//func TestGetPosts(t *testing.T) {
//	r := New(db)
//
//	p, err := r.GetPosts(1, &query.GetThreadPosts{
//		Limit: 10,
//		Since: 0,
//		Sort: "tree",
//		Desc: false,
//	})
//
//	assert.Nil(t, err)
//	for _, pp := range p {
//		fmt.Println(pp)
//	}
//}

func TestVoteThread(t *testing.T) {
	r := New(db)

	err := r.VoteThread(1, &domain.Vote{
		Voice: -1,
		Nickname: "username",
	})

	assert.Nil(t, err)

}