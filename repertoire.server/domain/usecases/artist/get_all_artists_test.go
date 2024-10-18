package artist

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

func TestGetAll_WhenGetArtistsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := &GetAllArtists{
		repository: artistRepository,
	}
	request := requests.GetArtistsRequest{
		UserID: uuid.New(),
	}

	internalError := errors.New("internal error")
	artistRepository.
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

	artistRepository.AssertExpectations(t)
}

func TestGetAll_WhenGetArtistsCountFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := &GetAllArtists{
		repository: artistRepository,
	}
	request := requests.GetArtistsRequest{
		UserID: uuid.New(),
	}

	expectedArtists := &[]models.Artist{
		{Name: "Some Artist"},
		{Name: "Some other Artist"},
	}

	artistRepository.
		On(
			"GetAllByUser",
			mock.IsType(expectedArtists),
			request.UserID,
			request.CurrentPage,
			request.PageSize,
		).
		Return(nil, expectedArtists).
		Once()

	internalError := errors.New("internal error")
	artistRepository.
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
	assert.Equal(t, expectedArtists, &result.Data)
	assert.Empty(t, result.TotalCount)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	artistRepository.AssertExpectations(t)
}

func TestGetAll_WhenSuccessful_ShouldReturnArtistsWithTotalCount(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := &GetAllArtists{
		repository: artistRepository,
	}
	request := requests.GetArtistsRequest{
		UserID: uuid.New(),
	}

	expectedArtists := &[]models.Artist{
		{Name: "Some Artist"},
		{Name: "Some other Artist"},
	}
	expectedTotalCount := &[]int64{20}[0]

	artistRepository.
		On(
			"GetAllByUser",
			mock.IsType(expectedArtists),
			request.UserID,
			request.CurrentPage,
			request.PageSize,
		).
		Return(nil, expectedArtists).
		Once()

	artistRepository.
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
	assert.Equal(t, expectedArtists, &result.Data)
	assert.Equal(t, expectedTotalCount, &result.TotalCount)
	assert.Nil(t, errCode)

	artistRepository.AssertExpectations(t)
}
