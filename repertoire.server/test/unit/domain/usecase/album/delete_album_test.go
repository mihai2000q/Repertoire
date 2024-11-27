package album

import (
	"errors"
	"net/http"
	"repertoire/server/domain/usecase/album"
	"repertoire/server/test/unit/data/repository"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeleteAlbum_WhenDeleteAlbumFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := album.NewDeleteAlbum(albumRepository)

	id := uuid.New()

	internalError := errors.New("internal error")
	albumRepository.On("Delete", id).Return(internalError).Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	albumRepository.AssertExpectations(t)
}

func TestDeleteAlbum_WhenSuccessful_ShouldReturnAlbums(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := album.NewDeleteAlbum(albumRepository)

	id := uuid.New()

	albumRepository.On("Delete", id).Return(nil).Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.Nil(t, errCode)

	albumRepository.AssertExpectations(t)
}
