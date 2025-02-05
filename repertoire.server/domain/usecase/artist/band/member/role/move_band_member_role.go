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
	repository repository.ArtistRepository
	jwtService service.JwtService
}

func NewMoveBandMemberRole(
	repository repository.ArtistRepository,
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

func (MoveBandMemberRole) getIndexes(types []model.BandMemberRole, id uuid.UUID, overID uuid.UUID) (int, int, error) {
	var index *int
	var overIndex *int
	for i := 0; i < len(types); i++ {
		if types[i].ID == id {
			index = &i
		} else if types[i].ID == overID {
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

func (MoveBandMemberRole) move(types []model.BandMemberRole, index int, overIndex int) []model.BandMemberRole {
	if index < overIndex {
		for i := index + 1; i <= overIndex; i++ {
			types[i].Order = uint(i - 1)
		}
	} else {
		for i := overIndex; i <= index; i++ {
			types[i].Order = uint(i + 1)
		}
	}

	types[index].Order = uint(overIndex)

	return types
}
