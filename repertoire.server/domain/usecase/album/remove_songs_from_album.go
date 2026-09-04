package album

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/database/transaction"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/httperror"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type RemoveSongsFromAlbum struct {
	albumRepository         repository.AlbumRepository
	transaction             transaction.Manager
	messagePublisherService service.MessagePublisherService
}

func NewRemoveSongsFromAlbum(
	albumRepository repository.AlbumRepository,
	transaction transaction.Manager,
	messagePublisherService service.MessagePublisherService,
) RemoveSongsFromAlbum {
	return RemoveSongsFromAlbum{
		albumRepository:         albumRepository,
		transaction:             transaction,
		messagePublisherService: messagePublisherService,
	}
}

func (r RemoveSongsFromAlbum) Handle(request requests.RemoveSongsFromAlbumRequest) *httperror.ErrorCode {
	var album model.Album
	if err := r.albumRepository.GetWithSongs(&album, request.ID); err != nil {
		return httperror.DatabaseError(err)
	}
	if reflect.ValueOf(album).IsZero() {
		return httperror.NotFoundError(errors.New("album not found"))
	}

	// map for easy lookup
	songIDsMap := make(map[uuid.UUID]bool)
	for _, songID := range request.SongIDs {
		songIDsMap[songID] = true
	}

	var songsToDelete []model.Song
	var songsToPreserve []model.Song
	albumTrackNo := uint(1)

	for _, song := range album.Songs {
		if songIDsMap[song.ID] {
			songsToDelete = append(songsToDelete, song)
		} else {
			// reorder preserved songs
			*song.AlbumTrackNo = albumTrackNo
			songsToPreserve = append(songsToPreserve, song)
			albumTrackNo++
		}
	}

	if len(songsToDelete) != len(request.SongIDs) {
		return httperror.NotFoundError(errors.New("songs not found"))
	}

	err := r.transaction.Execute(func(factory transaction.RepositoryFactory) error {
		txAlbumRepo := factory.NewAlbumRepository()

		if err := txAlbumRepo.RemoveSongs(&album, &songsToDelete); err != nil {
			return err
		}

		album.Songs = songsToPreserve
		if err := txAlbumRepo.UpdateWithAssociations(&album); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return httperror.DatabaseError(err)
	}

	if err = r.messagePublisherService.Publish(topics.SongsUpdatedTopic, request.SongIDs); err != nil {
		return httperror.MessagePublisherError(err)
	}

	return nil
}
