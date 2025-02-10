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

func TestRemoveSongsFromAlbum_WhenGetWithSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := album.NewRemoveSongsFromAlbum(albumRepository)

	request := requests.RemoveSongsFromAlbumRequest{
		ID:      uuid.New(),
		SongIDs: []uuid.UUID{uuid.New()},
	}

	// given - mocking
	internalError := errors.New("internal error")
	albumRepository.On("GetWithSongs", new(model.Album), request.ID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	albumRepository.AssertExpectations(t)
}

func TestRemoveSongsFromAlbum_WhenAlbumIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := album.NewRemoveSongsFromAlbum(albumRepository)

	request := requests.RemoveSongsFromAlbumRequest{
		ID:      uuid.New(),
		SongIDs: []uuid.UUID{uuid.New()},
	}

	// given - mocking
	albumRepository.On("GetWithSongs", new(model.Album), request.ID).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "album not found", errCode.Error.Error())

	albumRepository.AssertExpectations(t)
}

func TestRemoveSongsFromAlbum_WhenNotAllSongsFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := album.NewRemoveSongsFromAlbum(albumRepository)

	request := requests.RemoveSongsFromAlbumRequest{
		ID:      uuid.New(),
		SongIDs: []uuid.UUID{uuid.New(), uuid.New()},
	}

	// given - mocking
	mockAlbum := &model.Album{
		ID: request.ID,
		Songs: []model.Song{
			{ID: request.SongIDs[0], AlbumTrackNo: &[]uint{1}[0]},
			{ID: uuid.New(), AlbumTrackNo: &[]uint{2}[0]},
		},
	}
	albumRepository.On("GetWithSongs", new(model.Album), request.ID).
		Return(nil, mockAlbum).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "could not find all songs", errCode.Error.Error())

	albumRepository.AssertExpectations(t)
}

func TestRemoveSongsFromAlbum_WhenRemoveSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := album.NewRemoveSongsFromAlbum(albumRepository)

	request := requests.RemoveSongsFromAlbumRequest{
		ID:      uuid.New(),
		SongIDs: []uuid.UUID{uuid.New()},
	}

	// given - mocking
	mockAlbum := &model.Album{
		ID: request.ID,
		Songs: []model.Song{
			{ID: request.SongIDs[0], AlbumTrackNo: &[]uint{1}[0]},
		},
	}
	albumRepository.On("GetWithSongs", new(model.Album), request.ID).
		Return(nil, mockAlbum).
		Once()

	internalError := errors.New("internal error")
	albumRepository.On("RemoveSongs", mock.IsType(mockAlbum), mock.IsType(&mockAlbum.Songs)).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	albumRepository.AssertExpectations(t)
}

func TestRemoveSongsFromAlbum_WhenUpdateAlbumFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := album.NewRemoveSongsFromAlbum(albumRepository)

	request := requests.RemoveSongsFromAlbumRequest{
		ID:      uuid.New(),
		SongIDs: []uuid.UUID{uuid.New()},
	}

	// given - mocking
	mockAlbum := &model.Album{
		ID: request.ID,
		Songs: []model.Song{
			{ID: request.SongIDs[0], AlbumTrackNo: &[]uint{1}[0]},
		},
	}
	albumRepository.On("GetWithSongs", new(model.Album), request.ID).
		Return(nil, mockAlbum).
		Once()

	albumRepository.On("RemoveSongs", mock.IsType(mockAlbum), mock.IsType(&mockAlbum.Songs)).
		Return(nil).
		Once()

	internalError := errors.New("internal error")
	albumRepository.On("UpdateWithAssociations", mock.IsType(mockAlbum)).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	albumRepository.AssertExpectations(t)
}

func TestRemoveSongsFromAlbum_WhenIsValid_ShouldNotReturnAnyError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := album.NewRemoveSongsFromAlbum(albumRepository)

	request := requests.RemoveSongsFromAlbumRequest{
		ID:      uuid.New(),
		SongIDs: []uuid.UUID{uuid.New(), uuid.New()},
	}

	// given - mocking
	mockAlbum := &model.Album{
		ID: request.ID,
		Songs: []model.Song{
			{ID: uuid.New(), AlbumTrackNo: &[]uint{1}[0]},
			{ID: request.SongIDs[0], AlbumTrackNo: &[]uint{2}[0]},
			{ID: uuid.New(), AlbumTrackNo: &[]uint{3}[0]},
			{ID: request.SongIDs[1], AlbumTrackNo: &[]uint{4}[0]},
		},
	}
	albumRepository.On("GetWithSongs", new(model.Album), request.ID).
		Return(nil, mockAlbum).
		Once()

	albumRepository.On("RemoveSongs", mock.IsType(mockAlbum), mock.IsType(&mockAlbum.Songs)).
		Run(func(args mock.Arguments) {
			newAlbum := args.Get(0).(*model.Album)
			assert.Equal(t, request.ID, newAlbum.ID)

			songs := args.Get(1).(*[]model.Song)
			assert.Len(t, *songs, len(request.SongIDs))
			for _, song := range *songs {
				assert.Contains(t, request.SongIDs, song.ID)
			}
		}).
		Return(nil).
		Once()

	albumRepository.On("UpdateWithAssociations", mock.IsType(mockAlbum)).
		Run(func(args mock.Arguments) {
			newAlbum := args.Get(0).(*model.Album)
			assert.Len(t, newAlbum.Songs, len(mockAlbum.Songs)-len(request.SongIDs))
			for i, song := range newAlbum.Songs {
				assert.NotContains(t, request.SongIDs, song.ID)
				assert.Equal(t, uint(i)+1, *song.AlbumTrackNo)
			}
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)

	albumRepository.AssertExpectations(t)
}
