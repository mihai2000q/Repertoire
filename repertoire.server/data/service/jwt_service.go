package service

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"repertoire/server/internal"
	"repertoire/server/internal/wrapper"
)

type JwtService interface {
	Authorize(tokenString string) *wrapper.ErrorCode
	GetUserIdFromJwt(tokenString string) (uuid.UUID, *wrapper.ErrorCode)
}

type jwtService struct {
	env internal.Env
}

func NewJwtService(env internal.Env) JwtService {
	return &jwtService{
		env: env,
	}
}

func (j *jwtService) Authorize(tokenString string) *wrapper.ErrorCode {
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(j.env.JwtPublicKey))
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	token, _ := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})

	if token != nil && token.Valid {
		return nil
	}

	return wrapper.UnauthorizedError(errors.New("invalid token"))
}

func (j *jwtService) GetUserIdFromJwt(tokenString string) (uuid.UUID, *wrapper.ErrorCode) {
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(j.env.JwtPublicKey))
	if err != nil {
		return uuid.Nil, wrapper.InternalServerError(err)
	}
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	if err != nil {
		return uuid.Nil, wrapper.ForbiddenError(err)
	}

	sub, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, wrapper.ForbiddenError(err)
	}

	userID, err := uuid.Parse(sub)
	if err != nil {
		return uuid.Nil, wrapper.ForbiddenError(err)
	}

	return userID, nil
}
