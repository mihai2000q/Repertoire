package bandmember

import (
	"errors"
	"reflect"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/httperror"
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

func (d DeleteImageFromBandMember) Handle(id uuid.UUID) *httperror.ErrorCode {
	var member model.BandMember
	if err := d.repository.GetBandMember(&member, id); err != nil {
		return httperror.DatabaseError(err)
	}
	if reflect.ValueOf(member).IsZero() {
		return httperror.NotFoundError(errors.New("band member not found"))
	}
	if member.ImageURL == nil {
		return httperror.ConflictError(errors.New("band member does not have an image"))
	}

	if errCode := d.storageService.DeleteFile(*member.ImageURL); errCode != nil {
		return errCode
	}

	member.ImageURL = nil
	if err := d.repository.UpdateBandMember(&member); err != nil {
		return httperror.DatabaseError(err)
	}

	return nil
}
