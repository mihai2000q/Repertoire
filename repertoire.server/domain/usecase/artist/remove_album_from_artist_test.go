package artist

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"repertoire/server/data/repository"
	"repertoire/server/model"
	"testing"
)

func TestRemoveAlbumFromArtist_WhenGetAlbumFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := RemoveAlbumFromArtist{albumRepository: albumRepository}

	id := uuid.New()
	albumID := uuid.New()

	internalError := errors.New("internal error")
	albumRepository.On("GetWithSongs", mock.IsType(new(model.Album)), albumID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id, albumID)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	albumRepository.AssertExpectations(t)
}

func TestRemoveAlbumFromArtist_WhenAlbumIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := RemoveAlbumFromArtist{albumRepository: albumRepository}

	id := uuid.New()
	albumID := uuid.New()

	albumRepository.On("GetWithSongs", mock.IsType(new(model.Album)), albumID).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(id, albumID)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "album not found", errCode.Error.Error())

	albumRepository.AssertExpectations(t)
}

func TestRemoveAlbumFromArtist_WhenAlbumArtistDoesNotMatch_ShouldReturnBadRequestError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := RemoveAlbumFromArtist{albumRepository: albumRepository}

	id := uuid.New()
	albumID := uuid.New()

	album := &model.Album{ID: albumID}
	albumRepository.On("GetWithSongs", mock.IsType(new(model.Album)), albumID).
		Return(nil, album).
		Once()

	// when
	errCode := _uut.Handle(id, albumID)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusBadRequest, errCode.Code)
	assert.Equal(t, "album is not owned by this artist", errCode.Error.Error())

	albumRepository.AssertExpectations(t)
}

func TestRemoveAlbumFromArtist_WhenUpdateAlbumFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := RemoveAlbumFromArtist{albumRepository: albumRepository}

	id := uuid.New()
	albumID := uuid.New()

	album := &model.Album{
		ID:       albumID,
		ArtistID: &id,
	}
	albumRepository.On("GetWithSongs", mock.IsType(album), albumID).
		Return(nil, album).
		Once()

	internalError := errors.New("internal error")
	albumRepository.On("UpdateWithAssociations", mock.IsType(album)).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id, albumID)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	albumRepository.AssertExpectations(t)
}

func TestRemoveAlbumFromArtist_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := RemoveAlbumFromArtist{albumRepository: albumRepository}

	id := uuid.New()
	albumID := uuid.New()

	album := &model.Album{
		ID:       albumID,
		ArtistID: &id,
		Songs:    []model.Song{{ArtistID: &id}, {ArtistID: &id}, {ArtistID: &id}},
	}
	albumRepository.On("GetWithSongs", mock.IsType(album), albumID).
		Return(nil, album).
		Once()

	albumRepository.On("UpdateWithAssociations", mock.IsType(album)).
		Run(func(args mock.Arguments) {
			newAlbum := args.Get(0).(*model.Album)
			assert.Nil(t, newAlbum.ArtistID)
			for _, song := range newAlbum.Songs {
				assert.Nil(t, song.ArtistID)
			}
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(id, albumID)

	// then
	assert.Nil(t, errCode)

	albumRepository.AssertExpectations(t)
}
