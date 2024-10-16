package playlist

import (
	"errors"
	"net/http"
	"repertoire/api/request"
	"repertoire/data/repository"
	"repertoire/model"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdatePlaylist_WhenGetPlaylistFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := &UpdatePlaylist{
		repository: playlistRepository,
	}
	request := request.UpdatePlaylistRequest{
		ID:    uuid.New(),
		Title: "New Playlist",
	}

	internalError := errors.New("internal error")
	playlistRepository.On("Get", new(model.Playlist), request.ID).Return(internalError).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	playlistRepository.AssertExpectations(t)
}

func TestUpdatePlaylist_WhenPlaylistIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := &UpdatePlaylist{
		repository: playlistRepository,
	}
	request := request.UpdatePlaylistRequest{
		ID:    uuid.New(),
		Title: "New Playlist",
	}

	playlistRepository.On("Get", new(model.Playlist), request.ID).Return(nil).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "playlist not found", errCode.Error.Error())

	playlistRepository.AssertExpectations(t)
}

func TestUpdatePlaylist_WhenUpdatePlaylistFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := &UpdatePlaylist{
		repository: playlistRepository,
	}
	request := request.UpdatePlaylistRequest{
		ID:    uuid.New(),
		Title: "New Playlist",
	}

	playlist := &model.Playlist{
		ID:    request.ID,
		Title: "Some Playlist",
	}

	playlistRepository.On("Get", new(model.Playlist), request.ID).Return(nil, playlist).Once()
	internalError := errors.New("internal error")
	playlistRepository.On("Update", mock.IsType(playlist)).
		Run(func(args mock.Arguments) {
			newPlaylist := args.Get(0).(*model.Playlist)
			assert.Equal(t, request.Title, newPlaylist.Title)
			assert.Equal(t, request.Description, newPlaylist.Description)
		}).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	playlistRepository.AssertExpectations(t)
}

func TestUpdatePlaylist_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := &UpdatePlaylist{
		repository: playlistRepository,
	}
	request := request.UpdatePlaylistRequest{
		ID:    uuid.New(),
		Title: "New Playlist",
	}

	playlist := &model.Playlist{
		ID:    request.ID,
		Title: "Some Playlist",
	}

	playlistRepository.On("Get", new(model.Playlist), request.ID).Return(nil, playlist).Once()
	playlistRepository.On("Update", mock.IsType(playlist)).
		Run(func(args mock.Arguments) {
			newPlaylist := args.Get(0).(*model.Playlist)
			assert.Equal(t, request.Title, newPlaylist.Title)
			assert.Equal(t, request.Description, newPlaylist.Description)
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)

	playlistRepository.AssertExpectations(t)
}
