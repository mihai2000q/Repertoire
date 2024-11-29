package core

import (
	"bytes"
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"repertoire/server/internal"
	"repertoire/server/model"
	"time"
)

type TestHandler interface {
	WithoutAuthentication() TestHandler
	WithInvalidToken() TestHandler
	WithUser(user model.User) TestHandler
	GET(w http.ResponseWriter, url string)
	POST(w http.ResponseWriter, url string, body interface{})
	PUT(w http.ResponseWriter, url string, body interface{})
	DELETE(w http.ResponseWriter, url string)
}

type settings struct {
	authentication bool
	invalidToken   bool
	user           *model.User
}

type testHandler struct {
	httpServer *http.Server
	settings   *settings
}

func NewTestHandler() TestHandler {
	return &testHandler{
		getHttpServer(),
		&settings{
			authentication: true,
		},
	}
}

func (t *testHandler) WithoutAuthentication() TestHandler {
	t.settings.authentication = false
	return t
}

func (t *testHandler) WithInvalidToken() TestHandler {
	t.settings.invalidToken = true
	return t
}

func (t *testHandler) WithUser(user model.User) TestHandler {
	t.settings.user = &user
	return t
}

func (t *testHandler) GET(w http.ResponseWriter, url string) {
	req, _ := http.NewRequest("GET", url, nil)
	t.requestWithAuthentication(req)
	t.httpServer.Handler.ServeHTTP(w, req)
}

func (t *testHandler) POST(w http.ResponseWriter, url string, body interface{}) {
	jsonBody, _ := json.Marshal(body)
	reqBody := bytes.NewBuffer(jsonBody)
	req, _ := http.NewRequest("POST", url, reqBody)

	req.Header.Set("Content-Type", "application/json")
	t.requestWithAuthentication(req)
	t.httpServer.Handler.ServeHTTP(w, req)
}

func (t *testHandler) PUT(w http.ResponseWriter, url string, body interface{}) {
	jsonBody, _ := json.Marshal(body)
	reqBody := bytes.NewBuffer(jsonBody)
	req, _ := http.NewRequest("PUT", url, reqBody)

	req.Header.Set("Content-Type", "application/json")
	t.requestWithAuthentication(req)
	t.httpServer.Handler.ServeHTTP(w, req)
}

func (t *testHandler) DELETE(w http.ResponseWriter, url string) {
	req, _ := http.NewRequest("DELETE", url, nil)
	t.requestWithAuthentication(req)
	t.httpServer.Handler.ServeHTTP(w, req)
}

func (t *testHandler) requestWithAuthentication(req *http.Request) {
	if !t.settings.authentication {
		return
	}

	if t.settings.invalidToken {
		req.Header.Set("Authorization", "bearer "+t.createInvalidToken())
		return
	}

	var user model.User
	if t.settings.user == nil {
		db, _ := gorm.Open(postgres.Open(Dsn))
		db.First(&user)
	} else {
		user = *t.settings.user
	}

	token := t.createToken(user)
	req.Header.Set("Authorization", "bearer "+token)
}

func (t *testHandler) createInvalidToken() string {
	env := internal.NewEnv()

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"jti": uuid.New().String(),
		"sub": uuid.New().String(),
		"iss": env.JwtIssuer,
		"aud": env.JwtAudience,
		"iat": time.Now().UTC().Unix(),
		"exp": time.Now().UTC().Add(time.Hour).Unix(),
	})
	token, _ := claims.SignedString([]byte(env.JwtSecretKey))
	return token
}

func (t *testHandler) createToken(user model.User) string {
	env := internal.NewEnv()

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"jti": uuid.New().String(),
		"sub": user.ID.String(),
		"iss": env.JwtIssuer,
		"aud": env.JwtAudience,
		"iat": time.Now().UTC().Unix(),
		"exp": time.Now().UTC().Add(time.Hour).Unix(),
	})
	token, _ := claims.SignedString([]byte(env.JwtSecretKey))
	return token
}
