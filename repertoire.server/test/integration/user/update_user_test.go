package auth

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"repertoire/server/api/requests"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/core"
	userData "repertoire/server/test/integration/test/data/user"
	"repertoire/server/test/integration/test/utils"
	"testing"
)

func TestUpdateUser_WhenUserIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	request := requests.UpdateUserRequest{
		Name: "New Name",
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().
		WithInvalidToken().
		PUT(w, "/api/users", request)

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateUser_WhenSuccessful_ShouldUpdateUser(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, userData.Users, userData.SeedData)

	user := userData.Users[0]
	request := requests.UpdateUserRequest{
		Name: "New Name",
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().
		WithUser(user).
		PUT(w, "/api/users", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)
	assertUpdatedUser(t, request, user.ID)
}

func assertUpdatedUser(t *testing.T, request requests.UpdateUserRequest, userID uuid.UUID) {
	db := utils.GetDatabase()

	var user model.User
	db.Find(&user, userID)

	assert.Equal(t, request.Name, user.Name)
}
