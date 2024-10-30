package playlist

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/model"
	"repertoire/server/utils/wrapper"

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
	userID, errCode := c.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return errCode
	}

	playlist := model.Playlist{
		ID:          uuid.New(),
		Title:       request.Title,
		Description: request.Description,
		UserID:      userID,
	}
	err := c.repository.Create(&playlist)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	return nil
}
