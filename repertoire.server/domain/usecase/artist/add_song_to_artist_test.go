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

func TestAddSongToArtist_WhenGetSongWithSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := AddSongToArtist{songRepository: songRepository}

	request := requests.AddSongToArtistRequest{
		ID:     uuid.New(),
		SongID: uuid.New(),
	}

	internalError := errors.New("internal error")
	songRepository.On("GetWithSongs", mock.IsType(new(model.Song)), request.SongID).
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

func TestAddSongToArtist_WhenSongIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := AddSongToArtist{songRepository: songRepository}

	request := requests.AddSongToArtistRequest{
		ID:     uuid.New(),
		SongID: uuid.New(),
	}

	songRepository.On("GetWithSongs", mock.IsType(new(model.Song)), request.SongID).
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

func TestAddSongToArtist_WhenSongHasArtist_ShouldReturnBadRequestError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := AddSongToArtist{songRepository: songRepository}

	request := requests.AddSongToArtistRequest{
		ID:     uuid.New(),
		SongID: uuid.New(),
	}

	song := &model.Song{ArtistID: &[]uuid.UUID{uuid.New()}[0]}
	songRepository.On("GetWithSongs", mock.IsType(song), request.SongID).
		Return(nil, song).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusBadRequest, errCode.Code)
	assert.Equal(t, "song already has an artist", errCode.Error.Error())

	songRepository.AssertExpectations(t)
}

func TestAddSongToArtist_WhenUpdateWithAssociationsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := AddSongToArtist{songRepository: songRepository}

	request := requests.AddSongToArtistRequest{
		ID:     uuid.New(),
		SongID: uuid.New(),
	}

	song := &model.Song{ID: request.SongID}
	songRepository.On("GetWithSongs", mock.IsType(song), request.SongID).
		Return(nil, song).
		Once()

	internalError := errors.New("internal error")
	songRepository.On("UpdateWithAssociations", mock.IsType(song)).
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

func TestAddSongToArtist_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	songID := uuid.New()

	tests := []struct {
		name    string
		request requests.AddSongToArtistRequest
		song    *model.Song
	}{
		{
			"Song has album so it should update the whole album",
			requests.AddSongToArtistRequest{
				ID:     uuid.New(),
				SongID: songID,
			},
			&model.Song{
				ID: songID,
				Album: &model.Album{
					Songs: []model.Song{
						{ID: uuid.New()},
						{ID: songID},
						{ID: uuid.New()},
					},
				},
			},
		},
		{
			"Song does not have an album so it should update only this song",
			requests.AddSongToArtistRequest{
				ID:     uuid.New(),
				SongID: songID,
			},
			&model.Song{ID: songID},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			songRepository := new(repository.SongRepositoryMock)
			_uut := AddSongToArtist{songRepository: songRepository}

			songRepository.On("GetWithSongs", mock.IsType(tt.song), tt.request.SongID).
				Return(nil, tt.song).
				Once()
			songRepository.On("UpdateWithAssociations", mock.IsType(tt.song)).
				Run(func(args mock.Arguments) {
					newSong := args.Get(0).(*model.Song)
					if newSong.Album != nil {
						assert.Equal(t, tt.request.ID, *tt.song.Album.ArtistID)
						for _, song := range newSong.Album.Songs {
							assert.Equal(t, tt.request.ID, *song.ArtistID)
						}
					} else {
						assert.Equal(t, tt.request.ID, *newSong.ArtistID)
					}
				}).
				Return(nil).
				Once()

			// when
			errCode := _uut.Handle(tt.request)

			// then
			assert.Nil(t, errCode)

			songRepository.AssertExpectations(t)
		})
	}
}
