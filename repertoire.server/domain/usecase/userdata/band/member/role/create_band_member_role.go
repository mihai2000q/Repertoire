package role

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type CreateBandMemberRole struct {
	repository repository.UserDataRepository
	jwtService service.JwtService
}

func NewCreateBandMemberRole(
	repository repository.UserDataRepository,
	jwtService service.JwtService,
) CreateBandMemberRole {
	return CreateBandMemberRole{
		repository: repository,
		jwtService: jwtService,
	}
}

func (c CreateBandMemberRole) Handle(
	request requests.CreateBandMemberRoleRequest,
	token string,
) *wrapper.ErrorCode {
	userID, errCode := c.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return errCode
	}

	var count int64
	err := c.repository.CountBandMemberRoles(&count, userID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	bandMemberRole := model.BandMemberRole{
		ID:     uuid.New(),
		Name:   request.Name,
		Order:  uint(count),
		UserID: userID,
	}

	err = c.repository.CreateBandMemberRole(&bandMemberRole)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
