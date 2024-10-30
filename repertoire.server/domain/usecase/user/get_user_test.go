package user

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"repertoire/server/data/repository"
	"repertoire/server/model"
	"testing"
)

func TestGetUser_WhenGetUserFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	userRepository := new(repository.UserRepositoryMock)
	_uut := &GetUser{
		repository: userRepository,
	}
	id := uuid.New()

	internalError := errors.New("internal error")
	userRepository.On("Get", new(model.User), id).Return(internalError).Once()

	// when
	user, errCode := _uut.Handle(id)

	// then
	assert.Empty(t, user)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	userRepository.AssertExpectations(t)
}

func TestGetUser_WhenUserIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	userRepository := new(repository.UserRepositoryMock)
	_uut := &GetUser{
		repository: userRepository,
	}
	id := uuid.New()

	userRepository.On("Get", new(model.User), id).Return(nil).Once()

	// when
	user, errCode := _uut.Handle(id)

	// then
	assert.Empty(t, user)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "user not found", errCode.Error.Error())

	userRepository.AssertExpectations(t)
}

func TestGetUser_WhenSuccessful_ShouldReturnUser(t *testing.T) {
	// given
	userRepository := new(repository.UserRepositoryMock)
	_uut := &GetUser{
		repository: userRepository,
	}
	id := uuid.New()

	expectedUser := &model.User{
		ID:    id,
		Name:  "Samuel",
		Email: "samuel.samuel@gmail.com",
	}

	userRepository.On("Get", new(model.User), id).Return(nil, expectedUser).Once()

	// when
	user, errCode := _uut.Handle(id)

	// then
	assert.NotEmpty(t, user)
	assert.Equal(t, expectedUser, &user)
	assert.Nil(t, errCode)

	userRepository.AssertExpectations(t)
}
