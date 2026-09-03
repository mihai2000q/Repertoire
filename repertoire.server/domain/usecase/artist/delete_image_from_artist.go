package artist

import (
	"errors"
	"reflect"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/httperror"
	"repertoire/server/internal/message/topics"
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

func (d DeleteImageFromArtist) Handle(id uuid.UUID) *httperror.ErrorCode {
	var artist model.Artist
	if err := d.repository.Get(&artist, id); err != nil {
		return httperror.DatabaseError(err)
	}
	if reflect.ValueOf(artist).IsZero() {
		return httperror.NotFoundError(errors.New("artist not found"))
	}
	if artist.ImageURL == nil {
		return nil
	}

	if errCode := d.storageService.DeleteFile(*artist.ImageURL); errCode != nil {
		return errCode
	}

	artist.ImageURL = nil
	if err := d.repository.Update(&artist); err != nil {
		return httperror.DatabaseError(err)
	}

	if err := d.messagePublisherService.Publish(topics.ArtistUpdatedTopic, artist.ID); err != nil {
		return httperror.MessagePublisherError(err)
	}

	return nil
}
