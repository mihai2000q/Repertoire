package album

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/httperror"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type UpdateAlbum struct {
	songRepository          repository.SongRepository
	albumRepository         repository.AlbumRepository
	messagePublisherService service.MessagePublisherService
}

func NewUpdateAlbum(
	songRepository repository.SongRepository,
	albumRepository repository.AlbumRepository,
	messagePublisherService service.MessagePublisherService,
) UpdateAlbum {
	return UpdateAlbum{
		songRepository:          songRepository,
		albumRepository:         albumRepository,
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

	artistHasChanged := album.ArtistID != nil && request.ArtistID == nil ||
		album.ArtistID == nil && request.ArtistID != nil ||
		album.ArtistID != nil && request.ArtistID != nil && *album.ArtistID != *request.ArtistID

	album.Title = request.Title
	album.ReleaseDate = request.ReleaseDate
	album.ArtistID = request.ArtistID

	if err := u.repository.Update(&album); err != nil {
		return httperror.DatabaseError(err)
	}

	if artistHasChanged {
		if errCode := u.updateAlbumSongsArtist(request); errCode != nil {
			return errCode
		}
	}

	if err := u.messagePublisherService.Publish(topics.AlbumsUpdatedTopic, []uuid.UUID{album.ID}); err != nil {
		return httperror.MessagePublisherError(err)
	}

	return nil
}

func (u UpdateAlbum) updateAlbumSongsArtist(request requests.UpdateAlbumRequest) *httperror.ErrorCode {
	var songs []model.Song
	if err := u.songRepository.GetAllByAlbum(&songs, request.ID); err != nil {
		return httperror.DatabaseError(err)
	}

	for i := range songs {
		songs[i].ArtistID = request.ArtistID
	}
	if err := u.songRepository.UpdateAll(&songs); err != nil {
		return httperror.DatabaseError(err)
	}

	return nil
}
