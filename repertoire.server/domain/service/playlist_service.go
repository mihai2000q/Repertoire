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

func (s *playlistService) Get(id uuid.UUID) (models.Playlist, *utils.ErrorCode) {
	return s.getPlaylist.Handle(id)
}

func (s *playlistService) GetAll(request requests.GetPlaylistsRequest) ([]models.Playlist, *utils.ErrorCode) {
	return s.getAllPlaylists.Handle(request)
}

func (s *playlistService) Create(request requests.CreatePlaylistRequest, token string) *utils.ErrorCode {
	return s.createPlaylist.Handle(request, token)
}

func (s *playlistService) Update(request requests.UpdatePlaylistRequest) *utils.ErrorCode {
	return s.updatePlaylist.Handle(request)
}

func (s *playlistService) Delete(id uuid.UUID) *utils.ErrorCode {
	return s.deletePlaylist.Handle(id)
}
