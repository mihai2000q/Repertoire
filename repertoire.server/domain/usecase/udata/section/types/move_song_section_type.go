package types

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/reorder"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type MoveSongSectionType struct {
	repository repository.UserDataRepository
	jwtService service.JwtService
}

func NewMoveSongSectionType(
	repository repository.UserDataRepository,
	jwtService service.JwtService,
) MoveSongSectionType {
	return MoveSongSectionType{
		repository: repository,
		jwtService: jwtService,
	}
}

func (m MoveSongSectionType) Handle(request requests.MoveSongSectionTypeRequest, token string) *wrapper.ErrorCode {
	userID, errCode := m.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return errCode
	}

	var types []model.SongSectionType
	err := m.repository.GetSectionTypes(&types, userID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	errCode = reorder.MoveEntity(
		types,
		request.ID,
		request.OverID,
		&reorder.Config{
			EntityNotFoundMsg:     "band member not found",
			OverEntityNotFoundMsg: "over band member not found",
		},
	)
	if errCode != nil {
		return errCode
	}

	err = m.repository.UpdateAllSectionTypes(&types)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
