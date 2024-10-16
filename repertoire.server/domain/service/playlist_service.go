package service

import (
	"repertoire/api/requests"
	"repertoire/domain/usecases/playlist"
	"repertoire/models"
	"repertoire/utils"

	"github.com/google/uuid"
)

type PlaylistService interface {
	Get(id uuid.UUID) (playlist models.Playlist, e *utils.ErrorCode)
	GetAll(request requests.GetPlaylistsRequest) (playlists []models.Playlist, e *utils.ErrorCode)
	Create(request requests.CreatePlaylistRequest, token string) *utils.ErrorCode
	Update(request requests.UpdatePlaylistRequest) *utils.ErrorCode
	Delete(id uuid.UUID) *utils.ErrorCode
}

type playlistService struct {
	getPlaylist     playlist.GetPlaylist
	getAllPlaylists playlist.GetAllPlaylists
	createPlaylist  playlist.CreatePlaylist
	updatePlaylist  playlist.UpdatePlaylist
	deletePlaylist  playlist.DeletePlaylist
}

func NewPlaylistService(
	getPlaylist playlist.GetPlaylist,
	getAllPlaylists playlist.GetAllPlaylists,
	createPlaylist playlist.CreatePlaylist,
	updatePlaylist playlist.UpdatePlaylist,
	deletePlaylist playlist.DeletePlaylist,
) PlaylistService {
	return &playlistService{
		getPlaylist:     getPlaylist,
		getAllPlaylists: getAllPlaylists,
		createPlaylist:  createPlaylist,
		updatePlaylist:  updatePlaylist,
		deletePlaylist:  deletePlaylist,
	}
}

func (p *playlistService) Get(id uuid.UUID) (models.Playlist, *utils.ErrorCode) {
	return p.getPlaylist.Handle(id)
}

func (p *playlistService) GetAll(request requests.GetPlaylistsRequest) ([]models.Playlist, *utils.ErrorCode) {
	return p.getAllPlaylists.Handle(request)
}

func (p *playlistService) Create(request requests.CreatePlaylistRequest, token string) *utils.ErrorCode {
	return p.createPlaylist.Handle(request, token)
}

func (p *playlistService) Update(request requests.UpdatePlaylistRequest) *utils.ErrorCode {
	return p.updatePlaylist.Handle(request)
}

func (p *playlistService) Delete(id uuid.UUID) *utils.ErrorCode {
	return p.deletePlaylist.Handle(id)
}
