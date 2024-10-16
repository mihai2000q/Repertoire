package playlist

import (
	"repertoire/api/requests"
	"repertoire/data/repository"
	"repertoire/models"
	"repertoire/utils"
)

type GetAllPlaylists struct {
	repository repository.PlaylistRepository
}

func NewGetAllPlaylists(repository repository.PlaylistRepository) GetAllPlaylists {
	return GetAllPlaylists{
		repository: repository,
	}
}

func (g GetAllPlaylists) Handle(request requests.GetPlaylistsRequest) (playlists []models.Playlist, e *utils.ErrorCode) {
	err := g.repository.GetAllByUser(&playlists, request.UserID)
	if err != nil {
		return playlists, utils.InternalServerError(err)
	}
	return playlists, nil
}
