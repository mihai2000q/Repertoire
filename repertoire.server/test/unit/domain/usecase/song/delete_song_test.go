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

func TestDeleteSong_WhenUpdateAllFails_ShouldReturnInternalServerError(t *testing.T) {
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
	songRepository.On("Delete", id).Return(internalError).Once()

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

	songRepository.On("Delete", id).Return(nil).Once()

	internalError := errors.New("internal error")
	messagePublisherService.On("Publish", topics.SongDeletedTopic, *mockSong).
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

func TestDeleteSong_WhenSuccessful_ShouldDeleteSong(t *testing.T) {
	songID := uuid.New()

	tests := []struct {
		name       string
		song       model.Song
		albumSongs []model.Song
	}{
		{
			"Normal delete, without album",
			model.Song{ID: songID},
			[]model.Song{},
		},
		{
			"With Album",
			model.Song{
				ID:           songID,
				AlbumID:      &[]uuid.UUID{uuid.New()}[0],
				AlbumTrackNo: &[]uint{2}[0],
			},
			[]model.Song{
				{AlbumTrackNo: &[]uint{3}[0]},
				{AlbumTrackNo: &[]uint{4}[0]},
				{AlbumTrackNo: &[]uint{5}[0]},
			},
		},
		{
			"With Playlist",
			model.Song{
				ID:           songID,
				AlbumID:      &[]uuid.UUID{uuid.New()}[0],
				AlbumTrackNo: &[]uint{2}[0],
				Playlists: []model.Playlist{
					{
						Title: "Playlist 1",
						PlaylistSongs: []model.PlaylistSong{
							{
								SongID:      uuid.New(),
								SongTrackNo: 1,
							},
							{
								SongID:      songID,
								SongTrackNo: 2,
							},
							{
								SongID:      uuid.New(),
								SongTrackNo: 3,
							},
						},
					},
					{
						Title: "Playlist 2",
						PlaylistSongs: []model.PlaylistSong{
							{
								SongID:      songID,
								SongTrackNo: 1,
							},
							{
								SongID:      uuid.New(),
								SongTrackNo: 2,
							},
						},
					},
				},
			},
			[]model.Song{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			songRepository := new(repository.SongRepositoryMock)
			playlistRepository := new(repository.PlaylistRepositoryMock)
			messagePublisherService := new(service.MessagePublisherServiceMock)
			_uut := song.NewDeleteSong(songRepository, playlistRepository, messagePublisherService)

			id := tt.song.ID

			songRepository.On("GetWithPlaylistsAndSongs", mock.IsType(&tt.song), id).
				Return(nil, &tt.song).
				Once()
			songRepository.On("Delete", id).Return(nil).Once()

			if tt.song.AlbumID != nil {
				mockAlbumSongs := slices.Clone(tt.albumSongs)
				songRepository.
					On(
						"GetAllByAlbumAndTrackNo",
						mock.IsType(&tt.albumSongs),
						*tt.song.AlbumID,
						*tt.song.AlbumTrackNo,
					).
					Return(nil, &mockAlbumSongs).
					Once()

				songRepository.On("UpdateAll", mock.IsType(&tt.albumSongs)).
					Run(func(args mock.Arguments) {
						newSongs := args.Get(0).(*[]model.Song)
						for i := range *newSongs {
							assert.Equal(t, *tt.albumSongs[i].AlbumTrackNo-1, *(*newSongs)[i].AlbumTrackNo)
						}
					}).
					Return(nil).
					Once()
			}

			for _, playlist := range tt.song.Playlists {
				songRemovedIndex := slices.IndexFunc(playlist.PlaylistSongs, func(playlistSong model.PlaylistSong) bool {
					return playlistSong.SongID == tt.song.ID
				})
				playlistRepository.On("UpdateAllPlaylistSongs", mock.IsType(&playlist.PlaylistSongs)).
					Run(func(args mock.Arguments) {
						newPlaylistSongs := args.Get(0).(*[]model.PlaylistSong)
						for i := range *newPlaylistSongs {
							expectedTrackNo := playlist.PlaylistSongs[songRemovedIndex+1+i].SongTrackNo - 1
							assert.Equal(t, expectedTrackNo, (*newPlaylistSongs)[i].SongTrackNo)
						}
					}).
					Return(nil).
					Once()
			}

			messagePublisherService.On("Publish", topics.SongDeletedTopic, tt.song).
				Return(nil).
				Once()

			// when
			errCode := _uut.Handle(id)

			// then
			assert.Nil(t, errCode)

			songRepository.AssertExpectations(t)
			messagePublisherService.AssertExpectations(t)
		})
	}
}
