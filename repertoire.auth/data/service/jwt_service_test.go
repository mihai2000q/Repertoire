package service

import (
	"net/http"
	"repertoire/auth/internal"
	"repertoire/auth/model"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var privateKey = `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`

var publicKey = `-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAgUt4M5t0ARBS76u7Hb3l
AUywf9kTiCGNoWUQiILXHye+eq3HJeAJ+zQP7QZphy+IbKCk8OsnB37P7K6SQZQ1
79pk8pztQL48gPxX7F78dnBS05cpc4rqrJ/YyDzV4DR6hwYnQuOiHb9JSjv7/eUM
V5TykX6GXQHDZpBRt/aMoahexAqTgrYtk1o1cydcPY/K9c/TaBljpGdUz+7bNNNM
MWWpntE9QPkHJpEup1Ds6IHUNB7C2Mc8VXiDU8qNbpeeQTDhbkczrGq70QEfLwok
/MRCluZZolrBtbaXa+FuME6dXNCJd1i65nIbWpkMxtHRBeXeDDXt/SNlrAaljJDa
ECAMbyjEVf7rGB2aaFqIvFsje7aGBFy9F9XHtibj5h53oCeYoTUj31WTkhTHIdEx
wqF1nQPISHtf7gMHwGfD9+B1Swh+T4Yu+q+7Cyanc0Z1yFx2gh41FesiIHzXOLQw
xboUKRmfq4aNKpJULCj5qOOaw77NdipAUuiIW5Rxo/nwumGV8MUnUN37LY6g3Mfz
w7swJvlgJrW/AHBO3dhXsJJKyxSnvpd4Bpun/o5tj73YJDU5iddekR1pcBhhUPAg
vFFWwmdW7QpnhSTSjcBkVPj9oaXliIT3gfrqHBEgUyKAD5YYOjHDpVw/egs3XTMv
ciITn/KOz9ExLS671KFwLV8CAwEAAQ==
-----END PUBLIC KEY-----`

var otherPrivateKey = `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`

var otherPublicKey = `-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAkSBgbK0TxcPWnk4HYRdV
6zRV1RuoZwlI5J0RyH0cpYGJoJrOGubdFuwxgpGQdQ//z4U0oTNh8sEins/IRX++
BJSTIA9rzAQYQuAIMtK3kv8XAFZyUH4utr926U5ERMgaRC9EwoSAyBJjnurnky4L
40hvXi2ow0o55mGq2Bf89vCqKJwRYgG+PK/1+4Rkv5XZeFt5yyh0e+JhnoTTmVkQ
CmMMbNSfmPSGA8SsHnUApQYbD2oHubO/61hoT+PfAa7fcOqplXUdu4AAxIIVHaVf
2qjjZuTYisLFJdryCNh5vQ/VRF08oomzfj8Y4VYuVboX/Ij0zhsPae2Suqia0+kI
ijLBYEBE2l9Q6YXYVVpvhCORqU2nYiP8FTYnVdKH7mft8w11iKKsHm3oZOiYtpqz
/Q3ssPYeoS8LHI3pJLC98fW2N74FHJrn6IILK7tbRH9bhz0MdS9is46ZePx9oBCq
Y/LWTtPDVrdEYnYqWAQyQoaTWHvT3SlY3dSO4TjbExZ3ILpizw0NEKz9UX9bwvWP
XCucxHVaxbF4y6YoSR1efMeNJKevXutFkHABHHBWNcm6nmQxso2x9hQjw1iz6y7m
Ezfo7f3lOVGJYT1kmdCN8cCAHiY25fQAr5GGvx7nv721y3LsRC6ti02YNxEtbUNt
DcfoMwR1oEps66XlYRXd/z8CAwEAAQ==
-----END PUBLIC KEY-----`

func TestJwtService_Authorize_WhenPublicKeyIsNotAKey_ShouldReturnInternalServerError(t *testing.T) {
	// given
	env := internal.Env{JwtPublicKey: "not a key"}
	_uut := NewJwtService(env)

	claims := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{})
	key, _ := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	token, _ := claims.SignedString(key)

	// when
	errCode := _uut.Authorize(token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Error(t, errCode.Error)
}

func TestJwtService_Authorize_WhenTokenIsInvalid_ShouldReturnUnauthorizedError(t *testing.T) {
	env := internal.Env{
		JwtPrivateKey: privateKey,
		JwtPublicKey:  publicKey,
		JwtIssuer:     "Issuer",
		JwtAudience:   "Audience",
	}

	tests := []struct {
		name       string
		claims     *jwt.Token
		publicKey  string
		privateKey string
	}{
		{
			"when expiration time has elapsed",
			jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
				"exp": time.Now().UTC().Add(-time.Second).Unix(),
			}),
			env.JwtPublicKey,
			env.JwtPrivateKey,
		},
		{
			"when public key does not match",
			jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{}),
			otherPublicKey,
			env.JwtPrivateKey,
		},
		{
			"when private key does not match",
			jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{}),
			env.JwtPublicKey,
			otherPrivateKey,
		},
		{
			"when audience is missing",
			jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{}),
			env.JwtPublicKey,
			env.JwtPrivateKey,
		},
		{
			"when issuer is missing",
			jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
				"aud": env.JwtAudience,
			}),
			env.JwtPublicKey,
			env.JwtPrivateKey,
		},
		{
			"when expiration time is missing",
			jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
				"aud": env.JwtAudience,
				"iss": env.JwtIssuer,
			}),
			env.JwtPublicKey,
			env.JwtPrivateKey,
		},
		{
			"when sub is missing",
			jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
				"jti": uuid.New().String(),
				"aud": env.JwtAudience,
				"iss": env.JwtIssuer,
				"exp": time.Now().UTC().Add(time.Minute).Unix(),
			}),
			env.JwtPublicKey,
			env.JwtPrivateKey,
		},
		{
			"when sub is not uuid",
			jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
				"jti": uuid.New().String(),
				"sub": "this is not a uuid",
				"aud": env.JwtAudience,
				"iss": env.JwtIssuer,
				"exp": time.Now().UTC().Add(time.Minute).Unix(),
			}),
			env.JwtPublicKey,
			env.JwtPrivateKey,
		},
		{
			"when jti is missing",
			jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
				"aud": env.JwtAudience,
				"iss": env.JwtIssuer,
			}),
			env.JwtPublicKey,
			env.JwtPrivateKey,
		},
		{
			"when jti is not uuid",
			jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
				"jti": "this is not a uuid",
				"aud": env.JwtAudience,
				"iss": env.JwtIssuer,
				"exp": time.Now().UTC().Add(time.Minute).Unix(),
			}),
			env.JwtPublicKey,
			env.JwtPrivateKey,
		},
		{
			"when signing method is not the same",
			jwt.NewWithClaims(jwt.SigningMethodRS384, jwt.MapClaims{
				"jti": uuid.New().String(),
				"sub": uuid.New().String(),
				"aud": env.JwtAudience,
				"iss": env.JwtIssuer,
				"exp": time.Now().UTC().Add(time.Minute).Unix(),
			}),
			env.JwtPublicKey,
			env.JwtPrivateKey,
		},
		{
			"when audience is wrong",
			jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
				"jti": uuid.New().String(),
				"sub": uuid.New().String(),
				"aud": "something else",
				"iss": env.JwtIssuer,
				"exp": time.Now().UTC().Add(time.Minute).Unix(),
			}),
			env.JwtPublicKey,
			env.JwtPrivateKey,
		},
		{
			"when audience has too elements",
			jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
				"jti": uuid.New().String(),
				"sub": uuid.New().String(),
				"aud": []string{env.JwtAudience, "something else"},
				"iss": env.JwtIssuer,
				"exp": time.Now().UTC().Add(time.Minute).Unix(),
			}),
			env.JwtPublicKey,
			env.JwtPrivateKey,
		},
		{
			"when issuer is wrong",
			jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
				"jti": uuid.New().String(),
				"sub": uuid.New().String(),
				"aud": env.JwtAudience,
				"iss": "something else",
				"exp": time.Now().UTC().Add(time.Minute).Unix(),
			}),
			env.JwtPublicKey,
			env.JwtPrivateKey,
		},
		{
			"when jti is uuid nil",
			jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
				"jti": uuid.Nil.String(),
				"sub": uuid.New().String(),
				"aud": env.JwtAudience,
				"iss": env.JwtIssuer,
				"exp": time.Now().UTC().Add(time.Minute).Unix(),
			}),
			env.JwtPublicKey,
			env.JwtPrivateKey,
		},
		{
			"when sub is uuid nil",
			jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
				"jti": uuid.New().String(),
				"sub": uuid.Nil.String(),
				"aud": env.JwtAudience,
				"iss": env.JwtIssuer,
				"exp": time.Now().UTC().Add(time.Minute).Unix(),
			}),
			env.JwtPublicKey,
			env.JwtPrivateKey,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			env.JwtPublicKey = tt.publicKey
			_uut := NewJwtService(env)

			key, _ := jwt.ParseRSAPrivateKeyFromPEM([]byte(tt.privateKey))
			token, _ := tt.claims.SignedString(key)

			// when
			errCode := _uut.Authorize(token)

			// then
			assert.NotNil(t, errCode)
			assert.Equal(t, http.StatusUnauthorized, errCode.Code)
			assert.Error(t, errCode.Error)
		})
	}
}

func TestJwtService_Authorize_WhenTokenIsValid_ShouldNotReturnAnyError(t *testing.T) {
	// given
	env := internal.Env{
		JwtPrivateKey: privateKey,
		JwtPublicKey:  publicKey,
		JwtIssuer:     "Issuer",
		JwtAudience:   "Audience",
	}
	_uut := NewJwtService(env)

	claims := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"aud": env.JwtAudience,
		"iss": env.JwtIssuer,
		"jti": uuid.New().String(),
		"sub": uuid.New().String(),
		"iat": time.Now().UTC().Unix(),
		"exp": time.Now().UTC().Add(time.Hour).Unix(),
	})
	key, _ := jwt.ParseRSAPrivateKeyFromPEM([]byte(env.JwtPrivateKey))
	token, _ := claims.SignedString(key)

	// when
	errCode := _uut.Authorize(token)

	// then
	assert.Nil(t, errCode)
}

func TestJwtService_GetUserIDFromJwt_WhenPublicKeyIsNotAKey_ShouldReturnInternalServerError(t *testing.T) {
	// given
	env := internal.Env{
		JwtPublicKey:  "this is not a public key",
		JwtPrivateKey: privateKey,
	}
	_uut := NewJwtService(env)

	claims := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{})
	key, _ := jwt.ParseRSAPrivateKeyFromPEM([]byte(env.JwtPrivateKey))
	token, _ := claims.SignedString(key)

	// when
	userID, errCode := _uut.GetUserIDFromJwt(token)

	// then
	assert.Empty(t, userID)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Error(t, errCode.Error)
}

func TestJwtService_GetUserIDFromJwt_WhenPublicKeyIsNotMatchingWithPrivateKey_ShouldReturnForbiddenError(t *testing.T) {
	tests := []struct {
		name       string
		publicKey  string
		privateKey string
	}{
		{
			"When private key is not matching",
			publicKey,
			otherPrivateKey,
		},
		{
			"When public key is not matching",
			otherPublicKey,
			privateKey,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			env := internal.Env{
				JwtPublicKey: tt.publicKey,
			}
			_uut := NewJwtService(env)

			claims := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{})
			key, _ := jwt.ParseRSAPrivateKeyFromPEM([]byte(tt.privateKey))
			token, _ := claims.SignedString(key)

			// when
			userID, errCode := _uut.GetUserIDFromJwt(token)

			// then
			assert.Empty(t, userID)
			assert.NotNil(t, errCode)
			assert.Equal(t, http.StatusForbidden, errCode.Code)
			assert.Error(t, errCode.Error)
		})
	}
}

func TestJwtService_GetUserIDFromJwt_WhenSubIsMissing_ShouldReturnForbiddenError(t *testing.T) {
	// given
	env := internal.Env{
		JwtPrivateKey: privateKey,
		JwtPublicKey:  publicKey,
	}
	_uut := NewJwtService(env)

	claims := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{})
	key, _ := jwt.ParseRSAPrivateKeyFromPEM([]byte(env.JwtPrivateKey))
	token, _ := claims.SignedString(key)

	// when
	userID, errCode := _uut.GetUserIDFromJwt(token)

	// then
	assert.Empty(t, userID)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusForbidden, errCode.Code)
	assert.Error(t, errCode.Error)
}

func TestJwtService_GetUserIDFromJwt_WhenSubIsNotUUID_ShouldReturnForbiddenError(t *testing.T) {
	// given
	env := internal.Env{
		JwtPrivateKey: privateKey,
		JwtPublicKey:  publicKey,
	}
	_uut := NewJwtService(env)

	claims := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub": "something-else",
	})
	key, _ := jwt.ParseRSAPrivateKeyFromPEM([]byte(env.JwtPrivateKey))
	token, _ := claims.SignedString(key)

	// when
	userID, errCode := _uut.GetUserIDFromJwt(token)

	// then
	assert.Empty(t, userID)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusForbidden, errCode.Code)
	assert.Error(t, errCode.Error)
}

func TestJwtService_GetUserIDFromJwt_WhenSuccessful_ShouldReturnUserId(t *testing.T) {
	// given
	env := internal.Env{
		JwtPrivateKey: privateKey,
		JwtPublicKey:  publicKey,
	}
	_uut := NewJwtService(env)

	user := model.User{
		ID: uuid.New(),
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub": user.ID.String(),
	})
	key, _ := jwt.ParseRSAPrivateKeyFromPEM([]byte(env.JwtPrivateKey))
	token, _ := claims.SignedString(key)

	// when
	userID, errCode := _uut.GetUserIDFromJwt(token)

	// then
	assert.Nil(t, errCode)
	assert.Equal(t, userID, user.ID)
}

// Validations

func TestJwtService_Validate_WhenPublicKeyIsNotAKey_ShouldReturnUnauthorizedError(t *testing.T) {
	// given
	env := internal.Env{JwtPublicKey: "not a key"}
	_uut := NewJwtService(env)

	claims := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{})
	key, _ := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	token, _ := claims.SignedString(key)

	// when
	userID, errCode := _uut.Validate(token)

	// then
	assert.Empty(t, userID)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Error(t, errCode.Error)
}

func TestJwtService_Validate_WhenTokenIsInvalid_ShouldReturnUnauthorizedError(t *testing.T) {
	env := internal.Env{
		JwtPublicKey:  publicKey,
		JwtPrivateKey: privateKey,
		JwtAudience:   "Repertoire",
		JwtIssuer:     "Repertoire",
	}

	tests := []struct {
		name       string
		claims     *jwt.Token
		publicKey  string
		privateKey string
	}{
		{
			"when public key does not match",
			jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
				"jti": uuid.New().String(),
				"sub": uuid.New().String(),
				"iss": env.JwtIssuer,
				"aud": env.JwtAudience,
			}),
			otherPublicKey,
			env.JwtPrivateKey,
		},
		{
			"when private key does not match",
			jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
				"jti": uuid.New().String(),
				"sub": uuid.New().String(),
				"iss": env.JwtIssuer,
				"aud": env.JwtAudience,
			}),
			env.JwtPublicKey,
			otherPrivateKey,
		},
		{
			"when audience is missing",
			jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{}),
			env.JwtPublicKey,
			env.JwtPrivateKey,
		},
		{
			"when issuer is missing",
			jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
				"aud": env.JwtAudience,
			}),
			env.JwtPublicKey,
			env.JwtPrivateKey,
		},
		{
			"when jti is missing",
			jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
				"iss": env.JwtIssuer,
				"aud": env.JwtAudience,
			}),
			env.JwtPublicKey,
			env.JwtPrivateKey,
		},
		{
			"when Signing Method is not the same",
			jwt.NewWithClaims(jwt.SigningMethodRS384, jwt.MapClaims{
				"jti": uuid.New().String(),
				"iss": env.JwtIssuer,
				"aud": env.JwtAudience,
			}),
			env.JwtPublicKey,
			env.JwtPrivateKey,
		},
		{
			"when issuer is not matching",
			jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
				"jti": uuid.New().String(),
				"iss": "some issuer",
				"aud": env.JwtAudience,
			}),
			env.JwtPublicKey,
			env.JwtPrivateKey,
		},
		{
			"when audience is not matching",
			jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
				"jti": uuid.New().String(),
				"iss": env.JwtIssuer,
				"aud": "some audience",
			}),
			env.JwtPublicKey,
			env.JwtPrivateKey,
		},
		{
			"when audience has too many elements",
			jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
				"jti": uuid.New().String(),
				"iss": env.JwtIssuer,
				"aud": []string{env.JwtAudience, "some other audience"},
			}),
			env.JwtPublicKey,
			env.JwtPrivateKey,
		},
		{
			"when jti is not uuid",
			jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
				"jti": "something else",
				"iss": env.JwtIssuer,
				"aud": env.JwtAudience,
			}),
			env.JwtPublicKey,
			env.JwtPrivateKey,
		},
		{
			"when jti is uuid nil",
			jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
				"jti": uuid.Nil,
				"iss": env.JwtIssuer,
				"aud": env.JwtAudience,
			}),
			env.JwtPublicKey,
			env.JwtPrivateKey,
		},
		{
			"when sub is missing",
			jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
				"jti": uuid.New().String(),
				"iss": env.JwtIssuer,
				"aud": env.JwtAudience,
			}),
			env.JwtPublicKey,
			env.JwtPrivateKey,
		},
		{
			"when sub is not uuid",
			jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
				"jti": uuid.New().String(),
				"sub": "This is a sub",
				"iss": env.JwtIssuer,
				"aud": env.JwtAudience,
			}),
			env.JwtPublicKey,
			env.JwtPrivateKey,
		},
		{
			"when sub is uuid nil",
			jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
				"jti": uuid.New().String(),
				"sub": uuid.Nil.String(),
				"iss": env.JwtIssuer,
				"aud": env.JwtAudience,
			}),
			env.JwtPublicKey,
			env.JwtPrivateKey,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			env.JwtPublicKey = tt.publicKey
			env.JwtPrivateKey = tt.privateKey
			_uut := NewJwtService(env)

			key, _ := jwt.ParseRSAPrivateKeyFromPEM([]byte(tt.privateKey))
			token, _ := tt.claims.SignedString(key)

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
		JwtPublicKey: publicKey,
	}
	_uut := NewJwtService(env)

	user := model.User{
		ID: uuid.New(),
	}

	// expired token
	claims := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"jti": uuid.New().String(),
		"sub": user.ID.String(),
		"iss": env.JwtIssuer,
		"aud": env.JwtAudience,
		"iat": time.Now().UTC().Unix(),
		"exp": time.Now().UTC().Add(-time.Hour).Unix(),
	})
	key, _ := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	token, _ := claims.SignedString(key)

	// when
	userID, errCode := _uut.Validate(token)

	// then
	assert.Equal(t, userID, user.ID)
	assert.Nil(t, errCode)
}

func TestJwtService_ValidateCredentials_WhenInvalid_ShouldReturnError(t *testing.T) {
	env := internal.Env{
		ClientID:     "Some-id",
		ClientSecret: "Some-secret",
	}

	tests := []struct {
		name              string
		clientCredentials model.ClientCredentials
	}{
		{
			"when grant type is invalid",
			model.ClientCredentials{
				GrantType: "invalid",
			},
		},
		{
			"when client id does not match",
			model.ClientCredentials{
				GrantType: "client_credentials",
				ClientID:  "invalid",
			},
		},
		{
			"when client secret does not match",
			model.ClientCredentials{
				GrantType:    "client_credentials",
				ClientID:     env.ClientID,
				ClientSecret: "invalid",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			_uut := NewJwtService(env)

			// when
			errCode := _uut.ValidateCredentials(tt.clientCredentials)

			// then
			assert.NotNil(t, errCode)
			assert.Equal(t, http.StatusUnauthorized, errCode.Code)
			assert.Error(t, errCode.Error)
		})
	}
}

func TestJwtService_ValidateCredentials_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	// given
	env := internal.Env{
		ClientID:     "Some-id",
		ClientSecret: "Some-secret",
	}
	_uut := NewJwtService(env)

	credentials := model.ClientCredentials{
		GrantType:    "client_credentials",
		ClientID:     env.ClientID,
		ClientSecret: env.ClientSecret,
	}

	// when
	errCode := _uut.ValidateCredentials(credentials)

	// then
	assert.Nil(t, errCode)
}

// Create Tokens

func TestJwtService_CreateToken_WhenExpirationTimeIsInvalid_ShouldReturnInternalError(t *testing.T) {
	// given
	env := internal.Env{
		JwtExpirationTime: "something",
	}
	_uut := NewJwtService(env)

	user := model.User{}

	// when
	tokenString, errCode := _uut.CreateToken(user)

	// then
	assert.Empty(t, tokenString)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Error(t, errCode.Error)
}

func TestJwtService_CreateToken_WhenPrivateKeyIsInvalid_ShouldReturnInternalError(t *testing.T) {
	// given
	env := internal.Env{
		JwtPrivateKey:     "Invalid key",
		JwtExpirationTime: "1h",
	}
	_uut := NewJwtService(env)

	user := model.User{}

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
		JwtPublicKey:      publicKey,
		JwtPrivateKey:     privateKey,
		JwtIssuer:         "Repertoire",
		JwtAudience:       "Repertoire",
		JwtExpirationTime: "1h",
	}
	_uut := NewJwtService(env)

	user := model.User{
		ID: uuid.New(),
	}

	expiresIn, _ := time.ParseDuration(env.JwtExpirationTime)

	// when
	tokenString, errCode := _uut.CreateToken(user)

	// then
	assert.Nil(t, errCode)

	key, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return key, nil
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

	assert.Equal(t, jwt.SigningMethodRS256, token.Method)
	assert.NotEmpty(t, jti)
	assert.Equal(t, user.ID.String(), sub)
	assert.Len(t, aud, 1)
	assert.Equal(t, env.JwtAudience, aud[0])
	assert.Equal(t, env.JwtIssuer, iss)
	assert.WithinDuration(t, time.Now().UTC(), iat.Time, 10*time.Second)
	assert.WithinDuration(t, time.Now().Add(expiresIn).UTC(), exp.Time, 10*time.Second)
}

func TestJwtService_CreateCentrifugoToken_WhenExpirationTimeIsInvalid_ShouldReturnInternalError(t *testing.T) {
	// given
	env := internal.Env{
		CentrifugoJwtExpirationTime: "something",
	}
	_uut := NewJwtService(env)

	userID := uuid.New()

	// when
	tokenString, expiresIn, errCode := _uut.CreateCentrifugoToken(userID)

	// then
	assert.Empty(t, tokenString)
	assert.Empty(t, expiresIn)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Error(t, errCode.Error)
}

func TestJwtService_CreateCentrifugoToken_WhenSuccessful_ShouldReturnNewToken(t *testing.T) {
	// given
	env := internal.Env{
		JwtIssuer:                   "Repertoire",
		CentrifugoJwtSecretKey:      "Some-secret-duper-key-that-should-not-be-known-to-everyone",
		CentrifugoJwtAudience:       "Repertoire",
		CentrifugoJwtExpirationTime: "1h",
	}
	_uut := NewJwtService(env)

	userID := uuid.New()

	expiresInDuration, _ := time.ParseDuration(env.CentrifugoJwtExpirationTime)

	// when
	tokenString, expiresIn, errCode := _uut.CreateCentrifugoToken(userID)

	// then
	assert.Nil(t, errCode)
	assert.Equal(t, env.CentrifugoJwtExpirationTime, expiresIn)

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(env.CentrifugoJwtSecretKey), nil
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
	assert.Equal(t, userID.String(), sub)
	assert.Len(t, aud, 1)
	assert.Equal(t, env.CentrifugoJwtAudience, aud[0])
	assert.Equal(t, env.JwtIssuer, iss)
	assert.WithinDuration(t, time.Now().UTC(), iat.Time, 10*time.Second)
	assert.WithinDuration(t, time.Now().Add(expiresInDuration).UTC(), exp.Time.UTC(), 10*time.Second)
}

func TestJwtService_CreateStorageToken_WhenExpirationTimeIsInvalid_ShouldReturnInternalError(t *testing.T) {
	// given
	env := internal.Env{
		StorageJwtExpirationTime: "something",
	}
	_uut := NewJwtService(env)

	userID := uuid.New()

	// when
	tokenString, expiresIn, errCode := _uut.CreateStorageToken(userID)

	// then
	assert.Empty(t, tokenString)
	assert.Empty(t, expiresIn)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Error(t, errCode.Error)
}

func TestJwtService_CreateStorageToken_WhenSuccessful_ShouldReturnNewToken(t *testing.T) {
	// given
	env := internal.Env{
		JwtIssuer:                "Repertoire",
		StorageJwtSecretKey:      "Some-secret-duper-key-that-should-not-be-known-to-everyone",
		StorageJwtAudience:       "Repertoire",
		StorageJwtExpirationTime: "1h",
	}
	_uut := NewJwtService(env)

	userID := uuid.New()

	expiresInDuration, _ := time.ParseDuration(env.StorageJwtExpirationTime)

	// when
	tokenString, expiresIn, errCode := _uut.CreateStorageToken(userID)

	// then
	assert.Nil(t, errCode)
	assert.Equal(t, env.StorageJwtExpirationTime, expiresIn)

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(env.StorageJwtSecretKey), nil
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
	assert.Equal(t, userID.String(), sub)
	assert.Len(t, aud, 1)
	assert.Equal(t, env.StorageJwtAudience, aud[0])
	assert.Equal(t, env.JwtIssuer, iss)
	assert.WithinDuration(t, time.Now().UTC(), iat.Time, 10*time.Second)
	assert.WithinDuration(t, time.Now().Add(expiresInDuration).UTC(), exp.Time.UTC(), 10*time.Second)
}
