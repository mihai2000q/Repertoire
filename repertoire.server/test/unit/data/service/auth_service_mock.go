package service

import (
	"repertoire/server/internal/wrapper"

	"github.com/stretchr/testify/mock"
)

type AuthServiceMock struct {
	mock.Mock
}

func (a *AuthServiceMock) SignIn(email string, password string) (string, *wrapper.ErrorCode) {
	args := a.Called(email, password)

	var errCode *wrapper.ErrorCode
	if e := args.Get(1); e != nil {
		errCode = e.(*wrapper.ErrorCode)
	}

	return args.String(0), errCode
}
