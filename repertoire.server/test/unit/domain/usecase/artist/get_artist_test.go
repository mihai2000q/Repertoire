package artist

import (
	"errors"
	"net/http"
	"repertoire/server/domain/usecase/artist"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetArtist_WhenGetArtistFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := artist.NewGetArtist(artistRepository)

	id := uuid.New()

	internalError := errors.New("internal error")
	artistRepository.On("GetWithAssociations", new(model.Artist), id).
		Return(internalError).
		Once()

	// when
	resultArtist, errCode := _uut.Handle(id)

	// then
	assert.Empty(t, resultArtist)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	artistRepository.AssertExpectations(t)
}

func TestGetArtist_WhenArtistIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := artist.NewGetArtist(artistRepository)

	id := uuid.New()

	artistRepository.On("GetWithAssociations", new(model.Artist), id).
		Return(nil).
		Once()

	// when
	resultArtist, errCode := _uut.Handle(id)

	// then
	assert.Empty(t, resultArtist)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "artist not found", errCode.Error.Error())

	artistRepository.AssertExpectations(t)
}

func TestGetArtist_WhenSuccessful_ShouldReturnArtist(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := artist.NewGetArtist(artistRepository)

	id := uuid.New()

	expectedArtist := &model.Artist{
		ID:   id,
		Name: "Some Artist",
	}

	artistRepository.On("GetWithAssociations", new(model.Artist), id).
		Return(nil, expectedArtist).
		Once()

	// when
	resultArtist, errCode := _uut.Handle(id)

	// then
	assert.NotEmpty(t, resultArtist)
	assert.Equal(t, expectedArtist, &resultArtist)
	assert.Nil(t, errCode)

	artistRepository.AssertExpectations(t)
}
