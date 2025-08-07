package user

import (
	"errors"
	"net/http"
	"repertoire/server/domain/usecase/user"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetUser_WhenGetUserFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	userRepository := new(repository.UserRepositoryMock)
	_uut := user.NewGetUser(userRepository)

	id := uuid.New()

	internalError := errors.New("internal error")
	userRepository.On("Get", new(model.User), id).Return(internalError).Once()

	// when
	resultUser, errCode := _uut.Handle(id)

	// then
	assert.Empty(t, resultUser)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	userRepository.AssertExpectations(t)
}

func TestGetUser_WhenUserIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	userRepository := new(repository.UserRepositoryMock)
	_uut := user.NewGetUser(userRepository)

	id := uuid.New()

	userRepository.On("Get", new(model.User), id).Return(nil).Once()

	// when
	resultUser, errCode := _uut.Handle(id)

	// then
	assert.Empty(t, resultUser)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "user not found", errCode.Error.Error())

	userRepository.AssertExpectations(t)
}

func TestGetUser_WhenSuccessful_ShouldReturnUser(t *testing.T) {
	// given
	userRepository := new(repository.UserRepositoryMock)
	_uut := user.NewGetUser(userRepository)

	id := uuid.New()

	expectedUser := &model.User{
		ID:    id,
		Name:  "Samuel",
		Email: "samuel.samuel@gmail.com",
	}

	userRepository.On("Get", new(model.User), id).Return(nil, expectedUser).Once()

	// when
	resultUser, errCode := _uut.Handle(id)

	// then
	assert.NotEmpty(t, resultUser)
	assert.Equal(t, expectedUser, &resultUser)
	assert.Nil(t, errCode)

	userRepository.AssertExpectations(t)
}
