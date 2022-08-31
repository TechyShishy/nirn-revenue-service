package account

import (
	"fmt"
)

type AccountRegistry struct {
	accounts map[uint]*Account
}

func New() *AccountRegistry {
	return &AccountRegistry{accounts: make(map[uint]*Account)}
}

func (a *AccountRegistry) Add(id uint, name string) (*Account, error) {
	if a.accounts[id] != nil {
		if name != "" && a.accounts[id].Name != name {
			if a.accounts[id].Name != "" {
				return nil, fmt.Errorf("cowardly refusing to update conflicting account ids")
			}
			a.accounts[id].Name = name
		}
		return a.accounts[id], nil
	}
	account := &Account{Id: id, Name: name}
	a.accounts[id] = account
	return account, nil
}

type Account struct {
	Id   uint
	Name string
}
