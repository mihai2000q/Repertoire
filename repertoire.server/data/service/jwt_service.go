package service

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"repertoire/server/internal"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"time"
)

type JwtService interface {
	Authorize(tokenString string) *wrapper.ErrorCode
	CreateToken(user model.User) (string, *wrapper.ErrorCode)
	Validate(tokenString string) (uuid.UUID, *wrapper.ErrorCode)
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
	token, _ := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.env.JwtSecretKey), nil
	})

	if token != nil && token.Valid {
		return nil
	}

	return wrapper.UnauthorizedError(errors.New("invalid token"))
}

func (j *jwtService) CreateToken(user model.User) (string, *wrapper.ErrorCode) {
	expiresIn, err := time.ParseDuration(j.env.JwtExpirationTime)
	if err != nil {
		return "", wrapper.InternalServerError(err)
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"jti": uuid.New().String(),
		"sub": user.ID.String(),
		"iss": j.env.JwtIssuer,
		"aud": j.env.JwtAudience,
		"iat": time.Now().UTC().Unix(),
		"exp": time.Now().UTC().Add(expiresIn).Unix(),
	})
	token, err := claims.SignedString([]byte(j.env.JwtSecretKey))
	if err != nil {
		return "", wrapper.InternalServerError(err)
	}
	return token, nil
}

func (j *jwtService) Validate(tokenString string) (uuid.UUID, *wrapper.ErrorCode) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.env.JwtSecretKey), nil
	})
	if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
		return uuid.Nil, wrapper.UnauthorizedError(err)
	}

	aud, err := token.Claims.GetAudience()
	if err != nil {
		return uuid.Nil, wrapper.UnauthorizedError(err)
	}

	iss, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, wrapper.UnauthorizedError(err)
	}

	jtiClaim := token.Claims.(jwt.MapClaims)["jti"].(string)
	jti, err := uuid.Parse(jtiClaim)
	if err != nil {
		return uuid.Nil, wrapper.UnauthorizedError(err)
	}

	if token.Method != jwt.SigningMethodHS256 ||
		iss != j.env.JwtIssuer ||
		aud[0] != j.env.JwtAudience ||
		jti == uuid.Nil {
		return uuid.Nil, wrapper.UnauthorizedError(errors.New("invalid token"))
	}

	sub, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, wrapper.UnauthorizedError(err)
	}

	userID, err := uuid.Parse(sub)
	if err != nil {
		return uuid.Nil, wrapper.UnauthorizedError(err)
	}

	return userID, nil
}

func (j *jwtService) GetUserIdFromJwt(tokenString string) (uuid.UUID, *wrapper.ErrorCode) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.env.JwtSecretKey), nil
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
