package playlist

import (
	"errors"
	"net/http"
	"repertoire/server/data/repository"
	"repertoire/server/model"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetPlaylistQuery_WhenGetPlaylistFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := &GetPlaylist{
		repository: playlistRepository,
	}
	id := uuid.New()

	internalError := errors.New("internal error")
	playlistRepository.On("GetWithAssociations", new(model.Playlist), id).
		Return(internalError).
		Once()

	// when
	playlist, errCode := _uut.Handle(id)

	// then
	assert.Empty(t, playlist)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	playlistRepository.AssertExpectations(t)
}

func TestGetPlaylistQuery_WhenPlaylistIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := &GetPlaylist{
		repository: playlistRepository,
	}
	id := uuid.New()

	playlistRepository.On("GetWithAssociations", new(model.Playlist), id).
		Return(nil).
		Once()

	// when
	playlist, errCode := _uut.Handle(id)

	// then
	assert.Empty(t, playlist)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "playlist not found", errCode.Error.Error())

	playlistRepository.AssertExpectations(t)
}

func TestGetPlaylistQuery_WhenSuccessful_ShouldReturnPlaylist(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := &GetPlaylist{
		repository: playlistRepository,
	}
	id := uuid.New()

	expectedPlaylist := &model.Playlist{
		ID:    id,
		Title: "Some Playlist",
	}

	playlistRepository.On("GetWithAssociations", new(model.Playlist), id).
		Return(nil, expectedPlaylist).
		Once()

	// when
	playlist, errCode := _uut.Handle(id)

	// then
	assert.NotEmpty(t, playlist)
	assert.Equal(t, expectedPlaylist, &playlist)
	assert.Nil(t, errCode)

	playlistRepository.AssertExpectations(t)
}
