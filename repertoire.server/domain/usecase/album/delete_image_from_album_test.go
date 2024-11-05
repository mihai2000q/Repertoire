package album

import (
	"errors"
	"net/http"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal"
	"repertoire/server/model"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteImageFromAlbum_WhenGetAlbumFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := DeleteImageFromAlbum{repository: albumRepository}

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

func TestDeleteImageFromAlbum_WhenGetAlbumFails_ShouldReturnNotFoundError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := DeleteImageFromAlbum{repository: albumRepository}

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

func TestDeleteImageFromAlbum_WhenDeleteImageFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	storageService := new(service.StorageServiceMock)
	_uut := DeleteImageFromAlbum{
		repository:     albumRepository,
		storageService: storageService,
	}

	id := uuid.New()

	// given - mocking
	album := &model.Album{ID: id, ImageURL: &[]internal.FilePath{"This is some url"}[0]}
	albumRepository.On("Get", new(model.Album), id).Return(nil, album).Once()

	internalError := errors.New("internal error")
	storageService.On("Delete", string(*album.ImageURL)).Return(internalError).Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	albumRepository.AssertExpectations(t)
	storageService.AssertExpectations(t)
}

func TestDeleteImageFromAlbum_WhenUpdateAlbumFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	storageService := new(service.StorageServiceMock)
	_uut := DeleteImageFromAlbum{
		repository:     albumRepository,
		storageService: storageService,
	}

	id := uuid.New()

	// given - mocking
	album := &model.Album{ID: id, ImageURL: &[]internal.FilePath{"This is some url"}[0]}
	albumRepository.On("Get", new(model.Album), id).Return(nil, album).Once()

	storageService.On("Delete", string(*album.ImageURL)).Return(nil).Once()

	internalError := errors.New("internal error")
	albumRepository.On("Update", mock.IsType(album)).
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

func TestDeleteImageFromAlbum_WhenIsValid_ShouldNotReturnAnyError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	storageService := new(service.StorageServiceMock)
	_uut := DeleteImageFromAlbum{
		repository:     albumRepository,
		storageService: storageService,
	}

	id := uuid.New()

	// given - mocking
	album := &model.Album{ID: id, ImageURL: &[]internal.FilePath{"This is some url"}[0]}
	albumRepository.On("Get", new(model.Album), id).Return(nil, album).Once()

	storageService.On("Delete", string(*album.ImageURL)).Return(nil).Once()

	albumRepository.On("Update", mock.IsType(album)).
		Run(func(args mock.Arguments) {
			newAlbum := args.Get(0).(*model.Album)
			assert.Nil(t, newAlbum.ImageURL)
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.Nil(t, errCode)

	albumRepository.AssertExpectations(t)
	storageService.AssertExpectations(t)
}
