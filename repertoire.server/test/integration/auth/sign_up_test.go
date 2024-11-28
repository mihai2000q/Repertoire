package auth

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"repertoire/server/api/requests"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/core"
	"repertoire/server/test/integration/test/data/auth"
	"repertoire/server/test/integration/test/utils"
	"strings"
	"testing"
)

func TestSignUp_WhenUserAlreadyExists(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, auth.SeedData)

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

func TestSignUp(t *testing.T) {
	// given
	utils.CleanupData(t)

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
	//assert.JSONEq(t, `{"token": ""}`, w.Body.String())
	assertCreatedUser(t, request)
}

func assertCreatedUser(t *testing.T, request requests.SignUpRequest) {
	db := utils.GetDatabase()

	var user model.User
	db.Preload("SongSectionTypes").
		Preload("GuitarTunings").
		Find(&user, model.User{Email: strings.ToLower(request.Email)})

	assert.NotEmpty(t, user.ID)
	assert.Equal(t, request.Name, user.Name)
	assert.Equal(t, strings.ToLower(request.Email), user.Email)
	assert.NotEmpty(t, user.Password)
	assert.Nil(t, user.ProfilePictureURL)

	assert.Len(t, user.SongSectionTypes, len(model.DefaultSongSectionTypes))
	for i, sectionType := range user.SongSectionTypes {
		assert.NotEmpty(t, sectionType.ID)
		assert.Equal(t, model.DefaultSongSectionTypes[i], sectionType.Name)
		assert.Equal(t, uint(i), sectionType.Order)
	}

	assert.Len(t, user.GuitarTunings, len(model.DefaultGuitarTunings))
	for i, guitarTuning := range user.GuitarTunings {
		assert.NotEmpty(t, guitarTuning.ID)
		assert.Equal(t, model.DefaultGuitarTunings[i], guitarTuning.Name)
		assert.Equal(t, uint(i), guitarTuning.Order)
	}
}