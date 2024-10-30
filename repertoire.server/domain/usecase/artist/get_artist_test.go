package artist

import (
	"errors"
	"net/http"
	"repertoire/server/data/repository"
	"repertoire/server/model"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetArtistQuery_WhenGetArtistFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := &GetArtist{
		repository: artistRepository,
	}
	id := uuid.New()

	internalError := errors.New("internal error")
	artistRepository.On("GetWithAssociations", new(model.Artist), id).
		Return(internalError).
		Once()

	// when
	artist, errCode := _uut.Handle(id)

	// then
	assert.Empty(t, artist)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	artistRepository.AssertExpectations(t)
}

func TestGetArtistQuery_WhenArtistIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := &GetArtist{
		repository: artistRepository,
	}
	id := uuid.New()

	artistRepository.On("GetWithAssociations", new(model.Artist), id).
		Return(nil).
		Once()

	// when
	artist, errCode := _uut.Handle(id)

	// then
	assert.Empty(t, artist)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "artist not found", errCode.Error.Error())

	artistRepository.AssertExpectations(t)
}

func TestGetArtistQuery_WhenSuccessful_ShouldReturnArtist(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := &GetArtist{
		repository: artistRepository,
	}
	id := uuid.New()

	expectedArtist := &model.Artist{
		ID:   id,
		Name: "Some Artist",
	}

	artistRepository.On("GetWithAssociations", new(model.Artist), id).
		Return(nil, expectedArtist).
		Once()

	// when
	artist, errCode := _uut.Handle(id)

	// then
	assert.NotEmpty(t, artist)
	assert.Equal(t, expectedArtist, &artist)
	assert.Nil(t, errCode)

	artistRepository.AssertExpectations(t)
}
