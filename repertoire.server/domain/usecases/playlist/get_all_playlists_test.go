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

func TestGetAll_WhenGetPlaylistFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := &GetAllPlaylists{
		repository: playlistRepository,
	}
	request := requests.GetPlaylistsRequest{
		UserID: uuid.New(),
	}

	internalError := errors.New("internal error")
	playlistRepository.On("GetAllByUser", mock.Anything, request.UserID).
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

func TestGetAll_WhenSuccessful_ShouldReturnPlaylists(t *testing.T) {
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

	playlistRepository.On("GetAllByUser", mock.IsType(expectedPlaylists), request.UserID).
		Return(nil, expectedPlaylists).
		Once()

	// when
	playlists, errCode := _uut.Handle(request)

	// then
	assert.Equal(t, expectedPlaylists, &playlists)
	assert.Nil(t, errCode)

	playlistRepository.AssertExpectations(t)
}
