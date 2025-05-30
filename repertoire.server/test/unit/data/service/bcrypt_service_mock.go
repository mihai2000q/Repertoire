package service

import (
	"github.com/stretchr/testify/mock"
)

type BCryptServiceMock struct {
	mock.Mock
}

func (b *BCryptServiceMock) Hash(str string) (string, error) {
	args := b.Called(str)
	return args.String(0), args.Error(1)
}
