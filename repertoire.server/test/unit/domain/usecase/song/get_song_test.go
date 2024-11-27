package song

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"repertoire/server/domain/usecase/song"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"testing"
)

func TestGetSong_WhenGetSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewGetSong(songRepository)

	id := uuid.New()

	internalError := errors.New("internal error")
	songRepository.On("GetWithAssociations", new(model.Song), id).
		Return(internalError).
		Once()

	// when
	resultSong, errCode := _uut.Handle(id)

	// then
	assert.Empty(t, resultSong)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestGetSong_WhenSongIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewGetSong(songRepository)

	id := uuid.New()

	songRepository.On("GetWithAssociations", new(model.Song), id).
		Return(nil).
		Once()

	// when
	resultSong, errCode := _uut.Handle(id)

	// then
	assert.Empty(t, resultSong)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song not found", errCode.Error.Error())

	songRepository.AssertExpectations(t)
}

func TestGetSong_WhenSuccessful_ShouldReturnSong(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewGetSong(songRepository)

	id := uuid.New()

	expectedSong := &model.Song{
		ID:    id,
		Title: "Some Song",
	}

	songRepository.On("GetWithAssociations", new(model.Song), id).
		Return(nil, expectedSong).
		Once()

	// when
	resultSong, errCode := _uut.Handle(id)

	// then
	assert.NotEmpty(t, resultSong)
	assert.Equal(t, expectedSong, &resultSong)
	assert.Nil(t, errCode)

	songRepository.AssertExpectations(t)
}
