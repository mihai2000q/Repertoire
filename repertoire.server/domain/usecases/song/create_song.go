package song

import (
	"github.com/google/uuid"
	"repertoire/api/requests"
	"repertoire/data/repository"
	"repertoire/data/service"
	"repertoire/models"
	"repertoire/utils"
)

type CreateSong struct {
	jwtService service.JwtService
	repository repository.SongRepository
}

func NewCreateSong(jwtService service.JwtService, repository repository.SongRepository) CreateSong {
	return CreateSong{
		jwtService: jwtService,
		repository: repository,
	}
}

func (c CreateSong) Handle(request requests.CreateSongRequest, token string) *utils.ErrorCode {
	userId, errCode := c.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return errCode
	}

	song := models.Song{
		ID:         uuid.New(),
		Title:      request.Title,
		IsRecorded: request.IsRecorded,
		UserID:     userId,
	}
	err := c.repository.Create(&song)
	if err != nil {
		return utils.InternalServerError(err)
	}
	return nil
}
