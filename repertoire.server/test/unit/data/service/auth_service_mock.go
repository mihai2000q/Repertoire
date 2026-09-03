package service

import (
	"repertoire/server/internal/httperror"

	"github.com/stretchr/testify/mock"
)

type AuthServiceMock struct {
	mock.Mock
}

func (a *AuthServiceMock) SignIn(email string, password string) (string, *httperror.ErrorCode) {
	args := a.Called(email, password)

	var errCode *httperror.ErrorCode
	if e := args.Get(1); e != nil {
		errCode = e.(*httperror.ErrorCode)
	}

	return args.String(0), errCode
}
