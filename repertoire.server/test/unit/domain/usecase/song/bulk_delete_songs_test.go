package song

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/song"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBulkDeleteSongs_WhenGetSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewBulkDeleteSongs(songRepository, nil, nil)

	request := requests.BulkDeleteSongsRequest{
		IDs: []uuid.UUID{uuid.New()},
	}

	internalError := errors.New("internal error")
	songRepository.On("GetAllByIDsWithAlbumsAndPlaylists", mock.IsType(&[]model.Song{}), request.IDs).
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

func TestBulkDeleteSongs_WhenGetSongsAreEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewBulkDeleteSongs(songRepository, nil, nil)

	request := requests.BulkDeleteSongsRequest{
		IDs: []uuid.UUID{uuid.New()},
	}

	songRepository.On("GetAllByIDsWithAlbumsAndPlaylists", mock.IsType(&[]model.Song{}), request.IDs).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "songs not found", errCode.Error.Error())

	songRepository.AssertExpectations(t)
}

func TestBulkDeleteSongs_WhenUpdateAllAlbumSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewBulkDeleteSongs(songRepository, nil, nil)

	request := requests.BulkDeleteSongsRequest{
		IDs: []uuid.UUID{uuid.New()},
	}

	mockSongs := &[]model.Song{
		{
			ID:           request.IDs[0],
			AlbumID:      &[]uuid.UUID{uuid.New()}[0],
			AlbumTrackNo: &[]uint{1}[0],
			Album: &model.Album{
				Songs: []model.Song{
					{ID: request.IDs[0], AlbumTrackNo: &[]uint{1}[0]},
					{ID: uuid.New(), AlbumTrackNo: &[]uint{2}[0]},
				},
			},
		},
	}
	songRepository.On("GetAllByIDsWithAlbumsAndPlaylists", mock.IsType(mockSongs), request.IDs).
		Return(nil, mockSongs).
		Once()

	internalError := errors.New("internal error")
	songRepository.On("UpdateAll", mock.IsType(mockSongs)).Return(internalError).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestBulkDeleteSongs_WhenUpdateAllPlaylistsSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	playlistRepository := new(repository.PlaylistRepositoryMock)
	_uut := song.NewBulkDeleteSongs(songRepository, playlistRepository, nil)

	request := requests.BulkDeleteSongsRequest{
		IDs: []uuid.UUID{uuid.New()},
	}

	mockSongs := &[]model.Song{
		{
			ID: request.IDs[0],
			Playlists: []model.Playlist{
				{
					PlaylistSongs: []model.PlaylistSong{
						{SongID: request.IDs[0]},
						{SongID: uuid.New()},
					},
				},
			},
		},
	}
	songRepository.On("GetAllByIDsWithAlbumsAndPlaylists", mock.IsType(mockSongs), request.IDs).
		Return(nil, mockSongs).
		Once()

	internalError := errors.New("internal error")
	playlistRepository.On("UpdateAllPlaylistSongs", mock.IsType(new([]model.PlaylistSong))).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
	playlistRepository.AssertExpectations(t)
}

func TestBulkDeleteSongs_WhenDeleteSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewBulkDeleteSongs(songRepository, nil, nil)

	request := requests.BulkDeleteSongsRequest{
		IDs: []uuid.UUID{uuid.New()},
	}

	mockSongs := &[]model.Song{{ID: request.IDs[0]}}
	songRepository.On("GetAllByIDsWithAlbumsAndPlaylists", mock.IsType(mockSongs), request.IDs).
		Return(nil, mockSongs).
		Once()

	internalError := errors.New("internal error")
	songRepository.On("Delete", request.IDs).Return(internalError).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestBulkDeleteSongs_WhenPublishFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := song.NewBulkDeleteSongs(songRepository, nil, messagePublisherService)

	request := requests.BulkDeleteSongsRequest{
		IDs: []uuid.UUID{uuid.New()},
	}

	mockSongs := &[]model.Song{{ID: request.IDs[0]}}
	songRepository.On("GetAllByIDsWithAlbumsAndPlaylists", mock.IsType(mockSongs), request.IDs).
		Return(nil, mockSongs).
		Once()

	songRepository.On("Delete", request.IDs).Return(nil).Once()

	internalError := errors.New("internal error")
	messagePublisherService.On("Publish", topics.SongsDeletedTopic, *mockSongs).
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

func TestBulkDeleteSongs_WhenWithoutAlbumsOrPlaylists_ShouldDeleteSongs(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := song.NewBulkDeleteSongs(songRepository, nil, messagePublisherService)

	request := requests.BulkDeleteSongsRequest{
		IDs: []uuid.UUID{
			uuid.New(),
			uuid.New(),
			uuid.New(),
		},
	}

	mockSongs := &[]model.Song{
		{ID: request.IDs[0]},
		{ID: request.IDs[1]},
		{ID: request.IDs[2]},
	}

	songRepository.On("GetAllByIDsWithAlbumsAndPlaylists", mock.IsType(mockSongs), request.IDs).
		Return(nil, mockSongs).
		Once()
	songRepository.On("Delete", request.IDs).Return(nil).Once()

	messagePublisherService.On("Publish", topics.SongsDeletedTopic, *mockSongs).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)

	songRepository.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
}

func TestBulkDeleteSongs_WhenWithAlbums_ShouldDeleteSongsAndReorderAlbums(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := song.NewBulkDeleteSongs(songRepository, nil, messagePublisherService)

	request := requests.BulkDeleteSongsRequest{
		IDs: []uuid.UUID{
			uuid.New(),
			uuid.New(),
			uuid.New(),
		},
	}

	mockAlbums := []model.Album{
		{
			ID: uuid.New(),
			Songs: []model.Song{
				{ID: request.IDs[0], AlbumTrackNo: &[]uint{1}[0]},
				{ID: uuid.New(), AlbumTrackNo: &[]uint{2}[0]},
				{ID: uuid.New(), AlbumTrackNo: &[]uint{3}[0]},
				{ID: request.IDs[2], AlbumTrackNo: &[]uint{4}[0]},
				{ID: uuid.New(), AlbumTrackNo: &[]uint{5}[0]},
			},
		},
		{
			ID: uuid.New(),
			Songs: []model.Song{
				{ID: uuid.New(), AlbumTrackNo: &[]uint{1}[0]},
				{ID: request.IDs[1], AlbumTrackNo: &[]uint{2}[0]},
				{ID: uuid.New(), AlbumTrackNo: &[]uint{3}[0]},
			},
		},
	}

	mockSongs := &[]model.Song{
		{ID: request.IDs[0], Album: &mockAlbums[0], AlbumID: &mockAlbums[0].ID, AlbumTrackNo: &[]uint{1}[0]},
		{ID: request.IDs[1], Album: &mockAlbums[1], AlbumID: &mockAlbums[1].ID, AlbumTrackNo: &[]uint{2}[0]},
		{ID: request.IDs[2], Album: &mockAlbums[0], AlbumID: &mockAlbums[0].ID, AlbumTrackNo: &[]uint{3}[0]},
	}

	songRepository.On("GetAllByIDsWithAlbumsAndPlaylists", mock.IsType(mockSongs), request.IDs).
		Return(nil, mockSongs).
		Once()

	songRepository.On("UpdateAll", mock.IsType(mockSongs)).
		Run(func(args mock.Arguments) {
			newAlbumSongs := args.Get(0).(*[]model.Song)
			assert.NotEmpty(t, newAlbumSongs)
		}).
		Return(nil).
		Once()

	songRepository.On("Delete", request.IDs).Return(nil).Once()

	messagePublisherService.On("Publish", topics.SongsDeletedTopic, *mockSongs).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)

	songRepository.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
}

func TestBulkDeleteSongs_WhenWithPlaylists_ShouldDeleteSongsAndReorderPlaylists(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	playlistRepository := new(repository.PlaylistRepositoryMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := song.NewBulkDeleteSongs(songRepository, playlistRepository, messagePublisherService)

	request := requests.BulkDeleteSongsRequest{
		IDs: []uuid.UUID{
			uuid.New(),
			uuid.New(),
			uuid.New(),
		},
	}

	mockPlaylists := []model.Playlist{
		{
			ID: uuid.New(),
			PlaylistSongs: []model.PlaylistSong{
				{SongID: request.IDs[0], SongTrackNo: 1},
				{SongID: uuid.New(), SongTrackNo: 2},
				{SongID: request.IDs[0], SongTrackNo: 3},
				{SongID: uuid.New(), SongTrackNo: 4},
				{SongID: request.IDs[2], SongTrackNo: 5},
				{SongID: uuid.New(), SongTrackNo: 6},
			},
		},
		{
			ID: uuid.New(),
			PlaylistSongs: []model.PlaylistSong{
				{SongID: request.IDs[0], SongTrackNo: 1},
				{SongID: uuid.New(), SongTrackNo: 2},
				{SongID: request.IDs[1], SongTrackNo: 3},
				{SongID: uuid.New(), SongTrackNo: 4},
			},
		},
	}

	mockSongs := &[]model.Song{
		{ID: request.IDs[0], Playlists: mockPlaylists},
		{ID: request.IDs[1], Playlists: []model.Playlist{mockPlaylists[1]}},
		{ID: request.IDs[2], Playlists: []model.Playlist{mockPlaylists[0]}},
	}

	songRepository.On("GetAllByIDsWithAlbumsAndPlaylists", mock.IsType(mockSongs), request.IDs).
		Return(nil, mockSongs).
		Once()

	playlistRepository.On("UpdateAllPlaylistSongs", mock.IsType(new([]model.PlaylistSong))).
		Run(func(args mock.Arguments) {
			newPlaylistSongs := args.Get(0).(*[]model.PlaylistSong)
			assert.NotEmpty(t, newPlaylistSongs)
		}).
		Return(nil).
		Once()

	songRepository.On("Delete", request.IDs).Return(nil).Once()

	messagePublisherService.On("Publish", topics.SongsDeletedTopic, *mockSongs).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)

	songRepository.AssertExpectations(t)
	playlistRepository.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
}
