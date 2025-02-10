package section

import (
	"errors"
	"net/http"
	"repertoire/server/domain/usecase/song/section"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetSongSectionTypes_WhenGetUserIdFromJwtFails_ShouldReturnError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	_uut := section.NewGetSongSectionTypes(nil, jwtService)

	token := "this is a token"

	forbiddenError := wrapper.ForbiddenError(errors.New("forbidden error"))
	jwtService.On("GetUserIdFromJwt", token).Return(uuid.Nil, forbiddenError).Once()

	// when
	sectionTypes, errCode := _uut.Handle(token)

	// then
	assert.Empty(t, sectionTypes)
	assert.NotNil(t, errCode)
	assert.Equal(t, forbiddenError, errCode)

	jwtService.AssertExpectations(t)
}

func TestGetSongSectionTypes_WhenGetSectionTypesFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := section.NewGetSongSectionTypes(songRepository, jwtService)

	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	internalError := errors.New("internal error")
	songRepository.On("GetSectionTypes", new([]model.SongSectionType), userID).Return(internalError).Once()

	// when
	sectionTypes, errCode := _uut.Handle(token)

	// then
	assert.Empty(t, sectionTypes)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	jwtService.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}

func TestGetSongSectionTypes_WhenSuccessful_ShouldReturnSongSectionTypes(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := section.NewGetSongSectionTypes(songRepository, jwtService)

	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	expectedTypes := &[]model.SongSectionType{
		{ID: uuid.New(), Name: "Chorus"},
		{ID: uuid.New(), Name: "Interlude"},
	}
	songRepository.On("GetSectionTypes", mock.IsType(expectedTypes), userID).
		Return(nil, expectedTypes).
		Once()

	// when
	sectionTypes, errCode := _uut.Handle(token)

	// then
	assert.Equal(t, expectedTypes, &sectionTypes)
	assert.Nil(t, errCode)

	jwtService.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}
