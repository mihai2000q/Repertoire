package artist

import (
	"errors"
	"mime/multipart"
	"reflect"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/domain/provider"
	"repertoire/server/internal"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type SaveImageToArtist struct {
	repository              repository.ArtistRepository
	storageFilePathProvider provider.StorageFilePathProvider
	storageService          service.StorageService
}

func NewSaveImageToArtist(
	repository repository.ArtistRepository,
	storageFilePathProvider provider.StorageFilePathProvider,
	storageService service.StorageService,
) SaveImageToArtist {
	return SaveImageToArtist{
		repository:              repository,
		storageFilePathProvider: storageFilePathProvider,
		storageService:          storageService,
	}
}

func (s SaveImageToArtist) Handle(file *multipart.FileHeader, id uuid.UUID) *wrapper.ErrorCode {
	var artist model.Artist
	err := s.repository.Get(&artist, id)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(artist).IsZero() {
		return wrapper.NotFoundError(errors.New("artist not found"))
	}

	imagePath := s.storageFilePathProvider.GetArtistImagePath(file, artist)

	err = s.storageService.Upload(file, imagePath)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	artist.ImageURL = (*internal.FilePath)(&imagePath)
	err = s.repository.Update(&artist)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
