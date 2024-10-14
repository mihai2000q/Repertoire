package service

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"repertoire/models"
	"repertoire/utils"
)

type JwtServiceMock struct {
	mock.Mock
}

func (m *JwtServiceMock) Authorize(tokenString string) *utils.ErrorCode {
	args := m.Called(tokenString)

	var errCode *utils.ErrorCode
	if a := args.Get(0); a != nil {
		errCode = a.(*utils.ErrorCode)
	}

	return errCode
}

func (m *JwtServiceMock) CreateToken(user models.User) (string, *utils.ErrorCode) {
	args := m.Called(user)

	var errCode *utils.ErrorCode
	if a := args.Get(1); a != nil {
		errCode = a.(*utils.ErrorCode)
	}

	return args.String(0), errCode
}

func (m *JwtServiceMock) Validate(tokenString string) (uuid.UUID, *utils.ErrorCode) {
	args := m.Called(tokenString)

	var errCode *utils.ErrorCode
	if a := args.Get(1); a != nil {
		errCode = a.(*utils.ErrorCode)
	}

	return args.Get(0).(uuid.UUID), errCode
}

func (m *JwtServiceMock) GetUserIdFromJwt(token string) (uuid.UUID, *utils.ErrorCode) {
	args := m.Called(token)

	var errCode *utils.ErrorCode
	if a := args.Get(1); a != nil {
		errCode = a.(*utils.ErrorCode)
	}

	return args.Get(0).(uuid.UUID), errCode
}
