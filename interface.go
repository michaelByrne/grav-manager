package main

import (
	"github.com/eaigner/jet"
	"go.uber.org/zap"
)

// AccountService holds account methods
type AccountService interface {
	RegisterUser(user User) error
	UpgradePlan(limit int64, accountID string) error
	GetUserCount(id string) (int, error)
	GetAccount(id string) (Account, error)
	GetUsers(accountID string) ([]User, error)
}

// Account is an account
type Account struct {
	ID    string `json:"id"`
	Limit int64  `json:"limit"`
}

// User is a user
type User struct {
	ID        string `json:"user_id"`
	Timestamp string `json:"timestamp"`
	AccountID string `json:"account_id"`
}

type accountService struct {
	db     jet.Db
	logger *zap.SugaredLogger
}
