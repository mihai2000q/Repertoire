package album

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/album"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddSongToAlbum_WhenGetAllSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := album.NewAddSongsToAlbum(nil, songRepository)

	request := requests.AddSongsToAlbumRequest{
		ID:      uuid.New(),
		SongIDs: []uuid.UUID{uuid.New()},
	}

	internalError := errors.New("internal error")
	songRepository.On("GetAllByIDs", mock.Anything, request.SongIDs).
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

func TestAddSongToAlbum_WhenGetAlbumWithSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := album.NewAddSongsToAlbum(albumRepository, songRepository)

	request := requests.AddSongsToAlbumRequest{
		ID:      uuid.New(),
		SongIDs: []uuid.UUID{uuid.New()},
	}

	// given - mocking
	songs := &[]model.Song{{ID: request.SongIDs[0]}}
	songRepository.On("GetAllByIDs", mock.IsType(songs), request.SongIDs).
		Return(nil, songs).
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
	_uut := album.NewAddSongsToAlbum(albumRepository, songRepository)

	request := requests.AddSongsToAlbumRequest{
		ID:      uuid.New(),
		SongIDs: []uuid.UUID{uuid.New()},
	}

	// given - mocking
	songs := &[]model.Song{{ID: request.SongIDs[0]}}
	songRepository.On("GetAllByIDs", mock.IsType(songs), request.SongIDs).
		Return(nil, songs).
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
	_uut := album.NewAddSongsToAlbum(albumRepository, songRepository)

	request := requests.AddSongsToAlbumRequest{
		ID:      uuid.New(),
		SongIDs: []uuid.UUID{uuid.New()},
	}

	// given - mocking
	songs := &[]model.Song{{ID: request.SongIDs[0], ArtistID: &[]uuid.UUID{uuid.New()}[0]}}
	songRepository.On("GetAllByIDs", mock.IsType(songs), request.SongIDs).
		Return(nil, songs).
		Once()

	mockAlbum := &model.Album{ArtistID: &[]uuid.UUID{uuid.New()}[0]}
	albumRepository.On("GetWithSongs", mock.IsType(mockAlbum), request.ID).
		Return(nil, mockAlbum).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusBadRequest, errCode.Code)
	assert.Equal(t, "song "+request.SongIDs[0].String()+" and album do not share the same artist", errCode.Error.Error())

	albumRepository.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}

func TestAddSongToAlbum_WhenOneSongHasAlbum_ShouldReturnBadRequestError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := album.NewAddSongsToAlbum(albumRepository, songRepository)

	request := requests.AddSongsToAlbumRequest{
		ID:      uuid.New(),
		SongIDs: []uuid.UUID{uuid.New()},
	}

	songs := &[]model.Song{{ID: request.SongIDs[0], AlbumID: &[]uuid.UUID{uuid.New()}[0]}}
	songRepository.On("GetAllByIDs", mock.IsType(songs), request.SongIDs).
		Return(nil, songs).
		Once()

	mockAlbum := &model.Album{ArtistID: &[]uuid.UUID{uuid.New()}[0]}
	albumRepository.On("GetWithSongs", mock.IsType(mockAlbum), request.ID).
		Return(nil, mockAlbum).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusBadRequest, errCode.Code)
	assert.Equal(t, "song "+request.SongIDs[0].String()+" already has an album", errCode.Error.Error())

	albumRepository.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}

func TestAddSongToAlbum_WhenUpdateSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := album.NewAddSongsToAlbum(albumRepository, songRepository)

	request := requests.AddSongsToAlbumRequest{
		ID:      uuid.New(),
		SongIDs: []uuid.UUID{uuid.New()},
	}

	// given - mocking
	songs := &[]model.Song{{ID: request.SongIDs[0]}}
	songRepository.On("GetAllByIDs", mock.IsType(songs), request.SongIDs).
		Return(nil, songs).
		Once()

	mockAlbum := &model.Album{ID: request.ID}
	albumRepository.On("GetWithSongs", mock.IsType(mockAlbum), request.ID).
		Return(nil, mockAlbum).
		Once()

	internalError := errors.New("internal error")
	songRepository.On("UpdateAll", mock.IsType(songs)).
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
		request requests.AddSongsToAlbumRequest
		songs   *[]model.Song
		album   *model.Album
	}{
		{
			"when they bot got no artist",
			requests.AddSongsToAlbumRequest{
				ID:      id,
				SongIDs: []uuid.UUID{songID},
			},
			&[]model.Song{{ID: songID}},
			&model.Album{
				ID:    id,
				Songs: []model.Song{{}, {}, {}},
			},
		},
		{
			"when the album has an artist, but the song does not",
			requests.AddSongsToAlbumRequest{
				ID:      id,
				SongIDs: []uuid.UUID{songID},
			},
			&[]model.Song{{ID: songID}},
			&model.Album{
				ID:       id,
				Songs:    []model.Song{{}, {}, {}},
				ArtistID: &artistID,
			},
		},
		{
			"when the song has an artist, but the album does not",
			requests.AddSongsToAlbumRequest{
				ID:      id,
				SongIDs: []uuid.UUID{songID},
			},
			&[]model.Song{{
				ID:       songID,
				ArtistID: &artistID,
			}},
			&model.Album{
				ID:    id,
				Songs: []model.Song{{}, {}, {}},
			},
		},
		{
			"when both have the same artist",
			requests.AddSongsToAlbumRequest{
				ID:      id,
				SongIDs: []uuid.UUID{songID},
			},
			&[]model.Song{{
				ID:       songID,
				ArtistID: &artistID,
			}},
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
			_uut := album.NewAddSongsToAlbum(albumRepository, songRepository)

			// given - mocking
			songRepository.On("GetAllByIDs", mock.IsType(tt.songs), tt.request.SongIDs).
				Return(nil, tt.songs).
				Once()

			albumRepository.On("GetWithSongs", mock.IsType(tt.album), tt.request.ID).
				Return(nil, tt.album).
				Once()

			songRepository.On("UpdateAll", mock.IsType(tt.songs)).
				Run(func(args mock.Arguments) {
					newSongs := args.Get(0).(*[]model.Song)
					for i, song := range *newSongs {
						assert.Equal(t, *song.AlbumID, tt.request.ID)
						assert.Equal(t, *song.AlbumTrackNo, uint(len(tt.album.Songs)+i)+1)
						assert.Equal(t, song.ArtistID, tt.album.ArtistID)
					}
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
