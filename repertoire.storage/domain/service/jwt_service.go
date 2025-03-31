package service

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"repertoire/storage/internal"
)

type JwtService interface {
	Authorize(authToken string) error
}

type jwtService struct {
	env internal.Env
}

func NewJwtService(env internal.Env) JwtService {
	return jwtService{env: env}
}

func (j jwtService) Authorize(authToken string) error {
	token, _ := jwt.Parse(authToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.env.JwtSecretKey), nil
	})

	if token != nil && token.Valid {
		if err := j.validateToken(token); err != nil {
			return err
		}
		return nil
	}
	return errors.New("invalid token")
}

func (j jwtService) validateToken(token *jwt.Token) error {
	aud, err := token.Claims.GetAudience()
	if err != nil {
		return err
	}
	iss, err := token.Claims.GetIssuer()
	if err != nil {
		return err
	}
	_, err = token.Claims.GetExpirationTime()
	if err != nil {
		return err
	}

	jtiClaim, jtiFound := token.Claims.(jwt.MapClaims)["jti"]
	if !jtiFound {
		return errors.New("no jti found")
	}

	jti, err := uuid.Parse(jtiClaim.(string))
	if err != nil {
		return err
	}

	sub, err := token.Claims.GetSubject()
	if err != nil {
		return err
	}
	userID, err := uuid.Parse(sub)
	if err != nil {
		return err
	}

	if token.Method != jwt.SigningMethodHS256 ||
		iss != j.env.JwtIssuer ||
		aud[0] != j.env.JwtAudience ||
		len(aud) != 1 ||
		jti == uuid.Nil ||
		userID == uuid.Nil {
		return errors.New("invalid token")
	}

	return nil
}
