package artist

import (
	"errors"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/httperror"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type AddSongsToArtist struct {
	songRepository          repository.SongRepository
	messagePublisherService service.MessagePublisherService
}

func NewAddSongsToArtist(
	songRepository repository.SongRepository,
	messagePublisherService service.MessagePublisherService,
) AddSongsToArtist {
	return AddSongsToArtist{
		songRepository:          songRepository,
		messagePublisherService: messagePublisherService,
	}
}

func (a AddSongsToArtist) Handle(request requests.AddSongsToArtistRequest) *httperror.ErrorCode {
	var songs []model.Song
	if err := a.songRepository.GetAllByIDsWithAlbumSongs(&songs, request.SongIDs); err != nil {
		return httperror.DatabaseError(err)
	}
	if len(songs) != len(request.SongIDs) {
		return httperror.NotFoundError(errors.New("songs not found"))
	}

	var albumIDs []uuid.UUID
	var songIDs []uuid.UUID
	for i, song := range songs {
		if song.ArtistID != nil {
			return httperror.ConflictError(errors.New("song " + song.ID.String() + " already has an artist"))
		}

		songs[i].ArtistID = &request.ID

		// update the whole album's artist, including the other songs
		if song.Album != nil {
			songs[i].Album.ArtistID = &request.ID
			albumIDs = append(albumIDs, song.Album.ID)
			for j := range song.Album.Songs {
				songs[i].Album.Songs[j].ArtistID = &request.ID
			}
		} else {
			songIDs = append(songIDs, song.ID)
		}
	}

	if err := a.songRepository.UpdateAllWithAssociations(&songs); err != nil {
		return httperror.DatabaseError(err)
	}

	if len(songIDs) > 0 {
		if err := a.messagePublisherService.Publish(topics.SongsUpdatedTopic, songIDs); err != nil {
			return httperror.MessagePublisherError(err)
		}
	}
	if len(albumIDs) > 0 {
		if err := a.messagePublisherService.Publish(topics.AlbumsUpdatedTopic, albumIDs); err != nil {
			return httperror.MessagePublisherError(err)
		}
	}

	return nil
}
