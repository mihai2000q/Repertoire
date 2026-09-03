package role

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/httperror"
	"repertoire/server/internal/reorder"
	"repertoire/server/model"
)

type MoveBandMemberRole struct {
	repository repository.UserDataRepository
	jwtService service.JwtService
}

func NewMoveBandMemberRole(
	repository repository.UserDataRepository,
	jwtService service.JwtService,
) MoveBandMemberRole {
	return MoveBandMemberRole{
		repository: repository,
		jwtService: jwtService,
	}
}

func (m MoveBandMemberRole) Handle(request requests.MoveBandMemberRoleRequest, token string) *httperror.ErrorCode {
	userID, errCode := m.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return errCode
	}

	var roles []model.BandMemberRole
	if err := m.repository.GetBandMemberRoles(&roles, userID); err != nil {
		return httperror.DatabaseError(err)
	}

	errCode = reorder.MoveEntity(
		roles,
		request.ID,
		request.OverID,
		&reorder.Config{
			EntityNotFoundMsg:     "role not found",
			OverEntityNotFoundMsg: "over role not found",
		},
	)
	if errCode != nil {
		return errCode
	}

	if err := m.repository.UpdateAllBandMemberRoles(&roles); err != nil {
		return httperror.DatabaseError(err)
	}

	return nil
}
