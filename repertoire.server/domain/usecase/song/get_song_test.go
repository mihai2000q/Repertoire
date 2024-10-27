package song

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"repertoire/data/repository"
	"repertoire/model"
	"testing"
)

func TestGetSong_WhenGetSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := &GetSong{
		repository: songRepository,
	}
	id := uuid.New()

	internalError := errors.New("internal error")
	songRepository.On("Get", new(model.Song), id).Return(internalError).Once()

	// when
	song, errCode := _uut.Handle(id)

	// then
	assert.Empty(t, song)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestGetSong_WhenSongIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := &GetSong{
		repository: songRepository,
	}
	id := uuid.New()

	songRepository.On("Get", new(model.Song), id).Return(nil).Once()

	// when
	song, errCode := _uut.Handle(id)

	// then
	assert.Empty(t, song)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song not found", errCode.Error.Error())

	songRepository.AssertExpectations(t)
}

func TestGetSong_WhenSuccessful_ShouldReturnSong(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := &GetSong{
		repository: songRepository,
	}
	id := uuid.New()

	expectedSong := &model.Song{
		ID:    id,
		Title: "Some Song",
	}

	songRepository.On("Get", new(model.Song), id).Return(nil, expectedSong).Once()

	// when
	song, errCode := _uut.Handle(id)

	// then
	assert.NotEmpty(t, song)
	assert.Equal(t, expectedSong, &song)
	assert.Nil(t, errCode)

	songRepository.AssertExpectations(t)
}
