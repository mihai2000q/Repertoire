package user

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/assertion"
	"repertoire/server/test/integration/test/core"
	userData "repertoire/server/test/integration/test/data/user"
	"repertoire/server/test/integration/test/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCurrentUser_WhenUserIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().
		WithInvalidToken().
		GET(w, "/api/users/current")

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetCurrentUser_WhenSuccessful_ShouldReturnCurrentUser(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, userData.Users, userData.SeedData)

	user := userData.Users[0]

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().
		WithUser(user).
		GET(w, "/api/users/current")

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var responseUser model.User
	_ = json.Unmarshal(w.Body.Bytes(), &responseUser)

	db := utils.GetDatabase(t)
	db.Find(&user, user.ID)

	assertion.ResponseUser(t, user, responseUser)
}
