package artist

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/model"
	"testing"
)

func TestAddAlbumToArtist_WhenGetAlbumWithSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := AddAlbumToArtist{albumRepository: albumRepository}

	request := requests.AddAlbumToArtistRequest{
		ID:      uuid.New(),
		AlbumID: uuid.New(),
	}

	internalError := errors.New("internal error")
	albumRepository.On("GetWithSongs", mock.IsType(new(model.Album)), request.AlbumID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	albumRepository.AssertExpectations(t)
}

func TestAddAlbumToArtist_WhenAlbumIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := AddAlbumToArtist{albumRepository: albumRepository}

	request := requests.AddAlbumToArtistRequest{
		ID:      uuid.New(),
		AlbumID: uuid.New(),
	}

	albumRepository.On("GetWithSongs", mock.IsType(new(model.Album)), request.AlbumID).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "album not found", errCode.Error.Error())

	albumRepository.AssertExpectations(t)
}

func TestAddAlbumToArtist_WhenAlbumHasArtist_ShouldReturnBadRequestError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := AddAlbumToArtist{albumRepository: albumRepository}

	request := requests.AddAlbumToArtistRequest{
		ID:      uuid.New(),
		AlbumID: uuid.New(),
	}

	album := &model.Album{ArtistID: &[]uuid.UUID{uuid.New()}[0]}
	albumRepository.On("GetWithSongs", mock.IsType(album), request.AlbumID).
		Return(nil, album).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusBadRequest, errCode.Code)
	assert.Equal(t, "album already has an artist", errCode.Error.Error())

	albumRepository.AssertExpectations(t)
}

func TestAddAlbumToArtist_WhenUpdateWithAssociationsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := AddAlbumToArtist{albumRepository: albumRepository}

	request := requests.AddAlbumToArtistRequest{
		ID:      uuid.New(),
		AlbumID: uuid.New(),
	}

	album := &model.Album{ID: request.AlbumID}
	albumRepository.On("GetWithSongs", mock.IsType(album), request.AlbumID).
		Return(nil, album).
		Once()

	internalError := errors.New("internal error")
	albumRepository.On("UpdateWithAssociations", mock.IsType(album)).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	albumRepository.AssertExpectations(t)
}

func TestAddAlbumToArtist_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := AddAlbumToArtist{albumRepository: albumRepository}

	request := requests.AddAlbumToArtistRequest{
		ID:      uuid.New(),
		AlbumID: uuid.New(),
	}

	album := &model.Album{
		ID:       request.AlbumID,
		ArtistID: nil,
		Songs:    []model.Song{{}, {}, {}},
	}
	albumRepository.On("GetWithSongs", mock.IsType(album), request.AlbumID).
		Return(nil, album).
		Once()
	albumRepository.On("UpdateWithAssociations", mock.IsType(album)).
		Run(func(args mock.Arguments) {
			newAlbum := args.Get(0).(*model.Album)
			assert.Equal(t, request.ID, *newAlbum.ArtistID)
			for _, song := range newAlbum.Songs {
				assert.Equal(t, request.ID, *song.ArtistID)
			}
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)

	albumRepository.AssertExpectations(t)
}
