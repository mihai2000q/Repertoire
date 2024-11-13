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

func TestRemoveSongsFromArtist_WhenGetSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := RemoveSongFromArtist{songRepository: songRepository}

	request := requests.RemoveSongsFromArtistRequest{
		ID:      uuid.New(),
		SongIDs: []uuid.UUID{uuid.New(), uuid.New()},
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

func TestRemoveSongsFromArtist_WhenOneSongArtistDoesNotMatch_ShouldReturnBadRequestError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := RemoveSongFromArtist{songRepository: songRepository}

	request := requests.RemoveSongsFromArtistRequest{
		ID:      uuid.New(),
		SongIDs: []uuid.UUID{uuid.New(), uuid.New()},
	}

	songs := &[]model.Song{
		{
			ID:       request.SongIDs[0],
			ArtistID: nil,
		},
	}
	songRepository.On("GetAllByIDs", mock.IsType(songs), request.SongIDs).
		Return(nil, songs).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusBadRequest, errCode.Code)
	assert.Equal(t, "song "+request.SongIDs[0].String()+" is not owned by this artist", errCode.Error.Error())

	songRepository.AssertExpectations(t)
}

func TestRemoveSongsFromArtist_WhenUpdateSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := RemoveSongFromArtist{songRepository: songRepository}

	request := requests.RemoveSongsFromArtistRequest{
		ID:      uuid.New(),
		SongIDs: []uuid.UUID{uuid.New(), uuid.New()},
	}

	songs := []model.Song{
		{
			ID:       request.SongIDs[0],
			ArtistID: &request.ID,
		},
	}
	songRepository.On("GetAllByIDs", mock.IsType(&songs), request.SongIDs).
		Return(nil, &songs).
		Once()

	internalError := errors.New("internal error")
	songRepository.On("Update", mock.IsType(&songs[0])).
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

func TestRemoveSongsFromArtist_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := RemoveSongFromArtist{songRepository: songRepository}

	request := requests.RemoveSongsFromArtistRequest{
		ID:      uuid.New(),
		SongIDs: []uuid.UUID{uuid.New(), uuid.New()},
	}

	songs := []model.Song{
		{
			ID:       request.SongIDs[0],
			ArtistID: &request.ID,
		},
		{
			ID:       request.SongIDs[1],
			ArtistID: &request.ID,
		},
	}
	songRepository.On("GetAllByIDs", mock.IsType(&songs), request.SongIDs).
		Return(nil, &songs).
		Once()

	for _, song := range songs {
		songRepository.On("Update", mock.IsType(&song)).
			Run(func(args mock.Arguments) {
				newSong := args.Get(0).(*model.Song)
				assert.Nil(t, newSong.ArtistID)
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
