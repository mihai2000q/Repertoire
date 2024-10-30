package repository

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"repertoire/server/model"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (u *UserRepositoryMock) Get(user *model.User, id uuid.UUID) error {
	args := u.Called(user, id)

	if len(args) > 1 {
		*user = *args.Get(1).(*model.User)
	}

	return args.Error(0)
}

func (u *UserRepositoryMock) GetByEmail(user *model.User, email string) error {
	args := u.Called(user, email)

	if len(args) > 1 {
		*user = *args.Get(1).(*model.User)
	}

	return args.Error(0)
}

func (u *UserRepositoryMock) Create(user *model.User) error {
	args := u.Called(user)
	return args.Error(0)
}
