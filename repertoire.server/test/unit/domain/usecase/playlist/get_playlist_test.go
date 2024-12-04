package playlist

import (
	"errors"
	"net/http"
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

	id := uuid.New()

	internalError := errors.New("internal error")
	playlistRepository.On("GetWithAssociations", new(model.Playlist), id).
		Return(internalError).
		Once()

	// when
	resultPlaylist, errCode := _uut.Handle(id)

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

	id := uuid.New()

	playlistRepository.On("GetWithAssociations", new(model.Playlist), id).
		Return(nil).
		Once()

	// when
	resultPlaylist, errCode := _uut.Handle(id)

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

	id := uuid.New()

	mockPlaylist := &model.Playlist{
		ID:    id,
		Title: "Some Playlist",
		PlaylistSongs: []model.PlaylistSong{
			{
				Song:        model.Song{ID: uuid.New()},
				SongTrackNo: 1,
			},
			{
				Song:        model.Song{ID: uuid.New()},
				SongTrackNo: 2,
			},
		},
	}

	playlistRepository.On("GetWithAssociations", new(model.Playlist), id).
		Return(nil, mockPlaylist).
		Once()

	// when
	resultPlaylist, errCode := _uut.Handle(id)

	// then
	assert.Nil(t, errCode)
	assert.NotEmpty(t, resultPlaylist)

	assert.Equal(t, mockPlaylist.ID, resultPlaylist.ID)
	assert.Equal(t, mockPlaylist.Title, resultPlaylist.Title)
	assert.Equal(t, mockPlaylist.Description, resultPlaylist.Description)
	assert.Equal(t, mockPlaylist.ImageURL, resultPlaylist.ImageURL)

	for i := range resultPlaylist.Songs {
		assert.Equal(t, mockPlaylist.PlaylistSongs[i].Song.ID, resultPlaylist.Songs[i].ID)
		assert.Equal(t, mockPlaylist.PlaylistSongs[i].SongTrackNo, resultPlaylist.Songs[i].PlaylistTrackNo)
		assert.Equal(t, mockPlaylist.PlaylistSongs[i].CreatedAt, resultPlaylist.Songs[i].PlaylistCreatedAt)
	}

	playlistRepository.AssertExpectations(t)
}
