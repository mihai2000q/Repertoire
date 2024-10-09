package service

import (
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"repertoire/api/requests/auth"
	"repertoire/data/repository"
	"repertoire/data/service"
	"repertoire/models"
	"repertoire/utils"
	"strings"
	"time"
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

func (a *AuthService) Refresh(request auth.RefreshRequest) (string, *utils.ErrorCode) {
	// validate token
	userId, errCode := a.jwtService.Validate(request.AccessToken)
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

func (a *AuthService) SignIn(request auth.SignInRequest) (string, *utils.ErrorCode) {
	var user models.User

	// get user
	email := strings.ToLower(request.Email)
	err := a.userRepository.GetByEmail(&user, email)
	if err != nil {
		return "", utils.InternalServerError(err)
	}

	// check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return "", utils.InternalServerError(err)
	}

	return a.jwtService.CreateToken(user)
}

func (a *AuthService) SignUp(request auth.SignUpRequest) (string, *utils.ErrorCode) {
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
		ID:        uuid.New(),
		Name:      request.Name,
		Email:     email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	err = a.userRepository.Create(&user)
	if err != nil {
		return "", utils.InternalServerError(err)
	}

	return a.jwtService.CreateToken(user)
}
