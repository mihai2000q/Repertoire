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
)

func TestGetAllSongArrangements_WhenGetSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songArrangementRepository := new(repository.SongArrangementRepositoryMock)
	_uut := arrangement.NewGetAllSongArrangements(songArrangementRepository)

	request := requests.GetSongArrangementsRequest{
		SongID: uuid.New(),
	}

	internalError := errors.New("internal error")
	songArrangementRepository.On("GetAllBySong", new([]model.SongArrangement), request.SongID).
		Return(internalError).
		Once()

	// when
	resultArrangements, errCode := _uut.Handle(request)

	// then
	assert.Empty(t, resultArrangements)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songArrangementRepository.AssertExpectations(t)
}

func TestGetAllSongArrangements_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	// given
	songArrangementRepository := new(repository.SongArrangementRepositoryMock)
	_uut := arrangement.NewGetAllSongArrangements(songArrangementRepository)

	request := requests.GetSongArrangementsRequest{
		SongID: uuid.New(),
	}

	mockArrangements := []model.SongArrangement{
		{ID: uuid.New(), SongID: request.SongID, Order: 0},
		{ID: uuid.New(), SongID: request.SongID, Order: 1},
		{ID: uuid.New(), SongID: request.SongID, Order: 2},
	}
	songArrangementRepository.On("GetAllBySong", new([]model.SongArrangement), request.SongID).
		Return(nil, &mockArrangements).
		Once()

	// when
	resultArrangements, errCode := _uut.Handle(request)

	// then
	assert.Equal(t, mockArrangements, resultArrangements)
	assert.Nil(t, errCode)

	songArrangementRepository.AssertExpectations(t)
}
