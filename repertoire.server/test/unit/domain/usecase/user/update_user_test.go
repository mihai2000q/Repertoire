package user

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/user"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateUser_WhenGetUserIdFromJwtFails_ShouldReturnTheError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	_uut := user.NewUpdateUser(nil, jwtService)

	request := requests.UpdateUserRequest{
		Name: "New User Name",
	}

	token := "This is a token"

	// given - mocking
	forbiddenError := wrapper.ForbiddenError(errors.New("forbidden error"))
	jwtService.On("GetUserIdFromJwt", token).Return(uuid.Nil, forbiddenError).Once()

	// when
	errCode := _uut.Handle(request, token)

	// then
	assert.NotNil(t, errCode)
	assert.NotNil(t, forbiddenError, errCode)

	jwtService.AssertExpectations(t)
}

func TestUpdateUser_WhenGetUserFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userRepository := new(repository.UserRepositoryMock)
	_uut := user.NewUpdateUser(userRepository, jwtService)

	request := requests.UpdateUserRequest{
		Name: "New User Name",
	}

	token := "This is a token"

	// given - mocking
	id := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(id, nil).Once()

	internalError := errors.New("internal error")
	userRepository.On("Get", new(model.User), id).Return(internalError).Once()

	// when
	errCode := _uut.Handle(request, token)

	// then
	assert.NotNil(t, errCode)
	assert.NotNil(t, http.StatusInternalServerError, errCode.Code)
	assert.NotNil(t, internalError, errCode.Error)

	jwtService.AssertExpectations(t)
	userRepository.AssertExpectations(t)
}

func TestUpdateUser_WhenUserIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userRepository := new(repository.UserRepositoryMock)
	_uut := user.NewUpdateUser(userRepository, jwtService)

	request := requests.UpdateUserRequest{
		Name: "New User Name",
	}

	token := "This is a token"

	// given - mocking
	id := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(id, nil).Once()

	userRepository.On("Get", new(model.User), id).Return(nil).Once()

	// when
	errCode := _uut.Handle(request, token)

	// then
	assert.NotNil(t, errCode)
	assert.NotNil(t, http.StatusNotFound, errCode.Code)
	assert.NotNil(t, "user not found", errCode.Error.Error())

	jwtService.AssertExpectations(t)
	userRepository.AssertExpectations(t)
}

func TestUpdateUser_WhenUpdateUserFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userRepository := new(repository.UserRepositoryMock)
	_uut := user.NewUpdateUser(userRepository, jwtService)

	request := requests.UpdateUserRequest{
		Name: "New User Name",
	}

	token := "This is a token"

	// given - mocking
	id := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(id, nil).Once()

	mockUser := &model.User{ID: id}
	userRepository.On("Get", new(model.User), id).Return(nil, mockUser).Once()

	internalError := errors.New("internal error")
	userRepository.On("Update", mock.IsType(mockUser)).Return(internalError).Once()

	// when
	errCode := _uut.Handle(request, token)

	// then
	assert.NotNil(t, errCode)
	assert.NotNil(t, http.StatusInternalServerError, errCode.Code)
	assert.NotNil(t, internalError, errCode.Error)

	jwtService.AssertExpectations(t)
	userRepository.AssertExpectations(t)
}

func TestUpdateUser_WhenIsValid_ShouldNotReturnAnyError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userRepository := new(repository.UserRepositoryMock)
	_uut := user.NewUpdateUser(userRepository, jwtService)

	request := requests.UpdateUserRequest{
		Name: "New User Name",
	}

	token := "This is a token"

	// given - mocking
	id := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(id, nil).Once()

	mockUser := &model.User{
		ID:   id,
		Name: "Old name",
	}
	userRepository.On("Get", new(model.User), id).Return(nil, mockUser).Once()
	userRepository.On("Update", mock.IsType(mockUser)).
		Run(func(args mock.Arguments) {
			newUser := args.Get(0).(*model.User)
			assertUpdatedUser(t, *newUser, request)
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request, token)

	// then
	assert.Nil(t, errCode)

	jwtService.AssertExpectations(t)
	userRepository.AssertExpectations(t)
}

func assertUpdatedUser(t *testing.T, user model.User, request requests.UpdateUserRequest) {
	assert.Equal(t, request.Name, user.Name)
}
