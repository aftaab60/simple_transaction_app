package db_manager

import (
	"context"
	"errors"
	"log"
)

type key string

const (
	transactionKey key = "transaction_key"
)

// WrapInTransaction Receiving IDbTxBeginner interface. Same transacion wrapper can be used for other types of databases too
func WrapInTransaction(ctx context.Context, db IDbTxBeginner, f func(ctx context.Context) error, onRollback func(error)) (err error) {
	// create transaction if not present in context already
	tx := GetTransactionContext(ctx)
	if tx == nil {
		tx, err = db.BeginTxx(ctx, nil)
		if err != nil {
			return err
		}
		ctx = context.WithValue(ctx, transactionKey, tx)

		// commit or rollback transaction based on error, executed after function f(ctx)
		defer func() {
			if r := recover(); r != nil {
				RollbackTransaction(errors.New("panic error"), *tx, onRollback)
				panic(r)
			}
			if err != nil {
				RollbackTransaction(err, *tx, onRollback)
			} else {
				CommitTransaction(*tx)
			}
		}()
	}

	err = f(ctx)
	return err
}

func RollbackTransaction(err error, tx Tx, onRollback func(error)) {
	_err := tx.Tx.Rollback()
	if _err != nil {
		log.Println("WARN: Cannot rollback Transaction")
	}
	if onRollback != nil {
		onRollback(err)
	}
	log.Println("Transaction rollbacked")
}

func CommitTransaction(tx Tx) {
	err := tx.Tx.Commit()
	if err != nil {
		log.Println("WARN: Cannot commit Transaction")
	}
	log.Println("Transaction committed")
}

func GetTransactionContext(ctx context.Context) *Tx {
	if tx, ok := ctx.Value(transactionKey).(*Tx); ok {
		return tx
	}
	return nil
}
