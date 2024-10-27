package song

import (
	"repertoire/api/requests"
	"repertoire/data/repository"
	"repertoire/data/service"
	"repertoire/model"
	"repertoire/utils/wrapper"

	"github.com/google/uuid"
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

func (c CreateSong) Handle(request requests.CreateSongRequest, token string) *wrapper.ErrorCode {
	userID, errCode := c.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return errCode
	}

	song := model.Song{
		ID:             uuid.New(),
		Title:          request.Title,
		Description:    request.Description,
		Bpm:            request.Bpm,
		SongsterrLink:  request.SongsterrLink,
		GuitarTuningID: request.GuitarTuningID,
		UserID:         userID,
	}
	err := c.repository.Create(&song)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	return nil
}
