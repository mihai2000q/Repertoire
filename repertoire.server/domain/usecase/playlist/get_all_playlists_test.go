package playlist

import (
	"errors"
	"net/http"
	"repertoire/api/requests"
	"repertoire/data/repository"
	"repertoire/data/service"
	"repertoire/model"
	"repertoire/utils/wrapper"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAll_WhenGetUserIdFromJwtFails_ShouldReturnForbiddenError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	_uut := &GetAllPlaylists{
		jwtService: jwtService,
	}
	request := requests.GetPlaylistsRequest{}
	token := "This is a token"

	forbiddenError := wrapper.ForbiddenError(errors.New("forbidden error"))
	jwtService.On("GetUserIdFromJwt", token).Return(uuid.Nil, forbiddenError).Once()

	// when
	playlists, errCode := _uut.Handle(request, token)

	// then
	assert.Empty(t, playlists)
	assert.NotNil(t, errCode)
	assert.Equal(t, forbiddenError, errCode)

	jwtService.AssertExpectations(t)
}

func TestGetAll_WhenGetPlaylistsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	jwtService := new(service.JwtServiceMock)
	_uut := &GetAllPlaylists{
		repository: playlistRepository,
		jwtService: jwtService,
	}
	request := requests.GetPlaylistsRequest{}
	token := "This is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	internalError := errors.New("internal error")
	playlistRepository.
		On(
			"GetAllByUser",
			mock.Anything,
			userID,
			request.CurrentPage,
			request.PageSize,
			request.OrderBy,
		).
		Return(internalError).
		Once()

	// when
	playlists, errCode := _uut.Handle(request, token)

	// then
	assert.Empty(t, playlists)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	playlistRepository.AssertExpectations(t)
	jwtService.AssertExpectations(t)
}

func TestGetAll_WhenGetPlaylistsCountFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	jwtService := new(service.JwtServiceMock)
	_uut := &GetAllPlaylists{
		repository: playlistRepository,
		jwtService: jwtService,
	}
	request := requests.GetPlaylistsRequest{}
	token := "This is a token"

	expectedPlaylists := &[]model.Playlist{
		{Title: "Some Playlist"},
		{Title: "Some other Playlist"},
	}

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	playlistRepository.
		On(
			"GetAllByUser",
			mock.IsType(expectedPlaylists),
			userID,
			request.CurrentPage,
			request.PageSize,
			request.OrderBy,
		).
		Return(nil, expectedPlaylists).
		Once()

	internalError := errors.New("internal error")
	playlistRepository.
		On(
			"GetAllByUserCount",
			mock.Anything,
			userID,
		).
		Return(internalError).
		Once()

	// when
	result, errCode := _uut.Handle(request, token)

	// then
	assert.Equal(t, expectedPlaylists, &result.Models)
	assert.Empty(t, result.TotalCount)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	playlistRepository.AssertExpectations(t)
	jwtService.AssertExpectations(t)
}

func TestGetAll_WhenSuccessful_ShouldReturnPlaylistsWithTotalCount(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	jwtService := new(service.JwtServiceMock)
	_uut := &GetAllPlaylists{
		repository: playlistRepository,
		jwtService: jwtService,
	}
	request := requests.GetPlaylistsRequest{}
	token := "This is a token"

	expectedPlaylists := &[]model.Playlist{
		{Title: "Some Playlist"},
		{Title: "Some other Playlist"},
	}
	expectedTotalCount := &[]int64{20}[0]

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	playlistRepository.
		On(
			"GetAllByUser",
			mock.IsType(expectedPlaylists),
			userID,
			request.CurrentPage,
			request.PageSize,
			request.OrderBy,
		).
		Return(nil, expectedPlaylists).
		Once()

	playlistRepository.
		On(
			"GetAllByUserCount",
			mock.IsType(expectedTotalCount),
			userID,
		).
		Return(nil, expectedTotalCount).
		Once()

	// when
	result, errCode := _uut.Handle(request, token)

	// then
	assert.Equal(t, expectedPlaylists, &result.Models)
	assert.Equal(t, expectedTotalCount, &result.TotalCount)
	assert.Nil(t, errCode)

	playlistRepository.AssertExpectations(t)
	jwtService.AssertExpectations(t)
}
