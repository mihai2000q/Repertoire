package service

import (
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"repertoire/api/requests"
	"repertoire/data/repository"
	"repertoire/data/service"
	"repertoire/models"
	"repertoire/utils"
	"strings"
)

type AuthService struct {
	userRepository repository.UserRepository
	jwtService     service.JwtService
	env            utils.Env
}

func NewAuthService(
	userRepository repository.UserRepository,
	jwtService service.JwtService,
	env utils.Env,
) AuthService {
	return AuthService{
		userRepository: userRepository,
		jwtService:     jwtService,
		env:            env,
	}
}

func (a *AuthService) Refresh(request requests.RefreshRequest) (string, *utils.ErrorCode) {
	// validate token
	userId, errCode := a.jwtService.Validate(request.Token)
	if errCode != nil {
		return "", errCode
	}

	// get user
	var user models.User
	err := a.userRepository.Get(&user, userId)
	if err != nil {
		return "", utils.InternalServerError(err)
	}

	return a.jwtService.CreateToken(user)
}

func (a *AuthService) SignIn(request requests.SignInRequest) (string, *utils.ErrorCode) {
	var user models.User

	// get user
	email := strings.ToLower(request.Email)
	err := a.userRepository.GetByEmail(&user, email)
	if err != nil {
		return "", utils.InternalServerError(err)
	}
	if user.ID == uuid.Nil {
		return "", utils.UnauthorizedError(errors.New("invalid credentials"))
	}

	// check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return "", utils.UnauthorizedError(errors.New("invalid credentials"))
	}

	return a.jwtService.CreateToken(user)
}

func (a *AuthService) SignUp(request requests.SignUpRequest) (string, *utils.ErrorCode) {
	var user models.User

	// check if the user already exists
	email := strings.ToLower(request.Email)
	err := a.userRepository.GetByEmail(&user, email)
	if err != nil {
		return "", utils.InternalServerError(err)
	}
	if user.ID != uuid.Nil {
		return "", utils.BadRequestError(errors.New("user already exists"))
	}

	// hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", utils.InternalServerError(err)
	}

	// create user
	user = models.User{
		ID:       uuid.New(),
		Name:     request.Name,
		Email:    email,
		Password: string(hashedPassword),
	}
	err = a.userRepository.Create(&user)
	if err != nil {
		return "", utils.InternalServerError(err)
	}

	return a.jwtService.CreateToken(user)
}
