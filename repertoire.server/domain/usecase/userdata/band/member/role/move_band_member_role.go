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
	userDataRepository repository.UserDataRepository
	jwtService         service.JwtService
}

func NewMoveBandMemberRole(
	userDataRepository repository.UserDataRepository,
	jwtService service.JwtService,
) MoveBandMemberRole {
	return MoveBandMemberRole{
		userDataRepository: userDataRepository,
		jwtService:         jwtService,
	}
}

func (m MoveBandMemberRole) Handle(request requests.MoveBandMemberRoleRequest, token string) *httperror.ErrorCode {
	userID, errCode := m.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return errCode
	}

	var roles []model.BandMemberRole
	if err := m.userDataRepository.GetBandMemberRoles(&roles, userID); err != nil {
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

	if err := m.userDataRepository.UpdateAllBandMemberRoles(&roles); err != nil {
		return httperror.DatabaseError(err)
	}

	return nil
}
