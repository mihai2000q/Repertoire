package service

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"repertoire/data/repository"
	"repertoire/models"
	"testing"
)

// Get
func TestUserService_Get_WhenUserRepositoryReturnsError_ShouldReturnInternalServerError(t *testing.T) {
	// given
	userRepository := new(repository.UserRepositoryMock)
	_uut := &userService{
		repository: userRepository,
	}
	id := uuid.New()

	internalError := errors.New("internal error")
	userRepository.On("Get", &models.User{}, id).Return(internalError).Once()

	// when
	user, errCode := _uut.Get(id)

	// then
	assert.Empty(t, user)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	userRepository.AssertExpectations(t)
}

func TestUserService_Get_WhenUserIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	userRepository := new(repository.UserRepositoryMock)
	_uut := &userService{
		repository: userRepository,
	}
	id := uuid.New()

	userRepository.On("Get", &models.User{}, id).Return(nil).Once()

	// when
	user, errCode := _uut.Get(id)

	// then
	assert.Empty(t, user)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "user not found", errCode.Error.Error())

	userRepository.AssertExpectations(t)
}

func TestUserService_Get_WhenSuccessful_ShouldReturnUser(t *testing.T) {
	// given
	userRepository := new(repository.UserRepositoryMock)
	_uut := &userService{
		repository: userRepository,
	}
	id := uuid.New()

	expectedUser := &models.User{
		ID:    id,
		Name:  "Samuel",
		Email: "samuel.samuel@gmail.com",
	}

	userRepository.On("Get", &models.User{}, id).Return(nil, expectedUser).Once()

	// when
	user, errCode := _uut.Get(id)

	// then
	assert.NotEmpty(t, user)
	assert.Equal(t, expectedUser, &user)
	assert.Nil(t, errCode)

	userRepository.AssertExpectations(t)
}
