package bandmember

import (
	"errors"
	"mime/multipart"
	"reflect"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/domain/provider"
	"repertoire/server/internal"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"
	"time"

	"github.com/google/uuid"
)

type SaveImageToBandMember struct {
	repository              repository.ArtistRepository
	storageFilePathProvider provider.StorageFilePathProvider
	storageService          service.StorageService
}

func NewSaveImageToBandMember(
	repository repository.ArtistRepository,
	storageFilePathProvider provider.StorageFilePathProvider,
	storageService service.StorageService,
) SaveImageToBandMember {
	return SaveImageToBandMember{
		repository:              repository,
		storageFilePathProvider: storageFilePathProvider,
		storageService:          storageService,
	}
}

func (s SaveImageToBandMember) Handle(file *multipart.FileHeader, id uuid.UUID) *httperror.ErrorCode {
	var member model.BandMember
	if err := s.repository.GetBandMemberWithArtist(&member, id); err != nil {
		return httperror.DatabaseError(err)
	}
	if reflect.ValueOf(member).IsZero() {
		return httperror.NotFoundError(errors.New("band member not found"))
	}

	if member.ImageURL != nil {
		if errCode := s.storageService.DeleteFile(*member.ImageURL); errCode != nil {
			return errCode
		}
	}

	member.UpdatedAt = time.Now().UTC()
	imagePath := s.storageFilePathProvider.GetBandMemberImagePath(file, member)

	if errCode := s.storageService.Upload(file, imagePath); errCode != nil {
		return errCode
	}

	member.ImageURL = (*internal.FilePath)(&imagePath)
	if err := s.repository.UpdateBandMember(&member); err != nil {
		return httperror.DatabaseError(err)
	}

	return nil
}
