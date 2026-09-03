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
)

type AddSongsToAlbum struct {
	albumRepository         repository.AlbumRepository
	songRepository          repository.SongRepository
	messagePublisherService service.MessagePublisherService
}

func NewAddSongsToAlbum(
	albumRepository repository.AlbumRepository,
	songRepository repository.SongRepository,
	messagePublisherService service.MessagePublisherService,
) AddSongsToAlbum {
	return AddSongsToAlbum{
		albumRepository:         albumRepository,
		songRepository:          songRepository,
		messagePublisherService: messagePublisherService,
	}
}

func (a AddSongsToAlbum) Handle(request requests.AddSongsToAlbumRequest) *httperror.ErrorCode {
	var album model.Album
	if err := a.albumRepository.GetWithSongs(&album, request.ID); err != nil {
		return httperror.DatabaseError(err)
	}
	if reflect.ValueOf(album).IsZero() {
		return httperror.NotFoundError(errors.New("album not found"))
	}

	var songs []model.Song
	if err := a.songRepository.GetAllByIDs(&songs, request.SongIDs); err != nil {
		return httperror.DatabaseError(err)
	}
	if len(songs) != len(request.SongIDs) {
		return httperror.NotFoundError(errors.New("songs not found"))
	}

	songsLength := len(album.Songs) + 1
	for i, song := range songs {
		// if their artists don't match, or the song has an artist but the album doesn't, it results in failure
		// on the other hand, if the album has an artist and the song doesn't, it will inherit it (pass)
		if song.ArtistID != nil && (album.ArtistID == nil || *album.ArtistID != *song.ArtistID) {
			return httperror.ConflictError(errors.New("song " + song.ID.String() + song.Title + " and album do not share the same artist"))
		}
		if song.AlbumID != nil {
			return httperror.ConflictError(errors.New("song " + song.ID.String() + " already has an album"))
		}

		songs[i].AlbumID = &request.ID
		trackNo := uint(songsLength + i)
		songs[i].AlbumTrackNo = &trackNo
		songs[i].ArtistID = album.ArtistID // songs inherit album artist
		if songs[i].ReleaseDate == nil {
			songs[i].ReleaseDate = album.ReleaseDate // also inherit the release date if there is none
		}
	}

	if err := a.songRepository.UpdateAll(&songs); err != nil {
		return httperror.DatabaseError(err)
	}

	if err := a.messagePublisherService.Publish(topics.SongsUpdatedTopic, request.SongIDs); err != nil {
		return httperror.MessagePublisherError(err)
	}

	return nil
}
