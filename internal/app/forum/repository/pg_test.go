package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	config "github.com/paul-ss/forum-api/configs/go"
	"github.com/paul-ss/forum-api/internal/domain/query"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
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





//func TestStore(t *testing.T) {
//	r := New(db)
//
//	f := domain.Forum{
//		User: "username2",
//		Title: "f_title2",
//		Slug: "slug",
//	}
//
//	fres, err := r.StoreForum(&f)
//
//	assert.Nil(t, err)
//	fmt.Println(fres)
//
//}

//func TestGet(t *testing.T) {
//	r := New(db)
//
//
//	fres, err := r.GetForumBySlug("superslug")
//
//	assert.Nil(t, err)
//	fmt.Println(fres)
//
//}

//
//func TestStoreThread(t *testing.T) {
//	r := New(db)
//
//	th := domain.ThreadCreate{
//		Title: "thread_title",
//		Author: "username",
//		Message: "msg",
//		//Created: time.Now().Add(time.Hour),
//
//	}
//
//	thr, err := r.StoreThread("superslu", th)
//
//	assert.Nil(t, err)
//	fmt.Println(thr)
//}

//func TestGetUsers(t *testing.T) {
//	r := New(db)
//
//	q := query.GetForumUsers{
//		Limit: 0,
//		Since: "usernam",
//		Desc: true,
//	}
//
//
//	us, err := r.GetUsers("superslug", &q)
//
//
//	assert.Nil(t, err)
//	fmt.Println(us)
//}

func TestGetThreads(t *testing.T) {
	r := New(db)

	q := query.GetForumThreads{
		Limit: 10,
		Since: time.Now(),
		Desc: false,
	}

	us, err := r.GetThreads("superslu", &q)


	assert.Nil(t, err)
	fmt.Println(us)
}