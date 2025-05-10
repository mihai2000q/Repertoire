package album

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/database/transaction"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/message/topics"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"slices"
)

type RemoveSongsFromAlbum struct {
	repository              repository.AlbumRepository
	transaction             transaction.Manager
	messagePublisherService service.MessagePublisherService
}

func NewRemoveSongsFromAlbum(
	repository repository.AlbumRepository,
	messagePublisherService service.MessagePublisherService,
) RemoveSongsFromAlbum {
	return RemoveSongsFromAlbum{
		repository:              repository,
		messagePublisherService: messagePublisherService,
	}
}

func (r RemoveSongsFromAlbum) Handle(request requests.RemoveSongsFromAlbumRequest) *wrapper.ErrorCode {
	var album model.Album
	err := r.repository.GetWithSongs(&album, request.ID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(album).IsZero() {
		return wrapper.NotFoundError(errors.New("album not found"))
	}

	var songsToDelete []model.Song
	var songsToPreserve []model.Song

	albumTrackNo := uint(1)
	for _, song := range album.Songs {
		if slices.Contains(request.SongIDs, song.ID) {
			songsToDelete = append(songsToDelete, song)
		} else {
			// reorder preserved songs
			*song.AlbumTrackNo = albumTrackNo
			songsToPreserve = append(songsToPreserve, song)
			albumTrackNo++
		}
	}

	if len(songsToDelete) != len(request.SongIDs) {
		return wrapper.NotFoundError(errors.New("could not find all songs"))
	}

	err = r.transaction.Execute(func(factory transaction.RepositoryFactory) error {
		albumRepo := factory.NewAlbumRepository()

		if err := albumRepo.RemoveSongs(&album, &songsToDelete); err != nil {
			return err
		}

		album.Songs = songsToPreserve
		if err := albumRepo.UpdateWithAssociations(&album); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	err = r.messagePublisherService.Publish(topics.SongsUpdatedTopic, request.SongIDs)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
