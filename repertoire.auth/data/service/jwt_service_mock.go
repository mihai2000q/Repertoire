package service

import (
	"repertoire/auth/internal/wrapper"
	"repertoire/auth/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type JwtServiceMock struct {
	mock.Mock
}

func (j *JwtServiceMock) Authorize(authToken string) *wrapper.ErrorCode {
	args := j.Called(authToken)

	var errCode *wrapper.ErrorCode
	if a := args.Get(0); a != nil {
		errCode = a.(*wrapper.ErrorCode)
	}

	return errCode
}

func (j *JwtServiceMock) CreateToken(user model.User) (string, *wrapper.ErrorCode) {
	args := j.Called(user)

	var errCode *wrapper.ErrorCode
	if a := args.Get(1); a != nil {
		errCode = a.(*wrapper.ErrorCode)
	}

	return args.String(0), errCode
}

func (j *JwtServiceMock) CreateCentrifugoToken(userID uuid.UUID) (string, string, *wrapper.ErrorCode) {
	args := j.Called(userID)

	var errCode *wrapper.ErrorCode
	if a := args.Get(2); a != nil {
		errCode = a.(*wrapper.ErrorCode)
	}

	return args.String(0), args.String(1), errCode
}

func (j *JwtServiceMock) CreateStorageToken(userID uuid.UUID) (string, string, *wrapper.ErrorCode) {
	args := j.Called(userID)

	var errCode *wrapper.ErrorCode
	if a := args.Get(2); a != nil {
		errCode = a.(*wrapper.ErrorCode)
	}

	return args.String(0), args.String(1), errCode
}

func (j *JwtServiceMock) Validate(tokenString string) (uuid.UUID, *wrapper.ErrorCode) {
	args := j.Called(tokenString)

	var errCode *wrapper.ErrorCode
	if a := args.Get(1); a != nil {
		errCode = a.(*wrapper.ErrorCode)
	}

	return args.Get(0).(uuid.UUID), errCode
}

func (j *JwtServiceMock) ValidateCredentials(clientCredentials model.ClientCredentials) *wrapper.ErrorCode {
	args := j.Called(clientCredentials)

	var errCode *wrapper.ErrorCode
	if a := args.Get(0); a != nil {
		errCode = a.(*wrapper.ErrorCode)
	}

	return errCode
}

func (j *JwtServiceMock) GetUserIDFromJwt(token string) (uuid.UUID, *wrapper.ErrorCode) {
	args := j.Called(token)

	var errCode *wrapper.ErrorCode
	if a := args.Get(1); a != nil {
		errCode = a.(*wrapper.ErrorCode)
	}

	return args.Get(0).(uuid.UUID), errCode
}
