package types

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type CreateSongSectionType struct {
	userDataRepository repository.UserDataRepository
	jwtService         service.JwtService
}

func NewCreateSongSectionType(
	userDataRepository repository.UserDataRepository,
	jwtService service.JwtService,
) CreateSongSectionType {
	return CreateSongSectionType{
		userDataRepository: userDataRepository,
		jwtService:         jwtService,
	}
}

func (c CreateSongSectionType) Handle(
	request requests.CreateSongSectionTypeRequest,
	token string,
) *httperror.ErrorCode {
	userID, errCode := c.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return errCode
	}

	var count int64
	if err := c.userDataRepository.CountSectionTypes(&count, userID); err != nil {
		return httperror.DatabaseError(err)
	}

	sectionType := model.SongSectionType{
		ID:     uuid.New(),
		Name:   request.Name,
		Order:  uint(count),
		UserID: userID,
	}

	if err := c.userDataRepository.CreateSectionType(&sectionType); err != nil {
		return httperror.DatabaseError(err)
	}

	return nil
}
