package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/wire"
	"log"
	"simple_transaction_app/internal/db_manager"
	"simple_transaction_app/models"
	"simple_transaction_app/repositories"
	"simple_transaction_app/utils"
)

type ITransactionsService interface {
	GetTransactionById(ctx context.Context, transactionId int) (*models.TransactionDto, error)
	CreateTransaction(ctx context.Context, transactionCreate models.TransactionCreate) (*models.Transaction, error)
}

type TransactionsService struct {
	TransactionsRepository repositories.ITransactionsRepository
	AccountsRepository     repositories.IAccountsRepository
	// Notice this explicit transaction supported Db object composition in transaction service
	Db db_manager.IDbTxBeginner
}

var NewTransactionsService = wire.NewSet(
	wire.Struct(new(TransactionsService), "*"),
	wire.Bind(new(ITransactionsService), new(*TransactionsService)))

func (s *TransactionsService) GetTransactionById(ctx context.Context, transactionId int) (*models.TransactionDto, error) {
	transaction, err := s.TransactionsRepository.GetTransactionById(ctx, transactionId)
	if err != nil {
		return nil, err
	}
	return &models.TransactionDto{
		TransactionID:        transaction.TransactionID,
		SourceAccountID:      transaction.SourceAccountID,
		DestinationAccountID: transaction.DestinationAccountID,
		Amount:               transaction.Amount,
		CreatedAt:            transaction.CreatedAt,
		UpdatedAt:            transaction.UpdatedAt,
	}, nil
}

// create transaction needs to be under a transaction to make data consistency
func (s *TransactionsService) CreateTransaction(ctx context.Context, transactionCreate models.TransactionCreate) (*models.Transaction, error) {
	var transaction *models.Transaction

	onRollback := func(err error) {
		log.Printf("transaction is rolling back due to error: %v", err)
	}
	// wrap all operations in a traction
	err := db_manager.WrapInTransaction(ctx, s.Db, func(ctx context.Context) error {

		// get account info with source account id
		sourceAccount, err := s.AccountsRepository.GetAccountById(ctx, transactionCreate.SourceAccountID)
		if err != nil {
			log.Printf("Error: Invalid source account %d", transactionCreate.SourceAccountID)
			return err
		}
		// validate balance in source account
		if sourceAccount.Balance < transactionCreate.Amount {
			log.Printf("Error: insufficient fund in source account %d", transactionCreate.SourceAccountID)
			return errors.New(fmt.Sprintf("insufficient fund in source account %d", transactionCreate.SourceAccountID))
		}

		// get account info with destination account
		destinationAccount, err := s.AccountsRepository.GetAccountById(ctx, transactionCreate.DestinationAccountID)
		if err != nil {
			log.Printf("Error: Invalid destination account %d", transactionCreate.DestinationAccountID)
			return err
		}

		// update source & destination account balances
		sourceBalance := sourceAccount.Balance - transactionCreate.Amount
		destinationBalance := destinationAccount.Balance + transactionCreate.Amount
		ok, err := s.AccountsRepository.UpdateAccount(ctx, models.AccountUpdate{
			AccountID: sourceAccount.AccountID,
			Balance:   &sourceBalance,
			UpdatedAt: utils.GetCurrentTimePtr(),
		})
		if err != nil || !ok {
			return errors.New("error updating source account")
		}
		ok, err = s.AccountsRepository.UpdateAccount(ctx, models.AccountUpdate{
			AccountID: destinationAccount.AccountID,
			Balance:   &destinationBalance,
			UpdatedAt: utils.GetCurrentTimePtr(),
		})
		if err != nil || !ok {
			return errors.New("error updating destination account")
		}

		// insert a new record in transactions table
		if transaction, err = s.TransactionsRepository.CreateTransaction(ctx, transactionCreate); err != nil {
			return err
		}

		return nil
	}, onRollback)
	if err != nil {
		fmt.Printf("error creating account transaction: %v", err)
		return nil, err
	}

	return transaction, nil
}
