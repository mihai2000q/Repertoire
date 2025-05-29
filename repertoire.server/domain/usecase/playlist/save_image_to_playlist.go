package playlist

import (
	"errors"
	"mime/multipart"
	"reflect"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/domain/provider"
	"repertoire/server/internal"
	"repertoire/server/internal/message/topics"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"time"

	"github.com/google/uuid"
)

type SaveImageToPlaylist struct {
	repository              repository.PlaylistRepository
	storageFilePathProvider provider.StorageFilePathProvider
	storageService          service.StorageService
	messagePublisherService service.MessagePublisherService
}

func NewSaveImageToPlaylist(
	repository repository.PlaylistRepository,
	storageFilePathProvider provider.StorageFilePathProvider,
	storageService service.StorageService,
	messagePublisherService service.MessagePublisherService,
) SaveImageToPlaylist {
	return SaveImageToPlaylist{
		repository:              repository,
		storageFilePathProvider: storageFilePathProvider,
		storageService:          storageService,
		messagePublisherService: messagePublisherService,
	}
}

func (s SaveImageToPlaylist) Handle(file *multipart.FileHeader, id uuid.UUID) *wrapper.ErrorCode {
	var playlist model.Playlist
	err := s.repository.Get(&playlist, id)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(playlist).IsZero() {
		return wrapper.NotFoundError(errors.New("playlist not found"))
	}

	if playlist.ImageURL != nil {
		errCode := s.storageService.DeleteFile(*playlist.ImageURL)
		if errCode != nil {
			return errCode
		}
	}

	playlist.UpdatedAt = time.Now().UTC()
	imagePath := s.storageFilePathProvider.GetPlaylistImagePath(file, playlist)

	err = s.storageService.Upload(file, imagePath)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	playlist.ImageURL = (*internal.FilePath)(&imagePath)
	err = s.repository.Update(&playlist)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	err = s.messagePublisherService.Publish(topics.PlaylistUpdatedTopic, playlist)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
