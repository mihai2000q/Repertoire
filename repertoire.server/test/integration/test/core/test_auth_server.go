package core

import (
	"encoding/json"
	"net/http"
	"repertoire/server/data/http/auth"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type testAuthServer struct {
}

func (t *testAuthServer) handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if r.RequestURI == "/storage/token" {
			authResp := auth.TokenResponse{
				Token:     t.createToken(),
				ExpiresIn: "1h",
			}
			res, _ := json.Marshal(authResp)
			_, _ = w.Write(res)
		}
		if r.RequestURI == "/centrifugo/public-token" {
			authResp := auth.TokenResponse{
				Token:     t.createCentrifugoToken(),
				ExpiresIn: "1h",
			}
			res, _ := json.Marshal(authResp)
			_, _ = w.Write(res)
		}
		if r.RequestURI == "/sign-in" {
			token, _ := json.Marshal(t.createToken())
			_, _ = w.Write(token)
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
