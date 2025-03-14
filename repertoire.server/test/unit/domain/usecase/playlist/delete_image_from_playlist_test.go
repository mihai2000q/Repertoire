package playlist

import (
	"errors"
	"net/http"
	"repertoire/server/domain/usecase/playlist"
	"repertoire/server/internal"
	"repertoire/server/internal/message/topics"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteImageFromPlaylist_WhenGetPlaylistFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := playlist.NewDeleteImageFromPlaylist(playlistRepository, nil, nil)

	id := uuid.New()

	// given - mocking
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

func TestDeleteImageFromPlaylist_WhenPlaylistIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := playlist.NewDeleteImageFromPlaylist(playlistRepository, nil, nil)

	id := uuid.New()

	// given - mocking
	playlistRepository.On("Get", new(model.Playlist), id).Return(nil).Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "playlist not found", errCode.Error.Error())

	playlistRepository.AssertExpectations(t)
}

func TestDeleteImageFromPlaylist_WhenPlaylistHasNoImage_ShouldReturnBadRequestError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := playlist.NewDeleteImageFromPlaylist(playlistRepository, nil, nil)

	id := uuid.New()

	// given - mocking
	mockPlaylist := &model.Playlist{ID: id}
	playlistRepository.On("Get", new(model.Playlist), id).Return(nil, mockPlaylist).Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusBadRequest, errCode.Code)
	assert.Equal(t, "playlist does not have an image", errCode.Error.Error())

	playlistRepository.AssertExpectations(t)
}

func TestDeleteImageFromPlaylist_WhenDeleteImageFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	storageService := new(service.StorageServiceMock)
	_uut := playlist.NewDeleteImageFromPlaylist(playlistRepository, storageService, nil)

	id := uuid.New()

	// given - mocking
	mockPlaylist := &model.Playlist{ID: id, ImageURL: &[]internal.FilePath{"This is some url"}[0]}
	playlistRepository.On("Get", new(model.Playlist), id).Return(nil, mockPlaylist).Once()

	internalError := wrapper.InternalServerError(errors.New("internal error"))
	storageService.On("DeleteFile", *mockPlaylist.ImageURL).Return(internalError).Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, internalError, errCode)

	playlistRepository.AssertExpectations(t)
	storageService.AssertExpectations(t)
}

func TestDeleteImageFromPlaylist_WhenUpdatePlaylistFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	storageService := new(service.StorageServiceMock)
	_uut := playlist.NewDeleteImageFromPlaylist(playlistRepository, storageService, nil)

	id := uuid.New()

	// given - mocking
	mockPlaylist := &model.Playlist{ID: id, ImageURL: &[]internal.FilePath{"This is some url"}[0]}
	playlistRepository.On("Get", new(model.Playlist), id).Return(nil, mockPlaylist).Once()

	storageService.On("DeleteFile", *mockPlaylist.ImageURL).Return(nil).Once()

	internalError := errors.New("internal error")
	playlistRepository.On("Update", mock.IsType(mockPlaylist)).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	playlistRepository.AssertExpectations(t)
	storageService.AssertExpectations(t)
}

func TestDeleteImageFromPlaylist_WhenPublishFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	storageService := new(service.StorageServiceMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := playlist.NewDeleteImageFromPlaylist(playlistRepository, storageService, messagePublisherService)

	id := uuid.New()

	// given - mocking
	mockPlaylist := &model.Playlist{ID: id, ImageURL: &[]internal.FilePath{"This is some url"}[0]}
	playlistRepository.On("Get", new(model.Playlist), id).Return(nil, mockPlaylist).Once()

	storageService.On("DeleteFile", *mockPlaylist.ImageURL).Return(nil).Once()

	playlistRepository.On("Update", mock.IsType(mockPlaylist)).
		Return(nil).
		Once()

	internalError := errors.New("internal error")
	messagePublisherService.On("Publish", topics.PlaylistUpdatedTopic, mock.IsType(*mockPlaylist)).
		Return(internalError).
		Once()
	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	playlistRepository.AssertExpectations(t)
	storageService.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
}

func TestDeleteImageFromPlaylist_WhenIsValid_ShouldNotReturnAnyError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	storageService := new(service.StorageServiceMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := playlist.NewDeleteImageFromPlaylist(playlistRepository, storageService, messagePublisherService)

	id := uuid.New()

	// given - mocking
	mockPlaylist := &model.Playlist{ID: id, ImageURL: &[]internal.FilePath{"This is some url"}[0]}
	playlistRepository.On("Get", new(model.Playlist), id).Return(nil, mockPlaylist).Once()

	storageService.On("DeleteFile", *mockPlaylist.ImageURL).Return(nil).Once()

	var newPlaylist *model.Playlist
	playlistRepository.On("Update", mock.IsType(mockPlaylist)).
		Run(func(args mock.Arguments) {
			newPlaylist = args.Get(0).(*model.Playlist)
			assert.Nil(t, newPlaylist.ImageURL)
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
	errCode := _uut.Handle(id)

	// then
	assert.Nil(t, errCode)

	playlistRepository.AssertExpectations(t)
	storageService.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
}
