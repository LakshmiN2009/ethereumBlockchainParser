package driver

import (
	"errors"
	"ethereumBlockchainParser/client"
)

// Transaction ... inbound and outbound transactions for an address
type Transaction struct {
	InBound  []client.EthTransaction `json:"inbound"`
	OutBound []client.EthTransaction `json:"outbound"`
}

// InMemStore ... in memory store for subscribers and transactions
type InMemStore struct {
	Subscribers map[string]bool
	Transcation map[string]Transaction
}

// NewInMemStore ... create a new instance of InMemStore
func NewInMemStore() DB {
	return &InMemStore{
		Subscribers: make(map[string]bool),
		Transcation: make(map[string]Transaction),
	}
}

// Insert ... insert new transaction for a address
func (i *InMemStore) Insert(key string, value Transaction) error {
	i.Transcation[key] = value
	return nil
}

// Get ... get the transaction for the address
func (i *InMemStore) Get(key string) (Transaction, error) {
	if txn, ok := i.Transcation[key]; ok {
		return txn, nil
	}
	return Transaction{}, errors.New("transaction not found")
}
