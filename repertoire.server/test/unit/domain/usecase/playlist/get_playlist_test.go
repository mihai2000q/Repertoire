package playlist

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/playlist"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetPlaylist_WhenGetPlaylistFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := playlist.NewGetPlaylist(playlistRepository)

	request := requests.GetPlaylistRequest{ID: uuid.New()}

	internalError := errors.New("internal error")
	playlistRepository.On("Get", new(model.Playlist), request.ID).
		Return(internalError).
		Once()

	// when
	resultPlaylist, errCode := _uut.Handle(request)

	// then
	assert.Empty(t, resultPlaylist)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	playlistRepository.AssertExpectations(t)
}

func TestGetPlaylist_WhenPlaylistIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := playlist.NewGetPlaylist(playlistRepository)

	request := requests.GetPlaylistRequest{ID: uuid.New()}

	playlistRepository.On("Get", new(model.Playlist), request.ID).
		Return(nil).
		Once()

	// when
	resultPlaylist, errCode := _uut.Handle(request)

	// then
	assert.Empty(t, resultPlaylist)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "playlist not found", errCode.Error.Error())

	playlistRepository.AssertExpectations(t)
}

func TestGetPlaylist_WhenSuccessful_ShouldReturnPlaylist(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := playlist.NewGetPlaylist(playlistRepository)

	request := requests.GetPlaylistRequest{ID: uuid.New()}

	expectedPlaylist := &model.Playlist{
		ID:    request.ID,
		Title: "Some Playlist",
	}
	playlistRepository.On("Get", new(model.Playlist), request.ID).
		Return(nil, expectedPlaylist).
		Once()

	// when
	resultPlaylist, errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)
	assert.NotEmpty(t, resultPlaylist)

	assert.Equal(t, expectedPlaylist.ID, resultPlaylist.ID)
	assert.Equal(t, expectedPlaylist.Title, resultPlaylist.Title)
	assert.Equal(t, expectedPlaylist.Description, resultPlaylist.Description)
	assert.Equal(t, expectedPlaylist.ImageURL, resultPlaylist.ImageURL)

	playlistRepository.AssertExpectations(t)
}
