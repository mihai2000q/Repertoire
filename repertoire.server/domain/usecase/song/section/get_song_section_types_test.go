package section

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"repertoire/data/repository"
	"repertoire/data/service"
	"repertoire/model"
	"repertoire/utils/wrapper"
	"testing"
)

func TestGetSongSectionTypes_WhenGetUserIdFromJwtFails_ShouldReturnError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	_uut := &GetSongSectionTypes{
		jwtService: jwtService,
	}
	token := "this is a token"

	forbiddenError := wrapper.ForbiddenError(errors.New("forbidden error"))
	jwtService.On("GetUserIdFromJwt", token).Return(uuid.Nil, forbiddenError).Once()

	// when
	types, errCode := _uut.Handle(token)

	// then
	assert.Empty(t, types)
	assert.NotNil(t, errCode)
	assert.Equal(t, forbiddenError, errCode)

	jwtService.AssertExpectations(t)
}

func TestGetSongSectionTypes_WhenGetSectionTypesFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := &GetSongSectionTypes{
		repository: songRepository,
		jwtService: jwtService,
	}
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	internalError := errors.New("internal error")
	songRepository.On("GetSectionTypes", new([]model.SongSectionType), userID).Return(internalError).Once()

	// when
	types, errCode := _uut.Handle(token)

	// then
	assert.Empty(t, types)
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
	_uut := &GetSongSectionTypes{
		repository: songRepository,
		jwtService: jwtService,
	}
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
	types, errCode := _uut.Handle(token)

	// then
	assert.Equal(t, expectedTypes, &types)
	assert.Nil(t, errCode)

	jwtService.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}
