package song

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/song"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateSong_WhenGetSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewUpdateSong(songRepository)

	request := requests.UpdateSongRequest{
		ID:    uuid.New(),
		Title: "New Song",
	}

	internalError := errors.New("internal error")
	songRepository.On("Get", new(model.Song), request.ID).Return(internalError).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestUpdateSong_WhenSongIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewUpdateSong(songRepository)

	request := requests.UpdateSongRequest{
		ID:    uuid.New(),
		Title: "New Song",
	}

	songRepository.On("Get", new(model.Song), request.ID).Return(nil).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song not found", errCode.Error.Error())

	songRepository.AssertExpectations(t)
}

func TestUpdateSong_WhenUpdateSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewUpdateSong(songRepository)

	request := requests.UpdateSongRequest{
		ID:    uuid.New(),
		Title: "New Song",
	}

	mockSong := &model.Song{
		ID:    request.ID,
		Title: "Some Song",
	}
	songRepository.On("Get", new(model.Song), request.ID).
		Return(nil, mockSong).
		Once()

	internalError := errors.New("internal error")
	songRepository.On("Update", mock.IsType(mockSong)).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestUpdateSong_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewUpdateSong(songRepository)

	request := requests.UpdateSongRequest{
		ID:    uuid.New(),
		Title: "New Song",
	}

	mockSong := &model.Song{
		ID:    request.ID,
		Title: "Some Song",
	}
	songRepository.On("Get", new(model.Song), request.ID).
		Return(nil, mockSong).
		Once()

	songRepository.On("Update", mock.IsType(mockSong)).
		Run(func(args mock.Arguments) {
			newSong := args.Get(0).(*model.Song)
			assertUpdatedSong(t, request, newSong)
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)

	songRepository.AssertExpectations(t)
}

func assertUpdatedSong(t *testing.T, request requests.UpdateSongRequest, song *model.Song) {
	assert.Equal(t, request.Title, song.Title)
	assert.Equal(t, request.Description, song.Description)
	assert.Equal(t, request.IsRecorded, song.IsRecorded)
	assert.Equal(t, request.Bpm, song.Bpm)
	assert.Equal(t, request.SongsterrLink, song.SongsterrLink)
	assert.Equal(t, request.YoutubeLink, song.YoutubeLink)
	assert.Equal(t, request.ReleaseDate, song.ReleaseDate)
	assert.Equal(t, request.Difficulty, song.Difficulty)
	assert.Equal(t, request.GuitarTuningID, song.GuitarTuningID)
}
