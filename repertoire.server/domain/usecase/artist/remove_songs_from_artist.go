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

type RemoveSongsFromArtist struct {
	songRepository          repository.SongRepository
	messagePublisherService service.MessagePublisherService
}

func NewRemoveSongsFromArtist(
	songRepository repository.SongRepository,
	messagePublisherService service.MessagePublisherService,
) RemoveSongsFromArtist {
	return RemoveSongsFromArtist{
		songRepository:          songRepository,
		messagePublisherService: messagePublisherService,
	}
}

func (r RemoveSongsFromArtist) Handle(request requests.RemoveSongsFromArtistRequest) *httperror.ErrorCode {
	var songs []model.Song
	if err := r.songRepository.GetAllByIDs(&songs, request.SongIDs); err != nil {
		return httperror.DatabaseError(err)
	}
	if len(songs) != len(request.SongIDs) {
		return httperror.NotFoundError(errors.New("songs not found"))
	}

	for i, song := range songs {
		if song.ArtistID == nil || *song.ArtistID != request.ID {
			return httperror.ConflictError(errors.New("song " + song.ID.String() + " is not owned by this artist"))
		}

		songs[i].ArtistID = nil
	}

	if err := r.songRepository.UpdateAll(&songs); err != nil {
		return httperror.DatabaseError(err)
	}

	if err := r.messagePublisherService.Publish(topics.SongsUpdatedTopic, request.SongIDs); err != nil {
		return httperror.MessagePublisherError(err)
	}

	return nil
}
