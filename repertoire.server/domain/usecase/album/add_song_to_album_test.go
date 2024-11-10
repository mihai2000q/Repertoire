package album

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/model"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddSongToAlbum_WhenGetSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := AddSongToAlbum{songRepository: songRepository}

	request := requests.AddSongToAlbumRequest{
		ID:     uuid.New(),
		SongID: uuid.New(),
	}

	// given - mocking
	internalError := errors.New("internal error")
	songRepository.On("Get", mock.IsType(new(model.Song)), request.SongID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestAddSongToAlbum_WhenSongIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := AddSongToAlbum{songRepository: songRepository}

	request := requests.AddSongToAlbumRequest{
		ID:     uuid.New(),
		SongID: uuid.New(),
	}

	// given - mocking
	songRepository.On("Get", mock.IsType(new(model.Song)), request.SongID).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song not found", errCode.Error.Error())

	songRepository.AssertExpectations(t)
}

func TestAddSongToAlbum_WhenGetAlbumWithSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := AddSongToAlbum{
		repository:     albumRepository,
		songRepository: songRepository,
	}

	request := requests.AddSongToAlbumRequest{
		ID:     uuid.New(),
		SongID: uuid.New(),
	}

	// given - mocking
	song := &model.Song{ID: request.SongID}
	songRepository.On("Get", mock.IsType(song), request.SongID).
		Return(nil, song).
		Once()

	internalError := errors.New("internal error")
	albumRepository.On("GetWithSongs", mock.Anything, request.ID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	albumRepository.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}

func TestAddSongToAlbum_WhenAlbumIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := AddSongToAlbum{
		repository:     albumRepository,
		songRepository: songRepository,
	}

	request := requests.AddSongToAlbumRequest{
		ID:     uuid.New(),
		SongID: uuid.New(),
	}

	// given - mocking
	song := &model.Song{ID: request.SongID}
	songRepository.On("Get", mock.IsType(song), request.SongID).
		Return(nil, song).
		Once()

	albumRepository.On("GetWithSongs", mock.Anything, request.ID).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "album not found", errCode.Error.Error())

	albumRepository.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}

func TestAddSongToAlbum_WhenAlbumAndSongArtistDoesNotMatch_ShouldReturnBadRequestError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := AddSongToAlbum{
		repository:     albumRepository,
		songRepository: songRepository,
	}

	request := requests.AddSongToAlbumRequest{
		ID:     uuid.New(),
		SongID: uuid.New(),
	}

	// given - mocking
	song := &model.Song{ID: request.SongID}
	songRepository.On("Get", mock.IsType(song), request.SongID).
		Return(nil, song).
		Once()

	album := &model.Album{ArtistID: &[]uuid.UUID{uuid.New()}[0]}
	albumRepository.On("GetWithSongs", mock.IsType(album), request.ID).
		Return(nil, album).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusBadRequest, errCode.Code)
	assert.Equal(t, "song and album do not share the same artist", errCode.Error.Error())

	albumRepository.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}

func TestAddSongToAlbum_WhenUpdateSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := AddSongToAlbum{
		repository:     albumRepository,
		songRepository: songRepository,
	}

	request := requests.AddSongToAlbumRequest{
		ID:     uuid.New(),
		SongID: uuid.New(),
	}

	// given - mocking
	song := &model.Song{ID: request.SongID}
	songRepository.On("Get", mock.IsType(song), request.SongID).
		Return(nil, song).
		Once()

	album := &model.Album{ArtistID: &[]uuid.UUID{uuid.New()}[0]}
	albumRepository.On("GetWithSongs", mock.IsType(album), request.ID).
		Return(nil, album).
		Once()

	internalError := errors.New("internal error")
	songRepository.On("Update", mock.IsType(song)).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	albumRepository.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}

func TestAddSongToAlbum_WhenIsValid_ShouldNotReturnAnyError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := AddSongToAlbum{
		repository:     albumRepository,
		songRepository: songRepository,
	}

	request := requests.AddSongToAlbumRequest{
		ID:     uuid.New(),
		SongID: uuid.New(),
	}

	// given - mocking
	song := &model.Song{ID: request.SongID}
	songRepository.On("Get", mock.IsType(song), request.SongID).
		Return(nil, song).
		Once()

	album := &model.Album{
		ID:       request.ID,
		Songs:    []model.Song{{}, {}, {}},
		ArtistID: song.ArtistID,
	}
	albumRepository.On("GetWithSongs", mock.IsType(album), request.ID).
		Return(nil, album).
		Once()

	songRepository.On("Update", mock.IsType(song)).
		Run(func(args mock.Arguments) {
			newSong := args.Get(0).(*model.Song)
			assert.Equal(t, *newSong.AlbumID, request.ID)
			assert.Equal(t, *newSong.AlbumTrackNo, uint(len(album.Songs))+1)
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)

	albumRepository.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}
