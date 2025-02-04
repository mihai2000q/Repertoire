package artist

import (
	"errors"
	"mime/multipart"
	"net/http"
	"repertoire/server/domain/usecase/artist"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"repertoire/server/test/unit/domain/provider"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSaveImageToArtist_WhenGetArtistFails_ShouldReturnNotFoundError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := artist.NewSaveImageToArtist(artistRepository, nil, nil)

	file := new(multipart.FileHeader)
	id := uuid.New()

	// given - mocking
	internalError := errors.New("internal error")
	artistRepository.On("Get", new(model.Artist), id).Return(internalError).Once()

	// when
	errCode := _uut.Handle(file, id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	artistRepository.AssertExpectations(t)
}

func TestSaveImageToArtist_WhenArtistIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := artist.NewSaveImageToArtist(artistRepository, nil, nil)

	file := new(multipart.FileHeader)
	id := uuid.New()

	// given - mocking
	artistRepository.On("Get", new(model.Artist), id).Return(nil).Once()

	// when
	errCode := _uut.Handle(file, id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "artist not found", errCode.Error.Error())

	artistRepository.AssertExpectations(t)
}

func TestSaveImageToArtist_WhenStorageUploadFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	storageService := new(service.StorageServiceMock)
	_uut := artist.NewSaveImageToArtist(artistRepository, storageFilePathProvider, storageService)

	file := new(multipart.FileHeader)
	id := uuid.New()

	// given - mocking
	mockArtist := &model.Artist{ID: id, ImageURL: nil}
	artistRepository.On("Get", new(model.Artist), id).Return(nil, mockArtist).Once()

	imagePath := "artists file path"
	storageFilePathProvider.On("GetArtistImagePath", file, *mockArtist).Return(imagePath).Once()

	internalError := errors.New("internal error")
	storageService.On("Upload", file, imagePath).Return(internalError).Once()

	// when
	errCode := _uut.Handle(file, id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	artistRepository.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
	storageService.AssertExpectations(t)
}

func TestSaveImageToArtist_WhenUpdateArtistFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	storageService := new(service.StorageServiceMock)
	_uut := artist.NewSaveImageToArtist(artistRepository, storageFilePathProvider, storageService)

	file := new(multipart.FileHeader)
	id := uuid.New()

	// given - mocking
	mockArtist := &model.Artist{ID: id, ImageURL: nil}
	artistRepository.On("Get", new(model.Artist), id).Return(nil, mockArtist).Once()

	imagePath := "artists file path"
	storageFilePathProvider.On("GetArtistImagePath", file, *mockArtist).Return(imagePath).Once()

	storageService.On("Upload", file, imagePath).Return(nil).Once()

	internalError := errors.New("internal error")
	artistRepository.On("Update", mock.IsType(new(model.Artist))).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(file, id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	artistRepository.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
	storageService.AssertExpectations(t)
}

func TestSaveImageToArtist_WhenIsValid_ShouldNotReturnAnyError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	storageService := new(service.StorageServiceMock)
	_uut := artist.NewSaveImageToArtist(artistRepository, storageFilePathProvider, storageService)

	file := new(multipart.FileHeader)
	id := uuid.New()

	// given - mocking
	mockArtist := &model.Artist{ID: id, ImageURL: nil}
	artistRepository.On("Get", new(model.Artist), id).Return(nil, mockArtist).Once()

	imagePath := "artists file path"
	storageFilePathProvider.On("GetArtistImagePath", file, *mockArtist).Return(imagePath).Once()

	storageService.On("Upload", file, imagePath).Return(nil).Once()

	artistRepository.On("Update", mock.IsType(new(model.Artist))).
		Run(func(args mock.Arguments) {
			newArtist := args.Get(0).(*model.Artist)
			assert.Equal(t, imagePath, string(*newArtist.ImageURL))
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(file, id)

	// then
	assert.Nil(t, errCode)

	artistRepository.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
	storageService.AssertExpectations(t)
}
