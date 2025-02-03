package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"repertoire/server/api/requests"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/assertion"
	"repertoire/server/test/integration/test/core"
	"repertoire/server/test/integration/test/data/auth"
	"repertoire/server/test/integration/test/utils"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestSignUp_WhenUserAlreadyExists_ShouldReturnBadRequestError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, auth.Users, auth.SeedData)

	user := auth.Users[0]
	request := requests.SignUpRequest{
		Name:     "Nigel",
		Email:    user.Email,
		Password: "Password123",
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().
		WithoutAuthentication().
		POST(w, "/api/auth/sign-up", request)

	// then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSignUp_WhenSuccessful_ShouldCreateUserAndReturnToken(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, []model.User{}, func(*gorm.DB) {})

	request := requests.SignUpRequest{
		Name:     "Nigel",
		Email:    "Nigel@Herrison.com",
		Password: "Password123",
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().
		WithoutAuthentication().
		POST(w, "/api/auth/sign-up", request)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	user := assertCreatedUser(t, request)

	var response struct{ Token string }
	_ = json.Unmarshal(w.Body.Bytes(), &response)

	assertion.Token(t, response.Token, user)
}

func assertCreatedUser(t *testing.T, request requests.SignUpRequest) model.User {
	db := utils.GetDatabase(t)

	var user model.User
	db.Preload("SongSectionTypes").
		Preload("GuitarTunings").
		Find(&user, model.User{Email: strings.ToLower(request.Email)})

	assert.NotEmpty(t, user.ID)
	assert.Equal(t, request.Name, user.Name)
	assert.Equal(t, strings.ToLower(request.Email), user.Email)
	assert.NotEmpty(t, user.Password)
	assert.Nil(t, user.ProfilePictureURL)

	assert.Len(t, user.GuitarTunings, len(model.DefaultGuitarTunings))
	for i, guitarTuning := range user.GuitarTunings {
		assert.NotEmpty(t, guitarTuning.ID)
		assert.Equal(t, model.DefaultGuitarTunings[i], guitarTuning.Name)
		assert.Equal(t, uint(i), guitarTuning.Order)
	}

	assert.Len(t, user.SongSectionTypes, len(model.DefaultSongSectionTypes))
	for i, sectionType := range user.SongSectionTypes {
		assert.NotEmpty(t, sectionType.ID)
		assert.Equal(t, model.DefaultSongSectionTypes[i], sectionType.Name)
		assert.Equal(t, uint(i), sectionType.Order)
	}

	assert.Len(t, user.BandMemberRoles, len(model.DefaultBandMemberRoles))
	for i, bandMemberRole := range user.BandMemberRoles {
		assert.NotEmpty(t, bandMemberRole.ID)
		assert.Equal(t, model.DefaultSongSectionTypes[i], bandMemberRole.Name)
		assert.Equal(t, uint(i), bandMemberRole.Order)
	}

	return user
}
