package playlist

import (
	"errors"
	"net/http"
	"repertoire/server/domain/usecase/playlist"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeletePlaylist_WhenGetPlaylistFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := playlist.NewDeletePlaylist(playlistRepository, nil)

	id := uuid.New()

	internalError := errors.New("internal error")
	playlistRepository.On("Get", new(model.Playlist), id).Return(internalError).Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	playlistRepository.AssertExpectations(t)
}

func TestDeletePlaylist_WhenPlaylistIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := playlist.NewDeletePlaylist(playlistRepository, nil)

	id := uuid.New()

	playlistRepository.On("Get", new(model.Playlist), id).Return(nil).Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "playlist not found", errCode.Error.Error())

	playlistRepository.AssertExpectations(t)
}

func TestDeletePlaylist_WhenDeletePlaylistFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := playlist.NewDeletePlaylist(playlistRepository, nil)

	id := uuid.New()

	mockPlaylist := &model.Playlist{ID: id}
	playlistRepository.On("Get", new(model.Playlist), id).Return(nil, mockPlaylist).Once()

	internalError := errors.New("internal error")
	playlistRepository.On("Delete", id).Return(internalError).Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	playlistRepository.AssertExpectations(t)
}

func TestDeletePlaylist_WhenPublishFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := playlist.NewDeletePlaylist(playlistRepository, messagePublisherService)

	id := uuid.New()

	mockPlaylist := &model.Playlist{
		ID: id,
	}
	playlistRepository.On("Get", new(model.Playlist), id).Return(nil, mockPlaylist).Once()

	playlistRepository.On("Delete", id).Return(nil).Once()

	internalError := errors.New("internal error")
	messagePublisherService.On("Publish", topics.PlaylistsDeletedTopic, mock.IsType([]model.Playlist{})).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	playlistRepository.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
}

func TestDeletePlaylist_WhenSuccessful_ShouldDeletePlaylist(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := playlist.NewDeletePlaylist(playlistRepository, messagePublisherService)

	expectedPlaylist := model.Playlist{ID: uuid.New()}
	playlistRepository.On("Get", new(model.Playlist), expectedPlaylist.ID).
		Return(nil, &expectedPlaylist).
		Once()

	playlistRepository.On("Delete", expectedPlaylist.ID).Return(nil).Once()

	messagePublisherService.On("Publish", topics.PlaylistsDeletedTopic, []model.Playlist{expectedPlaylist}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(expectedPlaylist.ID)

	// then
	assert.Nil(t, errCode)

	playlistRepository.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
}
