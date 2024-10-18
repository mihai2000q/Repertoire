package service

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"repertoire/models"
	"repertoire/utils"
	"repertoire/utils/wrapper"
	"time"
)

type JwtService interface {
	Authorize(tokenString string) *wrapper.ErrorCode
	CreateToken(user models.User) (string, *wrapper.ErrorCode)
	Validate(tokenString string) (uuid.UUID, *wrapper.ErrorCode)
	GetUserIdFromJwt(tokenString string) (uuid.UUID, *wrapper.ErrorCode)
}

type jwtService struct {
	env utils.Env
}

func NewJwtService(env utils.Env) JwtService {
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

func (j *jwtService) CreateToken(user models.User) (string, *wrapper.ErrorCode) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"jti": uuid.New().String(),
		"sub": user.ID.String(),
		"iss": j.env.JwtIssuer,
		"aud": j.env.JwtAudience,
		"iat": time.Now().UTC().Unix(),
		"exp": time.Now().UTC().Add(time.Hour).Unix(),
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

	userId, err := uuid.Parse(sub)
	if err != nil {
		return uuid.Nil, wrapper.UnauthorizedError(err)
	}

	return userId, nil
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

	userId, err := uuid.Parse(sub)
	if err != nil {
		return uuid.Nil, wrapper.ForbiddenError(err)
	}

	return userId, nil
}
