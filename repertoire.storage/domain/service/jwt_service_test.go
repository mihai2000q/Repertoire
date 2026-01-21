package service

import (
	"repertoire/storage/data/logger"
	"repertoire/storage/internal"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// Tests

func TestJwtService_Authorize_WhenTokenIsInvalid_ShouldReturnError(t *testing.T) {
	env := internal.Env{
		JwtSecretKey: "This-is-a-very-long-secret-key-that-is-used-to-encrypt-the-token",
		JwtIssuer:    "JWTIssuer",
		JwtAudience:  "JWTAudience",
	}

	tests := []struct {
		name      string
		claims    *jwt.Token
		secretKey string
	}{
		{
			"When Secret key is wrong",
			jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{}),
			"wrong secret key",
		},
		{
			"When Token has expired",
			jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"exp": time.Now().UTC().Add(-1 * time.Second).Unix(),
			}),
			env.JwtSecretKey,
		},
		// signing method
		{
			"When signing method is wrong",
			jwt.NewWithClaims(jwt.SigningMethodHS384, jwt.MapClaims{}),
			env.JwtSecretKey,
		},
		// audience
		{
			"When Audience is missing",
			jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{}),
			env.JwtSecretKey,
		},
		{
			"When audience has too many elements",
			jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"aud": []string{env.JwtAudience, "some audience"},
			}),
			env.JwtSecretKey,
		},
		{
			"When audience is wrong",
			jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"aud": "some audience",
			}),
			env.JwtSecretKey,
		},
		// issuer
		{
			"When Issuer is missing",
			jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"aud": env.JwtAudience,
			}),
			env.JwtSecretKey,
		},
		{
			"When issuer is wrong",
			jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"aud": env.JwtAudience,
				"iss": "some issuer",
			}),
			env.JwtSecretKey,
		},
		// expiration time
		{
			"When Expiration Time is missing",
			jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"aud": env.JwtAudience,
				"iss": env.JwtIssuer,
			}),
			env.JwtSecretKey,
		},
		// jti
		{
			"When Jti is missing",
			jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"aud": env.JwtAudience,
				"iss": env.JwtIssuer,
				"exp": time.Now().UTC().Add(time.Minute).Unix(),
			}),
			env.JwtSecretKey,
		},
		{
			"When Jti is not uuid",
			jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"aud": env.JwtAudience,
				"iss": env.JwtIssuer,
				"exp": time.Now().UTC().Add(time.Minute).Unix(),
				"jti": "something",
			}),
			env.JwtSecretKey,
		},
		{
			"When jti is uuid nil",
			jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"aud": env.JwtAudience,
				"iss": env.JwtIssuer,
				"exp": time.Now().UTC().Add(time.Minute).Unix(),
				"jti": uuid.Nil.String(),
			}),
			env.JwtSecretKey,
		},
		// sub
		{
			"When sub is missing",
			jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"aud": env.JwtAudience,
				"iss": env.JwtIssuer,
				"exp": time.Now().UTC().Add(time.Minute).Unix(),
				"jti": uuid.New().String(),
			}),
			env.JwtSecretKey,
		},
		{
			"When sub is not uuid",
			jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"aud": env.JwtAudience,
				"iss": env.JwtIssuer,
				"exp": time.Now().UTC().Add(time.Minute).Unix(),
				"jti": uuid.New().String(),
				"sub": "something",
			}),
			env.JwtSecretKey,
		},
		{
			"When sub is uuid nil",
			jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"aud": env.JwtAudience,
				"iss": env.JwtIssuer,
				"exp": time.Now().UTC().Add(time.Minute).Unix(),
				"jti": uuid.New().String(),
				"sub": uuid.Nil.String(),
			}),
			env.JwtSecretKey,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			_uut := NewJwtService(env, logger.NewLoggerMock())

			token, _ := tt.claims.SignedString([]byte(tt.secretKey))

			// when
			err := _uut.Authorize(token)

			// then
			assert.Error(t, err)
		})
	}
}

func TestJwtService_Authorize_WhenSuccessful_ShouldNotReturnError(t *testing.T) {
	env := internal.Env{
		JwtSecretKey: "This-is-a-very-long-secret-key-that-is-used-to-encrypt-the-token",
		JwtIssuer:    "JWTIssuer",
		JwtAudience:  "JWTAudience",
	}

	tests := []struct {
		name   string
		claims *jwt.Token
	}{
		{
			"Normal Token",
			jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"jti": uuid.New().String(),
				"sub": uuid.New().String(),
				"iss": env.JwtIssuer,
				"aud": env.JwtAudience,
				"iat": time.Now().UTC().Unix(),
				"exp": time.Now().UTC().Add(time.Minute).Unix(),
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			_uut := NewJwtService(env, nil)

			token, _ := tt.claims.SignedString([]byte(env.JwtSecretKey))

			// when
			err := _uut.Authorize(token)

			// then
			assert.NoError(t, err)
		})
	}
}
