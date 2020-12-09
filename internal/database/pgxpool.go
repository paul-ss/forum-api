package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	config "github.com/paul-ss/forum-api/configs"
)

type DB struct {
	dbPool *pgxpool.Pool
	config *config.ConfDB
}

func NewDB() *DB {
	return &DB{
		config: &config.Conf.Db,
	}
}

func (db *DB) Open() error {
	conf, err := pgxpool.ParseConfig(fmt.Sprintf(
		"user=%s password=%s host=%s dbname=%s sslmode=%s pool_max_conns=%s",
		db.config.Postgres.Username,
		db.config.Postgres.Password,
		db.config.Postgres.Host,
		db.config.Postgres.DbName,
		db.config.Postgres.SslMode,
		db.config.Postgres.MaxConn,
	))
	if err != nil {
		return err
	}

	db.dbPool, err = pgxpool.ConnectConfig(context.Background(), conf)

	if err != nil {
		return err
	}
	return nil
}

func (db *DB) Close() {
	db.dbPool.Close()
}

