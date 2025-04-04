package artist

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/message/topics"
	"repertoire/server/internal/wrapper"
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

func (u UpdateArtist) Handle(request requests.UpdateArtistRequest) *wrapper.ErrorCode {
	var artist model.Artist
	err := u.repository.Get(&artist, request.ID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(artist).IsZero() {
		return wrapper.NotFoundError(errors.New("artist not found"))
	}

	artist.Name = request.Name
	artist.IsBand = request.IsBand

	err = u.repository.Update(&artist)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	err = u.messagePublisherService.Publish(topics.ArtistUpdatedTopic, artist.ID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
