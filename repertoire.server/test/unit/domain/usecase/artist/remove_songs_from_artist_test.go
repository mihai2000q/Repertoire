package artist

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/artist"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"testing"
)

func TestRemoveSongsFromArtist_WhenGetSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := artist.NewRemoveSongsFromArtist(songRepository, nil)

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

func TestRemoveSongsFromArtist_WhenOneSongArtistDoesNotMatch_ShouldReturnConflictError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := artist.NewRemoveSongsFromArtist(songRepository, nil)

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
	assert.Equal(t, http.StatusConflict, errCode.Code)
	assert.Equal(t, "song "+request.SongIDs[0].String()+" is not owned by this artist", errCode.Error.Error())

	songRepository.AssertExpectations(t)
}

func TestRemoveSongsFromArtist_WhenUpdateAllSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := artist.NewRemoveSongsFromArtist(songRepository, nil)

	request := requests.RemoveSongsFromArtistRequest{
		ID:      uuid.New(),
		SongIDs: []uuid.UUID{uuid.New(), uuid.New()},
	}

	songs := &[]model.Song{
		{
			ID:       request.SongIDs[0],
			ArtistID: &request.ID,
		},
	}
	songRepository.On("GetAllByIDs", mock.IsType(songs), request.SongIDs).
		Return(nil, songs).
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

	songRepository.AssertExpectations(t)
}

func TestRemoveSongsFromArtist_WhenPublishFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := artist.NewRemoveSongsFromArtist(songRepository, messagePublisherService)

	request := requests.RemoveSongsFromArtistRequest{
		ID:      uuid.New(),
		SongIDs: []uuid.UUID{uuid.New(), uuid.New()},
	}

	songs := &[]model.Song{
		{
			ID:       request.SongIDs[0],
			ArtistID: &request.ID,
		},
	}
	songRepository.On("GetAllByIDs", mock.IsType(songs), request.SongIDs).
		Return(nil, songs).
		Once()

	songRepository.On("UpdateAll", mock.IsType(songs)).
		Return(nil).
		Once()

	internalError := errors.New("internal error")
	messagePublisherService.On("Publish", topics.SongsUpdatedTopic, request.SongIDs).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
}

func TestRemoveSongsFromArtist_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := artist.NewRemoveSongsFromArtist(songRepository, messagePublisherService)

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

	songRepository.On("UpdateAll", mock.IsType(&songs)).
		Run(func(args mock.Arguments) {
			newSongs := args.Get(0).(*[]model.Song)
			for _, song := range *newSongs {
				assert.Nil(t, song.ArtistID)
			}
		}).
		Return(nil).
		Once()

	messagePublisherService.On("Publish", topics.SongsUpdatedTopic, request.SongIDs).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)

	songRepository.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
}
