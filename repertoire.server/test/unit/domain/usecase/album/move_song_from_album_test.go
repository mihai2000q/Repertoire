package album

import (
	"cmp"
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	album2 "repertoire/server/domain/usecase/album"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"slices"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMoveSongFromAlbum_WhenGetAlbumFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := album2.MoveSongFromAlbum{
		repository: albumRepository,
	}

	request := requests.MoveSongFromAlbumRequest{
		ID:         uuid.New(),
		SongID:     uuid.New(),
		OverSongID: uuid.New(),
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

func TestMoveSongFromAlbum_WhenAlbumIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := album2.MoveSongFromAlbum{
		repository: albumRepository,
	}

	request := requests.MoveSongFromAlbumRequest{
		ID:         uuid.New(),
		SongID:     uuid.New(),
		OverSongID: uuid.New(),
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

func TestMoveSongFromAlbum_WhenSongIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := album2.MoveSongFromAlbum{
		repository: albumRepository,
	}

	request := requests.MoveSongFromAlbumRequest{
		ID:         uuid.New(),
		SongID:     uuid.New(),
		OverSongID: uuid.New(),
	}

	// given - mocking
	album := &model.Album{ID: uuid.New()}
	albumRepository.On("GetWithSongs", new(model.Album), request.ID).
		Return(nil, album).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song not found", errCode.Error.Error())

	albumRepository.AssertExpectations(t)
}

func TestMoveSongFromAlbum_WhenOverSongIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := album2.MoveSongFromAlbum{
		repository: albumRepository,
	}

	album := &model.Album{
		ID: uuid.New(),
		Songs: []model.Song{
			{ID: uuid.New(), AlbumTrackNo: &[]uint{1}[0]},
		},
	}

	request := requests.MoveSongFromAlbumRequest{
		ID:         album.ID,
		SongID:     album.Songs[0].ID,
		OverSongID: uuid.New(),
	}

	// given - mocking
	albumRepository.On("GetWithSongs", new(model.Album), request.ID).
		Return(nil, album).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "over song not found", errCode.Error.Error())

	albumRepository.AssertExpectations(t)
}

func TestMoveSongFromAlbum_WhenUpdateFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := album2.MoveSongFromAlbum{
		repository: albumRepository,
	}

	album := &model.Album{
		ID: uuid.New(),
		Songs: []model.Song{
			{ID: uuid.New(), AlbumTrackNo: &[]uint{1}[0]},
			{ID: uuid.New(), AlbumTrackNo: &[]uint{2}[0]},
		},
	}

	request := requests.MoveSongFromAlbumRequest{
		ID:         album.ID,
		SongID:     album.Songs[0].ID,
		OverSongID: album.Songs[1].ID,
	}

	// given - mocking
	albumRepository.On("GetWithSongs", new(model.Album), request.ID).
		Return(nil, album).
		Once()

	internalError := errors.New("internal error")
	albumRepository.On("UpdateWithAssociations", mock.IsType(new(model.Album))).
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

func TestMoveSongFromAlbum_WhenIsValid_ShouldNotReturnAnyError(t *testing.T) {
	tests := []struct {
		name      string
		album     *model.Album
		index     uint
		overIndex uint
	}{
		{
			"Use case 1",
			&model.Album{
				ID: uuid.New(),
				Songs: []model.Song{
					{ID: uuid.New(), AlbumTrackNo: &[]uint{1}[0]},
					{ID: uuid.New(), AlbumTrackNo: &[]uint{2}[0]},
					{ID: uuid.New(), AlbumTrackNo: &[]uint{3}[0]},
					{ID: uuid.New(), AlbumTrackNo: &[]uint{4}[0]},
					{ID: uuid.New(), AlbumTrackNo: &[]uint{5}[0]},
				},
			},
			1,
			3,
		},
		{
			"Use case 2",
			&model.Album{
				ID: uuid.New(),
				Songs: []model.Song{
					{ID: uuid.New(), AlbumTrackNo: &[]uint{1}[0]},
					{ID: uuid.New(), AlbumTrackNo: &[]uint{2}[0]},
					{ID: uuid.New(), AlbumTrackNo: &[]uint{3}[0]},
					{ID: uuid.New(), AlbumTrackNo: &[]uint{4}[0]},
					{ID: uuid.New(), AlbumTrackNo: &[]uint{5}[0]},
				},
			},
			3,
			1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			albumRepository := new(repository.AlbumRepositoryMock)
			_uut := album2.MoveSongFromAlbum{
				repository: albumRepository,
			}

			request := requests.MoveSongFromAlbumRequest{
				ID:         tt.album.ID,
				SongID:     tt.album.Songs[tt.index].ID,
				OverSongID: tt.album.Songs[tt.overIndex].ID,
			}

			// given - mocking
			albumRepository.On("GetWithSongs", new(model.Album), request.ID).
				Return(nil, tt.album).
				Once()

			albumRepository.On("UpdateWithAssociations", mock.IsType(new(model.Album))).
				Run(func(args mock.Arguments) {
					album := args.Get(0).(*model.Album)
					songs := slices.Clone(album.Songs)
					slices.SortFunc(songs, func(a, b model.Song) int {
						return cmp.Compare(*a.AlbumTrackNo, *b.AlbumTrackNo)
					})
					if tt.index < tt.overIndex {
						assert.Equal(t, songs[tt.overIndex-1].ID, request.OverSongID)
					} else if tt.index > tt.overIndex {
						assert.Equal(t, songs[tt.overIndex+1].ID, request.OverSongID)
					}
					assert.Equal(t, songs[tt.overIndex].ID, request.SongID)
					for i, song := range songs {
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
		})
	}
}
