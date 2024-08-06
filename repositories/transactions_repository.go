package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/wire"
	"log"
	"simple_transaction_app/internal/db_manager"
	"simple_transaction_app/models"
	"strings"
)

type ITransactionsRepository interface {
	GetTransactionById(ctx context.Context, transactionId int) (*models.Transaction, error)
	CreateTransaction(ctx context.Context, c models.TransactionCreate) (*models.Transaction, error)
}

type TransactionsRepository struct {
	DB db_manager.IDB
}

var NewTransactionsRepository = wire.NewSet(
	wire.Struct(new(TransactionsRepository), "*"),
	wire.Bind(new(ITransactionsRepository), new(*TransactionsRepository)))

func (t *TransactionsRepository) GetTransactionById(ctx context.Context, transactionId int) (*models.Transaction, error) {
	cols := []string{
		"transaction_id",
		"source_account_id",
		"destination_account_id",
		"amount",
		"active",
		"created_at",
		"updated_at",
	}
	fields := strings.Join(cols, ",")
	query := fmt.Sprintf("SELECT %s FROM transactions WHERE transaction_id = $1", fields)

	transaction := models.Transaction{}
	err := t.DB.Get(ctx, &transaction, query, transactionId)
	if err != nil {
		log.Printf("Error: db error when getting transaction by id, %v", err)
		return nil, err
	}

	return &transaction, nil
}

func (t *TransactionsRepository) CreateTransaction(ctx context.Context, c models.TransactionCreate) (*models.Transaction, error) {
	// no need to insert values for active, created_at and updated_at. these are auto fields and gets added based on current datetime
	cols := []string{
		"source_account_id",
		"destination_account_id",
		"amount",
	}
	fields := strings.Join(cols, ",")
	query := fmt.Sprintf("INSERT INTO transactions(%s) VALUES ($1, $2, $3) RETURNING transaction_id", fields)

	var transactionId int
	err := t.DB.Get(ctx, &transactionId, query,
		c.SourceAccountID,
		c.DestinationAccountID,
		c.Amount,
	)
	if err != nil {
		log.Printf("Error: db error when inserting transaction, %v", err)
		return nil, err
	}
	if transactionId == 0 {
		return nil, errors.New("error inserting account transaction")
	}

	return t.GetTransactionById(ctx, transactionId)
}
