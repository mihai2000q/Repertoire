package auth

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/core"
	userData "repertoire/server/test/integration/test/data/user"
	"repertoire/server/test/integration/test/utils"
	"testing"
)

func TestGetUser_WhenUserIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, userData.Users, userData.SeedData)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().GET(w, "/api/users/"+uuid.New().String())

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetUser_WhenSuccessful_ShouldReturnUser(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, userData.Users, userData.SeedData)

	user := userData.Users[1]

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().GET(w, "/api/users/"+user.ID.String())

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var responseUser model.User
	_ = json.Unmarshal(w.Body.Bytes(), &responseUser)

	assert.Equal(t, user.ID, responseUser.ID)
	assert.Equal(t, user.Email, responseUser.Email)
	if user.ProfilePictureURL == nil {
		assert.Nil(t, responseUser.ProfilePictureURL)
	} else {
		assert.Equal(t, user.ProfilePictureURL.ToNullableFullURL(), responseUser.ProfilePictureURL)
	}
}
