package playlist

import (
	"errors"
	"mime/multipart"
	"net/http"
	"repertoire/server/domain/usecase/playlist"
	"repertoire/server/internal"
	"repertoire/server/internal/message/topics"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"repertoire/server/test/unit/domain/provider"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSaveImageToPlaylist_WhenGetPlaylistFails_ShouldReturnNotFoundError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := playlist.NewSaveImageToPlaylist(
		playlistRepository,
		nil,
		nil,
		nil,
	)

	file := new(multipart.FileHeader)
	id := uuid.New()

	// given - mocking
	internalError := errors.New("internal error")
	playlistRepository.On("Get", new(model.Playlist), id).Return(internalError).Once()

	// when
	errCode := _uut.Handle(file, id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	playlistRepository.AssertExpectations(t)
}

func TestSaveImageToPlaylist_WhenPlaylistIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := playlist.NewSaveImageToPlaylist(
		playlistRepository,
		nil,
		nil,
		nil,
	)

	file := new(multipart.FileHeader)
	id := uuid.New()

	// given - mocking
	playlistRepository.On("Get", new(model.Playlist), id).Return(nil).Once()

	// when
	errCode := _uut.Handle(file, id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "playlist not found", errCode.Error.Error())

	playlistRepository.AssertExpectations(t)
}

func TestSaveImageToPlaylist_WhenStorageDeleteFileFails_ShouldReturnError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	storageService := new(service.StorageServiceMock)
	_uut := playlist.NewSaveImageToPlaylist(
		playlistRepository,
		nil,
		storageService,
		nil,
	)

	file := new(multipart.FileHeader)
	id := uuid.New()

	// given - mocking
	mockPlaylist := &model.Playlist{ID: id, ImageURL: &[]internal.FilePath{"file_path"}[0]}
	playlistRepository.On("Get", new(model.Playlist), id).Return(nil, mockPlaylist).Once()

	internalError := wrapper.InternalServerError(errors.New("internal error"))
	storageService.On("DeleteFile", *mockPlaylist.ImageURL).Return(internalError).Once()

	// when
	errCode := _uut.Handle(file, id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, internalError, errCode)

	playlistRepository.AssertExpectations(t)
	storageService.AssertExpectations(t)
}

func TestSaveImageToPlaylist_WhenStorageUploadFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	storageService := new(service.StorageServiceMock)
	_uut := playlist.NewSaveImageToPlaylist(
		playlistRepository,
		storageFilePathProvider,
		storageService,
		nil,
	)

	file := new(multipart.FileHeader)
	id := uuid.New()

	// given - mocking
	mockPlaylist := &model.Playlist{ID: id, ImageURL: nil}
	playlistRepository.On("Get", new(model.Playlist), id).Return(nil, mockPlaylist).Once()

	imagePath := "playlists file path"
	storageFilePathProvider.On("GetPlaylistImagePath", file, mock.IsType(*mockPlaylist)).
		Return(imagePath).
		Once()

	internalError := errors.New("internal error")
	storageService.On("Upload", file, imagePath).Return(internalError).Once()

	// when
	errCode := _uut.Handle(file, id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	playlistRepository.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
	storageService.AssertExpectations(t)
}

func TestSaveImageToPlaylist_WhenUpdatePlaylistFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	storageService := new(service.StorageServiceMock)
	_uut := playlist.NewSaveImageToPlaylist(
		playlistRepository,
		storageFilePathProvider,
		storageService,
		nil,
	)

	file := new(multipart.FileHeader)
	id := uuid.New()

	// given - mocking
	mockPlaylist := &model.Playlist{ID: id, ImageURL: nil}
	playlistRepository.On("Get", new(model.Playlist), id).Return(nil, mockPlaylist).Once()

	imagePath := "playlists file path"
	storageFilePathProvider.On("GetPlaylistImagePath", file, mock.IsType(*mockPlaylist)).
		Return(imagePath).
		Once()

	storageService.On("Upload", file, imagePath).Return(nil).Once()

	internalError := errors.New("internal error")
	playlistRepository.On("Update", mock.IsType(new(model.Playlist))).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(file, id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	playlistRepository.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
	storageService.AssertExpectations(t)
}

func TestSaveImageToPlaylist_WhenPublishFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	storageService := new(service.StorageServiceMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := playlist.NewSaveImageToPlaylist(
		playlistRepository,
		storageFilePathProvider,
		storageService,
		messagePublisherService,
	)

	file := new(multipart.FileHeader)
	id := uuid.New()

	// given - mocking
	mockPlaylist := &model.Playlist{ID: id, ImageURL: nil}
	playlistRepository.On("Get", new(model.Playlist), id).Return(nil, mockPlaylist).Once()

	imagePath := "playlists file path"
	storageFilePathProvider.On("GetPlaylistImagePath", file, mock.IsType(*mockPlaylist)).
		Return(imagePath).
		Once()

	storageService.On("Upload", file, imagePath).Return(nil).Once()

	playlistRepository.On("Update", mock.IsType(new(model.Playlist))).
		Return(nil).
		Once()

	internalError := errors.New("internal error")
	messagePublisherService.On("Publish", topics.PlaylistUpdatedTopic, mock.IsType(*mockPlaylist)).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(file, id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	playlistRepository.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
	storageService.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
}

func TestSaveImageToPlaylist_WhenWithoutOldImage_ShouldSaveNewOneAndNotReturnAnyError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	storageService := new(service.StorageServiceMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := playlist.NewSaveImageToPlaylist(
		playlistRepository,
		storageFilePathProvider,
		storageService,
		messagePublisherService,
	)

	file := new(multipart.FileHeader)
	id := uuid.New()

	// given - mocking
	mockPlaylist := &model.Playlist{ID: id, ImageURL: nil}
	playlistRepository.On("Get", new(model.Playlist), id).Return(nil, mockPlaylist).Once()

	imagePath := "playlists file path"
	storageFilePathProvider.On("GetPlaylistImagePath", file, mock.IsType(*mockPlaylist)).
		Run(func(args mock.Arguments) {
			newPlaylist := args.Get(1).(model.Playlist)
			assert.Equal(t, newPlaylist.ID, mockPlaylist.ID)
			assert.WithinDuration(t, time.Now(), newPlaylist.UpdatedAt, time.Minute)
		}).
		Return(imagePath).
		Once()

	storageService.On("Upload", file, imagePath).Return(nil).Once()

	var newPlaylist *model.Playlist
	playlistRepository.On("Update", mock.IsType(new(model.Playlist))).
		Run(func(args mock.Arguments) {
			newPlaylist = args.Get(0).(*model.Playlist)
			assert.Equal(t, imagePath, string(*newPlaylist.ImageURL))
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
	errCode := _uut.Handle(file, id)

	// then
	assert.Nil(t, errCode)

	playlistRepository.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
	storageService.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
}

func TestSaveImageToPlaylist_WhenWithOldImage_ShouldDeleteOldImageSaveNewOneAndNotReturnAnyError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	storageService := new(service.StorageServiceMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := playlist.NewSaveImageToPlaylist(
		playlistRepository,
		storageFilePathProvider,
		storageService,
		messagePublisherService,
	)

	file := new(multipart.FileHeader)
	id := uuid.New()

	// given - mocking
	mockPlaylist := &model.Playlist{ID: id, ImageURL: &[]internal.FilePath{"file_path"}[0]}
	playlistRepository.On("Get", new(model.Playlist), id).Return(nil, mockPlaylist).Once()

	storageService.On("DeleteFile", *mockPlaylist.ImageURL).Return(nil).Once()

	imagePath := "playlists file path"
	storageFilePathProvider.On("GetPlaylistImagePath", file, mock.IsType(*mockPlaylist)).
		Run(func(args mock.Arguments) {
			newPlaylist := args.Get(1).(model.Playlist)
			assert.Equal(t, newPlaylist.ID, mockPlaylist.ID)
			assert.WithinDuration(t, time.Now(), newPlaylist.UpdatedAt, time.Minute)
		}).
		Return(imagePath).
		Once()

	storageService.On("Upload", file, imagePath).Return(nil).Once()

	var newPlaylist *model.Playlist
	playlistRepository.On("Update", mock.IsType(new(model.Playlist))).
		Run(func(args mock.Arguments) {
			newPlaylist = args.Get(0).(*model.Playlist)
			assert.Equal(t, imagePath, string(*newPlaylist.ImageURL))
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
	errCode := _uut.Handle(file, id)

	// then
	assert.Nil(t, errCode)

	playlistRepository.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
	storageService.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
}
