package service

import (
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/playlist"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type PlaylistService interface {
	Get(id uuid.UUID) (model.Playlist, *wrapper.ErrorCode)
	GetAll(request requests.GetPlaylistsRequest, token string) (wrapper.WithTotalCount[model.Playlist], *wrapper.ErrorCode)
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

func (p *playlistService) Get(id uuid.UUID) (model.Playlist, *wrapper.ErrorCode) {
	return p.getPlaylist.Handle(id)
}

func (p *playlistService) GetAll(request requests.GetPlaylistsRequest, token string) (wrapper.WithTotalCount[model.Playlist], *wrapper.ErrorCode) {
	return p.getAllPlaylists.Handle(request, token)
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
