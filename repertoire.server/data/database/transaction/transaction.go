package transaction

import (
	"gorm.io/gorm"
	"repertoire/server/data/database"
)

type functionType func(factory RepositoryFactory) error

type Manager interface {
	Execute(fn functionType) error
}

type transactionManager struct {
	client database.Client
}

func NewTransactionManager(client database.Client) Manager {
	return &transactionManager{client}
}

func (h *transactionManager) Execute(fn functionType) error {
	return h.client.Transaction(func(tx *gorm.DB) error {
		txFactory := repositoryFactory{client: database.Client{DB: tx}}
		return fn(txFactory)
	})
}
