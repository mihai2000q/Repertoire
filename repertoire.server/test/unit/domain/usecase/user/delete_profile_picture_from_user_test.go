package user

import (
	"errors"
	"net/http"
	user2 "repertoire/server/domain/usecase/user"
	"repertoire/server/internal"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	service2 "repertoire/server/test/unit/data/service"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteProfilePictureFromUser_WhenGetUserIdFromJwtFails_ShouldReturnTheError(t *testing.T) {
	// given
	jwtService := new(service2.JwtServiceMock)
	_uut := user2.DeleteProfilePictureFromUser{jwtService: jwtService}

	token := "this is a token"

	// given - mocking
	err := wrapper.ForbiddenError(errors.New("forbidden error"))
	jwtService.On("GetUserIdFromJwt", token).Return(uuid.Nil, err).Once()

	// when
	errCode := _uut.Handle(token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, err, errCode)

	jwtService.AssertExpectations(t)
}

func TestDeleteProfilePictureFromUser_WhenGetUserFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service2.JwtServiceMock)
	userRepository := new(repository.UserRepositoryMock)
	_uut := user2.DeleteProfilePictureFromUser{
		jwtService: jwtService,
		repository: userRepository,
	}

	token := "this is a token"

	// given - mocking
	id := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(id, nil).Once()

	internalError := errors.New("internal error")
	userRepository.On("Get", new(model.User), id).Return(internalError).Once()

	// when
	errCode := _uut.Handle(token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	jwtService.AssertExpectations(t)
	userRepository.AssertExpectations(t)
}

func TestDeleteProfilePictureFromUser_WhenGetUserFails_ShouldReturnNotFoundError(t *testing.T) {
	// given
	jwtService := new(service2.JwtServiceMock)
	userRepository := new(repository.UserRepositoryMock)
	_uut := user2.DeleteProfilePictureFromUser{
		jwtService: jwtService,
		repository: userRepository,
	}

	token := "this is a token"

	// given - mocking
	id := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(id, nil).Once()

	userRepository.On("Get", new(model.User), id).Return(nil).Once()

	// when
	errCode := _uut.Handle(token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "user not found", errCode.Error.Error())

	jwtService.AssertExpectations(t)
	userRepository.AssertExpectations(t)
}

func TestDeleteProfilePictureFromUser_WhenDeleteProfilePictureFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service2.JwtServiceMock)
	userRepository := new(repository.UserRepositoryMock)
	storageService := new(service2.StorageServiceMock)
	_uut := user2.DeleteProfilePictureFromUser{
		jwtService:     jwtService,
		repository:     userRepository,
		storageService: storageService,
	}

	token := "this is a token"

	// given - mocking
	id := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(id, nil).Once()

	user := &model.User{ID: id, ProfilePictureURL: &[]internal.FilePath{"This is some url"}[0]}
	userRepository.On("Get", new(model.User), id).Return(nil, user).Once()

	internalError := errors.New("internal error")
	storageService.On("Delete", string(*user.ProfilePictureURL)).Return(internalError).Once()

	// when
	errCode := _uut.Handle(token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	jwtService.AssertExpectations(t)
	userRepository.AssertExpectations(t)
	storageService.AssertExpectations(t)
}

func TestDeleteProfilePictureFromUser_WhenUpdateUserFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service2.JwtServiceMock)
	userRepository := new(repository.UserRepositoryMock)
	storageService := new(service2.StorageServiceMock)
	_uut := user2.DeleteProfilePictureFromUser{
		jwtService:     jwtService,
		repository:     userRepository,
		storageService: storageService,
	}

	token := "this is a token"

	// given - mocking
	id := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(id, nil).Once()

	user := &model.User{ID: id, ProfilePictureURL: &[]internal.FilePath{"This is some url"}[0]}
	userRepository.On("Get", new(model.User), id).Return(nil, user).Once()

	storageService.On("Delete", string(*user.ProfilePictureURL)).Return(nil).Once()

	internalError := errors.New("internal error")
	userRepository.On("Update", mock.IsType(user)).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	jwtService.AssertExpectations(t)
	userRepository.AssertExpectations(t)
	storageService.AssertExpectations(t)
}

func TestDeleteProfilePictureFromUser_WhenIsValid_ShouldNotReturnAnyError(t *testing.T) {
	// given
	jwtService := new(service2.JwtServiceMock)
	userRepository := new(repository.UserRepositoryMock)
	storageService := new(service2.StorageServiceMock)
	_uut := user2.DeleteProfilePictureFromUser{
		jwtService:     jwtService,
		repository:     userRepository,
		storageService: storageService,
	}

	token := "this is a token"

	// given - mocking
	id := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(id, nil).Once()

	user := &model.User{ID: id, ProfilePictureURL: &[]internal.FilePath{"This is some url"}[0]}
	userRepository.On("Get", new(model.User), id).Return(nil, user).Once()

	storageService.On("Delete", string(*user.ProfilePictureURL)).Return(nil).Once()

	userRepository.On("Update", mock.IsType(user)).
		Run(func(args mock.Arguments) {
			newUser := args.Get(0).(*model.User)
			assert.Nil(t, newUser.ProfilePictureURL)
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(token)

	// then
	assert.Nil(t, errCode)

	jwtService.AssertExpectations(t)
	userRepository.AssertExpectations(t)
	storageService.AssertExpectations(t)
}
