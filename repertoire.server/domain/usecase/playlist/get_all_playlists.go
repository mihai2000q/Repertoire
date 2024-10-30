package playlist

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/model"
	"repertoire/server/utils/wrapper"
)

type GetAllPlaylists struct {
	repository repository.PlaylistRepository
	jwtService service.JwtService
}

func NewGetAllPlaylists(repository repository.PlaylistRepository, jwtService service.JwtService) GetAllPlaylists {
	return GetAllPlaylists{
		repository: repository,
		jwtService: jwtService,
	}
}

func (g GetAllPlaylists) Handle(request requests.GetPlaylistsRequest, token string) (result wrapper.WithTotalCount[model.Playlist], e *wrapper.ErrorCode) {
	userID, errCode := g.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return result, errCode
	}

	err := g.repository.GetAllByUser(
		&result.Models,
		userID,
		request.CurrentPage,
		request.PageSize,
		request.OrderBy,
	)
	if err != nil {
		return result, wrapper.InternalServerError(err)
	}

	err = g.repository.GetAllByUserCount(&result.TotalCount, userID)
	if err != nil {
		return result, wrapper.InternalServerError(err)
	}

	return result, nil
}
