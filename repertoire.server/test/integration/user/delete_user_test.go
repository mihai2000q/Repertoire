package user

import (
	"net/http"
	"net/http/httptest"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/assertion"
	"repertoire/server/test/integration/test/core"
	userData "repertoire/server/test/integration/test/data/user"
	"repertoire/server/test/integration/test/utils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeleteUser_WhenSuccessful_ShouldDeleteUser(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, userData.Users, userData.SeedData)

	user := userData.Users[0]

	messages := utils.SubscribeToTopic(topics.UserDeletedTopic)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().
		WithUser(user).
		DELETE(w, "/api/users")

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)

	var deletedUser model.User
	db.Find(&deletedUser, user.ID)

	assert.Empty(t, deletedUser)

	assertion.AssertMessage(t, messages, func(userId uuid.UUID) {
		assert.Equal(t, user.ID, userId)
	})
}
