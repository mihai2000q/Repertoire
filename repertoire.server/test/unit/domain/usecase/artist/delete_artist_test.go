package artist

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/artist"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"repertoire/server/test/unit/domain/provider"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeleteArtist_WhenGetArtistFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := artist.NewDeleteArtist(artistRepository, nil, nil)

	request := requests.DeleteArtistRequest{
		ID: uuid.New(),
	}

	internalError := errors.New("internal error")
	artistRepository.On("Get", new(model.Artist), request.ID).Return(internalError).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	artistRepository.AssertExpectations(t)
}

func TestDeleteArtist_WhenArtistIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := artist.NewDeleteArtist(artistRepository, nil, nil)

	request := requests.DeleteArtistRequest{
		ID: uuid.New(),
	}

	artistRepository.On("Get", new(model.Artist), request.ID).Return(nil).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "artist not found", errCode.Error.Error())

	artistRepository.AssertExpectations(t)
}

func TestDeleteArtist_WhenDeleteDirectoryFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	storageService := new(service.StorageServiceMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	_uut := artist.NewDeleteArtist(artistRepository, storageService, storageFilePathProvider)

	request := requests.DeleteArtistRequest{
		ID: uuid.New(),
	}

	mockArtist := &model.Artist{
		ID: request.ID,
	}
	artistRepository.On("Get", new(model.Artist), request.ID).Return(nil, mockArtist).Once()

	storageFilePathProvider.On("HasArtistFiles", *mockArtist).Return(true).Once()

	directoryPath := "some directory path"
	storageFilePathProvider.On("GetArtistDirectoryPath", *mockArtist).Return(directoryPath).Once()

	internalError := wrapper.InternalServerError(errors.New("internal error"))
	storageService.On("DeleteDirectory", directoryPath).Return(internalError).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, internalError, errCode)

	artistRepository.AssertExpectations(t)
	storageService.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
}

func TestDeleteArtist_WhenDeleteArtistFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	_uut := artist.NewDeleteArtist(artistRepository, nil, storageFilePathProvider)

	request := requests.DeleteArtistRequest{
		ID: uuid.New(),
	}

	mockArtist := &model.Artist{
		ID: request.ID,
	}
	artistRepository.On("Get", new(model.Artist), request.ID).Return(nil, mockArtist).Once()

	storageFilePathProvider.On("HasArtistFiles", *mockArtist).Return(false).Once()

	internalError := errors.New("internal error")
	artistRepository.On("Delete", request.ID).Return(internalError).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	artistRepository.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
}
func TestDeleteArtist_WhenDeleteArtistAlbumsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	_uut := artist.NewDeleteArtist(artistRepository, nil, storageFilePathProvider)

	request := requests.DeleteArtistRequest{
		ID:         uuid.New(),
		WithAlbums: true,
	}

	mockArtist := &model.Artist{
		ID: request.ID,
	}
	artistRepository.On("Get", new(model.Artist), request.ID).Return(nil, mockArtist).Once()

	storageFilePathProvider.On("HasArtistFiles", *mockArtist).Return(false).Once()

	internalError := errors.New("internal error")
	artistRepository.On("DeleteAlbums", request.ID).Return(internalError).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	artistRepository.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
}

func TestDeleteArtist_WhenDeleteArtistSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	storageFilePathProvider := new(provider.StorageFilePathProviderMock)
	_uut := artist.NewDeleteArtist(artistRepository, nil, storageFilePathProvider)

	request := requests.DeleteArtistRequest{
		ID:        uuid.New(),
		WithSongs: true,
	}

	mockArtist := &model.Artist{
		ID: request.ID,
	}
	artistRepository.On("Get", new(model.Artist), request.ID).Return(nil, mockArtist).Once()

	storageFilePathProvider.On("HasArtistFiles", *mockArtist).Return(false).Once()

	internalError := errors.New("internal error")
	artistRepository.On("DeleteSongs", request.ID).Return(internalError).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	artistRepository.AssertExpectations(t)
	storageFilePathProvider.AssertExpectations(t)
}

func TestDeleteArtist_WhenSuccessful_ShouldDeleteArtist(t *testing.T) {
	tests := []struct {
		name       string
		artist     model.Artist
		hasFiles   bool
		withAlbums bool
		withSongs  bool
	}{
		{
			"Without Files",
			model.Artist{ID: uuid.New()},
			false,
			false,
			false,
		},
		{
			"With Files",
			model.Artist{ID: uuid.New()},
			true,
			false,
			false,
		},
		{
			"With Albums",
			model.Artist{ID: uuid.New()},
			true,
			true,
			false,
		},
		{
			"With Songs",
			model.Artist{ID: uuid.New()},
			true,
			false,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			artistRepository := new(repository.ArtistRepositoryMock)
			storageService := new(service.StorageServiceMock)
			storageFilePathProvider := new(provider.StorageFilePathProviderMock)
			_uut := artist.NewDeleteArtist(artistRepository, storageService, storageFilePathProvider)

			request := requests.DeleteArtistRequest{
				ID:         tt.artist.ID,
				WithAlbums: tt.withAlbums,
				WithSongs:  tt.withSongs,
			}

			artistRepository.On("Get", new(model.Artist), request.ID).Return(nil, &tt.artist).Once()

			storageFilePathProvider.On("HasArtistFiles", tt.artist).Return(tt.hasFiles).Once()

			if tt.hasFiles {
				directoryPath := "some directory path"
				storageFilePathProvider.On("GetArtistDirectoryPath", tt.artist).
					Return(directoryPath).
					Once()

				storageService.On("DeleteDirectory", directoryPath).
					Return(nil).
					Once()
			}

			if tt.withAlbums {
				artistRepository.On("DeleteAlbums", request.ID).Return(nil).Once()
			}
			if tt.withSongs {
				artistRepository.On("DeleteSongs", request.ID).Return(nil).Once()
			}

			artistRepository.On("Delete", request.ID).Return(nil).Once()

			// when
			errCode := _uut.Handle(request)

			// then
			assert.Nil(t, errCode)

			artistRepository.AssertExpectations(t)
			storageService.AssertExpectations(t)
			storageFilePathProvider.AssertExpectations(t)
		})
	}
}
