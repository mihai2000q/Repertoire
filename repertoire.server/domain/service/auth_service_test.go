package service

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"repertoire/api/requests"
	"repertoire/data/repository"
	"repertoire/data/service"
	"repertoire/models"
	"repertoire/utils"
	"strings"
	"testing"
)

// Refresh
func TestAuthService_Refresh_WhenJwtServiceReturnsError_ShouldReturnUnauthorizedError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	_uut := &authService{
		jwtService: jwtService,
	}
	request := requests.RefreshRequest{
		Token: "This is a token",
	}

	internalError := utils.InternalServerError(errors.New("internal error"))
	jwtService.On("Validate", request.Token).Return(uuid.Nil, internalError).Once()

	// when
	token, errCode := _uut.Refresh(request)

	// then
	assert.Empty(t, token)
	assert.NotNil(t, errCode)
	assert.Equal(t, internalError, errCode)

	jwtService.AssertExpectations(t)
}

func TestAuthService_Refresh_WhenUserRepositoryReturnsError_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userRepository := new(repository.UserRepositoryMock)
	_uut := &authService{
		jwtService:     jwtService,
		userRepository: userRepository,
	}
	request := requests.RefreshRequest{
		Token: "This is a token",
	}

	userID := uuid.New()
	jwtService.On("Validate", request.Token).Return(userID, nil).Once()

	internalError := errors.New("something went wrong")
	userRepository.On("Get", new(models.User), userID).Return(internalError).Once()

	// when
	token, errCode := _uut.Refresh(request)

	// then
	assert.Empty(t, token)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	jwtService.AssertExpectations(t)
	userRepository.AssertExpectations(t)
}

func TestAuthService_Refresh_WhenUserIsEmpty_ShouldReturnUnauthorizedError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userRepository := new(repository.UserRepositoryMock)
	_uut := &authService{
		jwtService:     jwtService,
		userRepository: userRepository,
	}
	request := requests.RefreshRequest{
		Token: "This is a token",
	}

	userID := uuid.New()
	jwtService.On("Validate", request.Token).Return(userID, nil).Once()
	userRepository.On("Get", new(models.User), userID).Return(nil).Once()

	// when
	token, errCode := _uut.Refresh(request)

	// then
	assert.Empty(t, token)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusUnauthorized, errCode.Code)
	assert.Equal(t, "not authorized", errCode.Error.Error())

	jwtService.AssertExpectations(t)
	userRepository.AssertExpectations(t)
}

func TestAuthService_Refresh_WhenCreateTokenReturnsError_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userRepository := new(repository.UserRepositoryMock)
	_uut := &authService{
		jwtService:     jwtService,
		userRepository: userRepository,
	}
	request := requests.RefreshRequest{
		Token: "This is a token",
	}

	userID := uuid.New()
	jwtService.On("Validate", request.Token).Return(userID, nil).Once()

	user := &models.User{ID: userID}
	userRepository.On("Get", new(models.User), userID).Return(nil, user).Once()

	internalError := utils.InternalServerError(errors.New("something went wrong"))
	jwtService.On("CreateToken", *user).Return("", internalError).Once()

	// when
	token, errCode := _uut.Refresh(request)

	// then
	assert.Empty(t, token)
	assert.NotNil(t, errCode)
	assert.Equal(t, internalError, errCode)

	jwtService.AssertExpectations(t)
	userRepository.AssertExpectations(t)
}

func TestAuthService_Refresh_WhenSuccessful_ShouldReturnNewToken(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	userRepository := new(repository.UserRepositoryMock)
	_uut := &authService{
		jwtService:     jwtService,
		userRepository: userRepository,
	}
	request := requests.RefreshRequest{
		Token: "This is a token",
	}

	userID := uuid.New()
	jwtService.On("Validate", request.Token).Return(userID, nil).Once()

	user := &models.User{ID: userID}
	userRepository.On("Get", new(models.User), userID).Return(nil, user).Once()

	expectedToken := "This is the new token"
	jwtService.On("CreateToken", *user).Return(expectedToken, nil).Once()

	// when
	token, errCode := _uut.Refresh(request)

	// then
	assert.Nil(t, errCode)
	assert.Equal(t, expectedToken, token)

	jwtService.AssertExpectations(t)
	userRepository.AssertExpectations(t)
}

// Sign In
func TestAuthService_SignIn_WhenUserRepositoryReturnsError_ShouldReturnInternalServerError(t *testing.T) {
	// given
	userRepository := new(repository.UserRepositoryMock)
	_uut := &authService{
		userRepository: userRepository,
	}
	request := requests.SignInRequest{
		Email:    "Samuel@yahoo.com",
		Password: "Password123",
	}

	internalError := errors.New("something went wrong")
	userRepository.On("GetByEmail", new(models.User), strings.ToLower(request.Email)).
		Return(internalError).
		Once()

	// when
	token, errCode := _uut.SignIn(request)

	// then
	assert.Empty(t, token)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	userRepository.AssertExpectations(t)
}

func TestAuthService_SignIn_WhenUserIsEmpty_ShouldReturnUnauthorizedError(t *testing.T) {
	// given
	userRepository := new(repository.UserRepositoryMock)
	_uut := &authService{
		userRepository: userRepository,
	}
	request := requests.SignInRequest{
		Email:    "Samuel@yahoo.com",
		Password: "Password123",
	}

	userRepository.On("GetByEmail", new(models.User), strings.ToLower(request.Email)).
		Return(nil).
		Once()

	// when
	token, errCode := _uut.SignIn(request)

	// then
	assert.Empty(t, token)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusUnauthorized, errCode.Code)
	assert.Equal(t, "invalid credentials", errCode.Error.Error())

	userRepository.AssertExpectations(t)
}

func TestAuthService_SignIn_WhenPasswordsAreNotTheSame_ShouldReturnUnauthorizedError(t *testing.T) {
	// given
	bCryptService := new(service.BCryptServiceMock)
	userRepository := new(repository.UserRepositoryMock)
	_uut := &authService{
		bCryptService:  bCryptService,
		userRepository: userRepository,
	}
	request := requests.SignInRequest{
		Email:    "Samuel@yahoo.com",
		Password: "Password123",
	}

	user := &models.User{
		ID:       uuid.New(),
		Email:    "samuel@yahoo.com",
		Password: "hashedPassword",
	}
	userRepository.On("GetByEmail", new(models.User), strings.ToLower(request.Email)).
		Return(nil, user).
		Once()

	bCryptService.On("CompareHash", user.Password, request.Password).Return(errors.New("")).Once()

	// when
	token, errCode := _uut.SignIn(request)

	// then
	assert.Empty(t, token)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusUnauthorized, errCode.Code)
	assert.Equal(t, "invalid credentials", errCode.Error.Error())

	userRepository.AssertExpectations(t)
	bCryptService.AssertExpectations(t)
}

func TestAuthService_SignIn_WhenCreateTokenReturnsError_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	bCryptService := new(service.BCryptServiceMock)
	userRepository := new(repository.UserRepositoryMock)
	_uut := &authService{
		jwtService:     jwtService,
		bCryptService:  bCryptService,
		userRepository: userRepository,
	}
	request := requests.SignInRequest{
		Email:    "Samuel@yahoo.com",
		Password: "Password123",
	}

	user := &models.User{
		ID:       uuid.New(),
		Email:    "samuel@yahoo.com",
		Password: "hashedPassword",
	}
	userRepository.On("GetByEmail", new(models.User), strings.ToLower(request.Email)).
		Return(nil, user).
		Once()

	bCryptService.On("CompareHash", user.Password, request.Password).Return(nil).Once()

	internalError := utils.InternalServerError(errors.New("something went wrong"))
	jwtService.On("CreateToken", *user).Return("", internalError).Once()

	// when
	token, errCode := _uut.SignIn(request)

	// then
	assert.Empty(t, token)
	assert.NotNil(t, errCode)
	assert.Equal(t, internalError, errCode)

	userRepository.AssertExpectations(t)
	bCryptService.AssertExpectations(t)
	jwtService.AssertExpectations(t)
}

func TestAuthService_SignIn_WhenSuccessful_ShouldReturnNewToken(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	bCryptService := new(service.BCryptServiceMock)
	userRepository := new(repository.UserRepositoryMock)
	_uut := &authService{
		jwtService:     jwtService,
		bCryptService:  bCryptService,
		userRepository: userRepository,
	}
	request := requests.SignInRequest{
		Email:    "Samuel@yahoo.com",
		Password: "Password123",
	}

	user := &models.User{
		ID:       uuid.New(),
		Email:    "samuel@yahoo.com",
		Password: "hashedPassword",
	}
	userRepository.On("GetByEmail", new(models.User), strings.ToLower(request.Email)).
		Return(nil, user).
		Once()

	bCryptService.On("CompareHash", user.Password, request.Password).Return(nil).Once()

	expectedToken := "This is the generated token"
	jwtService.On("CreateToken", *user).Return(expectedToken, nil).Once()

	// when
	token, errCode := _uut.SignIn(request)

	// then
	assert.Nil(t, errCode)
	assert.Equal(t, expectedToken, token)

	userRepository.AssertExpectations(t)
	bCryptService.AssertExpectations(t)
	jwtService.AssertExpectations(t)
}

// Sign Up
func TestAuthService_SignUp_WhenUserRepositoryReturnsError_ShouldReturnInternalServerError(t *testing.T) {
	// given
	userRepository := new(repository.UserRepositoryMock)
	_uut := &authService{
		userRepository: userRepository,
	}
	request := requests.SignUpRequest{
		Name:     "Samuel",
		Email:    "Samuel@yahoo.com",
		Password: "Password123",
	}

	internalError := errors.New("something went wrong")
	userRepository.On("GetByEmail", new(models.User), strings.ToLower(request.Email)).
		Return(internalError).
		Once()

	// when
	token, errCode := _uut.SignUp(request)

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
	_uut := &authService{
		userRepository: userRepository,
	}
	request := requests.SignUpRequest{
		Name:     "Samuel",
		Email:    "Samuel@yahoo.com",
		Password: "Password123",
	}

	userRepository.On("GetByEmail", new(models.User), strings.ToLower(request.Email)).
		Return(nil, &models.User{ID: uuid.New()}).
		Once()

	// when
	token, errCode := _uut.SignUp(request)

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
	_uut := &authService{
		bCryptService:  bCryptService,
		userRepository: userRepository,
	}
	request := requests.SignUpRequest{
		Name:     "Samuel",
		Email:    "Samuel@yahoo.com",
		Password: "Password123",
	}

	userRepository.On("GetByEmail", new(models.User), strings.ToLower(request.Email)).
		Return(nil).
		Once()

	internalError := errors.New("internal error")
	bCryptService.On("Hash", request.Password).Return("", internalError).Once()

	// when
	token, errCode := _uut.SignUp(request)

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
	_uut := &authService{
		bCryptService:  bCryptService,
		userRepository: userRepository,
	}
	request := requests.SignUpRequest{
		Name:     "Samuel",
		Email:    "Samuel@yahoo.com",
		Password: "Password123",
	}

	userRepository.On("GetByEmail", new(models.User), strings.ToLower(request.Email)).
		Return(nil).
		Once()

	hashedPassword := "hashed password"
	bCryptService.On("Hash", request.Password).Return(hashedPassword, nil).Once()

	var user *models.User
	internalError := errors.New("internal error")
	userRepository.On("Create", mock.IsType(&models.User{})).
		Run(func(args mock.Arguments) {
			user = args.Get(0).(*models.User)
			assert.NotEmpty(t, user.ID)
			assert.Equal(t, user.Name, request.Name)
			assert.Equal(t, user.Email, strings.ToLower(request.Email))
			assert.Equal(t, user.Password, hashedPassword)
		}).
		Return(internalError).
		Once()

	// when
	token, errCode := _uut.SignUp(request)

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
	_uut := &authService{
		jwtService:     jwtService,
		bCryptService:  bCryptService,
		userRepository: userRepository,
	}
	request := requests.SignUpRequest{
		Name:     "Samuel",
		Email:    "Samuel@yahoo.com",
		Password: "Password123",
	}

	userRepository.On("GetByEmail", new(models.User), strings.ToLower(request.Email)).
		Return(nil).
		Once()

	hashedPassword := "hashed password"
	bCryptService.On("Hash", request.Password).Return(hashedPassword, nil).Once()

	var user *models.User
	userRepository.On("Create", mock.IsType(&models.User{})).
		Run(func(args mock.Arguments) {
			user = args.Get(0).(*models.User)
			assert.NotEmpty(t, user.ID)
			assert.Equal(t, user.Name, request.Name)
			assert.Equal(t, user.Email, strings.ToLower(request.Email))
			assert.Equal(t, user.Password, hashedPassword)
		}).
		Return(nil).
		Once()

	internalError := utils.InternalServerError(errors.New("internal error"))
	jwtService.On("CreateToken", mock.IsType(models.User{})).
		Run(func(args mock.Arguments) {
			assert.Equal(t, *user, args.Get(0).(models.User))
		}).
		Return("", internalError).
		Once()

	// when
	token, errCode := _uut.SignUp(request)

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
	_uut := &authService{
		jwtService:     jwtService,
		bCryptService:  bCryptService,
		userRepository: userRepository,
	}
	request := requests.SignUpRequest{
		Name:     "Samuel",
		Email:    "Samuel@yahoo.com",
		Password: "Password123",
	}

	userRepository.On("GetByEmail", new(models.User), strings.ToLower(request.Email)).
		Return(nil).
		Once()

	hashedPassword := "hashed password"
	bCryptService.On("Hash", request.Password).Return(hashedPassword, nil).Once()

	var user *models.User
	userRepository.On("Create", mock.IsType(user)).
		Run(func(args mock.Arguments) {
			user = args.Get(0).(*models.User)
			assert.NotEmpty(t, user.ID)
			assert.Equal(t, user.Name, request.Name)
			assert.Equal(t, user.Email, strings.ToLower(request.Email))
			assert.Equal(t, user.Password, hashedPassword)
		}).
		Return(nil).
		Once()

	expectedToken := "This is the generated token"
	jwtService.On("CreateToken", mock.IsType(models.User{})).
		Run(func(args mock.Arguments) {
			assert.Equal(t, *user, args.Get(0).(models.User))
		}).
		Return(expectedToken, nil).
		Once()

	// when
	token, errCode := _uut.SignUp(request)

	// then
	assert.Nil(t, errCode)
	assert.Equal(t, expectedToken, token)

	userRepository.AssertExpectations(t)
	bCryptService.AssertExpectations(t)
	jwtService.AssertExpectations(t)
}
