package services

import (
	"context"
	"github.com/google/wire"
	"log"
	"simple_transaction_app/models"
	"simple_transaction_app/repositories"
	"time"
)

type IAccountsService interface {
	GetAccountByID(ctx context.Context, accountId int) (*models.AccountDto, error)
	CreateAccount(ctx context.Context, accountCreate models.AccountCreate) (*models.Account, error)
	UpdateAccount(ctx context.Context, accountUpdate models.AccountUpdate) (*models.Account, error)
}

type AccountsService struct {
	AccountsRepository repositories.IAccountsRepository
}

var NewAccountsService = wire.NewSet(
	wire.Struct(new(AccountsService), "*"),
	wire.Bind(new(IAccountsService), new(*AccountsService)))

func (s *AccountsService) GetAccountByID(ctx context.Context, accountId int) (*models.AccountDto, error) {
	account, err := s.AccountsRepository.GetAccountById(ctx, accountId)
	if err != nil {
		return nil, err
	}
	// convert model object to dto(that can be transported over network). We can create a common mapper to convert it.
	// In this conversion, return only required fields to routes. Doing it manually here instead of mapper.
	return &models.AccountDto{
		AccountID: account.AccountID,
		Balance:   account.Balance,
	}, err
}

func (s *AccountsService) CreateAccount(ctx context.Context, accountCreate models.AccountCreate) (*models.Account, error) {
	account, err := s.AccountsRepository.CreateAccount(ctx, accountCreate)
	if err != nil {
		log.Printf("Error creating account, error: %v", err)
		return nil, err
	}
	return account, nil
}

func (s *AccountsService) UpdateAccount(ctx context.Context, accountUpdate models.AccountUpdate) (*models.Account, error) {
	if accountUpdate.AccountID <= 0 {
		log.Printf("error: invalid account id")
	}
	if accountUpdate.UpdatedAt == nil {
		t := time.Now()
		accountUpdate.UpdatedAt = &t
	}
	_, err := s.AccountsRepository.UpdateAccount(ctx, accountUpdate)
	if err != nil {
		return nil, err
	}

	return s.AccountsRepository.GetAccountById(ctx, accountUpdate.AccountID)
}
