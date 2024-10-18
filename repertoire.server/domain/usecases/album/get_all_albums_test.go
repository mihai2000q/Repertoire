package album

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

func TestGetAll_WhenGetAlbumsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := &GetAllAlbums{
		repository: albumRepository,
	}
	request := requests.GetAlbumsRequest{
		UserID: uuid.New(),
	}

	internalError := errors.New("internal error")
	albumRepository.
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
	res, errCode := _uut.Handle(request)

	// then
	assert.Empty(t, res)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	albumRepository.AssertExpectations(t)
}

func TestGetAll_WhenGetAlbumsCountFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := &GetAllAlbums{
		repository: albumRepository,
	}
	request := requests.GetAlbumsRequest{
		UserID: uuid.New(),
	}

	expectedAlbums := &[]models.Album{
		{Title: "Some Album"},
		{Title: "Some other Album"},
	}

	albumRepository.
		On(
			"GetAllByUser",
			mock.IsType(expectedAlbums),
			request.UserID,
			request.CurrentPage,
			request.PageSize,
		).
		Return(nil, expectedAlbums).
		Once()

	internalError := errors.New("internal error")
	albumRepository.
		On(
			"GetAllByUserCount",
			mock.Anything,
			request.UserID,
		).
		Return(internalError).
		Once()

	// when
	res, errCode := _uut.Handle(request)

	// then
	assert.Equal(t, expectedAlbums, &res.Data)
	assert.Empty(t, res.TotalCount)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	albumRepository.AssertExpectations(t)
}

func TestGetAll_WhenSuccessful_ShouldReturnAlbumsWithTotalCount(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := &GetAllAlbums{
		repository: albumRepository,
	}
	request := requests.GetAlbumsRequest{
		UserID: uuid.New(),
	}

	expectedAlbums := &[]models.Album{
		{Title: "Some Album"},
		{Title: "Some other Album"},
	}
	expectedTotalCount := &[]int64{20}[0]

	albumRepository.
		On(
			"GetAllByUser",
			mock.IsType(expectedAlbums),
			request.UserID,
			request.CurrentPage,
			request.PageSize,
		).
		Return(nil, expectedAlbums).
		Once()

	albumRepository.
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
	assert.Equal(t, expectedAlbums, &result.Data)
	assert.Equal(t, expectedTotalCount, &result.TotalCount)
	assert.Nil(t, errCode)

	albumRepository.AssertExpectations(t)
}
