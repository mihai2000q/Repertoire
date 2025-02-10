package role

import (
	"errors"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
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

	var types []model.BandMemberRole
	err := m.repository.GetBandMemberRoles(&types, userID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	index, overIndex, err := m.getIndexes(types, request.ID, request.OverID)
	if err != nil {
		return wrapper.NotFoundError(err)
	}
	types = m.move(types, index, overIndex)

	err = m.repository.UpdateAllBandMemberRoles(&types)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}

func (MoveBandMemberRole) getIndexes(roles []model.BandMemberRole, id uuid.UUID, overID uuid.UUID) (int, int, error) {
	var index *int
	var overIndex *int
	for i := 0; i < len(roles); i++ {
		if roles[i].ID == id {
			index = &i
		} else if roles[i].ID == overID {
			overIndex = &i
		}
	}

	if index == nil {
		return -1, -1, errors.New("role not found")
	}
	if overIndex == nil {
		return -1, -1, errors.New("over role not found")
	}

	return *index, *overIndex, nil
}

func (MoveBandMemberRole) move(roles []model.BandMemberRole, index int, overIndex int) []model.BandMemberRole {
	if index < overIndex {
		for i := index + 1; i <= overIndex; i++ {
			roles[i].Order = uint(i - 1)
		}
	} else {
		for i := overIndex; i <= index; i++ {
			roles[i].Order = uint(i + 1)
		}
	}

	roles[index].Order = uint(overIndex)

	return roles
}
