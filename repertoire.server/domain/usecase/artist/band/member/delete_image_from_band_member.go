package member

import (
	"errors"
	"reflect"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type DeleteImageFromBandMember struct {
	repository     repository.ArtistRepository
	storageService service.StorageService
}

func NewDeleteImageFromBandMember(
	repository repository.ArtistRepository,
	storageService service.StorageService,
) DeleteImageFromBandMember {
	return DeleteImageFromBandMember{
		repository:     repository,
		storageService: storageService,
	}
}

func (d DeleteImageFromBandMember) Handle(id uuid.UUID) *wrapper.ErrorCode {
	var member model.BandMember
	err := d.repository.GetBandMember(&member, id)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(member).IsZero() {
		return wrapper.NotFoundError(errors.New("band member not found"))
	}
	if member.ImageURL == nil {
		return wrapper.BadRequestError(errors.New("band member does not have an image"))
	}

	errCode := d.storageService.DeleteFile(*member.ImageURL)
	if errCode != nil {
		return errCode
	}

	member.ImageURL = nil
	err = d.repository.UpdateBandMember(&member)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
