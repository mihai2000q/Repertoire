package album

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/model"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddSongToAlbum_WhenGetSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := AddSongToAlbum{repository: songRepository}

	request := requests.AddSongToAlbumRequest{
		ID:     uuid.New(),
		SongID: uuid.New(),
	}

	// given - mocking
	internalError := errors.New("internal error")
	songRepository.On("Get", mock.IsType(new(model.Song)), request.SongID).
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

func TestAddSongToAlbum_WhenCountSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := AddSongToAlbum{
		albumRepository: albumRepository,
		repository:      songRepository,
	}

	request := requests.AddSongToAlbumRequest{
		ID:     uuid.New(),
		SongID: uuid.New(),
	}

	// given - mocking
	song := &model.Song{ID: request.SongID}
	songRepository.On("Get", mock.IsType(song), request.SongID).
		Return(nil, song).
		Once()

	internalError := errors.New("internal error")
	albumRepository.On("CountSongs", mock.Anything, &request.ID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	albumRepository.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}

func TestAddSongToAlbum_WhenUpdateSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := AddSongToAlbum{
		albumRepository: albumRepository,
		repository:      songRepository,
	}

	request := requests.AddSongToAlbumRequest{
		ID:     uuid.New(),
		SongID: uuid.New(),
	}

	// given - mocking
	song := &model.Song{ID: request.SongID}
	songRepository.On("Get", mock.IsType(song), request.SongID).
		Return(nil, song).
		Once()

	var count = &[]int64{12}[0]
	albumRepository.On("CountSongs", mock.IsType(count), &request.ID).
		Return(nil, count).
		Once()

	internalError := errors.New("internal error")
	songRepository.On("Update", mock.IsType(song)).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	albumRepository.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}

func TestAddSongToAlbum_WhenIsValid_ShouldNotReturnAnyError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := AddSongToAlbum{
		albumRepository: albumRepository,
		repository:      songRepository,
	}

	request := requests.AddSongToAlbumRequest{
		ID:     uuid.New(),
		SongID: uuid.New(),
	}

	// given - mocking
	song := &model.Song{ID: request.SongID}
	songRepository.On("Get", mock.IsType(song), request.SongID).
		Return(nil, song).
		Once()

	var count = &[]int64{12}[0]
	albumRepository.On("CountSongs", mock.IsType(count), &request.ID).
		Return(nil, count).
		Once()

	songRepository.On("Update", mock.IsType(song)).
		Run(func(args mock.Arguments) {
			newSong := args.Get(0).(*model.Song)
			assert.Equal(t, *newSong.AlbumID, request.ID)
			assert.Equal(t, *newSong.AlbumTrackNo, uint(*count)+1)
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)

	albumRepository.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}
