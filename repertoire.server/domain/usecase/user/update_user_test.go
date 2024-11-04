package user

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/model"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateUser_WhenGetUserFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	userRepository := new(repository.UserRepositoryMock)
	_uut := UpdateUser{
		repository: userRepository,
	}

	request := requests.UpdateUserRequest{
		ID:   uuid.New(),
		Name: "New User Name",
	}

	// given - mocking
	internalError := errors.New("internal error")
	userRepository.On("Get", new(model.User), request.ID).Return(internalError).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.NotNil(t, http.StatusInternalServerError, errCode.Code)
	assert.NotNil(t, internalError, errCode.Error)

	userRepository.AssertExpectations(t)
}

func TestUpdateUser_WhenUserIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	userRepository := new(repository.UserRepositoryMock)
	_uut := UpdateUser{
		repository: userRepository,
	}

	request := requests.UpdateUserRequest{
		ID:   uuid.New(),
		Name: "New User Name",
	}

	// given - mocking
	userRepository.On("Get", new(model.User), request.ID).Return(nil).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.NotNil(t, http.StatusNotFound, errCode.Code)
	assert.NotNil(t, "user not found", errCode.Error.Error())

	userRepository.AssertExpectations(t)
}

func TestUpdateUser_WhenUpdateUserFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	userRepository := new(repository.UserRepositoryMock)
	_uut := UpdateUser{
		repository: userRepository,
	}

	request := requests.UpdateUserRequest{
		ID:   uuid.New(),
		Name: "New User Name",
	}

	// given - mocking
	user := &model.User{ID: request.ID}
	userRepository.On("Get", new(model.User), request.ID).Return(nil, user).Once()

	internalError := errors.New("internal error")
	userRepository.On("Update", mock.IsType(user)).Return(internalError).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.NotNil(t, http.StatusInternalServerError, errCode.Code)
	assert.NotNil(t, internalError, errCode.Error)

	userRepository.AssertExpectations(t)
}

func TestUpdateUser_WhenIsValid_ShouldNotReturnAnyError(t *testing.T) {
	// given
	userRepository := new(repository.UserRepositoryMock)
	_uut := UpdateUser{
		repository: userRepository,
	}

	request := requests.UpdateUserRequest{
		ID:   uuid.New(),
		Name: "New User Name",
	}

	// given - mocking
	user := &model.User{
		ID:   request.ID,
		Name: "Old name",
	}
	userRepository.On("Get", new(model.User), request.ID).Return(nil, user).Once()
	userRepository.On("Update", mock.IsType(user)).
		Run(func(args mock.Arguments) {
			newUser := args.Get(0).(*model.User)
			assertUpdatedUser(t, *newUser, request)
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)

	userRepository.AssertExpectations(t)
}

func assertUpdatedUser(t *testing.T, user model.User, request requests.UpdateUserRequest) {
	assert.Equal(t, request.Name, user.Name)
}
