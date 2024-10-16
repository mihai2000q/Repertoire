package artist

import (
	"errors"
	"net/http"
	"repertoire/api/requests"
	"repertoire/data/repository"
	"repertoire/models"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateArtist_WhenGetArtistFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := &UpdateArtist{
		repository: artistRepository,
	}
	request := requests.UpdateArtistRequest{
		ID:   uuid.New(),
		Name: "New Artist",
	}

	internalError := errors.New("internal error")
	artistRepository.On("Get", new(models.Artist), request.ID).Return(internalError).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	artistRepository.AssertExpectations(t)
}

func TestUpdateArtist_WhenArtistIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := &UpdateArtist{
		repository: artistRepository,
	}
	request := requests.UpdateArtistRequest{
		ID:   uuid.New(),
		Name: "New Artist",
	}

	artistRepository.On("Get", new(models.Artist), request.ID).Return(nil).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "artist not found", errCode.Error.Error())

	artistRepository.AssertExpectations(t)
}

func TestUpdateArtist_WhenUpdateArtistFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := &UpdateArtist{
		repository: artistRepository,
	}
	request := requests.UpdateArtistRequest{
		ID:   uuid.New(),
		Name: "New Artist",
	}

	artist := &models.Artist{
		ID:   request.ID,
		Name: "Some Artist",
	}

	artistRepository.On("Get", new(models.Artist), request.ID).Return(nil, artist).Once()
	internalError := errors.New("internal error")
	artistRepository.On("Update", mock.IsType(artist)).
		Run(func(args mock.Arguments) {
			newArtist := args.Get(0).(*models.Artist)
			assert.Equal(t, request.Name, newArtist.Name)
		}).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	artistRepository.AssertExpectations(t)
}

func TestUpdateArtist_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := &UpdateArtist{
		repository: artistRepository,
	}
	request := requests.UpdateArtistRequest{
		ID:   uuid.New(),
		Name: "New Artist",
	}

	artist := &models.Artist{
		ID:   request.ID,
		Name: "Some Artist",
	}

	artistRepository.On("Get", new(models.Artist), request.ID).Return(nil, artist).Once()
	artistRepository.On("Update", mock.IsType(artist)).
		Run(func(args mock.Arguments) {
			newArtist := args.Get(0).(*models.Artist)
			assert.Equal(t, request.Name, newArtist.Name)
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)

	artistRepository.AssertExpectations(t)
}
