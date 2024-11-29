package song

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"repertoire/server/domain/usecase/song"
	"repertoire/server/test/unit/data/repository"
	"testing"
)

func TestDeleteSong_WhenDeleteSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewDeleteSong(songRepository)

	id := uuid.New()

	internalError := errors.New("internal error")
	songRepository.On("Delete", id).Return(internalError).Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestDeleteSong_WhenSuccessful_ShouldDeleteSong(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewDeleteSong(songRepository)

	id := uuid.New()

	songRepository.On("Delete", id).Return(nil).Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.Nil(t, errCode)

	songRepository.AssertExpectations(t)
}
