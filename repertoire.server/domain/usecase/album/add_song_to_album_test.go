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

func TestAddSongToAlbum_WhenSongHasAlbum_ShouldReturnBadRequestError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := AddSongToAlbum{songRepository: songRepository}

	request := requests.AddSongToAlbumRequest{
		ID:     uuid.New(),
		SongID: uuid.New(),
	}

	song := &model.Song{AlbumID: &[]uuid.UUID{uuid.New()}[0]}
	songRepository.On("Get", mock.IsType(song), request.SongID).
		Return(nil, song).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusBadRequest, errCode.Code)
	assert.Equal(t, "song already has an album", errCode.Error.Error())

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
	song := &model.Song{ID: request.SongID, ArtistID: &[]uuid.UUID{uuid.New()}[0]}
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

	album := &model.Album{ID: request.ID}
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
	id := uuid.New()
	songID := uuid.New()
	artistID := uuid.New()

	tests := []struct {
		name    string
		request requests.AddSongToAlbumRequest
		song    *model.Song
		album   *model.Album
	}{
		{
			"when they bot got no artist",
			requests.AddSongToAlbumRequest{
				ID:     id,
				SongID: songID,
			},
			&model.Song{ID: songID},
			&model.Album{
				ID:    id,
				Songs: []model.Song{{}, {}, {}},
			},
		},
		{
			"when the album has an artist, but the song does not",
			requests.AddSongToAlbumRequest{
				ID:     id,
				SongID: songID,
			},
			&model.Song{ID: songID},
			&model.Album{
				ID:       id,
				Songs:    []model.Song{{}, {}, {}},
				ArtistID: &artistID,
			},
		},
		{
			"when the song has an artist, but the album does not",
			requests.AddSongToAlbumRequest{
				ID:     id,
				SongID: songID,
			},
			&model.Song{
				ID:       songID,
				ArtistID: &artistID,
			},
			&model.Album{
				ID:    id,
				Songs: []model.Song{{}, {}, {}},
			},
		},
		{
			"when both have the same artist",
			requests.AddSongToAlbumRequest{
				ID:     id,
				SongID: songID,
			},
			&model.Song{
				ID:       songID,
				ArtistID: &artistID,
			},
			&model.Album{
				ID:       id,
				Songs:    []model.Song{{}, {}, {}},
				ArtistID: &artistID,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			albumRepository := new(repository.AlbumRepositoryMock)
			songRepository := new(repository.SongRepositoryMock)
			_uut := AddSongToAlbum{
				repository:     albumRepository,
				songRepository: songRepository,
			}

			// given - mocking
			songRepository.On("Get", mock.IsType(tt.song), tt.request.SongID).
				Return(nil, tt.song).
				Once()

			albumRepository.On("GetWithSongs", mock.IsType(tt.album), tt.request.ID).
				Return(nil, tt.album).
				Once()

			songRepository.On("Update", mock.IsType(tt.song)).
				Run(func(args mock.Arguments) {
					newSong := args.Get(0).(*model.Song)
					assert.Equal(t, *newSong.AlbumID, tt.request.ID)
					assert.Equal(t, *newSong.AlbumTrackNo, uint(len(tt.album.Songs))+1)
					assert.Equal(t, newSong.ArtistID, tt.album.ArtistID)
				}).
				Return(nil).
				Once()

			// when
			errCode := _uut.Handle(tt.request)

			// then
			assert.Nil(t, errCode)

			albumRepository.AssertExpectations(t)
			songRepository.AssertExpectations(t)
		})
	}
}
