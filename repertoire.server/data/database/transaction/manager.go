package transaction

import (
	"gorm.io/gorm"
	"repertoire/server/data/database"
)

type FunctionWithFactories func(factory RepositoryFactory) error

type Manager interface {
	Execute(fn FunctionWithFactories) error
}

type manager struct {
	client database.Client
}

func NewManager(client database.Client) Manager {
	return &manager{client}
}

func (m *manager) Execute(fn FunctionWithFactories) error {
	return m.client.Transaction(func(tx *gorm.DB) error {
		txFactory := repositoryFactory{client: database.Client{DB: tx}}
		return fn(txFactory)
	})
}
