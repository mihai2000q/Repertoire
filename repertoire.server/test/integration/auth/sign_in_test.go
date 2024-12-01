package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"repertoire/server/api/requests"
	"repertoire/server/test/integration/test/assertion"
	"repertoire/server/test/integration/test/core"
	"repertoire/server/test/integration/test/data/auth"
	"repertoire/server/test/integration/test/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignIn_WhenUserIsNotFound_ShouldReturnUnauthorizedError(t *testing.T) {
	// given
	request := requests.SignInRequest{
		Email:    "random@email.com",
		Password: "Password123",
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().
		WithoutAuthentication().
		PUT(w, "/api/auth/sign-in", request)

	// then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestSignIn_WhenPasswordIsWrong_ShouldReturnUnauthorizedError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, auth.Users, auth.SeedData)

	request := requests.SignInRequest{
		Email:    "random@email.com",
		Password: "WrongPassword123",
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().
		WithoutAuthentication().
		PUT(w, "/api/auth/sign-in", request)

	// then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestSignIn_WhenSuccessful_ShouldReturnValidToken(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, auth.Users, auth.SeedData)

	user := auth.Users[0]
	request := requests.SignInRequest{
		Email:    user.Email,
		Password: "Password123",
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().
		WithoutAuthentication().
		PUT(w, "/api/auth/sign-in", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var response struct{ Token string }
	_ = json.Unmarshal(w.Body.Bytes(), &response)

	assertion.Token(t, response.Token, user)
}
