package service

import (
	"mime/multipart"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/playlist"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type PlaylistService interface {
	AddSong(request requests.AddSongToPlaylistRequest) *wrapper.ErrorCode
	Create(request requests.CreatePlaylistRequest, token string) (uuid.UUID, *wrapper.ErrorCode)
	Delete(id uuid.UUID) *wrapper.ErrorCode
	DeleteImage(id uuid.UUID) *wrapper.ErrorCode
	GetAll(request requests.GetPlaylistsRequest, token string) (wrapper.WithTotalCount[model.Playlist], *wrapper.ErrorCode)
	Get(id uuid.UUID) (model.Playlist, *wrapper.ErrorCode)
	RemoveSong(id uuid.UUID, songID uuid.UUID) *wrapper.ErrorCode
	SaveImage(file *multipart.FileHeader, id uuid.UUID) *wrapper.ErrorCode
	Update(request requests.UpdatePlaylistRequest) *wrapper.ErrorCode
}

type playlistService struct {
	addSongToPlaylist       playlist.AddSongToPlaylist
	createPlaylist          playlist.CreatePlaylist
	deletePlaylist          playlist.DeletePlaylist
	deleteImageFromPlaylist playlist.DeleteImageFromPlaylist
	getAllPlaylists         playlist.GetAllPlaylists
	getPlaylist             playlist.GetPlaylist
	removeSongFromPlaylist  playlist.RemoveSongFromPlaylist
	saveImageToPlaylist     playlist.SaveImageToPlaylist
	updatePlaylist          playlist.UpdatePlaylist
}

func NewPlaylistService(
	addSongToPlaylist playlist.AddSongToPlaylist,
	createPlaylist playlist.CreatePlaylist,
	deletePlaylist playlist.DeletePlaylist,
	deleteImageFromPlaylist playlist.DeleteImageFromPlaylist,
	getAllPlaylists playlist.GetAllPlaylists,
	getPlaylist playlist.GetPlaylist,
	removeSongFromPlaylist playlist.RemoveSongFromPlaylist,
	saveImageToPlaylist playlist.SaveImageToPlaylist,
	updatePlaylist playlist.UpdatePlaylist,
) PlaylistService {
	return &playlistService{
		addSongToPlaylist:       addSongToPlaylist,
		createPlaylist:          createPlaylist,
		deletePlaylist:          deletePlaylist,
		deleteImageFromPlaylist: deleteImageFromPlaylist,
		getAllPlaylists:         getAllPlaylists,
		getPlaylist:             getPlaylist,
		removeSongFromPlaylist:  removeSongFromPlaylist,
		saveImageToPlaylist:     saveImageToPlaylist,
		updatePlaylist:          updatePlaylist,
	}
}

func (p *playlistService) AddSong(request requests.AddSongToPlaylistRequest) *wrapper.ErrorCode {
	return p.addSongToPlaylist.Handle(request)
}

func (p *playlistService) Create(request requests.CreatePlaylistRequest, token string) (uuid.UUID, *wrapper.ErrorCode) {
	return p.createPlaylist.Handle(request, token)
}

func (p *playlistService) Delete(id uuid.UUID) *wrapper.ErrorCode {
	return p.deletePlaylist.Handle(id)
}

func (p *playlistService) DeleteImage(id uuid.UUID) *wrapper.ErrorCode {
	return p.deleteImageFromPlaylist.Handle(id)
}

func (p *playlistService) GetAll(request requests.GetPlaylistsRequest, token string) (wrapper.WithTotalCount[model.Playlist], *wrapper.ErrorCode) {
	return p.getAllPlaylists.Handle(request, token)
}

func (p *playlistService) Get(id uuid.UUID) (model.Playlist, *wrapper.ErrorCode) {
	return p.getPlaylist.Handle(id)
}

func (p *playlistService) RemoveSong(id uuid.UUID, songID uuid.UUID) *wrapper.ErrorCode {
	return p.removeSongFromPlaylist.Handle(id, songID)
}

func (p *playlistService) SaveImage(file *multipart.FileHeader, id uuid.UUID) *wrapper.ErrorCode {
	return p.saveImageToPlaylist.Handle(file, id)
}

func (p *playlistService) Update(request requests.UpdatePlaylistRequest) *wrapper.ErrorCode {
	return p.updatePlaylist.Handle(request)
}
