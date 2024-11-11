package artist

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"repertoire/server/data/repository"
	"repertoire/server/model"
	"testing"
)

func TestRemoveSongFromArtist_WhenGetSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := RemoveSongFromArtist{songRepository: songRepository}

	id := uuid.New()
	songID := uuid.New()

	internalError := errors.New("internal error")
	songRepository.On("Get", mock.IsType(new(model.Song)), songID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id, songID)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestRemoveSongFromArtist_WhenSongIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := RemoveSongFromArtist{songRepository: songRepository}

	id := uuid.New()
	songID := uuid.New()

	songRepository.On("Get", mock.IsType(new(model.Song)), songID).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(id, songID)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song not found", errCode.Error.Error())

	songRepository.AssertExpectations(t)
}

func TestRemoveSongFromArtist_WhenSongArtistDoesNotMatch_ShouldReturnBadRequestError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := RemoveSongFromArtist{songRepository: songRepository}

	id := uuid.New()
	songID := uuid.New()

	song := &model.Song{ID: songID}
	songRepository.On("Get", mock.IsType(new(model.Song)), songID).
		Return(nil, song).
		Once()

	// when
	errCode := _uut.Handle(id, songID)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusBadRequest, errCode.Code)
	assert.Equal(t, "song is not owned by this artist", errCode.Error.Error())

	songRepository.AssertExpectations(t)
}

func TestRemoveSongFromArtist_WhenUpdateSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := RemoveSongFromArtist{songRepository: songRepository}

	id := uuid.New()
	songID := uuid.New()

	song := &model.Song{
		ID:       songID,
		ArtistID: &id,
	}
	songRepository.On("Get", mock.IsType(song), songID).
		Return(nil, song).
		Once()

	internalError := errors.New("internal error")
	songRepository.On("Update", mock.IsType(song)).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id, songID)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestRemoveSongFromArtist_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := RemoveSongFromArtist{songRepository: songRepository}

	id := uuid.New()
	songID := uuid.New()

	song := &model.Song{
		ID:       songID,
		ArtistID: &id,
	}
	songRepository.On("Get", mock.IsType(song), songID).
		Return(nil, song).
		Once()

	songRepository.On("Update", mock.IsType(song)).
		Run(func(args mock.Arguments) {
			newSong := args.Get(0).(*model.Song)
			assert.Nil(t, newSong.ArtistID)
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(id, songID)

	// then
	assert.Nil(t, errCode)

	songRepository.AssertExpectations(t)
}
