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

func TestUpdateAlbum_WhenGetAlbumFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := &UpdateAlbum{
		repository: albumRepository,
	}
	request := requests.UpdateAlbumRequest{
		ID:    uuid.New(),
		Title: "New Album",
	}

	internalError := errors.New("internal error")
	albumRepository.On("Get", new(model.Album), request.ID).Return(internalError).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	albumRepository.AssertExpectations(t)
}

func TestUpdateAlbum_WhenAlbumIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := &UpdateAlbum{
		repository: albumRepository,
	}
	request := requests.UpdateAlbumRequest{
		ID:    uuid.New(),
		Title: "New Album",
	}

	albumRepository.On("Get", new(model.Album), request.ID).Return(nil).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "album not found", errCode.Error.Error())

	albumRepository.AssertExpectations(t)
}

func TestUpdateAlbum_WhenUpdateAlbumFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := &UpdateAlbum{
		repository: albumRepository,
	}
	request := requests.UpdateAlbumRequest{
		ID:    uuid.New(),
		Title: "New Album",
	}

	album := &model.Album{
		ID:    request.ID,
		Title: "Some Album",
	}

	albumRepository.On("Get", new(model.Album), request.ID).Return(nil, album).Once()
	internalError := errors.New("internal error")
	albumRepository.On("Update", mock.IsType(album)).
		Run(func(args mock.Arguments) {
			newAlbum := args.Get(0).(*model.Album)
			assert.Equal(t, request.Title, newAlbum.Title)
		}).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	albumRepository.AssertExpectations(t)
}

func TestUpdateAlbum_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := &UpdateAlbum{
		repository: albumRepository,
	}
	request := requests.UpdateAlbumRequest{
		ID:    uuid.New(),
		Title: "New Album",
	}

	album := &model.Album{
		ID:    request.ID,
		Title: "Some Album",
	}

	albumRepository.On("Get", new(model.Album), request.ID).Return(nil, album).Once()
	albumRepository.On("Update", mock.IsType(album)).
		Run(func(args mock.Arguments) {
			newAlbum := args.Get(0).(*model.Album)
			assert.Equal(t, request.Title, newAlbum.Title)
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)

	albumRepository.AssertExpectations(t)
}
