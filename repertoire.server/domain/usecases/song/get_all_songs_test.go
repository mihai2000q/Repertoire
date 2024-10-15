package song

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"repertoire/api/requests"
	"repertoire/data/repository"
	"repertoire/models"
	"testing"
)

func TestGetAll_WhenGetSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := &GetAllSongs{
		repository: songRepository,
	}
	request := requests.GetSongsRequest{
		UserID: uuid.New(),
	}

	internalError := errors.New("internal error")
	songRepository.On("GetAllByUser", mock.Anything, request.UserID).
		Return(internalError).
		Once()

	// when
	songs, errCode := _uut.Handle(request)

	// then
	assert.Empty(t, songs)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestGetAll_WhenSuccessful_ShouldReturnSongs(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := &GetAllSongs{
		repository: songRepository,
	}
	request := requests.GetSongsRequest{
		UserID: uuid.New(),
	}

	expectedSongs := &[]models.Song{
		{Title: "Some Song"},
		{Title: "Some other Song"},
	}

	songRepository.On("GetAllByUser", mock.IsType(expectedSongs), request.UserID).
		Return(nil, expectedSongs).
		Once()

	// when
	songs, errCode := _uut.Handle(request)

	// then
	assert.Equal(t, expectedSongs, &songs)
	assert.Nil(t, errCode)

	songRepository.AssertExpectations(t)
}
