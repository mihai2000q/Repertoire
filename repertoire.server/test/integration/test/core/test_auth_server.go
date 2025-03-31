package core

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type testAuthServer struct {
}

func (t *testAuthServer) handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		if r.RequestURI == "/storage/token" {
			_, _ = w.Write([]byte(t.createToken()))
		}
		if r.RequestURI == "/centrifugo/token" {
			_, _ = w.Write([]byte(t.createCentrifugoToken()))
		}
		if r.RequestURI == "/sign-in" {
			_, _ = w.Write([]byte(t.createToken()))
		}
	}
}

func (t *testAuthServer) createToken() string {
	claims := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"jti": uuid.New().String(),
		"sub": "Integration Testing",
		"iss": jwtInfo.Issuer,
		"aud": jwtInfo.Audience,
		"iat": time.Now().UTC().Unix(),
		"exp": time.Now().UTC().Add(time.Hour).Unix(),
	})
	privateKey, _ := jwt.ParseRSAPrivateKeyFromPEM([]byte(jwtInfo.PrivateKey))
	token, _ := claims.SignedString(privateKey)
	return token
}

func (t *testAuthServer) createCentrifugoToken() string {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"jti": uuid.New().String(),
		"sub": "Integration Testing",
		"iss": CentrifugoJwtInfo.Issuer,
		"aud": CentrifugoJwtInfo.Audience,
		"iat": time.Now().UTC().Unix(),
		"exp": time.Now().UTC().Add(time.Hour).Unix(),
	})
	token, _ := claims.SignedString([]byte(CentrifugoJwtInfo.SecretKey))
	return token
}
