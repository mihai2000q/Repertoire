package service

import (
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/playlist"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type PlaylistService interface {
	AddSong(request requests.AddSongToPlaylistRequest) *wrapper.ErrorCode
	Create(request requests.CreatePlaylistRequest, token string) *wrapper.ErrorCode
	Get(id uuid.UUID) (model.Playlist, *wrapper.ErrorCode)
	GetAll(request requests.GetPlaylistsRequest, token string) (wrapper.WithTotalCount[model.Playlist], *wrapper.ErrorCode)
	Update(request requests.UpdatePlaylistRequest) *wrapper.ErrorCode
	Delete(id uuid.UUID) *wrapper.ErrorCode
}

type playlistService struct {
	addSongToPlaylist playlist.AddSongToPlaylist
	createPlaylist    playlist.CreatePlaylist
	deletePlaylist    playlist.DeletePlaylist
	getPlaylist       playlist.GetPlaylist
	getAllPlaylists   playlist.GetAllPlaylists
	updatePlaylist    playlist.UpdatePlaylist
}

func NewPlaylistService(
	addSongToPlaylist playlist.AddSongToPlaylist,
	createPlaylist playlist.CreatePlaylist,
	deletePlaylist playlist.DeletePlaylist,
	getPlaylist playlist.GetPlaylist,
	getAllPlaylists playlist.GetAllPlaylists,
	updatePlaylist playlist.UpdatePlaylist,
) PlaylistService {
	return &playlistService{
		addSongToPlaylist: addSongToPlaylist,
		createPlaylist:    createPlaylist,
		deletePlaylist:    deletePlaylist,
		getPlaylist:       getPlaylist,
		getAllPlaylists:   getAllPlaylists,
		updatePlaylist:    updatePlaylist,
	}
}

func (p *playlistService) AddSong(request requests.AddSongToPlaylistRequest) *wrapper.ErrorCode {
	return p.addSongToPlaylist.Handle(request)
}

func (p *playlistService) Create(request requests.CreatePlaylistRequest, token string) *wrapper.ErrorCode {
	return p.createPlaylist.Handle(request, token)
}

func (p *playlistService) Delete(id uuid.UUID) *wrapper.ErrorCode {
	return p.deletePlaylist.Handle(id)
}

func (p *playlistService) Get(id uuid.UUID) (model.Playlist, *wrapper.ErrorCode) {
	return p.getPlaylist.Handle(id)
}

func (p *playlistService) GetAll(request requests.GetPlaylistsRequest, token string) (wrapper.WithTotalCount[model.Playlist], *wrapper.ErrorCode) {
	return p.getAllPlaylists.Handle(request, token)
}

func (p *playlistService) Update(request requests.UpdatePlaylistRequest) *wrapper.ErrorCode {
	return p.updatePlaylist.Handle(request)
}
