package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"repertoire/storage/domain/service"
	"repertoire/storage/internal"
	"testing"
)

// Utils

func getGinAuthHandler() (*gin.Engine, internal.Env, *service.JwtServiceMock) {
	gin.SetMode(gin.TestMode)
	engine := gin.Default()

	env := internal.Env{
		ClientID:          "client_id",
		ClientSecret:      "client_secret",
		JwtExpirationTime: "1h",
	}
	jwtService := new(service.JwtServiceMock)
	authHandler := AuthHandler{
		env:        env,
		jwtService: jwtService,
	}

	engine.POST("/token", authHandler.Token)

	return engine, env, jwtService
}

// Tests

func TestAuthHandler_Token_WhenRequestBodyValuesAreInvalid_ShouldReturnUnauthorizedError(t *testing.T) {
	engine, env, _ := getGinAuthHandler()

	tests := []struct {
		name         string
		grantType    string
		clientID     string
		clientSecret string
	}{
		{
			name:         "invalid grant type",
			grantType:    "invalid",
			clientID:     env.ClientID,
			clientSecret: env.ClientSecret,
		},
		{
			name:         "invalid client id",
			grantType:    "client_credentials",
			clientID:     "invalid_client_id",
			clientSecret: env.ClientSecret,
		},
		{
			name:         "invalid client secret",
			grantType:    "client_credentials",
			clientID:     env.ClientID,
			clientSecret: "invalid_client_secret",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			var requestBody bytes.Buffer
			multiWriter := multipart.NewWriter(&requestBody)
			_ = multiWriter.WriteField("grant_type", tt.grantType)
			_ = multiWriter.WriteField("client_id", tt.clientID)
			_ = multiWriter.WriteField("client_secret", tt.clientSecret)
			_ = multiWriter.Close()

			// when
			req := httptest.NewRequest(http.MethodPost, "/token", &requestBody)
			req.Header.Set("Content-Type", multiWriter.FormDataContentType())

			w := httptest.NewRecorder()
			engine.ServeHTTP(w, req)

			// then
			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})
	}
}

func TestAuthHandler_Token_WhenCreateTokenFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	engine, env, jwtService := getGinAuthHandler()

	var requestBody bytes.Buffer
	multiWriter := multipart.NewWriter(&requestBody)
	_ = multiWriter.WriteField("grant_type", "client_credentials")
	_ = multiWriter.WriteField("client_id", env.ClientID)
	_ = multiWriter.WriteField("client_secret", env.ClientSecret)
	_ = multiWriter.Close()

	jwtService.On("CreateToken").Return("", errors.New("internal error")).Once()

	// when
	req := httptest.NewRequest(http.MethodPost, "/token", &requestBody)
	req.Header.Set("Content-Type", multiWriter.FormDataContentType())

	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	// then
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	jwtService.AssertExpectations(t)
}

func TestAuthHandler_Token_WhenSuccessful_ShouldReturnTokenResult(t *testing.T) {
	// given
	engine, env, jwtService := getGinAuthHandler()

	var requestBody bytes.Buffer
	multiWriter := multipart.NewWriter(&requestBody)
	_ = multiWriter.WriteField("grant_type", "client_credentials")
	_ = multiWriter.WriteField("client_id", env.ClientID)
	_ = multiWriter.WriteField("client_secret", env.ClientSecret)
	_ = multiWriter.Close()

	token := "some token"
	jwtService.On("CreateToken").Return(token, nil).Once()

	// when
	req := httptest.NewRequest(http.MethodPost, "/token", &requestBody)
	req.Header.Set("Content-Type", multiWriter.FormDataContentType())

	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		AccessToken string
		TokenType   string
		ExpiresIn   string
	}
	_ = json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, token, response.AccessToken)
	assert.Equal(t, "Bearer", response.TokenType)
	assert.Equal(t, env.JwtExpirationTime, response.ExpiresIn)

	jwtService.AssertExpectations(t)
}
