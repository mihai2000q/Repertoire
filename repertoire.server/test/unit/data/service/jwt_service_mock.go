package service

import (
	"repertoire/server/internal/wrapper"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type JwtServiceMock struct {
	mock.Mock
}

func (m *JwtServiceMock) Authorize(tokenString string) *wrapper.ErrorCode {
	args := m.Called(tokenString)

	var errCode *wrapper.ErrorCode
	if a := args.Get(0); a != nil {
		errCode = a.(*wrapper.ErrorCode)
	}

	return errCode
}

func (m *JwtServiceMock) GetUserIdFromJwt(token string) (uuid.UUID, *wrapper.ErrorCode) {
	args := m.Called(token)

	var errCode *wrapper.ErrorCode
	if a := args.Get(1); a != nil {
		errCode = a.(*wrapper.ErrorCode)
	}

	return args.Get(0).(uuid.UUID), errCode
}
