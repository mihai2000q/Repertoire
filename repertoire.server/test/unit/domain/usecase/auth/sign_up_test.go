package auth

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/user"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthService_SignUp_WhenUserRepositoryReturnsError_ShouldReturnInternalServerError(t *testing.T) {
	// given
	userRepository := new(repository.UserRepositoryMock)
	_uut := user.NewSignUp(nil, nil, userRepository)

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
	_uut := user.NewSignUp(nil, nil, userRepository)

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
	_uut := user.NewSignUp(nil, bCryptService, userRepository)

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
	_uut := user.NewSignUp(nil, bCryptService, userRepository)

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

	internalError := errors.New("internal error")
	userRepository.On("Create", mock.IsType(&model.User{})).
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

func TestAuthService_SignUp_WhenSignInFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	authService := new(service.AuthServiceMock)
	bCryptService := new(service.BCryptServiceMock)
	userRepository := new(repository.UserRepositoryMock)
	_uut := user.NewSignUp(authService, bCryptService, userRepository)

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

	userRepository.On("Create", mock.IsType(&model.User{})).
		Return(nil).
		Once()

	internalError := wrapper.InternalServerError(errors.New("internal error"))
	authService.On("SignIn", strings.ToLower(request.Email), request.Password).
		Return("", internalError).
		Once()

	// when
	token, errCode := _uut.Handle(request)

	// then
	assert.Empty(t, token)
	assert.NotNil(t, errCode)
	assert.Equal(t, internalError, errCode)

	authService.AssertExpectations(t)
	userRepository.AssertExpectations(t)
	bCryptService.AssertExpectations(t)
}

func TestAuthService_SignUp_WhenSuccessful_ShouldReturnNewToken(t *testing.T) {
	// given
	authService := new(service.AuthServiceMock)
	bCryptService := new(service.BCryptServiceMock)
	userRepository := new(repository.UserRepositoryMock)
	_uut := user.NewSignUp(authService, bCryptService, userRepository)

	request := requests.SignUpRequest{
		Name:     "Samuel",
		Email:    "Samuel@yahoo.com",
		Password: "Password123",
	}

	// given - mocking
	userRepository.On("GetByEmail", new(model.User), strings.ToLower(request.Email)).
		Return(nil).
		Once()

	hashedPassword := "hashed password"
	bCryptService.On("Hash", request.Password).Return(hashedPassword, nil).Once()

	var expectedUser *model.User
	userRepository.On("Create", mock.IsType(expectedUser)).
		Run(func(args mock.Arguments) {
			expectedUser = args.Get(0).(*model.User)
			assertCreatedUser(t, *expectedUser, request, hashedPassword)
		}).
		Return(nil).
		Once()

	expectedToken := "This is the generated token"
	authService.On("SignIn", strings.ToLower(request.Email), request.Password).
		Return(expectedToken, nil).
		Once()

	// when
	token, errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)
	assert.Equal(t, expectedToken, token)

	authService.AssertExpectations(t)
	userRepository.AssertExpectations(t)
	bCryptService.AssertExpectations(t)
}

func assertCreatedUser(
	t *testing.T,
	user model.User,
	request requests.SignUpRequest,
	hashedPassword string,
) {
	assert.NotEmpty(t, user.ID)
	assert.Equal(t, request.Name, user.Name)
	assert.Equal(t, strings.ToLower(request.Email), user.Email)
	assert.Equal(t, hashedPassword, user.Password)
	assert.Nil(t, user.ProfilePictureURL)

	for i, guitarTuning := range user.GuitarTunings {
		assert.NotEmpty(t, guitarTuning.ID)
		assert.Equal(t, user.ID, guitarTuning.UserID)
		assert.Equal(t, model.DefaultGuitarTunings[i], guitarTuning.Name)
		assert.Equal(t, uint(i), guitarTuning.Order)
	}

	for i, songSectionType := range user.SongSectionTypes {
		assert.NotEmpty(t, songSectionType.ID)
		assert.Equal(t, user.ID, songSectionType.UserID)
		assert.Equal(t, model.DefaultSongSectionTypes[i], songSectionType.Name)
		assert.Equal(t, uint(i), songSectionType.Order)
	}

	for i, bandMemberRole := range user.BandMemberRoles {
		assert.NotEmpty(t, bandMemberRole.ID)
		assert.Equal(t, user.ID, bandMemberRole.UserID)
		assert.Equal(t, model.DefaultBandMemberRoles[i], bandMemberRole.Name)
		assert.Equal(t, uint(i), bandMemberRole.Order)
	}

	for i, instrument := range user.Instruments {
		assert.NotEmpty(t, instrument.ID)
		assert.Equal(t, user.ID, instrument.UserID)
		assert.Equal(t, model.DefaultInstruments[i], instrument.Name)
		assert.Equal(t, uint(i), instrument.Order)
	}
}
