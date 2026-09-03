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
	repository              repository.ArtistRepository
	messagePublisherService service.MessagePublisherService
}

func NewUpdateArtist(
	repository repository.ArtistRepository,
	messagePublisherService service.MessagePublisherService,
) UpdateArtist {
	return UpdateArtist{
		repository:              repository,
		messagePublisherService: messagePublisherService,
	}
}

func (u UpdateArtist) Handle(request requests.UpdateArtistRequest) *httperror.ErrorCode {
	var artist model.Artist
	if err := u.repository.Get(&artist, request.ID); err != nil {
		return httperror.DatabaseError(err)
	}
	if reflect.ValueOf(artist).IsZero() {
		return httperror.NotFoundError(errors.New("artist not found"))
	}

	artist.Name = request.Name
	artist.IsBand = request.IsBand

	if err := u.repository.Update(&artist); err != nil {
		return httperror.DatabaseError(err)
	}

	if err := u.messagePublisherService.Publish(topics.ArtistUpdatedTopic, artist.ID); err != nil {
		return httperror.MessagePublisherError(err)
	}

	return nil
}
