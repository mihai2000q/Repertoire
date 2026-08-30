package role

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/reorder"
	"repertoire/server/internal/wrapper"
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

func (m MoveBandMemberRole) Handle(request requests.MoveBandMemberRoleRequest, token string) *wrapper.ErrorCode {
	userID, errCode := m.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return errCode
	}

	var roles []model.BandMemberRole
	err := m.repository.GetBandMemberRoles(&roles, userID)
	if err != nil {
		return wrapper.InternalServerError(err)
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

	err = m.repository.UpdateAllBandMemberRoles(&roles)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
