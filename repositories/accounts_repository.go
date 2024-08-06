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

type IAccountsRepository interface {
	GetAccountById(ctx context.Context, accountId int) (*models.Account, error)
	CreateAccount(ctx context.Context, c models.AccountCreate) (*models.Account, error)
	UpdateAccount(ctx context.Context, accountUpdate models.AccountUpdate) (bool, error)
}

type AccountsRepository struct {
	DB db_manager.IDB
}

var NewAccountsRepository = wire.NewSet(
	wire.Struct(new(AccountsRepository), "*"),
	wire.Bind(new(IAccountsRepository), new(*AccountsRepository)))

func (a *AccountsRepository) GetAccountById(ctx context.Context, accountId int) (*models.Account, error) {
	cols := []string{
		"account_id",
		"balance",
		"active",
		"created_at",
		"updated_at",
	}
	fields := strings.Join(cols, ",")
	// upon deleting, mark active = false instead of making hard delete. Get calls filters only active accounts
	query := fmt.Sprintf("SELECT %s FROM accounts WHERE account_id = $1 and active = true", fields)

	account := models.Account{}
	err := a.DB.Get(ctx, &account, query, accountId)
	if err != nil {
		log.Printf("Error: db error when getting account by id, %v", err)
		return nil, err
	}

	return &account, nil
}

// Note: Ideally when inserting, we don't need to provide account_id as this is incremented and auto fields in Db already.
// But adding here as this is part of requirement
func (a *AccountsRepository) CreateAccount(ctx context.Context, c models.AccountCreate) (*models.Account, error) {
	// no need to insert values for active, created_at and updated_at. these are auto fields and gets added based on current datetime
	cols := []string{
		"account_id",
		"balance",
	}
	fields := strings.Join(cols, ",")
	query := fmt.Sprintf("INSERT INTO accounts(%s) VALUES ($1, $2)", fields)

	response, err := a.DB.Exec(ctx, query,
		c.AccountID,
		c.Balance,
	)
	if err != nil {
		log.Printf("Error: db error when inserting account, %v", err)
		return nil, err
	}
	numRows, err := response.RowsAffected()
	if err != nil {
		return nil, err
	}
	if numRows == 0 {
		return nil, errors.New("error inserting account")
	}

	return a.GetAccountById(ctx, c.AccountID)
}

func (a *AccountsRepository) UpdateAccount(ctx context.Context, u models.AccountUpdate) (bool, error) {
	cols := []string{
		"balance = COALESCE($1, balance)",
		"active = COALESCE($2, active)",
		"updated_at = COALESCE($3, updated_at)",
	}
	fields := strings.Join(cols, ",")
	query := fmt.Sprintf("UPDATE accounts SET %s WHERE account_id = $4", fields)

	response, err := a.DB.Exec(ctx, query,
		u.Balance,
		u.Active,
		u.UpdatedAt,
		u.AccountID,
	)
	if err != nil {
		log.Printf("Error: db error when updating account, %v", err)
		return false, err
	}
	updatedRows, err := response.RowsAffected()
	if err != nil {
		return false, err
	}

	return updatedRows > 0, nil
}
