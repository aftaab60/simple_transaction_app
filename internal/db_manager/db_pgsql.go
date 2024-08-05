package db_manager

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"sync"
)

var once sync.Once
var db *DB

func InitPgsqlConnection(cfg IConfig) *DB {
	once.Do(func() {
		connectStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=%s",
			cfg.GetHost(),
			cfg.GetPort(),
			cfg.GetUsername(),
			cfg.GetPassword(),
			cfg.GetDatabase(),
			"%27Asia/Singapore%27")

		sqlDb, err := sqlx.Open("postgres", connectStr)
		if err != nil {
			log.Fatal(err)
		}

		db = &DB{
			SqlDb: sqlDb,
			cfg:   cfg,
		}
	})
	return db
}

// map rows into destination (slice of elements)
func (db *DB) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return handleDBError(db.SqlDb.SelectContext(ctx, dest, query, args...))
}

// map row into destination
func (db *DB) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return handleDBError(db.SqlDb.GetContext(ctx, dest, query, args...))
}

func (db *DB) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	result, err := db.SqlDb.ExecContext(ctx, query, args...)
	return result, handleDBError(err)
}

func (db *DB) BeginTxx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := db.SqlDb.BeginTxx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &Tx{Tx: *tx}, nil
}
