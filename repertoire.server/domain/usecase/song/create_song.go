package song

import (
	"github.com/google/uuid"
	"repertoire/api/request"
	"repertoire/data/repository"
	"repertoire/data/service"
	"repertoire/model"
	"repertoire/utils/wrapper"
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

func (c CreateSong) Handle(request request.CreateSongRequest, token string) *wrapper.ErrorCode {
	userId, errCode := c.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return errCode
	}

	song := model.Song{
		ID:         uuid.New(),
		Title:      request.Title,
		IsRecorded: request.IsRecorded,
		UserID:     userId,
	}
	err := c.repository.Create(&song)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	return nil
}
