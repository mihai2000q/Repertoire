package transaction

import (
	"repertoire/server/data/database/transaction"

	"github.com/stretchr/testify/mock"
)

type ManagerMock struct {
	mock.Mock
}

func (m *ManagerMock) Execute(fn transaction.FunctionWithFactories) error {
	args := m.Called(fn)

	if len(args) > 1 {
		return fn(args.Get(1).(transaction.RepositoryFactory))
	}
	return args.Error(0)
}
