package song

import (
	"errors"
	"github.com/google/uuid"
	"mime/multipart"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/domain/provider"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type SaveImageToSong struct {
	repository              repository.SongRepository
	storageFilePathProvider provider.StorageFilePathProvider
	storageService          service.StorageService
}

func NewSaveImageToSong(
	repository repository.SongRepository,
	storageFilePathProvider provider.StorageFilePathProvider,
	storageService service.StorageService,
) SaveImageToSong {
	return SaveImageToSong{
		repository:              repository,
		storageFilePathProvider: storageFilePathProvider,
		storageService:          storageService,
	}
}

func (a SaveImageToSong) Handle(file *multipart.FileHeader, songID uuid.UUID, token string) *wrapper.ErrorCode {
	var song model.Song
	err := a.repository.Get(&song, songID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if song.ID == uuid.Nil {
		return wrapper.NotFoundError(errors.New("song not found"))
	}

	filePath, imageUrl := a.storageFilePathProvider.GetSongImagePathAndURL(file, songID)

	err = a.storageService.Upload(token, file, filePath)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	song.ImageURL = &imageUrl
	err = a.repository.Update(&song)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
