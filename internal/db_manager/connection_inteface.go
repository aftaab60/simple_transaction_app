package db_manager

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"log"
	"simple_transaction_app/consts"
)

type IConfig interface {
	GetHost() string
	GetPort() uint
	GetUsername() string
	GetPassword() string
	GetDatabase() string
}

type IDB interface {
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	//ExecContext(ctx context.Context, s string, i ...interface{}) (sql.Result, error)
	//QueryContext(ctx context.Context, s string, i ...interface{}) (*sql.Rows, error)
	//QueryRowContext(ctx context.Context, s string, i ...interface{}) *sql.Row
}

type DB struct {
	SqlDb *sqlx.DB
	cfg   IConfig
}

// Extended DB object with begin transaction method.
// This separates having two DB objects. One with transaction and other without transaction starter support
type IDbTxBeginner interface {
	IDB
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (*Tx, error)
}

// ITx is db transaction interface.
// Using interface, we can change transaction to any type like optimistic, pessimistic transaction locking or both later.
// In this project, using pessimistic locking based transaction.
//type ITx interface {
//	Commit() error
//	Rollback() error
//}

type Tx struct {
	sqlx.Tx
}

// handle DB errors based on error type and application need
// Here, logging and returning generic error to callers
// This avoids exposing db/table details in error message
func handleDBError(err error) error {
	if err == nil {
		return nil
	}
	log.Println("DB Error: ", err)
	return errors.New(consts.UnexpectedDbError)
}
