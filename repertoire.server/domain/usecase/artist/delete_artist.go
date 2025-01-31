package artist

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/domain/provider"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type DeleteArtist struct {
	repository              repository.ArtistRepository
	storageService          service.StorageService
	storageFilePathProvider provider.StorageFilePathProvider
}

func NewDeleteArtist(
	repository repository.ArtistRepository,
	storageService service.StorageService,
	storageFilePathProvider provider.StorageFilePathProvider,
) DeleteArtist {
	return DeleteArtist{
		repository:              repository,
		storageService:          storageService,
		storageFilePathProvider: storageFilePathProvider,
	}
}

func (d DeleteArtist) Handle(request requests.DeleteArtistRequest) *wrapper.ErrorCode {
	var artist model.Artist
	err := d.repository.Get(&artist, request.ID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(artist).IsZero() {
		return wrapper.NotFoundError(errors.New("artist not found"))
	}

	if d.storageFilePathProvider.HasArtistFiles(artist) {
		directoryPath := d.storageFilePathProvider.GetArtistDirectoryPath(artist)
		err = d.storageService.DeleteDirectory(directoryPath)
		if err != nil {
			return wrapper.InternalServerError(err)
		}
	}

	if request.WithAlbums {
		err = d.repository.DeleteAlbums(request.ID)
		if err != nil {
			return wrapper.InternalServerError(err)
		}
	}
	if request.WithSongs {
		err = d.repository.DeleteSongs(request.ID)
		if err != nil {
			return wrapper.InternalServerError(err)
		}
	}
	err = d.repository.Delete(request.ID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	return nil
}
