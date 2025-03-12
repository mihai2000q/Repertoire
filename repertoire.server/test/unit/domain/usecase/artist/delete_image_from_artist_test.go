package artist

import (
	"errors"
	"net/http"
	"repertoire/server/domain/usecase/artist"
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

func TestDeleteImageFromArtist_WhenGetArtistFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := artist.NewDeleteImageFromArtist(artistRepository, nil, nil)

	id := uuid.New()

	// given - mocking
	internalError := errors.New("internal error")
	artistRepository.On("Get", new(model.Artist), id).Return(internalError).Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	artistRepository.AssertExpectations(t)
}

func TestDeleteImageFromArtist_WhenArtistIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := artist.NewDeleteImageFromArtist(artistRepository, nil, nil)

	id := uuid.New()

	// given - mocking
	artistRepository.On("Get", new(model.Artist), id).Return(nil).Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "artist not found", errCode.Error.Error())

	artistRepository.AssertExpectations(t)
}

func TestDeleteImageFromArtist_WhenArtistHasNoImage_ShouldReturnBadRequestError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := artist.NewDeleteImageFromArtist(artistRepository, nil, nil)

	id := uuid.New()

	// given - mocking
	mockArtist := &model.Artist{ID: id}
	artistRepository.On("Get", new(model.Artist), id).Return(nil, mockArtist).Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusBadRequest, errCode.Code)
	assert.Equal(t, "artist does not have an image", errCode.Error.Error())

	artistRepository.AssertExpectations(t)
}

func TestDeleteImageFromArtist_WhenDeleteImageFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	storageService := new(service.StorageServiceMock)
	_uut := artist.NewDeleteImageFromArtist(artistRepository, storageService, nil)

	id := uuid.New()

	// given - mocking
	mockArtist := &model.Artist{ID: id, ImageURL: &[]internal.FilePath{"This is some url"}[0]}
	artistRepository.On("Get", new(model.Artist), id).Return(nil, mockArtist).Once()

	internalError := wrapper.InternalServerError(errors.New("internal error"))
	storageService.On("DeleteFile", *mockArtist.ImageURL).Return(internalError).Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode)

	artistRepository.AssertExpectations(t)
	storageService.AssertExpectations(t)
}

func TestDeleteImageFromArtist_WhenUpdateArtistFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	storageService := new(service.StorageServiceMock)
	_uut := artist.NewDeleteImageFromArtist(artistRepository, storageService, nil)

	id := uuid.New()

	// given - mocking
	mockArtist := &model.Artist{ID: id, ImageURL: &[]internal.FilePath{"This is some url"}[0]}
	artistRepository.On("Get", new(model.Artist), id).Return(nil, mockArtist).Once()

	storageService.On("DeleteFile", *mockArtist.ImageURL).Return(nil).Once()

	internalError := errors.New("internal error")
	artistRepository.On("Update", mock.IsType(mockArtist)).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	artistRepository.AssertExpectations(t)
	storageService.AssertExpectations(t)
}

func TestDeleteImageFromArtist_WhenPublishFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	storageService := new(service.StorageServiceMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := artist.NewDeleteImageFromArtist(artistRepository, storageService, messagePublisherService)

	id := uuid.New()

	// given - mocking
	mockArtist := &model.Artist{ID: id, ImageURL: &[]internal.FilePath{"This is some url"}[0]}
	artistRepository.On("Get", new(model.Artist), id).Return(nil, mockArtist).Once()

	storageService.On("DeleteFile", *mockArtist.ImageURL).Return(nil).Once()

	artistRepository.On("Update", mock.IsType(mockArtist)).
		Return(nil).
		Once()

	internalError := errors.New("internal error")
	messagePublisherService.On("Publish", topics.ArtistUpdatedTopic, mockArtist.ID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	artistRepository.AssertExpectations(t)
	storageService.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
}

func TestDeleteImageFromArtist_WhenIsValid_ShouldNotReturnAnyError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	storageService := new(service.StorageServiceMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := artist.NewDeleteImageFromArtist(artistRepository, storageService, messagePublisherService)

	id := uuid.New()

	// given - mocking
	mockArtist := &model.Artist{ID: id, ImageURL: &[]internal.FilePath{"This is some url"}[0]}
	artistRepository.On("Get", new(model.Artist), id).Return(nil, mockArtist).Once()

	storageService.On("DeleteFile", *mockArtist.ImageURL).Return(nil).Once()

	artistRepository.On("Update", mock.IsType(mockArtist)).
		Run(func(args mock.Arguments) {
			newArtist := args.Get(0).(*model.Artist)
			assert.Nil(t, newArtist.ImageURL)
		}).
		Return(nil).
		Once()

	messagePublisherService.On("Publish", topics.ArtistUpdatedTopic, mockArtist.ID).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.Nil(t, errCode)

	artistRepository.AssertExpectations(t)
	storageService.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
}
