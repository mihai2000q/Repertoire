package album

import (
	"errors"
	"net/http"
	"repertoire/server/domain/usecase/album"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetAlbum_WhenGetAlbumFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := album.NewGetAlbum(albumRepository)

	id := uuid.New()

	internalError := errors.New("internal error")
	albumRepository.On("GetWithAssociations", new(model.Album), id).
		Return(internalError).
		Once()

	// when
	resultAlbum, errCode := _uut.Handle(id)

	// then
	assert.Empty(t, resultAlbum)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	albumRepository.AssertExpectations(t)
}

func TestGetAlbum_WhenAlbumIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := album.NewGetAlbum(albumRepository)

	id := uuid.New()

	albumRepository.On("GetWithAssociations", new(model.Album), id).
		Return(nil).
		Once()

	// when
	resultAlbum, errCode := _uut.Handle(id)

	// then
	assert.Empty(t, resultAlbum)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "album not found", errCode.Error.Error())

	albumRepository.AssertExpectations(t)
}

func TestGetAlbum_WhenSuccessful_ShouldReturnAlbum(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := album.NewGetAlbum(albumRepository)

	id := uuid.New()

	expectedAlbum := &model.Album{
		ID:    id,
		Title: "Some Album",
	}

	albumRepository.On("GetWithAssociations", new(model.Album), id).
		Return(nil, expectedAlbum).
		Once()

	// when
	resultAlbum, errCode := _uut.Handle(id)

	// then
	assert.NotEmpty(t, resultAlbum)
	assert.Equal(t, expectedAlbum, &resultAlbum)
	assert.Nil(t, errCode)

	albumRepository.AssertExpectations(t)
}
