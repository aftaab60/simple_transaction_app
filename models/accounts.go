package models

import (
	"errors"
	"time"
)

// ignoring currency as this is static and same as per requirements
type Account struct {
	AccountID int       `json:"account_id" db:"account_id"`
	Balance   float64   `json:"balance" db:"balance"`
	Active    bool      `json:"active" db:"active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type AccountCreate struct {
	AccountID int     `json:"account_id" db:"id"`
	Balance   float64 `json:"initial_balance" db:"balance"`
}

func (a *AccountCreate) Validate() error {
	if a.AccountID <= 0 {
		return errors.New("invalid Account ID")
	}
	if a.Balance < 0 {
		return errors.New("balance must be greater than 0")
	}
	return nil
}

type AccountUpdate struct {
	AccountID int        `json:"account_id" db:"account_id"`
	Balance   *float64   `json:"balance" db:"balance"`
	Active    *bool      `json:"active" db:"active"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}

type AccountDto struct {
	AccountID int     `json:"account_id" db:"account_id"`
	Balance   float64 `json:"balance" db:"balance"`
}
