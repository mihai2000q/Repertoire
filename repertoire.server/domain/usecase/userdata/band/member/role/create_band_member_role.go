package role

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type CreateBandMemberRole struct {
	userDataRepository repository.UserDataRepository
	jwtService         service.JwtService
}

func NewCreateBandMemberRole(
	userDataRepository repository.UserDataRepository,
	jwtService service.JwtService,
) CreateBandMemberRole {
	return CreateBandMemberRole{
		userDataRepository: userDataRepository,
		jwtService:         jwtService,
	}
}

func (c CreateBandMemberRole) Handle(
	request requests.CreateBandMemberRoleRequest,
	token string,
) *httperror.ErrorCode {
	userID, errCode := c.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return errCode
	}

	var count int64
	if err := c.userDataRepository.CountBandMemberRoles(&count, userID); err != nil {
		return httperror.DatabaseError(err)
	}

	bandMemberRole := model.BandMemberRole{
		ID:     uuid.New(),
		Name:   request.Name,
		Order:  uint(count),
		UserID: userID,
	}

	if err := c.userDataRepository.CreateBandMemberRole(&bandMemberRole); err != nil {
		return httperror.DatabaseError(err)
	}

	return nil
}
