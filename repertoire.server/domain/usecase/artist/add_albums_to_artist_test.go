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

func TestAddAlbumsToArtist_WhenGetAlbumsWithSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := AddAlbumsToArtist{albumRepository: albumRepository}

	request := requests.AddAlbumsToArtistRequest{
		ID:       uuid.New(),
		AlbumIDs: []uuid.UUID{uuid.New()},
	}

	internalError := errors.New("internal error")
	albumRepository.On("GetAllByIDsWithSongs", mock.Anything, request.AlbumIDs).
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

func TestAddAlbumsToArtist_WhenOneAlbumAlreadyHasArtist_ShouldReturnBadRequestError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := AddAlbumsToArtist{albumRepository: albumRepository}

	request := requests.AddAlbumsToArtistRequest{
		ID:       uuid.New(),
		AlbumIDs: []uuid.UUID{uuid.New()},
	}

	albums := &[]model.Album{
		{
			ID:       request.AlbumIDs[0],
			ArtistID: &[]uuid.UUID{uuid.New()}[0],
		},
	}
	albumRepository.On("GetAllByIDsWithSongs", mock.IsType(albums), request.AlbumIDs).
		Return(nil, albums).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusBadRequest, errCode.Code)
	assert.Equal(t, "album "+request.AlbumIDs[0].String()+" already has an artist", errCode.Error.Error())

	albumRepository.AssertExpectations(t)
}

func TestAddAlbumsToArtist_WhenUpdateAllAlbumsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := AddAlbumsToArtist{albumRepository: albumRepository}

	request := requests.AddAlbumsToArtistRequest{
		ID:       uuid.New(),
		AlbumIDs: []uuid.UUID{uuid.New()},
	}

	albums := &[]model.Album{{ID: request.AlbumIDs[0]}}
	albumRepository.On("GetAllByIDsWithSongs", mock.IsType(albums), request.AlbumIDs).
		Return(nil, albums).
		Once()

	internalError := errors.New("internal error")
	albumRepository.On("UpdateAllWithAssociations", mock.IsType(albums)).
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

func TestAddAlbumsToArtist_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := AddAlbumsToArtist{albumRepository: albumRepository}

	request := requests.AddAlbumsToArtistRequest{
		ID:       uuid.New(),
		AlbumIDs: []uuid.UUID{uuid.New(), uuid.New()},
	}

	albums := &[]model.Album{
		{
			ID:       request.AlbumIDs[0],
			ArtistID: nil,
			Songs:    []model.Song{{}, {}, {}},
		},
		{
			ID:       request.AlbumIDs[1],
			ArtistID: nil,
		},
	}
	albumRepository.On("GetAllByIDsWithSongs", mock.IsType(albums), request.AlbumIDs).
		Return(nil, albums).
		Once()

	albumRepository.On("UpdateAllWithAssociations", mock.IsType(albums)).
		Run(func(args mock.Arguments) {
			newAlbums := args.Get(0).(*[]model.Album)
			for _, album := range *newAlbums {
				assert.Equal(t, &request.ID, album.ArtistID)
				for _, song := range album.Songs {
					assert.Equal(t, &request.ID, song.ArtistID)
				}
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
