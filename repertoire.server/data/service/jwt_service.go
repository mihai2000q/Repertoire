package service

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"repertoire/config"
	"repertoire/models"
	"repertoire/utils"
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

func (j *JwtService) Authorize(tokenString string) *utils.ErrorCode {
	token, _ := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.env.JwtSecretKey), nil
	})

	if token != nil && token.Valid {
		return nil
	}

	return utils.UnauthorizedError(errors.New("invalid token"))
}

func (j *JwtService) CreateToken(user models.User) (string, *utils.ErrorCode) {
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
		return "", utils.InternalServerError(err)
	}
	return token, nil
}

func (j *JwtService) Validate(tokenString string) (uuid.UUID, *utils.ErrorCode) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.env.JwtSecretKey), nil
	})
	if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
		return uuid.Nil, utils.UnauthorizedError(err)
	}

	aud, err := token.Claims.GetAudience()
	if err != nil {
		return uuid.Nil, utils.UnauthorizedError(err)
	}

	iss, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, utils.UnauthorizedError(err)
	}

	jtiClaim := token.Claims.(jwt.MapClaims)["jti"].(string)
	jti, err := uuid.Parse(jtiClaim)
	if err != nil {
		return uuid.Nil, utils.UnauthorizedError(err)
	}

	if token.Method != jwt.SigningMethodHS256 ||
		iss != j.env.JwtIssuer ||
		aud[0] != j.env.JwtAudience ||
		jti == uuid.Nil {
		return uuid.Nil, utils.UnauthorizedError(errors.New("invalid token"))
	}

	sub, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, utils.UnauthorizedError(err)
	}

	userId, err := uuid.Parse(sub)
	if err != nil {
		return uuid.Nil, utils.UnauthorizedError(err)
	}

	return userId, nil
}

func (j *JwtService) GetUserIdFromJwt(tokenString string) (uuid.UUID, *utils.ErrorCode) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.env.JwtSecretKey), nil
	})
	if err != nil {
		return uuid.Nil, utils.UnauthorizedError(err)
	}

	sub, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, utils.UnauthorizedError(err)
	}

	userId, err := uuid.Parse(sub)
	if err != nil {
		return uuid.Nil, utils.UnauthorizedError(err)
	}

	return userId, nil
}
