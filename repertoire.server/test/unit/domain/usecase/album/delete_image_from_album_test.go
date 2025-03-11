package album

import (
	"errors"
	"net/http"
	"repertoire/server/domain/usecase/album"
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

func TestDeleteImageFromAlbum_WhenGetAlbumFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := album.NewDeleteImageFromAlbum(albumRepository, nil, nil)

	id := uuid.New()

	// given - mocking
	internalError := errors.New("internal error")
	albumRepository.On("Get", new(model.Album), id).Return(internalError).Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	albumRepository.AssertExpectations(t)
}

func TestDeleteImageFromAlbum_WhenAlbumIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := album.NewDeleteImageFromAlbum(albumRepository, nil, nil)

	id := uuid.New()

	// given - mocking
	albumRepository.On("Get", new(model.Album), id).Return(nil).Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "album not found", errCode.Error.Error())

	albumRepository.AssertExpectations(t)
}

func TestDeleteImageFromAlbum_WhenAlbumHasNoImage_ShouldReturnBadRequestError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := album.NewDeleteImageFromAlbum(albumRepository, nil, nil)

	id := uuid.New()

	// given - mocking
	mockAlbum := &model.Album{ID: id}
	albumRepository.On("Get", new(model.Album), id).Return(nil, mockAlbum).Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusBadRequest, errCode.Code)
	assert.Equal(t, "album does not have an image", errCode.Error.Error())

	albumRepository.AssertExpectations(t)
}

func TestDeleteImageFromAlbum_WhenDeleteImageFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	storageService := new(service.StorageServiceMock)
	_uut := album.NewDeleteImageFromAlbum(albumRepository, storageService, nil)

	id := uuid.New()

	// given - mocking
	mockAlbum := &model.Album{ID: id, ImageURL: &[]internal.FilePath{"This is some url"}[0]}
	albumRepository.On("Get", new(model.Album), id).Return(nil, mockAlbum).Once()

	internalError := wrapper.InternalServerError(errors.New("internal error"))
	storageService.On("DeleteFile", *mockAlbum.ImageURL).Return(internalError).Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, internalError, errCode)

	albumRepository.AssertExpectations(t)
	storageService.AssertExpectations(t)
}

func TestDeleteImageFromAlbum_WhenUpdateAlbumFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	storageService := new(service.StorageServiceMock)
	_uut := album.NewDeleteImageFromAlbum(albumRepository, storageService, nil)

	id := uuid.New()

	// given - mocking
	mockAlbum := &model.Album{ID: id, ImageURL: &[]internal.FilePath{"This is some url"}[0]}
	albumRepository.On("Get", new(model.Album), id).Return(nil, mockAlbum).Once()

	storageService.On("DeleteFile", *mockAlbum.ImageURL).Return(nil).Once()

	internalError := errors.New("internal error")
	albumRepository.On("Update", mock.IsType(mockAlbum)).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	albumRepository.AssertExpectations(t)
	storageService.AssertExpectations(t)
}

func TestDeleteImageFromAlbum_WhenPublishFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	storageService := new(service.StorageServiceMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := album.NewDeleteImageFromAlbum(albumRepository, storageService, messagePublisherService)

	id := uuid.New()

	// given - mocking
	mockAlbum := &model.Album{ID: id, ImageURL: &[]internal.FilePath{"This is some url"}[0]}
	albumRepository.On("Get", new(model.Album), id).Return(nil, mockAlbum).Once()

	storageService.On("DeleteFile", *mockAlbum.ImageURL).Return(nil).Once()

	albumRepository.On("Update", mock.IsType(mockAlbum)).
		Return(nil).
		Once()

	internalError := errors.New("internal error")
	messagePublisherService.On("Publish", topics.AlbumUpdatedTopic, mockAlbum.ID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	albumRepository.AssertExpectations(t)
	storageService.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
}

func TestDeleteImageFromAlbum_WhenIsValid_ShouldNotReturnAnyError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	storageService := new(service.StorageServiceMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := album.NewDeleteImageFromAlbum(albumRepository, storageService, messagePublisherService)

	id := uuid.New()

	// given - mocking
	mockAlbum := &model.Album{ID: id, ImageURL: &[]internal.FilePath{"This is some url"}[0]}
	albumRepository.On("Get", new(model.Album), id).Return(nil, mockAlbum).Once()

	storageService.On("DeleteFile", *mockAlbum.ImageURL).Return(nil).Once()

	albumRepository.On("Update", mock.IsType(mockAlbum)).
		Run(func(args mock.Arguments) {
			newAlbum := args.Get(0).(*model.Album)
			assert.Nil(t, newAlbum.ImageURL)
		}).
		Return(nil).
		Once()

	messagePublisherService.On("Publish", topics.AlbumUpdatedTopic, mockAlbum.ID).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.Nil(t, errCode)

	albumRepository.AssertExpectations(t)
	storageService.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
}
