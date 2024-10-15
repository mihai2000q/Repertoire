package playlist

import (
	"errors"
	"net/http"
	"repertoire/data/repository"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeletePlaylist_WhenDeletePlaylistFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := &DeletePlaylist{
		repository: playlistRepository,
	}
	id := uuid.New()

	internalError := errors.New("internal error")
	playlistRepository.On("Delete", id).Return(internalError).Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	playlistRepository.AssertExpectations(t)
}

func TestDeletePlaylist_WhenSuccessful_ShouldReturnPlaylists(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := &DeletePlaylist{
		repository: playlistRepository,
	}
	id := uuid.New()

	playlistRepository.On("Delete", id).Return(nil).Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.Nil(t, errCode)

	playlistRepository.AssertExpectations(t)
}
