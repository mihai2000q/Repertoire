package service

import "github.com/stretchr/testify/mock"

type SearchTaskTrackerServiceMock struct {
	mock.Mock
}

func (a *SearchTaskTrackerServiceMock) Track(taskID string, userID string) {
	a.Called(taskID, userID)
}

func (a *SearchTaskTrackerServiceMock) GetUserID(taskID string) (string, bool) {
	args := a.Called(taskID)
	return args.String(0), args.Bool(1)
}
