package artist

import (
	"errors"
	"net/http"
	artist2 "repertoire/server/domain/usecase/artist"
	"repertoire/server/test/unit/data/repository"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeleteArtist_WhenDeleteArtistFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := &artist2.DeleteArtist{
		repository: artistRepository,
	}
	id := uuid.New()

	internalError := errors.New("internal error")
	artistRepository.On("Delete", id).Return(internalError).Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	artistRepository.AssertExpectations(t)
}

func TestDeleteArtist_WhenSuccessful_ShouldReturnArtists(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := &artist2.DeleteArtist{
		repository: artistRepository,
	}
	id := uuid.New()

	artistRepository.On("Delete", id).Return(nil).Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.Nil(t, errCode)

	artistRepository.AssertExpectations(t)
}
