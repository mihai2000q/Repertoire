package service

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"repertoire/config"
	"repertoire/models"
	"time"
)

type JwtService struct {
	env config.Env
}

func NewJwtService(env config.Env) JwtService {
	return JwtService{
		env: env,
	}
}

func (j *JwtService) Authorize(tokenString string) error {
	token, _ := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.env.JwtSecretKey), nil
	})
	if token.Valid {
		return nil
	}

	return errors.New("invalid token")
}

func (j *JwtService) CreateToken(user models.User) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"jti":   uuid.New().String(),
		"sub":   user.ID.String(),
		"email": user.Email,
		"iss":   j.env.JwtIssuer,
		"aud":   j.env.JwtAudience,
		"iat":   time.Now().UTC().Unix(),
		"exp":   time.Now().UTC().Add(time.Hour).Unix(),
	})
	token, err := claims.SignedString([]byte(j.env.JwtSecretKey))
	if err != nil {
		return "", err
	}
	return token, nil
}
