package types

import (
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"slices"

	"github.com/google/uuid"
)

type DeleteSongSectionType struct {
	repository repository.SongRepository
	jwtService service.JwtService
}

func NewDeleteSongSectionType(
	repository repository.SongRepository,
	jwtService service.JwtService,
) DeleteSongSectionType {
	return DeleteSongSectionType{
		repository: repository,
		jwtService: jwtService,
	}
}

func (d DeleteSongSectionType) Handle(id uuid.UUID, token string) *wrapper.ErrorCode {
	userID, errCode := d.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return errCode
	}

	var types []model.SongSectionType
	err := d.repository.GetSectionTypes(&types, userID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	index := slices.IndexFunc(types, func(s model.SongSectionType) bool {
		return s.ID == id
	})

	for i := index + 1; i < len(types); i++ {
		types[i].Order = uint(types[i].Order - 1)
	}

	err = d.repository.UpdateAllSectionTypes(&types)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	err = d.repository.DeleteSectionType(id)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}