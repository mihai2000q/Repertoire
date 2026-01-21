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

func TestBulkDeletePlaylists_WhenGetPlaylistsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := playlist.NewBulkDeletePlaylists(playlistRepository, nil)

	request := requests.BulkDeletePlaylistsRequest{
		IDs: []uuid.UUID{uuid.New()},
	}

	internalError := errors.New("internal error")
	playlistRepository.On("GetAllByIDs", mock.IsType(&[]model.Playlist{}), request.IDs).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	playlistRepository.AssertExpectations(t)
}

func TestBulkDeletePlaylists_WhenPlaylistsAreEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := playlist.NewBulkDeletePlaylists(playlistRepository, nil)

	request := requests.BulkDeletePlaylistsRequest{
		IDs: []uuid.UUID{uuid.New()},
	}

	playlistRepository.On("GetAllByIDs", mock.IsType(&[]model.Playlist{}), request.IDs).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "playlists not found", errCode.Error.Error())

	playlistRepository.AssertExpectations(t)
}

func TestBulkDeletePlaylists_WhenDeletePlaylistsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := playlist.NewBulkDeletePlaylists(playlistRepository, nil)

	request := requests.BulkDeletePlaylistsRequest{
		IDs: []uuid.UUID{uuid.New()},
	}

	mockPlaylists := &[]model.Playlist{{ID: request.IDs[0]}}
	playlistRepository.On("GetAllByIDs", mock.IsType(mockPlaylists), request.IDs).
		Return(nil, mockPlaylists).
		Once()

	internalError := errors.New("internal error")
	playlistRepository.On("Delete", request.IDs).Return(internalError).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	playlistRepository.AssertExpectations(t)
}

func TestBulkDeletePlaylists_WhenPublishFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := playlist.NewBulkDeletePlaylists(playlistRepository, messagePublisherService)

	request := requests.BulkDeletePlaylistsRequest{
		IDs: []uuid.UUID{uuid.New()},
	}

	mockPlaylists := &[]model.Playlist{{ID: request.IDs[0]}}
	playlistRepository.On("GetAllByIDs", mock.IsType(mockPlaylists), request.IDs).
		Return(nil, mockPlaylists).
		Once()

	playlistRepository.On("Delete", request.IDs).Return(nil).Once()

	internalError := errors.New("internal error")
	messagePublisherService.On("Publish", topics.PlaylistsDeletedTopic, *mockPlaylists).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	playlistRepository.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
}

func TestBulkDeletePlaylists_WhenIsSuccessful_ShouldDeletePlaylists(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := playlist.NewBulkDeletePlaylists(playlistRepository, messagePublisherService)

	request := requests.BulkDeletePlaylistsRequest{
		IDs: []uuid.UUID{
			uuid.New(),
			uuid.New(),
			uuid.New(),
		},
	}

	mockPlaylists := &[]model.Playlist{
		{ID: request.IDs[0]},
		{ID: request.IDs[1]},
		{ID: request.IDs[2]},
	}

	playlistRepository.On("GetAllByIDs", mock.IsType(mockPlaylists), request.IDs).
		Return(nil, mockPlaylists).
		Once()
	playlistRepository.On("Delete", request.IDs).Return(nil).Once()

	messagePublisherService.On("Publish", topics.PlaylistsDeletedTopic, *mockPlaylists).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)

	playlistRepository.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
}
