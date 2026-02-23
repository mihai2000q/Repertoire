package arrangement

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/song/arrangement"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateDefaultSongArrangement_WhenGetSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := arrangement.NewUpdateDefaultSongArrangement(songRepository)

	request := requests.UpdateDefaultSongArrangementRequest{
		ID:     uuid.New(),
		SongID: uuid.New(),
	}

	internalError := errors.New("internal error")
	songRepository.On("Get", new(model.Song), request.SongID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestUpdateDefaultSongArrangement_WhenSongIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := arrangement.NewUpdateDefaultSongArrangement(songRepository)

	request := requests.UpdateDefaultSongArrangementRequest{
		ID:     uuid.New(),
		SongID: uuid.New(),
	}

	songRepository.On("Get", new(model.Song), request.SongID).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song not found", errCode.Error.Error())

	songRepository.AssertExpectations(t)
}

func TestUpdateDefaultSongArrangement_WhenUpdateSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := arrangement.NewUpdateDefaultSongArrangement(songRepository)

	request := requests.UpdateDefaultSongArrangementRequest{
		ID:     uuid.New(),
		SongID: uuid.New(),
	}

	mockSong := model.Song{
		ID: request.SongID,
	}
	songRepository.On("Get", new(model.Song), request.SongID).
		Return(nil, &mockSong).
		Once()

	internalError := errors.New("internal error")
	songRepository.On("Update", mock.IsType(new(model.Song))).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestUpdateDefaultSongArrangement_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := arrangement.NewUpdateDefaultSongArrangement(songRepository)

	request := requests.UpdateDefaultSongArrangementRequest{
		ID:     uuid.New(),
		SongID: uuid.New(),
	}

	mockSong := model.Song{
		ID: request.SongID,
	}
	songRepository.On("Get", new(model.Song), request.SongID).
		Return(nil, &mockSong).
		Once()

	songRepository.On("Update", mock.IsType(new(model.Song))).
		Run(func(args mock.Arguments) {
			newSong := *args.Get(0).(*model.Song)
			assert.Equal(t, request.ID, *newSong.DefaultArrangementID)
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)

	songRepository.AssertExpectations(t)
}
