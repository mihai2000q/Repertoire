package song

import (
	"errors"
	"mime/multipart"
	"net/http"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/domain/provider"
	"repertoire/server/model"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSaveImageToSong_WhenGetSongFails_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := SaveImageToSong{
		repository: songRepository,
	}

	file := new(multipart.FileHeader)
	id := uuid.New()

	// given - mocking
	internalError := errors.New("internal error")
	songRepository.On("Get", new(model.Song), id).Return(internalError).Once()

	// when
	errCode := _uut.Handle(file, id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestSaveImageToSong_WhenSongIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := SaveImageToSong{
		repository: songRepository,
	}

	file := new(multipart.FileHeader)
	id := uuid.New()

	// given - mocking
	songRepository.On("Get", new(model.Song), id).Return(nil).Once()

	// when
	errCode := _uut.Handle(file, id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song not found", errCode.Error.Error())

	songRepository.AssertExpectations(t)
}

func TestSaveImageToSong_WhenStorageUploadFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	storageService := new(service.StorageServiceMock)
	_uut := SaveImageToSong{
		repository:              songRepository,
		storageFilePathProvider: storageFilePathProvider,
		storageService:          storageService,
	}

	file := new(multipart.FileHeader)
	id := uuid.New()

	// given - mocking
	song := &model.Song{ID: id, ImageURL: nil}
	songRepository.On("Get", new(model.Song), id).Return(nil, song).Once()

	imagePath := "songs file path"
	storageFilePathProvider.On("GetSongImagePath", file, id).Return(imagePath).Once()

	internalError := errors.New("internal error")
	storageService.On("Upload", file, imagePath).Return(internalError).Once()

	// when
	errCode := _uut.Handle(file, id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
	storageService.AssertExpectations(t)
}

func TestSaveImageToSong_WhenUpdateSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	storageService := new(service.StorageServiceMock)
	_uut := SaveImageToSong{
		repository:              songRepository,
		storageFilePathProvider: storageFilePathProvider,
		storageService:          storageService,
	}

	file := new(multipart.FileHeader)
	id := uuid.New()

	// given - mocking
	song := &model.Song{ID: id, ImageURL: nil}
	songRepository.On("Get", new(model.Song), id).Return(nil, song).Once()

	imagePath := "songs file path"
	storageFilePathProvider.On("GetSongImagePath", file, id).Return(imagePath).Once()

	storageService.On("Upload", file, imagePath).Return(nil).Once()

	internalError := errors.New("internal error")
	songRepository.On("Update", mock.IsType(new(model.Song))).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(file, id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
	storageService.AssertExpectations(t)
}

func TestSaveImageToSong_WhenIsValid_ShouldNotReturnAnyError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	storageService := new(service.StorageServiceMock)
	_uut := SaveImageToSong{
		repository:              songRepository,
		storageFilePathProvider: storageFilePathProvider,
		storageService:          storageService,
	}

	file := new(multipart.FileHeader)
	id := uuid.New()

	// given - mocking
	song := &model.Song{ID: id, ImageURL: nil}
	songRepository.On("Get", new(model.Song), id).Return(nil, song).Once()

	imagePath := "songs file path"
	storageFilePathProvider.On("GetSongImagePath", file, id).Return(imagePath).Once()

	storageService.On("Upload", file, imagePath).Return(nil).Once()

	songRepository.On("Update", mock.IsType(new(model.Song))).
		Run(func(args mock.Arguments) {
			newSong := args.Get(0).(*model.Song)
			assert.Equal(t, imagePath, string(*newSong.ImageURL))
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(file, id)

	// then
	assert.Nil(t, errCode)

	songRepository.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
	storageService.AssertExpectations(t)
}
