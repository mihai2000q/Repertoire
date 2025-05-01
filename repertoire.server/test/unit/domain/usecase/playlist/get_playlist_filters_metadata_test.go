package playlist

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/playlist"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetPlaylistFiltersMetadata_WhenGetUserIdFromJwtFails_ShouldReturnInternalServerErrorError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := playlist.NewGetPlaylistFiltersMetadata(jwtService, playlistRepository)

	request := requests.GetPlaylistFiltersMetadataRequest{}
	token := "some token"

	internalError := wrapper.InternalServerError(errors.New("some internal error"))
	jwtService.On("GetUserIdFromJwt", token).Return(uuid.Nil, internalError).Once()

	// when
	result, errCode := _uut.Handle(request, token)

	// then
	assert.Empty(t, result)
	assert.NotNil(t, errCode)
	assert.Equal(t, internalError, errCode)

	jwtService.AssertExpectations(t)
	playlistRepository.AssertExpectations(t)
}

func TestGetPlaylistFiltersMetadata_WhenGetFails_ShouldReturnInternalServerErrorError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := playlist.NewGetPlaylistFiltersMetadata(jwtService, playlistRepository)

	request := requests.GetPlaylistFiltersMetadataRequest{}
	token := "some token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	internalError := errors.New("some internal error")
	playlistRepository.On("GetFiltersMetadata", new(model.PlaylistFiltersMetadata), userID, request.SearchBy).
		Return(internalError).
		Once()

	// when
	result, errCode := _uut.Handle(request, token)

	// then
	assert.Empty(t, result)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	jwtService.AssertExpectations(t)
	playlistRepository.AssertExpectations(t)
}

func TestGetPlaylistFiltersMetadata_WhenSuccessful_ShouldReturnPlaylistFiltersMetadata(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := playlist.NewGetPlaylistFiltersMetadata(jwtService, playlistRepository)

	request := requests.GetPlaylistFiltersMetadataRequest{}
	token := "some token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	metadata := model.PlaylistFiltersMetadata{
		MinSongsCount: 0,
	}
	playlistRepository.On("GetFiltersMetadata", new(model.PlaylistFiltersMetadata), userID, request.SearchBy).
		Return(nil, &metadata).
		Once()

	// when
	result, errCode := _uut.Handle(request, token)

	// then
	assert.Nil(t, errCode)
	assert.Equal(t, metadata, result)

	jwtService.AssertExpectations(t)
	playlistRepository.AssertExpectations(t)
}
