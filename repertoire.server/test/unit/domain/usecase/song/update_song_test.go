package song

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/song"
	"repertoire/server/internal/enums"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"slices"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateSong_WhenGetSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewUpdateSong(songRepository, nil, nil)

	request := requests.UpdateSongRequest{
		ID:    uuid.New(),
		Title: "New Song",
	}

	internalError := errors.New("internal error")
	songRepository.On("Get", new(model.Song), request.ID).Return(internalError).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestUpdateSong_WhenSongIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewUpdateSong(songRepository, nil, nil)

	request := requests.UpdateSongRequest{
		ID:    uuid.New(),
		Title: "New Song",
	}

	songRepository.On("Get", new(model.Song), request.ID).Return(nil).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song not found", errCode.Error.Error())

	songRepository.AssertExpectations(t)
}

func TestUpdateSong_WhenGetAlbumFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := song.NewUpdateSong(songRepository, albumRepository, nil)

	request := requests.UpdateSongRequest{
		ID:      uuid.New(),
		Title:   "New Song",
		AlbumID: &[]uuid.UUID{uuid.New()}[0],
	}

	mockSong := model.Song{ID: uuid.New()}

	songRepository.On("Get", new(model.Song), request.ID).
		Return(nil, &mockSong).
		Once()

	internalError := errors.New("internal error")
	albumRepository.On("Get", new(model.Album), *request.AlbumID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
	albumRepository.AssertExpectations(t)
}

func TestUpdateSong_WhenAlbumIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := song.NewUpdateSong(songRepository, albumRepository, nil)

	request := requests.UpdateSongRequest{
		ID:      uuid.New(),
		Title:   "New Song",
		AlbumID: &[]uuid.UUID{uuid.New()}[0],
	}

	mockSong := model.Song{ID: uuid.New()}

	songRepository.On("Get", new(model.Song), request.ID).
		Return(nil, &mockSong).
		Once()

	albumRepository.On("Get", new(model.Album), *request.AlbumID).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "album not found", errCode.Error.Error())

	songRepository.AssertExpectations(t)
	albumRepository.AssertExpectations(t)
}

func TestUpdateSong_WhenAlbumArtistAndRequestArtistDoNotMatch_ShouldReturnBadRequestError(t *testing.T) {
	albumID := uuid.New()
	artistID := uuid.New()

	tests := []struct {
		name    string
		request requests.UpdateSongRequest
		song    model.Song
		album   model.Album
	}{
		{
			"Artist Changes",
			requests.UpdateSongRequest{
				ID:       uuid.New(),
				Title:    "New Song",
				ArtistID: &[]uuid.UUID{uuid.New()}[0],
				AlbumID:  &albumID,
			},
			model.Song{ID: uuid.New()},
			model.Album{ID: albumID},
		},
		{
			"Album Changes",
			requests.UpdateSongRequest{
				ID:       uuid.New(),
				Title:    "New Song",
				AlbumID:  &[]uuid.UUID{uuid.New()}[0],
				ArtistID: &artistID,
			},
			model.Song{ID: uuid.New(), ArtistID: &artistID},
			model.Album{ID: uuid.New()},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			songRepository := new(repository.SongRepositoryMock)
			albumRepository := new(repository.AlbumRepositoryMock)
			_uut := song.NewUpdateSong(songRepository, albumRepository, nil)

			songRepository.On("Get", new(model.Song), tt.request.ID).
				Return(nil, &tt.song).
				Once()

			albumRepository.On("Get", new(model.Album), *tt.request.AlbumID).
				Return(nil, &tt.album).
				Once()

			// when
			errCode := _uut.Handle(tt.request)

			// then
			assert.NotNil(t, errCode)
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
			assert.Equal(t, "album's artist does not match the request's artist", errCode.Error.Error())

			songRepository.AssertExpectations(t)
			albumRepository.AssertExpectations(t)
		})
	}
}

func TestUpdateSong_WhenGetAllSongsByAlbumFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewUpdateSong(songRepository, nil, nil)

	request := requests.UpdateSongRequest{
		ID: uuid.New(),
	}

	mockSong := &model.Song{
		ID:           request.ID,
		AlbumID:      &[]uuid.UUID{uuid.New()}[0],
		AlbumTrackNo: &[]uint{1}[0],
	}
	songRepository.On("Get", new(model.Song), request.ID).
		Return(nil, mockSong).
		Once()

	internalError := errors.New("internal error")
	songRepository.
		On(
			"GetAllByAlbumAndTrackNo",
			new([]model.Song),
			*mockSong.AlbumID,
			*mockSong.AlbumTrackNo,
		).
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

func TestUpdateSong_WhenUpdateAllSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewUpdateSong(songRepository, nil, nil)

	request := requests.UpdateSongRequest{
		ID: uuid.New(),
	}

	mockSong := &model.Song{
		ID:           request.ID,
		AlbumID:      &[]uuid.UUID{uuid.New()}[0],
		AlbumTrackNo: &[]uint{1}[0],
	}
	songRepository.On("Get", new(model.Song), request.ID).
		Return(nil, mockSong).
		Once()

	songRepository.
		On(
			"GetAllByAlbumAndTrackNo",
			new([]model.Song),
			*mockSong.AlbumID,
			*mockSong.AlbumTrackNo,
		).
		Return(nil).
		Once()

	internalError := errors.New("internal error")
	songRepository.On("UpdateAll", new([]model.Song)).
		Return(internalError, mockSong).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestUpdateSong_WhenCountSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := song.NewUpdateSong(songRepository, albumRepository, nil)

	request := requests.UpdateSongRequest{
		ID:       uuid.New(),
		Title:    "New Song",
		ArtistID: &[]uuid.UUID{uuid.New()}[0],
		AlbumID:  &[]uuid.UUID{uuid.New()}[0],
	}
	mockSong := model.Song{ID: uuid.New()}
	mockAlbum := model.Album{ID: *request.AlbumID, ArtistID: request.ArtistID}

	songRepository.On("Get", new(model.Song), request.ID).
		Return(nil, &mockSong).
		Once()

	albumRepository.On("Get", new(model.Album), *request.AlbumID).
		Return(nil, &mockAlbum).
		Once()

	internalError := errors.New("internal error")
	songRepository.On("CountByAlbum", new(int64), *request.AlbumID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
	albumRepository.AssertExpectations(t)
}

func TestUpdateSong_WhenUpdateSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewUpdateSong(songRepository, nil, nil)

	request := requests.UpdateSongRequest{
		ID:    uuid.New(),
		Title: "New Song",
	}

	mockSong := &model.Song{
		ID:    request.ID,
		Title: "Some Song",
	}
	songRepository.On("Get", new(model.Song), request.ID).
		Return(nil, mockSong).
		Once()

	internalError := errors.New("internal error")
	songRepository.On("Update", mock.IsType(mockSong)).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestUpdateSong_WhenPublishFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := song.NewUpdateSong(songRepository, nil, messagePublisherService)

	request := requests.UpdateSongRequest{
		ID:    uuid.New(),
		Title: "New Song",
	}

	mockSong := &model.Song{
		ID:    request.ID,
		Title: "Some Song",
	}
	songRepository.On("Get", new(model.Song), request.ID).
		Return(nil, mockSong).
		Once()

	songRepository.On("Update", mock.IsType(mockSong)).
		Return(nil).
		Once()

	internalError := errors.New("internal error")
	messagePublisherService.On("Publish", topics.SongUpdatedTopic, mock.IsType(*mockSong)).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
}

func TestUpdateSong_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	songID := uuid.New()
	albumID := uuid.New()
	artistID := uuid.New()

	tests := []struct {
		name    string
		request requests.UpdateSongRequest
		song    *model.Song
	}{
		{
			"Non Optional",
			requests.UpdateSongRequest{
				ID:    songID,
				Title: "New Song",
			},
			&model.Song{
				ID:    songID,
				Title: "Some Song",
			},
		},
		{
			"With Artist - add it",
			requests.UpdateSongRequest{
				ID:       songID,
				Title:    "New Song",
				ArtistID: &[]uuid.UUID{uuid.New()}[0],
			},
			&model.Song{
				ID:    songID,
				Title: "Some Song",
			},
		},
		{
			"With Artist - change it",
			requests.UpdateSongRequest{
				ID:       songID,
				Title:    "New Song",
				ArtistID: &[]uuid.UUID{uuid.New()}[0],
			},
			&model.Song{
				ID:       songID,
				Title:    "Some Song",
				ArtistID: &[]uuid.UUID{uuid.New()}[0],
			},
		},
		{
			"With Artist - remove it",
			requests.UpdateSongRequest{
				ID:    songID,
				Title: "New Song",
			},
			&model.Song{
				ID:       songID,
				Title:    "Some Song",
				ArtistID: &[]uuid.UUID{uuid.New()}[0],
			},
		},
		{
			"With Album - remove it",
			requests.UpdateSongRequest{
				ID:    songID,
				Title: "New Song",
			},
			&model.Song{
				ID:           songID,
				Title:        "Some Song",
				AlbumID:      &[]uuid.UUID{uuid.New()}[0],
				AlbumTrackNo: &[]uint{2}[0],
			},
		},
		{
			"All filled",
			requests.UpdateSongRequest{
				ID:             songID,
				Title:          "New Song",
				Description:    "This is a nice description",
				IsRecorded:     true,
				Bpm:            &[]uint{123}[0],
				GuitarTuningID: &[]uuid.UUID{uuid.New()}[0],
				ReleaseDate:    &[]time.Time{time.Now()}[0],
				Difficulty:     &[]enums.Difficulty{enums.Hard}[0],
				SongsterrLink:  &[]string{"http://songsterr.com/some-song"}[0],
				YoutubeLink:    &[]string{"https://www.youtube.com/watch?v=IHgFJEJgUrg"}[0],
				ArtistID:       &artistID,
				AlbumID:        &albumID,
			},
			&model.Song{
				ID:           songID,
				Title:        "Some Song",
				ArtistID:     &artistID,
				AlbumID:      &albumID,
				AlbumTrackNo: &[]uint{2}[0],
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			songRepository := new(repository.SongRepositoryMock)
			messagePublisherService := new(service.MessagePublisherServiceMock)
			_uut := song.NewUpdateSong(songRepository, nil, messagePublisherService)

			songRepository.On("Get", new(model.Song), tt.request.ID).
				Return(nil, tt.song).
				Once()

			assertOldAlbumReordered(t, songRepository, *tt.song, tt.request.AlbumID)

			var newSong *model.Song
			songRepository.On("Update", mock.IsType(tt.song)).
				Run(func(args mock.Arguments) {
					newSong = args.Get(0).(*model.Song)
					assertUpdatedSong(t, tt.request, newSong)
				}).
				Return(nil).
				Once()

			messagePublisherService.On("Publish", topics.SongUpdatedTopic, mock.IsType(*tt.song)).
				Run(func(args mock.Arguments) {
					assert.Equal(t, *newSong, args.Get(1).(model.Song))
				}).
				Return(nil).
				Once()

			// when
			errCode := _uut.Handle(tt.request)

			// then
			assert.Nil(t, errCode)

			songRepository.AssertExpectations(t)
			messagePublisherService.AssertExpectations(t)
		})
	}
}

func TestUpdateSong_WhenRequestHasAlbum_ShouldNotReturnAnyError(t *testing.T) {
	albumID := uuid.New()
	artistID := uuid.New()

	tests := []struct {
		name    string
		request requests.UpdateSongRequest
		song    model.Song
		album   model.Album
	}{
		{
			"Change Album - without artist",
			requests.UpdateSongRequest{
				ID:      uuid.New(),
				Title:   "New Song",
				AlbumID: &albumID,
			},
			model.Song{
				ID:           uuid.New(),
				AlbumID:      &[]uuid.UUID{uuid.New()}[0],
				AlbumTrackNo: &[]uint{2}[0],
			},
			model.Album{
				ID: albumID,
			},
		},
		{
			"Add Album - without artist",
			requests.UpdateSongRequest{
				ID:      uuid.New(),
				Title:   "New Song",
				AlbumID: &albumID,
			},
			model.Song{ID: uuid.New()},
			model.Album{ID: albumID},
		},
		{
			"Add Album - with unchanged artist",
			requests.UpdateSongRequest{
				ID:       uuid.New(),
				Title:    "New Song",
				AlbumID:  &albumID,
				ArtistID: &artistID,
			},
			model.Song{
				ID:       uuid.New(),
				ArtistID: &artistID,
			},
			model.Album{
				ID:       albumID,
				ArtistID: &artistID,
			},
		},
		{
			"Add Album - with changing artist (when song didn't have one)",
			requests.UpdateSongRequest{
				ID:       uuid.New(),
				Title:    "New Song",
				AlbumID:  &albumID,
				ArtistID: &artistID,
			},
			model.Song{ID: uuid.New()},
			model.Album{
				ID:       albumID,
				ArtistID: &artistID,
			},
		},
		{
			"Add Album - with changing artist (when song used to have one)",
			requests.UpdateSongRequest{
				ID:       uuid.New(),
				Title:    "New Song",
				AlbumID:  &albumID,
				ArtistID: &artistID,
			},
			model.Song{
				ID:       uuid.New(),
				ArtistID: &[]uuid.UUID{uuid.New()}[0],
			},
			model.Album{
				ID:       albumID,
				ArtistID: &artistID,
			},
		},
		{
			"Change Album - with changing artist (when song didn't have one)",
			requests.UpdateSongRequest{
				ID:       uuid.New(),
				Title:    "New Song",
				AlbumID:  &albumID,
				ArtistID: &artistID,
			},
			model.Song{
				ID:           uuid.New(),
				AlbumID:      &[]uuid.UUID{uuid.New()}[0],
				AlbumTrackNo: &[]uint{2}[0],
			},
			model.Album{
				ID:       albumID,
				ArtistID: &artistID,
			},
		},
		{
			"Change Album - with changing artist (when song used to have one)",
			requests.UpdateSongRequest{
				ID:       uuid.New(),
				Title:    "New Song",
				AlbumID:  &albumID,
				ArtistID: &artistID,
			},
			model.Song{
				ID:           uuid.New(),
				ArtistID:     &[]uuid.UUID{uuid.New()}[0],
				AlbumID:      &[]uuid.UUID{uuid.New()}[0],
				AlbumTrackNo: &[]uint{2}[0],
			},
			model.Album{
				ID:       albumID,
				ArtistID: &artistID,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			songRepository := new(repository.SongRepositoryMock)
			albumRepository := new(repository.AlbumRepositoryMock)
			messagePublisherService := new(service.MessagePublisherServiceMock)
			_uut := song.NewUpdateSong(songRepository, albumRepository, messagePublisherService)

			songRepository.On("Get", new(model.Song), tt.request.ID).
				Return(nil, &tt.song).
				Once()

			albumRepository.On("Get", new(model.Album), *tt.request.AlbumID).
				Return(nil, &tt.album).
				Once()

			assertOldAlbumReordered(t, songRepository, tt.song, tt.request.AlbumID)

			var songsCount int64 = 12
			songRepository.On("CountByAlbum", new(int64), *tt.request.AlbumID).
				Return(nil, &songsCount).
				Once()

			var newSong *model.Song
			songRepository.On("Update", mock.IsType(&tt.song)).
				Run(func(args mock.Arguments) {
					newSong = args.Get(0).(*model.Song)
					assertUpdatedSong(t, tt.request, newSong)
					assert.Equal(t, uint(songsCount+1), *newSong.AlbumTrackNo)
				}).
				Return(nil).
				Once()

			messagePublisherService.On("Publish", topics.SongUpdatedTopic, mock.IsType(tt.song)).
				Run(func(args mock.Arguments) {
					assert.Equal(t, *newSong, args.Get(1).(model.Song))
				}).
				Return(nil).
				Once()

			// when
			errCode := _uut.Handle(tt.request)

			// then
			assert.Nil(t, errCode)

			songRepository.AssertExpectations(t)
			albumRepository.AssertExpectations(t)
			messagePublisherService.AssertExpectations(t)
		})
	}
}

func assertUpdatedSong(t *testing.T, request requests.UpdateSongRequest, song *model.Song) {
	assert.Equal(t, request.Title, song.Title)
	assert.Equal(t, request.Description, song.Description)
	assert.Equal(t, request.IsRecorded, song.IsRecorded)
	assert.Equal(t, request.Bpm, song.Bpm)
	assert.Equal(t, request.SongsterrLink, song.SongsterrLink)
	assert.Equal(t, request.YoutubeLink, song.YoutubeLink)
	assert.Equal(t, request.ReleaseDate, song.ReleaseDate)
	assert.Equal(t, request.Difficulty, song.Difficulty)
	assert.Equal(t, request.GuitarTuningID, song.GuitarTuningID)
	assert.Equal(t, request.ArtistID, song.ArtistID)
	assert.Equal(t, request.AlbumID, song.AlbumID)
}

func assertOldAlbumReordered(
	t *testing.T,
	songRepository *repository.SongRepositoryMock,
	song model.Song,
	newAlbumID *uuid.UUID,
) {
	if song.AlbumID == nil || song.AlbumID == newAlbumID {
		return
	}

	songs := []model.Song{
		{ID: uuid.New(), AlbumTrackNo: &[]uint{3}[0]},
		{ID: uuid.New(), AlbumTrackNo: &[]uint{4}[0]},
		{ID: uuid.New(), AlbumTrackNo: &[]uint{5}[0]},
	}
	oldSongs := slices.Clone(songs)

	songRepository.
		On(
			"GetAllByAlbumAndTrackNo",
			new([]model.Song),
			*song.AlbumID,
			*song.AlbumTrackNo,
		).
		Return(nil, &songs).
		Once()

	songRepository.On("UpdateAll", &songs).
		Run(func(args mock.Arguments) {
			newSongs := args.Get(0).(*[]model.Song)
			for i := range *newSongs {
				assert.Equal(t, oldSongs[i].ID, (*newSongs)[i].ID)
				assert.Equal(t, *oldSongs[i].AlbumTrackNo-1, *(*newSongs)[i].AlbumTrackNo)
			}
		}).
		Return(nil).
		Once()
}
