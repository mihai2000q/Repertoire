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

type UpdateAlbum struct {
	albumRepository         repository.AlbumRepository
	transactionManager      transaction.Manager
	messagePublisherService service.MessagePublisherService

	txSongRepo repository.SongRepository
}

func NewUpdateAlbum(
	albumRepository repository.AlbumRepository,
	transactionManager transaction.Manager,
	messagePublisherService service.MessagePublisherService,
) UpdateAlbum {
	return UpdateAlbum{
		albumRepository:         albumRepository,
		transactionManager:      transactionManager,
		messagePublisherService: messagePublisherService,
	}
}

func (u UpdateAlbum) Handle(request requests.UpdateAlbumRequest) *httperror.ErrorCode {
	var album model.Album
	if err := u.albumRepository.Get(&album, request.ID); err != nil {
		return httperror.DatabaseError(err)
	}
	if reflect.ValueOf(album).IsZero() {
		return httperror.NotFoundError(errors.New("album not found"))
	}

	err := u.transactionManager.Execute(func(factory transaction.RepositoryFactory) error {
		txAlbumRepo := factory.NewAlbumRepository()

		artistHasChanged := album.ArtistID != nil && request.ArtistID == nil ||
			album.ArtistID == nil && request.ArtistID != nil ||
			album.ArtistID != nil && request.ArtistID != nil && *album.ArtistID != *request.ArtistID

		album.Title = request.Title
		album.ReleaseDate = request.ReleaseDate
		album.ArtistID = request.ArtistID

		if err := txAlbumRepo.Update(&album); err != nil {
			return err
		}

		if artistHasChanged {
			u.txSongRepo = factory.NewSongRepository()
			if errCode := u.updateAlbumSongsArtist(request); errCode != nil {
				return errCode
			}
		}

		return nil
	})
	if err != nil {
		return httperror.DatabaseError(err)
	}

	if err := u.messagePublisherService.Publish(topics.AlbumsUpdatedTopic, []uuid.UUID{album.ID}); err != nil {
		return httperror.MessagePublisherError(err)
	}

	return nil
}

func (u UpdateAlbum) updateAlbumSongsArtist(request requests.UpdateAlbumRequest) error {
	var songs []model.Song
	if err := u.txSongRepo.GetAllByAlbum(&songs, request.ID); err != nil {
		return err
	}

	for i := range songs {
		songs[i].ArtistID = request.ArtistID
	}
	if err := u.txSongRepo.UpdateAll(&songs); err != nil {
		return err
	}

	return nil
}
