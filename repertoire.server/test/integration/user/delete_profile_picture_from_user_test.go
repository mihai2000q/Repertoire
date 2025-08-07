package user

import (
	"net/http"
	"net/http/httptest"
	"repertoire/server/test/integration/test/core"
	userData "repertoire/server/test/integration/test/data/user"
	"repertoire/server/test/integration/test/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteProfilePictureFromUser_WhenUserIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, userData.Users, userData.SeedData)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().
		WithInvalidToken().
		DELETE(w, "/api/users/pictures")

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteProfilePictureFromUser_WhenUserHasNoProfilePicture_ShouldReturnConflictError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, userData.Users, userData.SeedData)

	user := userData.Users[0]

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().
		WithUser(user).
		DELETE(w, "/api/users/pictures")

	// then
	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestDeleteProfilePictureFromUser_WhenSuccessful_ShouldUpdateUserAndDeleteProfilePicture(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, userData.Users, userData.SeedData)

	user := userData.Users[1]

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().
		WithUser(user).
		DELETE(w, "/api/users/pictures")

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)
	db.Find(&user, user.ID)

	assert.Nil(t, user.ProfilePictureURL)
}
