package song

import (
	"errors"
	"net/http"
	"repertoire/server/domain/usecase/song"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"slices"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteSong_WhenGetSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewDeleteSong(songRepository, nil, nil)

	id := uuid.New()

	internalError := errors.New("internal error")
	songRepository.On("GetWithPlaylistsAndSongs", mock.IsType(new(model.Song)), id).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestDeleteSong_WhenGetSongIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewDeleteSong(songRepository, nil, nil)

	id := uuid.New()

	songRepository.On("GetWithPlaylistsAndSongs", mock.IsType(new(model.Song)), id).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song not found", errCode.Error.Error())

	songRepository.AssertExpectations(t)
}

func TestDeleteSong_WhenGetAllByAlbumAndTrackNoFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewDeleteSong(songRepository, nil, nil)

	id := uuid.New()

	mockSong := &model.Song{
		ID:           id,
		AlbumID:      &[]uuid.UUID{uuid.New()}[0],
		AlbumTrackNo: &[]uint{1}[0],
	}
	songRepository.On("GetWithPlaylistsAndSongs", mock.IsType(mockSong), id).
		Return(nil, mockSong).
		Once()

	internalError := errors.New("internal error")
	songRepository.
		On(
			"GetAllByAlbumAndTrackNo",
			mock.IsType(&[]model.Song{}),
			*mockSong.AlbumID,
			*mockSong.AlbumTrackNo,
		).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestDeleteSong_WhenUpdateAllAlbumSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewDeleteSong(songRepository, nil, nil)

	id := uuid.New()

	mockSong := &model.Song{
		ID:           id,
		AlbumID:      &[]uuid.UUID{uuid.New()}[0],
		AlbumTrackNo: &[]uint{1}[0],
	}
	songRepository.On("GetWithPlaylistsAndSongs", mock.IsType(mockSong), id).
		Return(nil, mockSong).
		Once()

	mockAlbumSongs := &[]model.Song{
		{AlbumTrackNo: &[]uint{1}[0]},
	}
	songRepository.
		On(
			"GetAllByAlbumAndTrackNo",
			mock.IsType(&[]model.Song{}),
			*mockSong.AlbumID,
			*mockSong.AlbumTrackNo,
		).
		Return(nil, mockAlbumSongs).
		Once()

	internalError := errors.New("internal error")
	songRepository.On("UpdateAll", mock.IsType(mockAlbumSongs)).Return(internalError).Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestDeleteSong_WhenUpdateAllPlaylistsSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := song.NewDeleteSong(songRepository, playlistRepository, nil)

	id := uuid.New()

	mockSong := &model.Song{
		ID: id,
		Playlists: []model.Playlist{
			{
				PlaylistSongs: []model.PlaylistSong{
					{SongID: id},
					{SongID: uuid.New()},
				},
			},
		},
	}
	songRepository.On("GetWithPlaylistsAndSongs", mock.IsType(mockSong), id).
		Return(nil, mockSong).
		Once()

	internalError := errors.New("internal error")
	playlistRepository.On("UpdateAllPlaylistSongs", mock.IsType(new([]model.PlaylistSong))).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
	playlistRepository.AssertExpectations(t)
}

func TestDeleteSong_WhenDeleteSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewDeleteSong(songRepository, nil, nil)

	id := uuid.New()

	mockSong := &model.Song{ID: id}
	songRepository.On("GetWithPlaylistsAndSongs", mock.IsType(mockSong), id).
		Return(nil, mockSong).
		Once()

	internalError := errors.New("internal error")
	songRepository.On("Delete", []uuid.UUID{id}).Return(internalError).Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestDeleteSong_WhenPublishFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := song.NewDeleteSong(songRepository, nil, messagePublisherService)

	id := uuid.New()

	mockSong := &model.Song{ID: id}
	songRepository.On("GetWithPlaylistsAndSongs", mock.IsType(mockSong), id).
		Return(nil, mockSong).
		Once()

	songRepository.On("Delete", []uuid.UUID{id}).Return(nil).Once()

	internalError := errors.New("internal error")
	messagePublisherService.On("Publish", topics.SongsDeletedTopic, []model.Song{*mockSong}).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
}

func TestDeleteSong_WhenWithoutAlbumOrPlaylists_ShouldDeleteSong(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := song.NewDeleteSong(songRepository, nil, messagePublisherService)

	mockSong := model.Song{ID: uuid.New()}

	id := mockSong.ID

	songRepository.On("GetWithPlaylistsAndSongs", mock.IsType(&mockSong), id).
		Return(nil, &mockSong).
		Once()
	songRepository.On("Delete", []uuid.UUID{id}).Return(nil).Once()

	messagePublisherService.On("Publish", topics.SongsDeletedTopic, []model.Song{mockSong}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.Nil(t, errCode)

	songRepository.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
}

func TestDeleteSong_WhenWithAlbum_ShouldDeleteSongAndReorderAlbum(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := song.NewDeleteSong(songRepository, nil, messagePublisherService)

	mockSong := model.Song{
		ID:           uuid.New(),
		AlbumID:      &[]uuid.UUID{uuid.New()}[0],
		AlbumTrackNo: &[]uint{2}[0],
	}

	id := mockSong.ID

	songRepository.On("GetWithPlaylistsAndSongs", mock.IsType(&mockSong), id).
		Return(nil, &mockSong).
		Once()

	mockAlbumSongs := &[]model.Song{
		{AlbumTrackNo: &[]uint{3}[0]},
		{AlbumTrackNo: &[]uint{4}[0]},
		{AlbumTrackNo: &[]uint{5}[0]},
	}
	oldAlbumSongs := slices.Clone(*mockAlbumSongs)
	songRepository.
		On(
			"GetAllByAlbumAndTrackNo",
			mock.IsType(mockAlbumSongs),
			*mockSong.AlbumID,
			*mockSong.AlbumTrackNo,
		).
		Return(nil, mockAlbumSongs).
		Once()

	songRepository.On("UpdateAll", mock.IsType(mockAlbumSongs)).
		Run(func(args mock.Arguments) {
			newSongs := args.Get(0).(*[]model.Song)
			for i := range *newSongs {
				assert.Equal(t, *oldAlbumSongs[i].AlbumTrackNo-1, *(*newSongs)[i].AlbumTrackNo)
			}
		}).
		Return(nil).
		Once()

	songRepository.On("Delete", []uuid.UUID{id}).Return(nil).Once()

	messagePublisherService.On("Publish", topics.SongsDeletedTopic, []model.Song{mockSong}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.Nil(t, errCode)

	songRepository.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
}

func TestDeleteSong_WhenWithPlaylists_ShouldDeleteSongAndReorderPlaylists(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	playlistRepository := new(repository.PlaylistRepositoryMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := song.NewDeleteSong(songRepository, playlistRepository, messagePublisherService)

	id := uuid.New()

	mockPlaylists := []model.Playlist{
		{
			Title: "Playlist 1",
			PlaylistSongs: []model.PlaylistSong{
				{
					SongID:      uuid.New(),
					SongTrackNo: 1,
				},
				{
					SongID:      id,
					SongTrackNo: 2,
				},
				{
					SongID:      uuid.New(),
					SongTrackNo: 3,
				},
				{
					SongID:      id,
					SongTrackNo: 4,
				},
				{
					SongID:      uuid.New(),
					SongTrackNo: 5,
				},
			},
		},
		{
			Title: "Playlist 2",
			PlaylistSongs: []model.PlaylistSong{
				{
					SongID:      id,
					SongTrackNo: 1,
				},
				{
					SongID:      uuid.New(),
					SongTrackNo: 2,
				},
			},
		},
	}
	mockSong := model.Song{
		ID:        id,
		Playlists: mockPlaylists,
	}

	expectedOrderedPlaylistSongs := &[]model.PlaylistSong{
		// Playlist 1
		{
			SongID:      mockPlaylists[0].PlaylistSongs[0].SongID,
			SongTrackNo: 1,
		},
		{
			SongID:      mockPlaylists[0].PlaylistSongs[2].SongID,
			SongTrackNo: 2,
		},
		{
			SongID:      mockPlaylists[0].PlaylistSongs[4].SongID,
			SongTrackNo: 3,
		},
		// Playlist 2
		{
			SongID:      mockPlaylists[1].PlaylistSongs[1].SongID,
			SongTrackNo: 1,
		},
	}

	songRepository.On("GetWithPlaylistsAndSongs", mock.IsType(&mockSong), id).
		Return(nil, &mockSong).
		Once()

	playlistRepository.On("UpdateAllPlaylistSongs", mock.IsType(expectedOrderedPlaylistSongs)).
		Return(nil).
		Once()

	songRepository.On("Delete", []uuid.UUID{id}).Return(nil).Once()

	messagePublisherService.On("Publish", topics.SongsDeletedTopic, []model.Song{mockSong}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(id)

	// then
	assert.Nil(t, errCode)

	songRepository.AssertExpectations(t)
	playlistRepository.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
}
