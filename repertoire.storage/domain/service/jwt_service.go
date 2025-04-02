package service

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"repertoire/storage/data/logger"
	"repertoire/storage/internal"
)

type JwtService interface {
	Authorize(authToken string) error
}

type jwtService struct {
	env    internal.Env
	logger *logger.Logger
}

func NewJwtService(env internal.Env, logger *logger.Logger) JwtService {
	return jwtService{
		env:    env,
		logger: logger,
	}
}

func (j jwtService) Authorize(authToken string) error {
	token, _ := jwt.Parse(authToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.env.JwtSecretKey), nil
	})

	if token != nil && token.Valid {
		if err := j.validateToken(token); err != nil {
			j.logger.Warn("Invalid Token", zap.Error(err), zap.String("token", authToken))
			return err
		}
		return nil
	}
	return errors.New("invalid token")
}

func (j jwtService) validateToken(token *jwt.Token) error {
	// signing method
	if token.Method != jwt.SigningMethodHS256 {
		return errors.New("invalid signing method")
	}

	// audience
	aud, err := token.Claims.GetAudience()
	if err != nil {
		return err
	}
	if len(aud) != 1 || aud[0] != j.env.JwtAudience {
		return errors.New("wrong audience")
	}

	// issuer
	iss, err := token.Claims.GetIssuer()
	if err != nil {
		return err
	}
	if iss != j.env.JwtIssuer {
		return errors.New("wrong issuer")
	}

	// expiration time
	exp, err := token.Claims.GetExpirationTime()
	if err != nil {
		return err
	}
	if exp == nil {
		return errors.New("no expiration time found")
	}

	// jti
	jtiClaim, jtiFound := token.Claims.(jwt.MapClaims)["jti"]
	if !jtiFound {
		return errors.New("missing jti")
	}

	jti, err := uuid.Parse(jtiClaim.(string))
	if err != nil {
		return err
	}
	if jti == uuid.Nil {
		return errors.New("invalid jti")
	}

	// sub
	sub, err := token.Claims.GetSubject()
	if err != nil {
		return err
	}
	userID, err := uuid.Parse(sub)
	if err != nil {
		return err
	}
	if userID == uuid.Nil {
		return errors.New("invalid sub")
	}

	return nil
}
