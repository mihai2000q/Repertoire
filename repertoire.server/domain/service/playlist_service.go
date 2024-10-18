package service

import (
	"repertoire/api/requests"
	"repertoire/domain/usecases/playlist"
	"repertoire/models"
	"repertoire/utils/wrapper"

	"github.com/google/uuid"
)

type PlaylistService interface {
	Get(id uuid.UUID) (playlist models.Playlist, e *wrapper.ErrorCode)
	GetAll(request requests.GetPlaylistsRequest) (playlists []models.Playlist, e *wrapper.ErrorCode)
	Create(request requests.CreatePlaylistRequest, token string) *wrapper.ErrorCode
	Update(request requests.UpdatePlaylistRequest) *wrapper.ErrorCode
	Delete(id uuid.UUID) *wrapper.ErrorCode
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

func (p *playlistService) Get(id uuid.UUID) (models.Playlist, *wrapper.ErrorCode) {
	return p.getPlaylist.Handle(id)
}

func (p *playlistService) GetAll(request requests.GetPlaylistsRequest) ([]models.Playlist, *wrapper.ErrorCode) {
	return p.getAllPlaylists.Handle(request)
}

func (p *playlistService) Create(request requests.CreatePlaylistRequest, token string) *wrapper.ErrorCode {
	return p.createPlaylist.Handle(request, token)
}

func (p *playlistService) Update(request requests.UpdatePlaylistRequest) *wrapper.ErrorCode {
	return p.updatePlaylist.Handle(request)
}

func (p *playlistService) Delete(id uuid.UUID) *wrapper.ErrorCode {
	return p.deletePlaylist.Handle(id)
}
