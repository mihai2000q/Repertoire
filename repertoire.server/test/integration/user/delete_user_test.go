package auth

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"repertoire/server/test/integration/test/core"
	userData "repertoire/server/test/integration/test/data/user"
	"repertoire/server/test/integration/test/utils"
	"testing"
)

func TestDeleteUser_WhenSuccessful_ShouldReturnUser(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, userData.SeedData)

	user := userData.Users[1]

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().
		WithUser(user).
		DELETE(w, "/api/users")

	// then
	assert.Equal(t, http.StatusOK, w.Code)
}
