package models

import (
	"errors"
	"time"
)

type Transaction struct {
	TransactionID        int       `json:"transaction_id" db:"transaction_id"`
	SourceAccountID      int       `json:"source_account_id" db:"source_account_id"`
	DestinationAccountID int       `json:"destination_account_id" db:"destination_account_id"`
	Amount               float64   `json:"amount" db:"amount"`
	Active               bool      `json:"active" db:"active"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}

type TransactionCreate struct {
	SourceAccountID      int     `json:"source_account_id" db:"source_account_id"`
	DestinationAccountID int     `json:"destination_account_id" db:"destination_account_id"`
	Amount               float64 `json:"amount" db:"amount"`
}

func (t *TransactionCreate) Validate() error {
	if t.SourceAccountID <= 0 || t.DestinationAccountID <= 0 {
		return errors.New("invalid source or destination Account ID")
	}
	// validate source and destination accounts are not same
	if t.SourceAccountID == t.DestinationAccountID {
		return errors.New("source and destination Account IDs can not be same")
	}
	if t.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	return nil
}

type TransactionDto struct {
	TransactionID        int       `json:"transaction_id" db:"transaction_id"`
	SourceAccountID      int       `json:"source_account_id" db:"source_account_id"`
	DestinationAccountID int       `json:"destination_account_id" db:"destination_account_id"`
	Amount               float64   `json:"amount" db:"amount"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}
