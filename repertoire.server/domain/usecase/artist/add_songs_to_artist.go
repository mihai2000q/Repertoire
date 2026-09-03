package artist

import (
	"errors"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/httperror"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
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

	for i, song := range songs {
		if song.ArtistID != nil {
			return httperror.ConflictError(errors.New("song " + song.ID.String() + "already has an artist"))
		}

		// update the whole album's artist, including the other songs
		if song.Album != nil {
			songs[i].Album.ArtistID = &request.ID
			for j := range song.Album.Songs {
				songs[i].Album.Songs[j].ArtistID = &request.ID
			}
		} else {
			songs[i].ArtistID = &request.ID
		}
	}

	if err := a.songRepository.UpdateAllWithAssociations(&songs); err != nil {
		return httperror.DatabaseError(err)
	}

	if err := a.messagePublisherService.Publish(topics.SongsUpdatedTopic, request.SongIDs); err != nil {
		return httperror.MessagePublisherError(err)
	}

	return nil
}
