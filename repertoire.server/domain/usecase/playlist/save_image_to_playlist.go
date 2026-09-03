package playlist

import (
	"errors"
	"mime/multipart"
	"reflect"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/domain/provider"
	"repertoire/server/internal"
	"repertoire/server/internal/httperror"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"time"

	"github.com/google/uuid"
)

type SaveImageToPlaylist struct {
	playlistRepository      repository.PlaylistRepository
	storageFilePathProvider provider.StorageFilePathProvider
	storageService          service.StorageService
	messagePublisherService service.MessagePublisherService
}

func NewSaveImageToPlaylist(
	playlistRepository repository.PlaylistRepository,
	storageFilePathProvider provider.StorageFilePathProvider,
	storageService service.StorageService,
	messagePublisherService service.MessagePublisherService,
) SaveImageToPlaylist {
	return SaveImageToPlaylist{
		playlistRepository:      playlistRepository,
		storageFilePathProvider: storageFilePathProvider,
		storageService:          storageService,
		messagePublisherService: messagePublisherService,
	}
}

func (s SaveImageToPlaylist) Handle(file *multipart.FileHeader, id uuid.UUID) *httperror.ErrorCode {
	var playlist model.Playlist
	if err := s.playlistRepository.Get(&playlist, id); err != nil {
		return httperror.DatabaseError(err)
	}
	if reflect.ValueOf(playlist).IsZero() {
		return httperror.NotFoundError(errors.New("playlist not found"))
	}

	if playlist.ImageURL != nil {
		if errCode := s.storageService.DeleteFile(*playlist.ImageURL); errCode != nil {
			return errCode
		}
	}

	playlist.UpdatedAt = time.Now().UTC()
	imagePath := s.storageFilePathProvider.GetPlaylistImagePath(file, playlist)

	if errCode := s.storageService.Upload(file, imagePath); errCode != nil {
		return errCode
	}

	playlist.ImageURL = (*internal.FilePath)(&imagePath)
	if err := s.playlistRepository.Update(&playlist); err != nil {
		return httperror.DatabaseError(err)
	}

	if err := s.messagePublisherService.Publish(topics.PlaylistUpdatedTopic, playlist); err != nil {
		return httperror.MessagePublisherError(err)
	}

	return nil
}
