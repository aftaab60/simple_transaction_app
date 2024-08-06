package db_manager

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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
	sqldb := db.SqlDb
	if tx := GetTransactionContext(ctx); tx != nil {
		sqldb = tx.DB.SqlDb
	}
	return handleDBError(sqldb.SelectContext(ctx, dest, query, args...))
}

// map row into destination
func (db *DB) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	sqldb := db.SqlDb
	if tx := GetTransactionContext(ctx); tx != nil {
		sqldb = tx.DB.SqlDb
	}
	return handleDBError(sqldb.GetContext(ctx, dest, query, args...))
}

func (db *DB) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	sqldb := db.SqlDb
	if tx := GetTransactionContext(ctx); tx != nil {
		sqldb = tx.DB.SqlDb
	}
	result, err := sqldb.ExecContext(ctx, query, args...)
	return result, handleDBError(err)
}

func (db *DB) BeginTxx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := db.SqlDb.BeginTxx(ctx, opts)
	if err != nil {
		return nil, errors.New("can not start transaction")
	}
	return &Tx{
		DB: db,
		Tx: tx,
	}, nil
}
