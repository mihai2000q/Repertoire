package artist

import (
	"errors"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/message/topics"
	"repertoire/server/internal/wrapper"
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

func (r RemoveSongsFromArtist) Handle(request requests.RemoveSongsFromArtistRequest) *wrapper.ErrorCode {
	var songs []model.Song
	err := r.songRepository.GetAllByIDs(&songs, request.SongIDs)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	for i, song := range songs {
		if song.ArtistID == nil || *song.ArtistID != request.ID {
			return wrapper.ConflictError(errors.New("song " + song.ID.String() + " is not owned by this artist"))
		}

		songs[i].ArtistID = nil
	}

	err = r.songRepository.UpdateAll(&songs)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	err = r.messagePublisherService.Publish(topics.SongsUpdatedTopic, request.SongIDs)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
