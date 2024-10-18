package playlist

import (
	"repertoire/api/requests"
	"repertoire/data/repository"
	"repertoire/models"
	"repertoire/utils/wrapper"
)

type GetAllPlaylists struct {
	repository repository.PlaylistRepository
}

func NewGetAllPlaylists(repository repository.PlaylistRepository) GetAllPlaylists {
	return GetAllPlaylists{
		repository: repository,
	}
}

func (g GetAllPlaylists) Handle(request requests.GetPlaylistsRequest) (result wrapper.WithTotalCount[models.Playlist], e *wrapper.ErrorCode) {
	err := g.repository.GetAllByUser(&result.Data, request.UserID, request.CurrentPage, request.PageSize)
	if err != nil {
		return result, wrapper.InternalServerError(err)
	}
	err = g.repository.GetAllByUserCount(&result.TotalCount, request.UserID)
	if err != nil {
		return result, wrapper.InternalServerError(err)
	}
	return result, nil
}
