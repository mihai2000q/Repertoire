package part

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/song/part"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAllSongParts_WhenGetAllBySongFails_ShouldReturnDatabaseError(t *testing.T) {
	// given
	songPartRepository := new(repository.SongPartRepositoryMock)
	_uut := part.NewGetAllSongParts(songPartRepository)

	request := requests.GetSongPartsRequest{
		SongID: uuid.New(),
	}

	internalError := errors.New("internal error")
	songPartRepository.On("GetAllBySong", new([]model.SongPart), request.SongID).
		Return(internalError).
		Once()

	// when
	parts, errCode := _uut.Handle(request)

	// then
	assert.Empty(t, parts)
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songPartRepository.AssertExpectations(t)
}

func TestGetAllSongParts_WhenSuccessful_ShouldReturnParts(t *testing.T) {
	// given
	songPartRepository := new(repository.SongPartRepositoryMock)
	_uut := part.NewGetAllSongParts(songPartRepository)

	request := requests.GetSongPartsRequest{
		SongID: uuid.New(),
	}

	expectedParts := []model.SongPart{
		{ID: uuid.New(), Name: "Verse Riff"},
		{ID: uuid.New(), Name: "Chorus Riff"},
		{ID: uuid.New(), Name: "Solo"},
	}

	songPartRepository.On("GetAllBySong", new([]model.SongPart), request.SongID).
		Return(nil, &expectedParts).
		Once()

	// when
	parts, errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)
	assert.Equal(t, expectedParts, parts)

	songPartRepository.AssertExpectations(t)
}
