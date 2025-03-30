package service

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"repertoire/storage/internal"
	"testing"
	"time"
)

// Utils

func createToken(
	signingMethod jwt.SigningMethod,
	jwtExpirationTime string,
	jwtIssuer string,
	jwtAudience string,
	jwtSecretKey string,
	jti string,
) string {
	expiresIn, _ := time.ParseDuration(jwtExpirationTime)

	claims := jwt.NewWithClaims(signingMethod, jwt.MapClaims{
		"jti": jti,
		"iss": jwtIssuer,
		"aud": jwtAudience,
		"iat": time.Now().UTC().Unix(),
		"exp": time.Now().UTC().Add(expiresIn).Unix(),
	})
	token, _ := claims.SignedString([]byte(jwtSecretKey))
	return token
}

func validateToken(t *testing.T, tokenString string, env internal.Env) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(env.JwtSecretKey), nil
	})
	assert.NoError(t, err)

	aud, err := token.Claims.GetAudience()
	assert.NoError(t, err)

	iss, err := token.Claims.GetIssuer()
	assert.NoError(t, err)

	jtiClaim := token.Claims.(jwt.MapClaims)["jti"].(string)
	jti, err := uuid.Parse(jtiClaim)
	assert.NoError(t, err)

	assert.Equal(t, jwt.SigningMethodHS256, token.Method)
	assert.Equal(t, env.JwtIssuer, iss)
	assert.Len(t, aud, 1)
	assert.Equal(t, env.JwtAudience, aud[0])
	assert.NotEmpty(t, jti)
}

// Tests

func TestJwtService_Authorize_WhenTokenIsInvalid_ShouldReturnError(t *testing.T) {
	env := internal.Env{
		JwtSecretKey: "This-is-a-very-long-secret-key-that-is-used-to-encrypt-the-token",
		JwtIssuer:    "JWTIssuer",
		JwtAudience:  "JWTAudience",
	}

	tests := []struct {
		name  string
		token string
	}{
		{
			"When Secret key is wrong",
			createToken(
				jwt.SigningMethodHS256,
				"1h",
				env.JwtIssuer,
				env.JwtAudience,
				"secret is wrong",
				uuid.New().String(),
			),
		},
		{
			"When Token has expired",
			createToken(
				jwt.SigningMethodHS256,
				"-1s",
				env.JwtIssuer,
				env.JwtAudience,
				env.JwtSecretKey,
				uuid.New().String(),
			),
		},
		{
			"When Signing Method is wrong",
			createToken(
				jwt.SigningMethodEdDSA,
				"1h",
				env.JwtIssuer,
				env.JwtAudience,
				env.JwtSecretKey,
				uuid.New().String(),
			),
		},
		{
			"When Issuer is wrong",
			createToken(
				jwt.SigningMethodHS256,
				"1h",
				"Random issuer",
				env.JwtAudience,
				env.JwtSecretKey,
				uuid.New().String(),
			),
		},
		{
			"When Audience is wrong",
			createToken(
				jwt.SigningMethodHS256,
				"1h",
				env.JwtIssuer,
				"audience",
				env.JwtSecretKey,
				uuid.New().String(),
			),
		},
		{
			"When jti is wrong",
			createToken(
				jwt.SigningMethodHS256,
				"1h",
				env.JwtIssuer,
				env.JwtAudience,
				env.JwtSecretKey,
				"Random jti",
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			_uut := jwtService{
				env: env,
			}

			// when
			err := _uut.Authorize(tt.token)

			// then
			assert.Error(t, err)
		})
	}
}

func TestJwtService_Authorize_WhenSuccessful_ShouldNotReturnError(t *testing.T) {
	// given
	env := internal.Env{
		JwtSecretKey: "This-is-a-very-long-secret-key-that-is-used-to-encrypt-the-token",
		JwtIssuer:    "JWTIssuer",
		JwtAudience:  "JWTAudience",
	}
	_uut := jwtService{
		env: env,
	}

	token := createToken(
		jwt.SigningMethodHS256,
		"1h",
		env.JwtIssuer,
		env.JwtAudience,
		env.JwtSecretKey,
		uuid.New().String(),
	)

	// when
	err := _uut.Authorize(token)

	// then
	assert.NoError(t, err)
}
