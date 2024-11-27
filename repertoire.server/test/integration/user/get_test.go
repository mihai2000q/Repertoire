package auth

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/core"
	userData "repertoire/server/test/integration/test/data/user"
	"repertoire/server/test/integration/test/utils"
	"testing"
)

func TestGetUser_WhenSuccessful_ShouldReturnUser(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, userData.SeedData)

	user := userData.Users[1]

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().GET(w, "/api/users/"+user.ID.String())

	// then
	var returnedUser model.User
	_ = json.Unmarshal(w.Body.Bytes(), &returnedUser)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, user.ID, returnedUser.ID)
	assert.Equal(t, user.Email, returnedUser.Email)
	if user.ProfilePictureURL == nil {
		assert.Nil(t, returnedUser.ProfilePictureURL)
	} else {
		assert.Equal(t, user.ProfilePictureURL.ToNullableFullURL(), returnedUser.ProfilePictureURL)
	}
}
