package album

import (
	"cmp"
	"errors"
	"net/http"
	"repertoire/server/data/repository"
	"repertoire/server/model"
	"slices"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRemoveSongFromAlbum_WhenGetWithSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := RemoveSongFromAlbum{
		repository: albumRepository,
	}

	id := uuid.New()
	songID := uuid.New()

	// given - mocking
	internalError := errors.New("internal error")
	albumRepository.On("GetWithSongs", new(model.Album), id).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id, songID)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	albumRepository.AssertExpectations(t)
}

func TestRemoveSongFromAlbum_WhenAlbumIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := RemoveSongFromAlbum{
		repository: albumRepository,
	}

	id := uuid.New()
	songID := uuid.New()

	// given - mocking
	albumRepository.On("GetWithSongs", new(model.Album), id).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(id, songID)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "album not found", errCode.Error.Error())

	albumRepository.AssertExpectations(t)
}

func TestRemoveSongFromAlbum_WhenSongIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := RemoveSongFromAlbum{
		repository: albumRepository,
	}

	id := uuid.New()
	songID := uuid.New()

	// given - mocking
	album := &model.Album{
		ID: id,
		Songs: []model.Song{
			{ID: uuid.New(), AlbumTrackNo: &[]uint{1}[0]},
		},
	}
	albumRepository.On("GetWithSongs", new(model.Album), id).
		Return(nil, album).
		Once()

	// when
	errCode := _uut.Handle(id, songID)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song not found", errCode.Error.Error())

	albumRepository.AssertExpectations(t)
}

func TestRemoveSongFromAlbum_WhenUpdateWithAssociationsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := RemoveSongFromAlbum{
		repository: albumRepository,
	}

	id := uuid.New()
	songID := uuid.New()

	// given - mocking
	album := &model.Album{
		ID: id,
		Songs: []model.Song{
			{ID: songID, AlbumTrackNo: &[]uint{1}[0]},
		},
	}
	albumRepository.On("GetWithSongs", new(model.Album), id).
		Return(nil, album).
		Once()

	internalError := errors.New("internal error")
	albumRepository.On("UpdateWithAssociations", mock.IsType(album)).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id, songID)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	albumRepository.AssertExpectations(t)
}

func TestRemoveSongFromAlbum_WhenRemoveSongFromAlbumFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := RemoveSongFromAlbum{
		repository: albumRepository,
	}

	id := uuid.New()
	songID := uuid.New()

	// given - mocking
	album := &model.Album{
		ID: id,
		Songs: []model.Song{
			{ID: songID, AlbumTrackNo: &[]uint{1}[0]},
		},
	}
	albumRepository.On("GetWithSongs", new(model.Album), id).
		Return(nil, album).
		Once()

	albumRepository.On("UpdateWithAssociations", mock.IsType(album)).
		Return(nil).
		Once()

	internalError := errors.New("internal error")
	albumRepository.On("RemoveSong", mock.IsType(album), mock.IsType(new(model.Song))).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id, songID)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	albumRepository.AssertExpectations(t)
}

func TestRemoveSongFromAlbum_WhenIsValid_ShouldNotReturnAnyError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := RemoveSongFromAlbum{
		repository: albumRepository,
	}

	id := uuid.New()
	songID := uuid.New()

	// given - mocking
	album := &model.Album{
		ID: id,
		Songs: []model.Song{
			{ID: uuid.New(), AlbumTrackNo: &[]uint{1}[0]},
			{ID: songID, AlbumTrackNo: &[]uint{2}[0]},
			{ID: uuid.New(), AlbumTrackNo: &[]uint{3}[0]},
			{ID: uuid.New(), AlbumTrackNo: &[]uint{4}[0]},
		},
	}
	albumRepository.On("GetWithSongs", new(model.Album), id).
		Return(nil, album).
		Once()

	albumRepository.On("UpdateWithAssociations", mock.IsType(album)).
		Run(func(args mock.Arguments) {
			newAlbum := args.Get(0).(*model.Album)
			songs := slices.Clone(newAlbum.Songs)

			songs = slices.DeleteFunc(songs, func(s model.Song) bool {
				return s.ID == songID
			})

			slices.SortFunc(songs, func(a, b model.Song) int {
				return cmp.Compare(*a.AlbumTrackNo, *b.AlbumTrackNo)
			})
			for i, song := range songs {
				assert.Equal(t, uint(i)+1, *song.AlbumTrackNo)
			}
		}).
		Return(nil).
		Once()

	albumRepository.On("RemoveSong", mock.IsType(album), mock.IsType(new(model.Song))).
		Run(func(args mock.Arguments) {
			album := args.Get(0).(*model.Album)
			assert.Equal(t, id, album.ID)

			song := args.Get(1).(*model.Song)
			assert.Equal(t, songID, song.ID)
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(id, songID)

	// then
	assert.Nil(t, errCode)

	albumRepository.AssertExpectations(t)
}
