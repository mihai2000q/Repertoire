package service

import (
	"repertoire/server/internal/httperror"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type JwtServiceMock struct {
	mock.Mock
}

func (m *JwtServiceMock) Authorize(tokenString string) *httperror.ErrorCode {
	args := m.Called(tokenString)

	var errCode *httperror.ErrorCode
	if a := args.Get(0); a != nil {
		errCode = a.(*httperror.ErrorCode)
	}

	return errCode
}

func (m *JwtServiceMock) GetUserIdFromJwt(token string) (uuid.UUID, *httperror.ErrorCode) {
	args := m.Called(token)

	var errCode *httperror.ErrorCode
	if a := args.Get(1); a != nil {
		errCode = a.(*httperror.ErrorCode)
	}

	return args.Get(0).(uuid.UUID), errCode
}
