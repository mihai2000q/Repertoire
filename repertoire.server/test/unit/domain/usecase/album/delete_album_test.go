package album

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/album"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"repertoire/server/test/unit/domain/provider"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeleteAlbum_WhenGetAlbumFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := album.NewDeleteAlbum(albumRepository, nil, nil)

	request := requests.DeleteAlbumRequest{
		ID: uuid.New(),
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

func TestDeleteAlbum_WhenAlbumIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := album.NewDeleteAlbum(albumRepository, nil, nil)

	request := requests.DeleteAlbumRequest{
		ID: uuid.New(),
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

func TestDeleteAlbum_WhenDeleteDirectoryFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	storageService := new(service.StorageServiceMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	_uut := album.NewDeleteAlbum(albumRepository, storageService, storageFilePathProvider)

	request := requests.DeleteAlbumRequest{
		ID: uuid.New(),
	}

	mockAlbum := &model.Album{
		ID: request.ID,
	}
	albumRepository.On("Get", new(model.Album), request.ID).Return(nil, mockAlbum).Once()

	storageFilePathProvider.On("HasAlbumFiles", *mockAlbum).Return(true).Once()

	directoryPath := "some directory path"
	storageFilePathProvider.On("GetAlbumDirectoryPath", *mockAlbum).Return(directoryPath).Once()

	internalError := wrapper.InternalServerError(errors.New("internal error"))
	storageService.On("DeleteDirectory", directoryPath).Return(internalError).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, internalError, errCode)

	albumRepository.AssertExpectations(t)
	storageService.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
}

func TestDeleteAlbum_WhenDeleteAlbumFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	_uut := album.NewDeleteAlbum(albumRepository, nil, storageFilePathProvider)

	request := requests.DeleteAlbumRequest{
		ID: uuid.New(),
	}

	mockAlbum := &model.Album{
		ID: request.ID,
	}
	albumRepository.On("Get", new(model.Album), request.ID).Return(nil, mockAlbum).Once()

	storageFilePathProvider.On("HasAlbumFiles", *mockAlbum).Return(false).Once()

	internalError := errors.New("internal error")
	albumRepository.On("Delete", request.ID).Return(internalError).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	albumRepository.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
}

func TestDeleteAlbum_WhenDeleteAlbumWithSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	_uut := album.NewDeleteAlbum(albumRepository, nil, storageFilePathProvider)

	request := requests.DeleteAlbumRequest{
		ID:        uuid.New(),
		WithSongs: true,
	}

	mockAlbum := &model.Album{
		ID: request.ID,
	}
	albumRepository.On("Get", new(model.Album), request.ID).Return(nil, mockAlbum).Once()

	storageFilePathProvider.On("HasAlbumFiles", *mockAlbum).Return(false).Once()

	internalError := errors.New("internal error")
	albumRepository.On("DeleteWithSongs", request.ID).Return(internalError).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	albumRepository.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
}

func TestDeleteAlbum_WhenSuccessful_ShouldDeleteAlbum(t *testing.T) {
	tests := []struct {
		name      string
		album     model.Album
		hasFiles  bool
		withSongs bool
	}{
		{
			"Without Files",
			model.Album{ID: uuid.New()},
			false,
			false,
		},
		{
			"With Files",
			model.Album{ID: uuid.New()},
			true,
			false,
		},
		{
			"With Songs",
			model.Album{ID: uuid.New()},
			false,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			albumRepository := new(repository.AlbumRepositoryMock)
			storageService := new(service.StorageServiceMock)
			storageFilePathProvider := new(provider.StorageFilePathProviderMock)
			_uut := album.NewDeleteAlbum(albumRepository, storageService, storageFilePathProvider)

			request := requests.DeleteAlbumRequest{
				ID:        tt.album.ID,
				WithSongs: tt.withSongs,
			}

			albumRepository.On("Get", new(model.Album), request.ID).Return(nil, &tt.album).Once()

			storageFilePathProvider.On("HasAlbumFiles", tt.album).Return(tt.hasFiles).Once()

			if tt.hasFiles {
				directoryPath := "some directory path"
				storageFilePathProvider.On("GetAlbumDirectoryPath", tt.album).
					Return(directoryPath).
					Once()

				storageService.On("DeleteDirectory", directoryPath).
					Return(nil).
					Once()
			}

			if tt.withSongs {
				albumRepository.On("DeleteWithSongs", request.ID).Return(nil).Once()
			} else {
				albumRepository.On("Delete", request.ID).Return(nil).Once()
			}

			// when
			errCode := _uut.Handle(request)

			// then
			assert.Nil(t, errCode)

			albumRepository.AssertExpectations(t)
			storageService.AssertExpectations(t)
			storageFilePathProvider.AssertExpectations(t)
		})
	}
}
