package playlist

import (
	"errors"
	"net/http"
	"repertoire/api/requests"
	"repertoire/data/repository"
	"repertoire/models"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAll_WhenGetPlaylistsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := &GetAllPlaylists{
		repository: playlistRepository,
	}
	request := requests.GetPlaylistsRequest{
		UserID: uuid.New(),
	}

	internalError := errors.New("internal error")
	playlistRepository.
		On(
			"GetAllByUser",
			mock.Anything,
			request.UserID,
			request.CurrentPage,
			request.PageSize,
		).
		Return(internalError).
		Once()

	// when
	playlists, errCode := _uut.Handle(request)

	// then
	assert.Empty(t, playlists)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	playlistRepository.AssertExpectations(t)
}

func TestGetAll_WhenGetPlaylistsCountFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := &GetAllPlaylists{
		repository: playlistRepository,
	}
	request := requests.GetPlaylistsRequest{
		UserID: uuid.New(),
	}

	expectedPlaylists := &[]models.Playlist{
		{Title: "Some Playlist"},
		{Title: "Some other Playlist"},
	}

	playlistRepository.
		On(
			"GetAllByUser",
			mock.IsType(expectedPlaylists),
			request.UserID,
			request.CurrentPage,
			request.PageSize,
		).
		Return(nil, expectedPlaylists).
		Once()

	internalError := errors.New("internal error")
	playlistRepository.
		On(
			"GetAllByUserCount",
			mock.Anything,
			request.UserID,
		).
		Return(internalError).
		Once()

	// when
	result, errCode := _uut.Handle(request)

	// then
	assert.Equal(t, expectedPlaylists, &result.Data)
	assert.Empty(t, result.TotalCount)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	playlistRepository.AssertExpectations(t)
}

func TestGetAll_WhenSuccessful_ShouldReturnPlaylistsWithTotalCount(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := &GetAllPlaylists{
		repository: playlistRepository,
	}
	request := requests.GetPlaylistsRequest{
		UserID: uuid.New(),
	}

	expectedPlaylists := &[]models.Playlist{
		{Title: "Some Playlist"},
		{Title: "Some other Playlist"},
	}
	expectedTotalCount := &[]int64{20}[0]

	playlistRepository.
		On(
			"GetAllByUser",
			mock.IsType(expectedPlaylists),
			request.UserID,
			request.CurrentPage,
			request.PageSize,
		).
		Return(nil, expectedPlaylists).
		Once()

	playlistRepository.
		On(
			"GetAllByUserCount",
			mock.IsType(expectedTotalCount),
			request.UserID,
		).
		Return(nil, expectedTotalCount).
		Once()

	// when
	result, errCode := _uut.Handle(request)

	// then
	assert.Equal(t, expectedPlaylists, &result.Data)
	assert.Equal(t, expectedTotalCount, &result.TotalCount)
	assert.Nil(t, errCode)

	playlistRepository.AssertExpectations(t)
}
