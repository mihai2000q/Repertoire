package album

import (
	"errors"
	"net/http"
	"repertoire/server/data/repository"
	"repertoire/server/model"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetAlbum_WhenGetAlbumFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := &GetAlbum{
		repository: albumRepository,
	}
	id := uuid.New()

	internalError := errors.New("internal error")
	albumRepository.On("GetWithAssociations", new(model.Album), id).
		Return(internalError).
		Once()

	// when
	album, errCode := _uut.Handle(id)

	// then
	assert.Empty(t, album)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	albumRepository.AssertExpectations(t)
}

func TestGetAlbum_WhenAlbumIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := &GetAlbum{
		repository: albumRepository,
	}
	id := uuid.New()

	albumRepository.On("GetWithAssociations", new(model.Album), id).
		Return(nil).
		Once()

	// when
	album, errCode := _uut.Handle(id)

	// then
	assert.Empty(t, album)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "album not found", errCode.Error.Error())

	albumRepository.AssertExpectations(t)
}

func TestGetAlbum_WhenSuccessful_ShouldReturnAlbum(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := &GetAlbum{
		repository: albumRepository,
	}
	id := uuid.New()

	expectedAlbum := &model.Album{
		ID:    id,
		Title: "Some Album",
	}

	albumRepository.On("GetWithAssociations", new(model.Album), id).
		Return(nil, expectedAlbum).
		Once()

	// when
	album, errCode := _uut.Handle(id)

	// then
	assert.NotEmpty(t, album)
	assert.Equal(t, expectedAlbum, &album)
	assert.Nil(t, errCode)

	albumRepository.AssertExpectations(t)
}
