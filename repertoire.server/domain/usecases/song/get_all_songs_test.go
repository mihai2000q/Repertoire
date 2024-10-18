package song

import (
	"errors"
	"net/http"
	"repertoire/api/requests"
	"repertoire/data/repository"
	"repertoire/models"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAll_WhenGetSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := &GetAllSongs{
		repository: songRepository,
	}
	request := requests.GetSongsRequest{
		UserID: uuid.New(),
	}

	internalError := errors.New("internal error")
	songRepository.
		On(
			"GetAllByUser",
			mock.Anything,
			request.UserID,
			request.CurrentPage,
			request.PageSize,
		).
		Return(internalError).
		Once()

	// when
	result, errCode := _uut.Handle(request)

	// then
	assert.Empty(t, result)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestGetAll_WhenGetSongsCountFails_ShouldReturnInternalServerError(t *testing.T) {
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

	songRepository.
		On(
			"GetAllByUser",
			mock.IsType(expectedSongs),
			request.UserID,
			request.CurrentPage,
			request.PageSize,
		).
		Return(nil, expectedSongs).
		Once()

	internalError := errors.New("internal error")
	songRepository.
		On(
			"GetAllByUserCount",
			mock.Anything,
			request.UserID,
		).
		Return(internalError).
		Once()

	// when
	result, errCode := _uut.Handle(request)

	// then
	assert.Equal(t, expectedSongs, &result.Data)
	assert.Empty(t, result.TotalCount)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestGetAll_WhenSuccessful_ShouldReturnSongsWithTotalCount(t *testing.T) {
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
	expectedTotalCount := &[]int64{20}[0]

	songRepository.
		On(
			"GetAllByUser",
			mock.IsType(expectedSongs),
			request.UserID,
			request.CurrentPage,
			request.PageSize,
		).
		Return(nil, expectedSongs).
		Once()

	songRepository.
		On(
			"GetAllByUserCount",
			mock.IsType(expectedTotalCount),
			request.UserID,
		).
		Return(nil, expectedTotalCount).
		Once()

	// when
	result, errCode := _uut.Handle(request)

	// then
	assert.Equal(t, expectedSongs, &result.Data)
	assert.Equal(t, expectedTotalCount, &result.TotalCount)
	assert.Nil(t, errCode)

	songRepository.AssertExpectations(t)
}
