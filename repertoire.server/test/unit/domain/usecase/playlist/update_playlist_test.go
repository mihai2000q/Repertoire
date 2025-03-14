package playlist

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/playlist"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdatePlaylist_WhenGetPlaylistFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := playlist.NewUpdatePlaylist(playlistRepository, nil)

	request := requests.UpdatePlaylistRequest{
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
	_uut := playlist.NewUpdatePlaylist(playlistRepository, nil)

	request := requests.UpdatePlaylistRequest{
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
	_uut := playlist.NewUpdatePlaylist(playlistRepository, nil)

	request := requests.UpdatePlaylistRequest{
		ID:    uuid.New(),
		Title: "New Playlist",
	}

	mockPlaylist := &model.Playlist{
		ID:    request.ID,
		Title: "Some Playlist",
	}
	playlistRepository.On("Get", new(model.Playlist), request.ID).
		Return(nil, mockPlaylist).
		Once()

	internalError := errors.New("internal error")
	playlistRepository.On("Update", mock.IsType(mockPlaylist)).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	playlistRepository.AssertExpectations(t)
}

func TestUpdatePlaylist_WhenPublishFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := playlist.NewUpdatePlaylist(playlistRepository, messagePublisherService)

	request := requests.UpdatePlaylistRequest{
		ID:    uuid.New(),
		Title: "New Playlist",
	}

	mockPlaylist := &model.Playlist{
		ID:    request.ID,
		Title: "Some Playlist",
	}
	playlistRepository.On("Get", new(model.Playlist), request.ID).
		Return(nil, mockPlaylist).
		Once()

	playlistRepository.On("Update", mock.IsType(mockPlaylist)).
		Return(nil).
		Once()

	internalError := errors.New("internal error")
	messagePublisherService.On("Publish", topics.PlaylistUpdatedTopic, mock.IsType(*mockPlaylist)).
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
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := playlist.NewUpdatePlaylist(playlistRepository, messagePublisherService)

	request := requests.UpdatePlaylistRequest{
		ID:    uuid.New(),
		Title: "New Playlist",
	}

	mockPlaylist := &model.Playlist{
		ID:    request.ID,
		Title: "Some Playlist",
	}
	playlistRepository.On("Get", new(model.Playlist), request.ID).
		Return(nil, mockPlaylist).
		Once()

	var newPlaylist *model.Playlist
	playlistRepository.On("Update", mock.IsType(mockPlaylist)).
		Run(func(args mock.Arguments) {
			newPlaylist = args.Get(0).(*model.Playlist)
			assertUpdatedPlaylist(t, *newPlaylist, request)
		}).
		Return(nil).
		Once()

	messagePublisherService.On("Publish", topics.PlaylistUpdatedTopic, mock.IsType(*mockPlaylist)).
		Run(func(args mock.Arguments) {
			assert.Equal(t, *newPlaylist, args.Get(1).(model.Playlist))
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)

	playlistRepository.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
}

func assertUpdatedPlaylist(
	t *testing.T,
	playlist model.Playlist,
	request requests.UpdatePlaylistRequest,
) {
	assert.Equal(t, request.Title, playlist.Title)
	assert.Equal(t, request.Description, playlist.Description)
}
