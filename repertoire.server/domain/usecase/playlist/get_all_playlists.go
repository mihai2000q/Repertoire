package playlist

import (
	"repertoire/api/request"
	"repertoire/data/repository"
	"repertoire/data/service"
	"repertoire/model"
	"repertoire/utils/wrapper"
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

func (g GetAllPlaylists) Handle(request request.GetPlaylistsRequest, token string) (result wrapper.WithTotalCount[model.Playlist], e *wrapper.ErrorCode) {
	userId, errCode := g.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return result, errCode
	}

	err := g.repository.GetAllByUser(&result.Models, userId, request.CurrentPage, request.PageSize)
	if err != nil {
		return result, wrapper.InternalServerError(err)
	}

	err = g.repository.GetAllByUserCount(&result.TotalCount, userId)
	if err != nil {
		return result, wrapper.InternalServerError(err)
	}

	return result, nil
}
