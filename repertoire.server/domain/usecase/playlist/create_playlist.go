package playlist

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type CreatePlaylist struct {
	jwtService          service.JwtService
	repository          repository.PlaylistRepository
	searchEngineService service.SearchEngineService
}

func NewCreatePlaylist(
	jwtService service.JwtService,
	repository repository.PlaylistRepository,
	searchEngineService service.SearchEngineService,
) CreatePlaylist {
	return CreatePlaylist{
		jwtService:          jwtService,
		repository:          repository,
		searchEngineService: searchEngineService,
	}
}

func (c CreatePlaylist) Handle(request requests.CreatePlaylistRequest, token string) (uuid.UUID, *wrapper.ErrorCode) {
	userID, errCode := c.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return uuid.Nil, errCode
	}

	playlist := model.Playlist{
		ID:          uuid.New(),
		Title:       request.Title,
		Description: request.Description,
		UserID:      userID,
	}
	err := c.repository.Create(&playlist)
	if err != nil {
		return uuid.Nil, wrapper.InternalServerError(err)
	}

	errCode = c.searchEngineService.Add([]any{playlist.ToSearch()})
	if errCode != nil {
		return uuid.Nil, errCode
	}

	return playlist.ID, nil
}
