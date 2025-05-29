package member

import (
	"errors"
	"mime/multipart"
	"reflect"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/domain/provider"
	"repertoire/server/internal"
	"repertoire/server/internal/wrapper"
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

func (s SaveImageToBandMember) Handle(file *multipart.FileHeader, id uuid.UUID) *wrapper.ErrorCode {
	var member model.BandMember
	err := s.repository.GetBandMemberWithArtist(&member, id)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(member).IsZero() {
		return wrapper.NotFoundError(errors.New("band member not found"))
	}

	if member.ImageURL != nil {
		errCode := s.storageService.DeleteFile(*member.ImageURL)
		if errCode != nil {
			return errCode
		}
	}

	member.UpdatedAt = time.Now().UTC()
	imagePath := s.storageFilePathProvider.GetBandMemberImagePath(file, member)

	err = s.storageService.Upload(file, imagePath)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	member.ImageURL = (*internal.FilePath)(&imagePath)
	err = s.repository.UpdateBandMember(&member)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
