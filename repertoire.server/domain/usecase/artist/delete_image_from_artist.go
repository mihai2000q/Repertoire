package artist

import (
	"errors"
	"reflect"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/message/topics"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type DeleteImageFromArtist struct {
	repository              repository.ArtistRepository
	storageService          service.StorageService
	messagePublisherService service.MessagePublisherService
}

func NewDeleteImageFromArtist(
	repository repository.ArtistRepository,
	storageService service.StorageService,
	messagePublisherService service.MessagePublisherService,
) DeleteImageFromArtist {
	return DeleteImageFromArtist{
		repository:              repository,
		storageService:          storageService,
		messagePublisherService: messagePublisherService,
	}
}

func (d DeleteImageFromArtist) Handle(id uuid.UUID) *wrapper.ErrorCode {
	var artist model.Artist
	err := d.repository.Get(&artist, id)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(artist).IsZero() {
		return wrapper.NotFoundError(errors.New("artist not found"))
	}
	if artist.ImageURL == nil {
		return wrapper.ConflictError(errors.New("artist does not have an image"))
	}

	errCode := d.storageService.DeleteFile(*artist.ImageURL)
	if errCode != nil {
		return errCode
	}

	artist.ImageURL = nil
	err = d.repository.Update(&artist)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	err = d.messagePublisherService.Publish(topics.ArtistUpdatedTopic, artist.ID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
