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

func TestAddSongsToArtist_WhenGetSongWithSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := AddSongsToArtist{songRepository: songRepository}

	request := requests.AddSongsToArtistRequest{
		ID:      uuid.New(),
		SongIDs: []uuid.UUID{uuid.New()},
	}

	internalError := errors.New("internal error")
	songRepository.On("GetAllByIDsWithSongs", mock.Anything, request.SongIDs).
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

func TestAddSongsToArtist_WhenOneSongHasArtist_ShouldReturnBadRequestError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := AddSongsToArtist{songRepository: songRepository}

	request := requests.AddSongsToArtistRequest{
		ID:      uuid.New(),
		SongIDs: []uuid.UUID{uuid.New()},
	}

	songs := []model.Song{
		{
			ID:       request.SongIDs[0],
			ArtistID: &[]uuid.UUID{uuid.New()}[0],
		},
	}
	songRepository.On("GetAllByIDsWithSongs", mock.IsType(&songs), request.SongIDs).
		Return(nil, &songs).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusBadRequest, errCode.Code)
	assert.Equal(t, "song "+songs[0].ID.String()+"already has an artist", errCode.Error.Error())

	songRepository.AssertExpectations(t)
}

func TestAddSongsToArtist_WhenUpdateWithAssociationsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := AddSongsToArtist{songRepository: songRepository}

	request := requests.AddSongsToArtistRequest{
		ID:      uuid.New(),
		SongIDs: []uuid.UUID{uuid.New()},
	}

	songs := []model.Song{
		{
			ID:       request.SongIDs[0],
			ArtistID: nil,
		},
	}
	songRepository.On("GetAllByIDsWithSongs", mock.IsType(&songs), request.SongIDs).
		Return(nil, &songs).
		Once()

	internalError := errors.New("internal error")
	songRepository.On("UpdateWithAssociations", mock.IsType(&songs[0])).
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

func TestAddSongsToArtist_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := AddSongsToArtist{songRepository: songRepository}

	request := requests.AddSongsToArtistRequest{
		ID:      uuid.New(),
		SongIDs: []uuid.UUID{uuid.New(), uuid.New()},
	}

	songs := []model.Song{
		{
			ID:       request.SongIDs[0],
			ArtistID: nil,
		},
		{
			ID:       request.SongIDs[1],
			ArtistID: nil,
			Album: &model.Album{
				ArtistID: nil,
				Songs: []model.Song{
					{ID: uuid.New()},
					{ID: uuid.New()},
					{ID: request.SongIDs[1]},
					{ID: uuid.New()},
				},
			},
		},
	}
	songRepository.On("GetAllByIDsWithSongs", mock.IsType(&songs), request.SongIDs).
		Return(nil, &songs).
		Once()

	for _, song := range songs {
		songRepository.On("UpdateWithAssociations", mock.IsType(&song)).
			Run(func(args mock.Arguments) {
				newSong := args.Get(0).(*model.Song)
				if newSong.Album != nil {
					assert.Equal(t, request.ID, *song.Album.ArtistID)
					for _, s := range newSong.Album.Songs {
						assert.Equal(t, request.ID, *s.ArtistID)
					}
				} else {
					assert.Equal(t, request.ID, *newSong.ArtistID)
				}
			}).
			Return(nil).
			Once()
	}

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)

	songRepository.AssertExpectations(t)
}
