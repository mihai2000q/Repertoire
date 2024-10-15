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

func TestGetAll_WhenGetArtistFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := &GetAllArtists{
		repository: artistRepository,
	}
	request := requests.GetArtistsRequest{
		UserID: uuid.New(),
	}

	internalError := errors.New("internal error")
	artistRepository.On("GetAllByUser", mock.Anything, request.UserID).
		Return(internalError).
		Once()

	// when
	artists, errCode := _uut.Handle(request)

	// then
	assert.Empty(t, artists)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	artistRepository.AssertExpectations(t)
}

func TestGetAll_WhenSuccessful_ShouldReturnArtists(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := &GetAllArtists{
		repository: artistRepository,
	}
	request := requests.GetArtistsRequest{
		UserID: uuid.New(),
	}

	expectedArtists := &[]models.Artist{
		{Title: "Some Artist"},
		{Title: "Some other Artist"},
	}

	artistRepository.On("GetAllByUser", mock.IsType(expectedArtists), request.UserID).
		Return(nil, expectedArtists).
		Once()

	// when
	artists, errCode := _uut.Handle(request)

	// then
	assert.Equal(t, expectedArtists, &artists)
	assert.Nil(t, errCode)

	artistRepository.AssertExpectations(t)
}
