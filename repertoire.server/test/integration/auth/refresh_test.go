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

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRefresh_WhenTokenIsInvalid_ShouldReturnUnauthorizedError(t *testing.T) {
	tests := []struct {
		name    string
		request requests.RefreshRequest
	}{
		{
			"when token is empty",
			requests.RefreshRequest{
				Token: "",
			},
		},
		{
			"when jti is invalid",
			requests.RefreshRequest{
				Token: utils.CreateCustomToken(uuid.New().String(), "something"),
			},
		},
		{
			"when sub is invalid",
			requests.RefreshRequest{
				Token: utils.CreateCustomToken("something", uuid.New().String()),
			},
		},
		{
			"when sub does not represent a user",
			requests.RefreshRequest{
				Token: utils.CreateCustomToken(uuid.New().String(), uuid.New().String()),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// given
			utils.SeedAndCleanupData(t, auth.Users, auth.SeedData)

			// when
			w := httptest.NewRecorder()
			core.NewTestHandler().
				WithoutAuthentication().
				PUT(w, "/api/auth/refresh", test.request)

			// then
			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})
	}
}

func TestRefresh_WhenSuccessful_ShouldReturnValidToken(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, auth.Users, auth.SeedData)

	user := auth.Users[0]
	request := requests.RefreshRequest{
		Token: utils.CreateValidToken(user),
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
