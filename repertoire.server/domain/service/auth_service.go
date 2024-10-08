package service

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"repertoire/api/contracts/auth"
	"repertoire/config"
	"repertoire/data/repository"
	"repertoire/models"
	"strings"
	"time"
)

type AuthService struct {
	userRepository repository.UserRepository
	env            config.Env
}

func NewAuthService(userRepository repository.UserRepository, env config.Env) AuthService {
	return AuthService{
		userRepository: userRepository,
		env:            env,
	}
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

	// create, sign and return token
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"jti":   uuid.New().String(),
		"sub":   user.ID.String(),
		"email": user.Email,
		"iss":   a.env.JwtIssuer,
		"aud":   a.env.JwtAudience,
		"iat":   time.Now().UTC().Unix(),
		"exp":   time.Now().UTC().Add(time.Hour).Unix(),
	})
	token, err := claims.SignedString([]byte(a.env.JwtSecretKey))
	if err != nil {
		return "", err
	}
	return token, nil
}
