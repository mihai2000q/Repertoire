package playlist

import (
	"repertoire/api/requests"
	"repertoire/data/repository"
	"repertoire/data/service"
	"repertoire/models"
	"repertoire/utils/wrapper"

	"github.com/google/uuid"
)

type CreatePlaylist struct {
	jwtService service.JwtService
	repository repository.PlaylistRepository
}

func NewCreatePlaylist(jwtService service.JwtService, repository repository.PlaylistRepository) CreatePlaylist {
	return CreatePlaylist{
		jwtService: jwtService,
		repository: repository,
	}
}

func (c CreatePlaylist) Handle(request requests.CreatePlaylistRequest, token string) *wrapper.ErrorCode {
	userId, errCode := c.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return errCode
	}

	playlist := models.Playlist{
		ID:          uuid.New(),
		Title:       request.Title,
		Description: request.Description,
		UserID:      userId,
	}
	err := c.repository.Create(&playlist)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	return nil
}
