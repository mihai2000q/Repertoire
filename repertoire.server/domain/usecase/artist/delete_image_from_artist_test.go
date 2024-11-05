package artist

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

func TestDeleteImageFromArtist_WhenGetArtistFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := DeleteImageFromArtist{repository: artistRepository}

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

func TestDeleteImageFromArtist_WhenGetArtistFails_ShouldReturnNotFoundError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := DeleteImageFromArtist{repository: artistRepository}

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

func TestDeleteImageFromArtist_WhenDeleteImageFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	storageService := new(service.StorageServiceMock)
	_uut := DeleteImageFromArtist{
		repository:     artistRepository,
		storageService: storageService,
	}

	id := uuid.New()

	// given - mocking
	artist := &model.Artist{ID: id, ImageURL: &[]internal.FilePath{"This is some url"}[0]}
	artistRepository.On("Get", new(model.Artist), id).Return(nil, artist).Once()

	internalError := errors.New("internal error")
	storageService.On("Delete", string(*artist.ImageURL)).Return(internalError).Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	artistRepository.AssertExpectations(t)
	storageService.AssertExpectations(t)
}

func TestDeleteImageFromArtist_WhenUpdateArtistFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	storageService := new(service.StorageServiceMock)
	_uut := DeleteImageFromArtist{
		repository:     artistRepository,
		storageService: storageService,
	}

	id := uuid.New()

	// given - mocking
	artist := &model.Artist{ID: id, ImageURL: &[]internal.FilePath{"This is some url"}[0]}
	artistRepository.On("Get", new(model.Artist), id).Return(nil, artist).Once()

	storageService.On("Delete", string(*artist.ImageURL)).Return(nil).Once()

	internalError := errors.New("internal error")
	artistRepository.On("Update", mock.IsType(artist)).
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

func TestDeleteImageFromArtist_WhenIsValid_ShouldNotReturnAnyError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	storageService := new(service.StorageServiceMock)
	_uut := DeleteImageFromArtist{
		repository:     artistRepository,
		storageService: storageService,
	}

	id := uuid.New()

	// given - mocking
	artist := &model.Artist{ID: id, ImageURL: &[]internal.FilePath{"This is some url"}[0]}
	artistRepository.On("Get", new(model.Artist), id).Return(nil, artist).Once()

	storageService.On("Delete", string(*artist.ImageURL)).Return(nil).Once()

	artistRepository.On("Update", mock.IsType(artist)).
		Run(func(args mock.Arguments) {
			newArtist := args.Get(0).(*model.Artist)
			assert.Nil(t, newArtist.ImageURL)
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.Nil(t, errCode)

	artistRepository.AssertExpectations(t)
	storageService.AssertExpectations(t)
}
