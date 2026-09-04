package service

import (
	"repertoire/server/internal/env"
	"repertoire/server/internal/httperror"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JwtService interface {
	Authorize(tokenString string) *httperror.ErrorCode
	GetUserIdFromJwt(tokenString string) (uuid.UUID, *httperror.ErrorCode)
}

type jwtService struct {
	env env.Env
}

func NewJwtService(env env.Env) JwtService {
	return &jwtService{
		env: env,
	}
}

func (j *jwtService) Authorize(tokenString string) *httperror.ErrorCode {
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(j.env.JwtPublicKey))
	if err != nil {
		return httperror.InternalServerError(err)
	}
	_, err = jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	if err != nil {
		return httperror.UnauthorizedError(err)
	}

	return nil
}

func (j *jwtService) GetUserIdFromJwt(tokenString string) (uuid.UUID, *httperror.ErrorCode) {
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(j.env.JwtPublicKey))
	if err != nil {
		return uuid.Nil, httperror.InternalServerError(err)
	}
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	if err != nil {
		return uuid.Nil, httperror.ForbiddenError(err)
	}

	sub, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, httperror.ForbiddenError(err)
	}

	userID, err := uuid.Parse(sub)
	if err != nil {
		return uuid.Nil, httperror.ForbiddenError(err)
	}

	return userID, nil
}
