package artist

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

type UpdateArtist struct {
	artistRepository        repository.ArtistRepository
	messagePublisherService service.MessagePublisherService
}

func NewUpdateArtist(
	artistRepository repository.ArtistRepository,
	messagePublisherService service.MessagePublisherService,
) UpdateArtist {
	return UpdateArtist{
		artistRepository:        artistRepository,
		messagePublisherService: messagePublisherService,
	}
}

func (u UpdateArtist) Handle(request requests.UpdateArtistRequest) *httperror.ErrorCode {
	var artist model.Artist
	if err := u.artistRepository.Get(&artist, request.ID); err != nil {
		return httperror.DatabaseError(err)
	}
	if reflect.ValueOf(artist).IsZero() {
		return httperror.NotFoundError(errors.New("artist not found"))
	}

	artist.Name = request.Name
	artist.IsBand = request.IsBand

	if err := u.artistRepository.Update(&artist); err != nil {
		return httperror.DatabaseError(err)
	}

	if err := u.messagePublisherService.Publish(topics.ArtistUpdatedTopic, artist.ID); err != nil {
		return httperror.MessagePublisherError(err)
	}

	return nil
}
