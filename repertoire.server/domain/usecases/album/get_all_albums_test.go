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

func TestGetAll_WhenGetAlbumFails_ShouldReturnInternalServerError(t *testing.T) {
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
	albums, errCode := _uut.Handle(request)

	// then
	assert.Empty(t, albums)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	albumRepository.AssertExpectations(t)
}

func TestGetAll_WhenSuccessful_ShouldReturnAlbums(t *testing.T) {
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
			mock.Anything,
			request.UserID,
			request.CurrentPage,
			request.PageSize,
		).
		Return(nil, expectedAlbums).
		Once()

	// when
	albums, errCode := _uut.Handle(request)

	// then
	assert.Equal(t, expectedAlbums, &albums)
	assert.Nil(t, errCode)

	albumRepository.AssertExpectations(t)
}
