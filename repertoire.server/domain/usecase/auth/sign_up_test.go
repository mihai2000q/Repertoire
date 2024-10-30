package auth

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"strings"
	"testing"
)

func TestAuthService_SignUp_WhenUserRepositoryReturnsError_ShouldReturnInternalServerError(t *testing.T) {
	// given
	userRepository := new(repository.UserRepositoryMock)
	_uut := &SignUp{
		userRepository: userRepository,
	}
	request := requests.SignUpRequest{
		Name:     "Samuel",
		Email:    "Samuel@yahoo.com",
		Password: "Password123",
	}

	internalError := errors.New("something went wrong")
	userRepository.On("GetByEmail", new(model.User), strings.ToLower(request.Email)).
		Return(internalError).
		Once()

	// when
	token, errCode := _uut.Handle(request)

	// then
	assert.Empty(t, token)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	userRepository.AssertExpectations(t)
}

func TestAuthService_SignUp_WhenUserIsNotEmpty_ShouldReturnUnauthorizedError(t *testing.T) {
	// given
	userRepository := new(repository.UserRepositoryMock)
	_uut := &SignUp{
		userRepository: userRepository,
	}
	request := requests.SignUpRequest{
		Name:     "Samuel",
		Email:    "Samuel@yahoo.com",
		Password: "Password123",
	}

	userRepository.On("GetByEmail", new(model.User), strings.ToLower(request.Email)).
		Return(nil, &model.User{ID: uuid.New()}).
		Once()

	// when
	token, errCode := _uut.Handle(request)

	// then
	assert.Empty(t, token)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusBadRequest, errCode.Code)
	assert.Equal(t, "user already exists", errCode.Error.Error())

	userRepository.AssertExpectations(t)
}

func TestAuthService_SignUp_WhenHashPasswordFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	bCryptService := new(service.BCryptServiceMock)
	userRepository := new(repository.UserRepositoryMock)
	_uut := &SignUp{
		bCryptService:  bCryptService,
		userRepository: userRepository,
	}
	request := requests.SignUpRequest{
		Name:     "Samuel",
		Email:    "Samuel@yahoo.com",
		Password: "Password123",
	}

	userRepository.On("GetByEmail", new(model.User), strings.ToLower(request.Email)).
		Return(nil).
		Once()

	internalError := errors.New("internal error")
	bCryptService.On("Hash", request.Password).Return("", internalError).Once()

	// when
	token, errCode := _uut.Handle(request)

	// then
	assert.Empty(t, token)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	userRepository.AssertExpectations(t)
	bCryptService.AssertExpectations(t)
}

func TestAuthService_SignUp_WhenCreateUserFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	bCryptService := new(service.BCryptServiceMock)
	userRepository := new(repository.UserRepositoryMock)
	_uut := &SignUp{
		bCryptService:  bCryptService,
		userRepository: userRepository,
	}
	request := requests.SignUpRequest{
		Name:     "Samuel",
		Email:    "Samuel@yahoo.com",
		Password: "Password123",
	}

	userRepository.On("GetByEmail", new(model.User), strings.ToLower(request.Email)).
		Return(nil).
		Once()

	hashedPassword := "hashed password"
	bCryptService.On("Hash", request.Password).Return(hashedPassword, nil).Once()

	var user *model.User
	internalError := errors.New("internal error")
	userRepository.On("Create", mock.IsType(&model.User{})).
		Run(func(args mock.Arguments) {
			user = args.Get(0).(*model.User)
			assert.NotEmpty(t, user.ID)
			assert.Equal(t, user.Name, request.Name)
			assert.Equal(t, user.Email, strings.ToLower(request.Email))
			assert.Equal(t, user.Password, hashedPassword)

			for i, guitarTuning := range user.GuitarTunings {
				assert.NotEmpty(t, guitarTuning.ID)
				assert.Equal(t, user.ID, guitarTuning.UserID)
				assert.Equal(t, model.DefaultGuitarTuning[i], guitarTuning.Name)
				assert.Equal(t, uint(i), guitarTuning.Order)
			}

			for i, songSectionType := range user.SongSectionTypes {
				assert.NotEmpty(t, songSectionType.ID)
				assert.Equal(t, user.ID, songSectionType.UserID)
				assert.Equal(t, model.DefaultSongSectionTypes[i], songSectionType.Name)
				assert.Equal(t, uint(i), songSectionType.Order)
			}
		}).
		Return(internalError).
		Once()

	// when
	token, errCode := _uut.Handle(request)

	// then
	assert.Empty(t, token)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	userRepository.AssertExpectations(t)
	bCryptService.AssertExpectations(t)
}

func TestAuthService_SignUp_WhenCreateTokenFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	bCryptService := new(service.BCryptServiceMock)
	userRepository := new(repository.UserRepositoryMock)
	_uut := &SignUp{
		jwtService:     jwtService,
		bCryptService:  bCryptService,
		userRepository: userRepository,
	}
	request := requests.SignUpRequest{
		Name:     "Samuel",
		Email:    "Samuel@yahoo.com",
		Password: "Password123",
	}

	userRepository.On("GetByEmail", new(model.User), strings.ToLower(request.Email)).
		Return(nil).
		Once()

	hashedPassword := "hashed password"
	bCryptService.On("Hash", request.Password).Return(hashedPassword, nil).Once()

	var user *model.User
	userRepository.On("Create", mock.IsType(&model.User{})).
		Run(func(args mock.Arguments) {
			user = args.Get(0).(*model.User)
			assert.NotEmpty(t, user.ID)
			assert.Equal(t, user.Name, request.Name)
			assert.Equal(t, user.Email, strings.ToLower(request.Email))
			assert.Equal(t, user.Password, hashedPassword)

			for i, guitarTuning := range user.GuitarTunings {
				assert.NotEmpty(t, guitarTuning.ID)
				assert.Equal(t, user.ID, guitarTuning.UserID)
				assert.Equal(t, model.DefaultGuitarTuning[i], guitarTuning.Name)
				assert.Equal(t, uint(i), guitarTuning.Order)
			}

			for i, songSectionType := range user.SongSectionTypes {
				assert.NotEmpty(t, songSectionType.ID)
				assert.Equal(t, user.ID, songSectionType.UserID)
				assert.Equal(t, model.DefaultSongSectionTypes[i], songSectionType.Name)
				assert.Equal(t, uint(i), songSectionType.Order)
			}
		}).
		Return(nil).
		Once()

	internalError := wrapper.InternalServerError(errors.New("internal error"))
	jwtService.On("CreateToken", mock.IsType(model.User{})).
		Run(func(args mock.Arguments) {
			assert.Equal(t, *user, args.Get(0).(model.User))
		}).
		Return("", internalError).
		Once()

	// when
	token, errCode := _uut.Handle(request)

	// then
	assert.Empty(t, token)
	assert.NotNil(t, errCode)
	assert.Equal(t, internalError, errCode)

	userRepository.AssertExpectations(t)
	bCryptService.AssertExpectations(t)
	jwtService.AssertExpectations(t)
}

func TestAuthService_SignUp_WhenSuccessful_ShouldReturnNewToken(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	bCryptService := new(service.BCryptServiceMock)
	userRepository := new(repository.UserRepositoryMock)
	_uut := &SignUp{
		jwtService:     jwtService,
		bCryptService:  bCryptService,
		userRepository: userRepository,
	}
	request := requests.SignUpRequest{
		Name:     "Samuel",
		Email:    "Samuel@yahoo.com",
		Password: "Password123",
	}

	userRepository.On("GetByEmail", new(model.User), strings.ToLower(request.Email)).
		Return(nil).
		Once()

	hashedPassword := "hashed password"
	bCryptService.On("Hash", request.Password).Return(hashedPassword, nil).Once()

	var user *model.User
	userRepository.On("Create", mock.IsType(user)).
		Run(func(args mock.Arguments) {
			user = args.Get(0).(*model.User)
			assert.NotEmpty(t, user.ID)
			assert.Equal(t, user.Name, request.Name)
			assert.Equal(t, user.Email, strings.ToLower(request.Email))
			assert.Equal(t, user.Password, hashedPassword)

			for i, guitarTuning := range user.GuitarTunings {
				assert.NotEmpty(t, guitarTuning.ID)
				assert.Equal(t, user.ID, guitarTuning.UserID)
				assert.Equal(t, model.DefaultGuitarTuning[i], guitarTuning.Name)
				assert.Equal(t, uint(i), guitarTuning.Order)
			}

			for i, songSectionType := range user.SongSectionTypes {
				assert.NotEmpty(t, songSectionType.ID)
				assert.Equal(t, user.ID, songSectionType.UserID)
				assert.Equal(t, model.DefaultSongSectionTypes[i], songSectionType.Name)
				assert.Equal(t, uint(i), songSectionType.Order)
			}
		}).
		Return(nil).
		Once()

	expectedToken := "This is the generated token"
	jwtService.On("CreateToken", mock.IsType(model.User{})).
		Run(func(args mock.Arguments) {
			assert.Equal(t, *user, args.Get(0).(model.User))
		}).
		Return(expectedToken, nil).
		Once()

	// when
	token, errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)
	assert.Equal(t, expectedToken, token)

	userRepository.AssertExpectations(t)
	bCryptService.AssertExpectations(t)
	jwtService.AssertExpectations(t)
}
