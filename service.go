package main

import (
	"errors"
	"log"

	"github.com/eaigner/jet"
	"github.com/lib/pq"
)

// NewAccountService returns a new account service
func NewAccountService() (AccountService, error) {
	// TODO: make URL an environment variable
	pgURL, err := pq.ParseURL("postgres://rexawhzp:Falpwbfj0BlAO-4Mund6QGSCHGV4Jsb9@salt.db.elephantsql.com:5432/rexawhzp")
	if err != nil {
		log.Fatal(err)
	}

	db, err := jet.Open("postgres", pgURL)
	if err != nil {
		log.Fatal(err)
	}

	return &accountService{db: *db}, nil
}

func (a *accountService) RegisterUser(user User) error {
	if exists, _ := a.UserExists(user.ID); exists {
		return nil
	}

	acct, err := a.GetAccount(user.AccountID)
	if err != nil {
		return err
	}

	count, err := a.GetUserCount(user.AccountID)
	if err != nil {
		return err
	}

	if int64(count)+1 >= acct.Limit {
		return errors.New("user quota reached")
	}

	err = a.db.Query(`INSERT INTO users (id, timestamp, accountid) VALUES ($1, $2, $3)`, user.ID, user.Timestamp, user.AccountID).Run()
	if err != nil {
		return err
	}

	return nil
}

func (a *accountService) UpgradePlan(limit int64, accountID string) error {
	err := a.db.Query(`UPDATE account SET "limit" = 1000 WHERE id = $1`, accountID).Run()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (a *accountService) GetUserCount(accountID string) (int, error) {
	users, err := a.GetUsers(accountID)
	if err != nil {
		return 0, err
	}

	return len(users), nil
}

func (a *accountService) GetUsers(accountID string) ([]User, error) {
	var users []*struct {
		Id        string
		Timestamp string
		Accountid string
	}

	err := a.db.Query(`SELECT * FROM users WHERE accountid = $1`, accountID).Rows(&users)
	if err != nil {
		log.Fatal(err)
	}

	outUsers := []User{}
	for _, user := range users {
		currUser := User{
			ID:        user.Id,
			Timestamp: user.Timestamp,
			AccountID: user.Accountid,
		}

		outUsers = append(outUsers, currUser)
	}
	return outUsers, nil
}

// GetAccount fetches an account from the db based on an id
func (a accountService) GetAccount(id string) (Account, error) {

	var account []*struct {
		Id    string
		Limit int64
	}
	err := a.db.Query(`SELECT * FROM account WHERE id = $1`, id).Rows(&account)
	if err != nil {
		log.Fatal(err)
	}

	if len(account) == 0 {
		return Account{}, errors.New("account not found")
	}

	return Account{
		ID:    account[0].Id,
		Limit: account[0].Limit,
	}, nil
}

func (a accountService) UserExists(id string) (bool, error) {
	_, err := a.GetUser(id)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (a accountService) GetUser(id string) (User, error) {
	var user []*struct {
		Id        string
		AccountId string
		Timestamp string
	}

	err := a.db.Query(`SELECT * FROM "users" WHERE id = $1`, id).Rows(&user)
	if err != nil {
		log.Fatal(err)
	}

	if len(user) == 0 {
		return User{}, errors.New("user not found")
	}

	return User{
		ID:        user[0].Id,
		AccountID: user[0].AccountId,
		Timestamp: user[0].Timestamp,
	}, nil
}
