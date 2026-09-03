package types

import (
	"errors"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"
	"slices"

	"github.com/google/uuid"
)

type DeleteSongSectionType struct {
	repository repository.UserDataRepository
	jwtService service.JwtService
}

func NewDeleteSongSectionType(
	repository repository.UserDataRepository,
	jwtService service.JwtService,
) DeleteSongSectionType {
	return DeleteSongSectionType{
		repository: repository,
		jwtService: jwtService,
	}
}

func (d DeleteSongSectionType) Handle(id uuid.UUID, token string) *httperror.ErrorCode {
	userID, errCode := d.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return errCode
	}

	var types []model.SongSectionType
	if err := d.repository.GetSectionTypes(&types, userID); err != nil {
		return httperror.DatabaseError(err)
	}

	index := slices.IndexFunc(types, func(s model.SongSectionType) bool {
		return s.ID == id
	})
	if index == -1 {
		return httperror.NotFoundError(errors.New("song section type not found"))
	}

	for i := index + 1; i < len(types); i++ {
		types[i].Order = types[i].Order - 1
	}

	if err := d.repository.UpdateAllSectionTypes(&types); err != nil {
		return httperror.DatabaseError(err)
	}

	if err := d.repository.DeleteSectionType(id); err != nil {
		return httperror.DatabaseError(err)
	}

	return nil
}
