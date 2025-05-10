package transaction

import (
	"gorm.io/gorm"
	"repertoire/server/data/database"
)

type functionType func(factory RepositoryFactory) error

type Handler interface {
	Execute(fn functionType) error
}

type gormTransactionHandler struct {
	client database.Client
}

func NewTransactionHandler(client database.Client) Handler {
	return &gormTransactionHandler{client}
}

func (h *gormTransactionHandler) Execute(fn functionType) error {
	return h.client.Transaction(func(tx *gorm.DB) error {
		txFactory := repositoryFactory{client: database.Client{DB: tx}}
		return fn(txFactory)
	})
}
