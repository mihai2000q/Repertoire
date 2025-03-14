package song

import (
	"errors"
	"net/http"
	"repertoire/server/domain/usecase/song"
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

func TestDeleteImageFromSong_WhenGetSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewDeleteImageFromSong(songRepository, nil, nil)

	id := uuid.New()

	// given - mocking
	internalError := errors.New("internal error")
	songRepository.On("Get", new(model.Song), id).Return(internalError).Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestDeleteImageFromSong_WhenSongIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewDeleteImageFromSong(songRepository, nil, nil)

	id := uuid.New()

	// given - mocking
	songRepository.On("Get", new(model.Song), id).Return(nil).Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song not found", errCode.Error.Error())

	songRepository.AssertExpectations(t)
}

func TestDeleteImageFromSong_WhenSongHasNoImage_ShouldReturnBadRequestError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewDeleteImageFromSong(songRepository, nil, nil)

	id := uuid.New()

	// given - mocking
	mockSong := &model.Song{ID: id}
	songRepository.On("Get", new(model.Song), id).Return(nil, mockSong).Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusBadRequest, errCode.Code)
	assert.Equal(t, "song does not have an image", errCode.Error.Error())

	songRepository.AssertExpectations(t)
}

func TestDeleteImageFromSong_WhenDeleteImageFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	storageService := new(service.StorageServiceMock)
	_uut := song.NewDeleteImageFromSong(songRepository, storageService, nil)

	id := uuid.New()

	// given - mocking
	mockSong := &model.Song{ID: id, ImageURL: &[]internal.FilePath{"This is some url"}[0]}
	songRepository.On("Get", new(model.Song), id).Return(nil, mockSong).Once()

	internalError := wrapper.InternalServerError(errors.New("internal error"))
	storageService.On("DeleteFile", *mockSong.ImageURL).Return(internalError).Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, internalError, errCode)

	songRepository.AssertExpectations(t)
	storageService.AssertExpectations(t)
}

func TestDeleteImageFromSong_WhenUpdateSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	storageService := new(service.StorageServiceMock)
	_uut := song.NewDeleteImageFromSong(songRepository, storageService, nil)

	id := uuid.New()

	// given - mocking
	mockSong := &model.Song{ID: id, ImageURL: &[]internal.FilePath{"This is some url"}[0]}
	songRepository.On("Get", new(model.Song), id).Return(nil, mockSong).Once()

	storageService.On("DeleteFile", *mockSong.ImageURL).Return(nil).Once()

	internalError := errors.New("internal error")
	songRepository.On("Update", mock.IsType(mockSong)).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
	storageService.AssertExpectations(t)
}

func TestDeleteImageFromSong_WhenPublishFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	storageService := new(service.StorageServiceMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := song.NewDeleteImageFromSong(songRepository, storageService, messagePublisherService)

	id := uuid.New()

	// given - mocking
	mockSong := &model.Song{ID: id, ImageURL: &[]internal.FilePath{"This is some url"}[0]}
	songRepository.On("Get", new(model.Song), id).Return(nil, mockSong).Once()

	storageService.On("DeleteFile", *mockSong.ImageURL).Return(nil).Once()

	songRepository.On("Update", mock.IsType(mockSong)).
		Return(nil).
		Once()

	internalError := errors.New("internal error")
	messagePublisherService.On("Publish", topics.SongsUpdatedTopic, []uuid.UUID{mockSong.ID}).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
	storageService.AssertExpectations(t)
	storageService.AssertExpectations(t)
}

func TestDeleteImageFromSong_WhenIsValid_ShouldNotReturnAnyError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	storageService := new(service.StorageServiceMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := song.NewDeleteImageFromSong(songRepository, storageService, messagePublisherService)

	id := uuid.New()

	// given - mocking
	mockSong := &model.Song{ID: id, ImageURL: &[]internal.FilePath{"This is some url"}[0]}
	songRepository.On("Get", new(model.Song), id).Return(nil, mockSong).Once()

	storageService.On("DeleteFile", *mockSong.ImageURL).Return(nil).Once()

	songRepository.On("Update", mock.IsType(mockSong)).
		Run(func(args mock.Arguments) {
			newSong := args.Get(0).(*model.Song)
			assert.Nil(t, newSong.ImageURL)
		}).
		Return(nil).
		Once()

	messagePublisherService.On("Publish", topics.SongsUpdatedTopic, []uuid.UUID{mockSong.ID}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.Nil(t, errCode)

	songRepository.AssertExpectations(t)
	storageService.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
}
