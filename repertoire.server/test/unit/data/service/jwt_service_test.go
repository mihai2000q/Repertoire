package service

import (
	"net/http"
	"repertoire/server/data/service"
	"repertoire/server/internal"
	"repertoire/server/model"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestJwtService_Authorize_WhenTokenIsInvalid_ShouldReturnUnauthorizedError(t *testing.T) {
	env := internal.Env{
		JwtSecretKey: "This is a secret key that is used to encrypt the token",
	}

	tests := []struct {
		name      string
		claims    *jwt.Token
		secretKey string
	}{
		{
			"when expiration time has elapsed",
			jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"iat": time.Now().UTC().Unix(),
				"exp": time.Now().UTC().Add(-time.Second).Unix(),
			}),
			env.JwtSecretKey,
		},
		{
			"when signing method is not the same",
			jwt.NewWithClaims(jwt.SigningMethodES512, jwt.MapClaims{
				"iat": time.Now().UTC().Unix(),
				"exp": time.Now().UTC().Add(time.Minute).Unix(),
			}),
			env.JwtSecretKey,
		},
		{
			"when secret key does not match",
			jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"iat": time.Now().UTC().Unix(),
				"exp": time.Now().UTC().Add(time.Minute).Unix(),
			}),
			"This is another secret key, not the original one and the test shall not pass",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			_uut := service.NewJwtService(env)

			token, _ := tt.claims.SignedString([]byte(tt.secretKey))

			// when
			errCode := _uut.Authorize(token)

			// then
			assert.NotNil(t, errCode)
			assert.Equal(t, http.StatusUnauthorized, errCode.Code)
			assert.Equal(t, "invalid token", errCode.Error.Error())
		})
	}
}

func TestJwtService_Authorize_WhenTokenIsValid_ShouldNotReturnAnyError(t *testing.T) {
	// given
	env := internal.Env{
		JwtSecretKey: "This is a secret key that is used to encrypt the token",
	}
	_uut := service.NewJwtService(env)

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iat": time.Now().UTC().Unix(),
		"exp": time.Now().UTC().Add(time.Hour).Unix(),
	})
	token, _ := claims.SignedString([]byte(env.JwtSecretKey))

	// when
	errCode := _uut.Authorize(token)

	// then
	assert.Nil(t, errCode)
}

func TestJwtService_CreateToken_WhenItFails_ShouldReturnInternalError(t *testing.T) {
	// given
	env := internal.Env{
		JwtSecretKey:      "This is a secret key that is used to encrypt the token",
		JwtIssuer:         "Repertoire",
		JwtAudience:       "Repertoire",
		JwtExpirationTime: "something",
	}
	_uut := service.NewJwtService(env)

	user := model.User{
		ID: uuid.New(),
	}

	// when
	tokenString, errCode := _uut.CreateToken(user)

	// then
	assert.Empty(t, tokenString)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Error(t, errCode.Error)
}

func TestJwtService_CreateToken_WhenSuccessful_ShouldReturnNewToken(t *testing.T) {
	// given
	env := internal.Env{
		JwtSecretKey:      "This is a secret key that is used to encrypt the token",
		JwtIssuer:         "Repertoire",
		JwtAudience:       "Repertoire",
		JwtExpirationTime: "1h",
	}
	_uut := service.NewJwtService(env)

	user := model.User{
		ID: uuid.New(),
	}

	expiresIn, _ := time.ParseDuration(env.JwtExpirationTime)

	// when
	tokenString, errCode := _uut.CreateToken(user)

	// then
	assert.Nil(t, errCode)

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(env.JwtSecretKey), nil
	})
	assert.NoError(t, err)

	jtiClaim := token.Claims.(jwt.MapClaims)["jti"].(string)
	jti, err := uuid.Parse(jtiClaim)
	assert.NoError(t, err)
	sub, err := token.Claims.GetSubject()
	assert.NoError(t, err)
	aud, err := token.Claims.GetAudience()
	assert.NoError(t, err)
	iss, err := token.Claims.GetIssuer()
	assert.NoError(t, err)
	iat, err := token.Claims.GetIssuedAt()
	assert.NoError(t, err)
	exp, err := token.Claims.GetExpirationTime()
	assert.NoError(t, err)

	assert.Equal(t, jwt.SigningMethodHS256, token.Method)
	assert.NotEmpty(t, jti)
	assert.Equal(t, user.ID.String(), sub)
	assert.Len(t, aud, 1)
	assert.Equal(t, env.JwtAudience, aud[0])
	assert.Equal(t, env.JwtIssuer, iss)
	assert.WithinDuration(t, time.Now().UTC(), iat.Time, 10*time.Second)
	assert.WithinDuration(t, time.Now().Add(expiresIn).UTC(), exp.Time, 10*time.Second)
}

func TestJwtService_Validate_WhenTokenIsInvalid_ShouldReturnUnauthorizedError(t *testing.T) {
	env := internal.Env{
		JwtSecretKey: "This is a secret key that is used to encrypt the token",
		JwtAudience:  "Repertoire",
		JwtIssuer:    "Repertoire",
	}

	tests := []struct {
		name      string
		claims    *jwt.Token
		secretKey string
	}{
		{
			"when audience is missing",
			jwt.NewWithClaims(jwt.SigningMethodES384, jwt.MapClaims{
				"jti": uuid.New().String(),
				"sub": uuid.New().String(),
				"iss": env.JwtIssuer,
				"aud": env.JwtAudience,
			}),
			"Some different key",
		},
		{
			"when audience is missing",
			jwt.NewWithClaims(jwt.SigningMethodES384, jwt.MapClaims{}),
			env.JwtSecretKey,
		},
		{
			"when issuer is missing",
			jwt.NewWithClaims(jwt.SigningMethodES384, jwt.MapClaims{
				"aud": env.JwtAudience,
			}),
			env.JwtSecretKey,
		},
		{
			"when jti is missing",
			jwt.NewWithClaims(jwt.SigningMethodES384, jwt.MapClaims{
				"iss": env.JwtIssuer,
				"aud": env.JwtAudience,
			}),
			env.JwtSecretKey,
		},
		{
			"when Signing Method is not the same",
			jwt.NewWithClaims(jwt.SigningMethodES384, jwt.MapClaims{
				"jti": uuid.New().String(),
				"iss": env.JwtIssuer,
				"aud": env.JwtAudience,
			}),
			env.JwtSecretKey,
		},
		{
			"when issuer is not matching",
			jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"jti": uuid.New().String(),
				"iss": "some issuer",
				"aud": env.JwtAudience,
			}),
			env.JwtSecretKey,
		},
		{
			"when audience is not matching",
			jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"jti": uuid.New().String(),
				"iss": env.JwtIssuer,
				"aud": "some audience",
			}),
			env.JwtSecretKey,
		},
		{
			"when jti is uuid nil",
			jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"jti": uuid.Nil,
				"iss": env.JwtIssuer,
				"aud": env.JwtAudience,
			}),
			env.JwtSecretKey,
		},
		{
			"when sub is missing",
			jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"jti": uuid.New().String(),
				"iss": env.JwtIssuer,
				"aud": env.JwtAudience,
			}),
			env.JwtSecretKey,
		},
		{
			"when sub is not uuid",
			jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"jti": uuid.New().String(),
				"sub": "This is a sub",
				"iss": env.JwtIssuer,
				"aud": env.JwtAudience,
			}),
			env.JwtSecretKey,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			_uut := service.NewJwtService(env)

			token, _ := tt.claims.SignedString([]byte(tt.secretKey))

			// when
			userID, errCode := _uut.Validate(token)

			// then
			assert.Empty(t, userID)
			assert.NotNil(t, errCode)
			assert.Equal(t, http.StatusUnauthorized, errCode.Code)
			assert.Error(t, errCode.Error)
		})
	}
}

func TestJwtService_Validate_WhenSuccessful_ShouldReturnUserID(t *testing.T) {
	// given
	env := internal.Env{
		JwtSecretKey: "This is a secret key that is used to encrypt the token",
	}
	_uut := service.NewJwtService(env)

	user := model.User{
		ID: uuid.New(),
	}

	// expired token
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"jti": uuid.New().String(),
		"sub": user.ID.String(),
		"iss": env.JwtIssuer,
		"aud": env.JwtAudience,
		"iat": time.Now().UTC().Unix(),
		"exp": time.Now().UTC().Add(-time.Hour).Unix(),
	})
	token, _ := claims.SignedString([]byte(env.JwtSecretKey))

	// when
	userID, errCode := _uut.Validate(token)

	// then
	assert.Equal(t, userID, user.ID)
	assert.Nil(t, errCode)
}

func TestJwtService_GetUserIdFromJwt_WhenSecretKeyIsNotMatching_ShouldReturnForbiddenError(t *testing.T) {
	// given
	env := internal.Env{
		JwtSecretKey: "This is a secret key that is used to encrypt the token",
	}
	_uut := service.NewJwtService(env)

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{})
	token, _ := claims.SignedString([]byte("some other key"))

	// when
	userID, errCode := _uut.GetUserIdFromJwt(token)

	// then
	assert.Empty(t, userID)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusForbidden, errCode.Code)
	assert.Error(t, errCode.Error)
}

func TestJwtService_GetUserIdFromJwt_WhenSubIsMissing_ShouldReturnForbiddenError(t *testing.T) {
	// given
	env := internal.Env{
		JwtSecretKey: "This is a secret key that is used to encrypt the token",
	}
	_uut := service.NewJwtService(env)

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{})
	token, _ := claims.SignedString([]byte(env.JwtSecretKey))

	// when
	userID, errCode := _uut.GetUserIdFromJwt(token)

	// then
	assert.Empty(t, userID)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusForbidden, errCode.Code)
	assert.Error(t, errCode.Error)
}

func TestJwtService_GetUserIdFromJwt_WhenSubIsNotUUID_ShouldReturnForbiddenError(t *testing.T) {
	// given
	env := internal.Env{
		JwtSecretKey: "This is a secret key that is used to encrypt the token",
	}
	_uut := service.NewJwtService(env)

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "something-else",
	})
	token, _ := claims.SignedString([]byte(env.JwtSecretKey))

	// when
	userID, errCode := _uut.GetUserIdFromJwt(token)

	// then
	assert.Empty(t, userID)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusForbidden, errCode.Code)
	assert.Error(t, errCode.Error)
}

func TestJwtService_GetUserIdFromJwt_WhenSuccessful_ShouldReturnUserId(t *testing.T) {
	// given
	env := internal.Env{
		JwtSecretKey: "This is a secret key that is used to encrypt the token",
	}
	_uut := service.NewJwtService(env)

	user := model.User{
		ID: uuid.New(),
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID.String(),
	})
	token, _ := claims.SignedString([]byte(env.JwtSecretKey))

	// when
	userID, errCode := _uut.GetUserIdFromJwt(token)

	// then
	assert.Equal(t, userID, user.ID)
	assert.Nil(t, errCode)
}
