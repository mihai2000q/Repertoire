package types

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/httperror"
	"repertoire/server/internal/reorder"
	"repertoire/server/model"
)

type MoveSongSectionType struct {
	userDataRepository repository.UserDataRepository
	jwtService         service.JwtService
}

func NewMoveSongSectionType(
	userDataRepository repository.UserDataRepository,
	jwtService service.JwtService,
) MoveSongSectionType {
	return MoveSongSectionType{
		userDataRepository: userDataRepository,
		jwtService:         jwtService,
	}
}

func (m MoveSongSectionType) Handle(request requests.MoveSongSectionTypeRequest, token string) *httperror.ErrorCode {
	userID, errCode := m.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return errCode
	}

	var types []model.SongSectionType
	if err := m.userDataRepository.GetSectionTypes(&types, userID); err != nil {
		return httperror.DatabaseError(err)
	}

	errCode = reorder.MoveEntity(
		types,
		request.ID,
		request.OverID,
		&reorder.Config{
			EntityNotFoundMsg:     "type not found",
			OverEntityNotFoundMsg: "over type not found",
		},
	)
	if errCode != nil {
		return errCode
	}

	if err := m.userDataRepository.UpdateAllSectionTypes(&types); err != nil {
		return httperror.DatabaseError(err)
	}

	return nil
}
