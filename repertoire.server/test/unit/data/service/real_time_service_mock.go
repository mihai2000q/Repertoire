package service

import "github.com/stretchr/testify/mock"

type RealTimeServiceMock struct {
	mock.Mock
}

func (r *RealTimeServiceMock) Publish(channel string, userID string, payload any) error {
	args := r.Called(channel, userID, payload)
	return args.Error(0)
}

func (r *RealTimeServiceMock) CreateToken(userID string) string {
	args := r.Called(userID)
	return args.String(0)
}
