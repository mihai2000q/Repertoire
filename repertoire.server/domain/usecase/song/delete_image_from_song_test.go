package song

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

func TestDeleteImageFromSong_WhenGetSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := DeleteImageFromSong{repository: songRepository}

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

func TestDeleteImageFromSong_WhenGetSongFails_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := DeleteImageFromSong{repository: songRepository}

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

func TestDeleteImageFromSong_WhenDeleteImageFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	storageService := new(service.StorageServiceMock)
	_uut := DeleteImageFromSong{
		repository:     songRepository,
		storageService: storageService,
	}

	id := uuid.New()

	// given - mocking
	song := &model.Song{ID: id, ImageURL: &[]internal.FilePath{"This is some url"}[0]}
	songRepository.On("Get", new(model.Song), id).Return(nil, song).Once()

	internalError := errors.New("internal error")
	storageService.On("Delete", string(*song.ImageURL)).Return(internalError).Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
	storageService.AssertExpectations(t)
}

func TestDeleteImageFromSong_WhenUpdateSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	storageService := new(service.StorageServiceMock)
	_uut := DeleteImageFromSong{
		repository:     songRepository,
		storageService: storageService,
	}

	id := uuid.New()

	// given - mocking
	song := &model.Song{ID: id, ImageURL: &[]internal.FilePath{"This is some url"}[0]}
	songRepository.On("Get", new(model.Song), id).Return(nil, song).Once()

	storageService.On("Delete", string(*song.ImageURL)).Return(nil).Once()

	internalError := errors.New("internal error")
	songRepository.On("Update", mock.IsType(song)).
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

func TestDeleteImageFromSong_WhenIsValid_ShouldNotReturnAnyError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	storageService := new(service.StorageServiceMock)
	_uut := DeleteImageFromSong{
		repository:     songRepository,
		storageService: storageService,
	}

	id := uuid.New()

	// given - mocking
	song := &model.Song{ID: id, ImageURL: &[]internal.FilePath{"This is some url"}[0]}
	songRepository.On("Get", new(model.Song), id).Return(nil, song).Once()

	storageService.On("Delete", string(*song.ImageURL)).Return(nil).Once()

	songRepository.On("Update", mock.IsType(song)).
		Run(func(args mock.Arguments) {
			newSong := args.Get(0).(*model.Song)
			assert.Nil(t, newSong.ImageURL)
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.Nil(t, errCode)

	songRepository.AssertExpectations(t)
	storageService.AssertExpectations(t)
}
