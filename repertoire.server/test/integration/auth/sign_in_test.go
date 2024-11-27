package auth

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"repertoire/server/api/requests"
	"repertoire/server/test/integration/test/core"
	"repertoire/server/test/integration/test/data/auth"
	"repertoire/server/test/integration/test/utils"
	"testing"
)

func TestSignIn_WhenUserIsNotFound_ShouldReturnInvalidCredentials(t *testing.T) {
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
	//assert.JSONEq(t, `{"error": ""}`, w.Body.String())
}

func TestSignIn_WhenValid_ShouldReturnToken(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, auth.SeedData)

	user := auth.Users[0]
	request := requests.SignInRequest{
		Email:    user.Email,
		Password: "Password123",
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().PUT(w, "/api/auth/sign-in", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)
	// assert.JSONEq(t, `{"token": ""}`, w.Body.String())
}
