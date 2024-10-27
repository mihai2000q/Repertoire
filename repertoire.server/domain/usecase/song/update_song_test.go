package song

import (
	"errors"
	"net/http"
	"repertoire/api/request"
	"repertoire/data/repository"
	"repertoire/model"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateSong_WhenGetSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := &UpdateSong{
		repository: songRepository,
	}
	request := request.UpdateSongRequest{
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
	_uut := &UpdateSong{
		repository: songRepository,
	}
	request := request.UpdateSongRequest{
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
	_uut := &UpdateSong{
		repository: songRepository,
	}
	request := request.UpdateSongRequest{
		ID:    uuid.New(),
		Title: "New Song",
	}

	song := &model.Song{
		ID:    request.ID,
		Title: "Some Song",
	}

	songRepository.On("Get", new(model.Song), request.ID).Return(nil, song).Once()
	internalError := errors.New("internal error")
	songRepository.On("Update", mock.IsType(song)).
		Run(func(args mock.Arguments) {
			newSong := args.Get(0).(*model.Song)
			assert.Equal(t, request.Title, newSong.Title)
			assert.Equal(t, request.IsRecorded, newSong.IsRecorded)
			assert.Equal(t, request.Reharsals, newSong.Reharsals)
			assert.Equal(t, request.Bpm, newSong.Bpm)
			assert.Equal(t, request.SongsterrLink, newSong.SongsterrLink)
			assert.Equal(t, request.GuitarTuningID, newSong.GuitarTuningID)
		}).
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
	_uut := &UpdateSong{
		repository: songRepository,
	}
	request := request.UpdateSongRequest{
		ID:    uuid.New(),
		Title: "New Song",
	}

	song := &model.Song{
		ID:    request.ID,
		Title: "Some Song",
	}

	songRepository.On("Get", new(model.Song), request.ID).Return(nil, song).Once()
	songRepository.On("Update", mock.IsType(song)).
		Run(func(args mock.Arguments) {
			newSong := args.Get(0).(*model.Song)
			assert.Equal(t, request.Title, newSong.Title)
			assert.Equal(t, request.IsRecorded, newSong.IsRecorded)
			assert.Equal(t, request.Reharsals, newSong.Reharsals)
			assert.Equal(t, request.Bpm, newSong.Bpm)
			assert.Equal(t, request.SongsterrLink, newSong.SongsterrLink)
			assert.Equal(t, request.GuitarTuningID, newSong.GuitarTuningID)
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)

	songRepository.AssertExpectations(t)
}
