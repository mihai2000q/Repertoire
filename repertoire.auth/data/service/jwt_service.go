package service

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"repertoire/auth/data/logger"
	"repertoire/auth/internal"
	"repertoire/auth/internal/wrapper"
	"repertoire/auth/model"
	"time"
)

type JwtService interface {
	Authorize(authToken string) *wrapper.ErrorCode
	GetUserIDFromJwt(tokenString string) (uuid.UUID, *wrapper.ErrorCode)

	Validate(tokenString string) (uuid.UUID, *wrapper.ErrorCode)
	ValidateCredentials(clientCredentials model.ClientCredentials) *wrapper.ErrorCode

	CreateToken(user model.User) (string, *wrapper.ErrorCode)
	CreateCentrifugoToken(userID uuid.UUID) (string, string, *wrapper.ErrorCode)
	CreateStorageToken(userID uuid.UUID) (string, string, *wrapper.ErrorCode)
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

func (j jwtService) Authorize(authToken string) *wrapper.ErrorCode {
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(j.env.JwtPublicKey))
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	token, err := jwt.Parse(authToken, func(t *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	if err != nil {
		return wrapper.UnauthorizedError(err)
	}

	if token != nil && token.Valid {
		return nil
	}
	return wrapper.UnauthorizedError(errors.New("invalid token"))
}

func (j jwtService) GetUserIDFromJwt(tokenString string) (uuid.UUID, *wrapper.ErrorCode) {
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

// Validation

func (j jwtService) Validate(tokenString string) (uuid.UUID, *wrapper.ErrorCode) {
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(j.env.JwtPublicKey))
	if err != nil {
		return uuid.Nil, wrapper.InternalServerError(err)
	}
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
		return uuid.Nil, wrapper.UnauthorizedError(err)
	}

	// signing method
	if token.Method != jwt.SigningMethodRS256 {
		return uuid.Nil, wrapper.UnauthorizedError(errors.New("invalid signing method"))
	}

	// audience
	aud, err := token.Claims.GetAudience()
	if err != nil {
		return uuid.Nil, wrapper.UnauthorizedError(err)
	}
	if len(aud) != 1 || aud[0] != j.env.JwtAudience {
		return uuid.Nil, wrapper.UnauthorizedError(errors.New("wrong audience"))
	}

	// issuer
	iss, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, wrapper.UnauthorizedError(err)
	}
	if iss != j.env.JwtIssuer {
		return uuid.Nil, wrapper.UnauthorizedError(errors.New("wrong issuer"))
	}

	// expiration time
	exp, err := token.Claims.GetExpirationTime()
	if err != nil {
		return uuid.Nil, wrapper.UnauthorizedError(err)
	}
	if exp == nil {
		return uuid.Nil, wrapper.UnauthorizedError(errors.New("missing expiration time"))
	}

	// jti
	jtiClaim, jtiFound := token.Claims.(jwt.MapClaims)["jti"]
	if !jtiFound {
		return uuid.Nil, wrapper.UnauthorizedError(errors.New("missing jti"))
	}

	jti, err := uuid.Parse(jtiClaim.(string))
	if err != nil {
		return uuid.Nil, wrapper.UnauthorizedError(err)
	}
	if jti == uuid.Nil {
		return uuid.Nil, wrapper.UnauthorizedError(errors.New("invalid jti"))
	}

	// sub
	sub, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, wrapper.UnauthorizedError(err)
	}

	userID, err := uuid.Parse(sub)
	if err != nil {
		return uuid.Nil, wrapper.UnauthorizedError(err)
	}
	if userID == uuid.Nil {
		return uuid.Nil, wrapper.UnauthorizedError(errors.New("invalid sub"))
	}

	return userID, nil
}

func (j jwtService) ValidateCredentials(clientCredentials model.ClientCredentials) *wrapper.ErrorCode {
	if clientCredentials.GrantType != "client_credentials" ||
		clientCredentials.ClientID != j.env.ClientID ||
		clientCredentials.ClientSecret != j.env.ClientSecret {
		return wrapper.UnauthorizedError(errors.New("you are not authorized"))
	}
	return nil
}

// Create Tokens

func (j jwtService) CreateToken(user model.User) (string, *wrapper.ErrorCode) {
	expiresIn, err := time.ParseDuration(j.env.JwtExpirationTime)
	if err != nil {
		return "", wrapper.InternalServerError(err)
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(j.env.JwtPrivateKey))
	if err != nil {
		return "", wrapper.InternalServerError(err)
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"jti": uuid.New().String(),
		"sub": user.ID.String(),
		"iss": j.env.JwtIssuer,
		"aud": "Repertoire",
		"iat": time.Now().UTC().Unix(),
		"exp": time.Now().UTC().Add(expiresIn).Unix(),
	})
	token, err := claims.SignedString(privateKey)
	if err != nil {
		return "", wrapper.InternalServerError(err)
	}
	return token, nil
}

func (j jwtService) CreateCentrifugoToken(userID uuid.UUID) (string, string, *wrapper.ErrorCode) {
	expiresIn, err := time.ParseDuration(j.env.CentrifugoJwtExpirationTime)
	if err != nil {
		return "", "", wrapper.InternalServerError(err)
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"jti": uuid.New().String(),
		"sub": userID,
		"iss": j.env.JwtIssuer,
		"aud": j.env.CentrifugoJwtAudience,
		"iat": time.Now().UTC().Unix(),
		"exp": time.Now().UTC().Add(expiresIn).Unix(),
	})
	token, err := claims.SignedString([]byte(j.env.CentrifugoJwtSecretKey))
	if err != nil {
		return "", "", wrapper.InternalServerError(err)
	}
	return token, j.env.CentrifugoJwtExpirationTime, nil
}

func (j jwtService) CreateStorageToken(userID uuid.UUID) (string, string, *wrapper.ErrorCode) {
	expiresIn, err := time.ParseDuration(j.env.StorageJwtExpirationTime)
	if err != nil {
		return "", "", wrapper.InternalServerError(err)
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"jti": uuid.New().String(),
		"sub": userID,
		"iss": j.env.JwtIssuer,
		"aud": j.env.StorageJwtAudience,
		"iat": time.Now().UTC().Unix(),
		"exp": time.Now().UTC().Add(expiresIn).Unix(),
	})
	token, err := claims.SignedString([]byte(j.env.StorageJwtSecretKey))
	if err != nil {
		return "", "", wrapper.InternalServerError(err)
	}
	return token, j.env.StorageJwtExpirationTime, nil
}
