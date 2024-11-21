package types

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type CreateSongSectionType struct {
	repository repository.SongRepository
	jwtService service.JwtService
}

func NewCreateSongSectionType(
	repository repository.SongRepository,
	jwtService service.JwtService,
) CreateSongSectionType {
	return CreateSongSectionType{
		repository: repository,
		jwtService: jwtService,
	}
}

func (c CreateSongSectionType) Handle(
	request requests.CreateSongSectionTypeRequest,
	token string,
) *wrapper.ErrorCode {
	userID, errCode := c.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return errCode
	}

	var count int64
	err := c.repository.CountSectionTypes(&count, userID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	sectionType := model.SongSectionType{
		ID:     uuid.New(),
		Name:   request.Name,
		Order:  uint(count),
		UserID: userID,
	}

	err = c.repository.CreateSectionType(&sectionType)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
