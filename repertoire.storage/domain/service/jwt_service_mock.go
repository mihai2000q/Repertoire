package service

import "github.com/stretchr/testify/mock"

type JwtServiceMock struct {
	mock.Mock
}

func (j *JwtServiceMock) Authorize(authToken string) error {
	args := j.Called(authToken)
	return args.Error(0)
}

func (j *JwtServiceMock) CreateToken() (string, error) {
	args := j.Called()
	return args.String(0), args.Error(1)
}
