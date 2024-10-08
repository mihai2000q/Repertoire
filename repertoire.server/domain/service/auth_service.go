package service

import (
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"repertoire/api/contracts/auth"
	"repertoire/config"
	"repertoire/data/repository"
	"repertoire/data/service"
	"repertoire/models"
	"strings"
	"time"
)

type AuthService struct {
	userRepository repository.UserRepository
	jwtService     service.JwtService
	env            config.Env
}

func NewAuthService(
	userRepository repository.UserRepository,
	jwtService service.JwtService,
	env config.Env,
) AuthService {
	return AuthService{
		userRepository: userRepository,
		jwtService:     jwtService,
		env:            env,
	}
}

func (a *AuthService) SignIn(request auth.SignInRequest) (string, error) {
	var user models.User

	// get user
	email := strings.ToLower(request.Email)
	err := a.userRepository.GetByEmail(&user, email)
	if err != nil {
		return "", err
	}

	// check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return "", err
	}

	return a.jwtService.CreateToken(user)
}

func (a *AuthService) SignUp(request auth.SignUpRequest) (string, error) {
	var user models.User

	// check if the user already exists
	email := strings.ToLower(request.Email)
	err := a.userRepository.GetByEmail(&user, email)
	if err != nil {
		return "", err
	}
	if user.ID != uuid.Nil {
		return "", errors.New("user already exists")
	}

	// hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
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
		return "", err
	}

	return a.jwtService.CreateToken(user)
}
