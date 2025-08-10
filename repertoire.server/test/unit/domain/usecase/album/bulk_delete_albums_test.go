package album

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/album"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBulkDeleteAlbums_WhenGetAlbumsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := album.NewBulkDeleteAlbums(albumRepository, nil)

	request := requests.BulkDeleteAlbumsRequest{
		IDs: []uuid.UUID{uuid.New()},
	}

	internalError := errors.New("internal error")
	albumRepository.On("GetAllByIDs", new([]model.Album), request.IDs).
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

func TestBulkDeleteAlbums_WhenGetAlbumsWithSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := album.NewBulkDeleteAlbums(albumRepository, nil)

	request := requests.BulkDeleteAlbumsRequest{
		IDs:       []uuid.UUID{uuid.New()},
		WithSongs: true,
	}

	internalError := errors.New("internal error")
	albumRepository.On("GetAllByIDsWithSongs", new([]model.Album), request.IDs).
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

func TestBulkDeleteAlbums_WhenAlbumsAreLen0_ShouldReturnNotFoundError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := album.NewBulkDeleteAlbums(albumRepository, nil)

	request := requests.BulkDeleteAlbumsRequest{
		IDs: []uuid.UUID{uuid.New()},
	}

	albumRepository.On("GetAllByIDs", new([]model.Album), request.IDs).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "albums not found", errCode.Error.Error())

	albumRepository.AssertExpectations(t)
}

func TestBulkDeleteAlbums_WhenDeleteAlbumsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := album.NewBulkDeleteAlbums(albumRepository, nil)

	request := requests.BulkDeleteAlbumsRequest{
		IDs: []uuid.UUID{uuid.New()},
	}

	mockAlbums := &[]model.Album{{ID: request.IDs[0]}}
	albumRepository.On("GetAllByIDs", new([]model.Album), request.IDs).
		Return(nil, mockAlbums).
		Once()

	internalError := errors.New("internal error")
	albumRepository.On("Delete", request.IDs).Return(internalError).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	albumRepository.AssertExpectations(t)
}

func TestBulkDeleteAlbums_WhenDeleteAlbumsWithSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := album.NewBulkDeleteAlbums(albumRepository, nil)

	request := requests.BulkDeleteAlbumsRequest{
		IDs:       []uuid.UUID{uuid.New()},
		WithSongs: true,
	}

	mockAlbums := &[]model.Album{{ID: request.IDs[0]}}
	albumRepository.On("GetAllByIDsWithSongs", new([]model.Album), request.IDs).
		Return(nil, mockAlbums).
		Once()

	internalError := errors.New("internal error")
	albumRepository.On("DeleteWithSongs", request.IDs).Return(internalError).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	albumRepository.AssertExpectations(t)
}

func TestBulkDeleteAlbums_WhenPublishFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := album.NewBulkDeleteAlbums(albumRepository, messagePublisherService)

	request := requests.BulkDeleteAlbumsRequest{
		IDs: []uuid.UUID{uuid.New()},
	}

	mockAlbums := &[]model.Album{{ID: request.IDs[0]}}
	albumRepository.On("GetAllByIDs", new([]model.Album), request.IDs).
		Return(nil, mockAlbums).
		Once()

	albumRepository.On("Delete", request.IDs).Return(nil).Once()

	internalError := errors.New("internal error")
	messagePublisherService.On("Publish", topics.AlbumsDeletedTopic, *mockAlbums).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	albumRepository.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
}

func TestBulkDeleteAlbums_WhenSuccessful_ShouldDeleteAlbum(t *testing.T) {
	tests := []struct {
		name      string
		albums    []model.Album
		withSongs bool
	}{
		{
			"Without Songs",
			[]model.Album{{ID: uuid.New()}},
			false,
		},
		{
			"With Songs",
			[]model.Album{{ID: uuid.New()}},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			albumRepository := new(repository.AlbumRepositoryMock)
			messagePublisherService := new(service.MessagePublisherServiceMock)
			_uut := album.NewBulkDeleteAlbums(albumRepository, messagePublisherService)

			var ids []uuid.UUID
			for _, a := range tt.albums {
				ids = append(ids, a.ID)
			}
			request := requests.BulkDeleteAlbumsRequest{
				IDs:       ids,
				WithSongs: tt.withSongs,
			}

			if request.WithSongs {
				albumRepository.On("GetAllByIDsWithSongs", new([]model.Album), request.IDs).
					Return(nil, &tt.albums).
					Once()
			} else {
				albumRepository.On("GetAllByIDs", new([]model.Album), request.IDs).
					Return(nil, &tt.albums).
					Once()
			}

			if tt.withSongs {
				albumRepository.On("DeleteWithSongs", request.IDs).Return(nil).Once()
			} else {
				albumRepository.On("Delete", request.IDs).Return(nil).Once()
			}

			messagePublisherService.On("Publish", topics.AlbumsDeletedTopic, tt.albums).
				Return(nil).
				Once()

			// when
			errCode := _uut.Handle(request)

			// then
			assert.Nil(t, errCode)

			albumRepository.AssertExpectations(t)
			messagePublisherService.AssertExpectations(t)
		})
	}
}
