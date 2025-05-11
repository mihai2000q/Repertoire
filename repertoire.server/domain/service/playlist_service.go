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
	AddAlbums(request requests.AddAlbumsToPlaylistRequest) *wrapper.ErrorCode
	AddArtists(request requests.AddArtistsToPlaylistRequest) *wrapper.ErrorCode
	AddSongs(request requests.AddSongsToPlaylistRequest) *wrapper.ErrorCode
	Create(request requests.CreatePlaylistRequest, token string) (uuid.UUID, *wrapper.ErrorCode)
	Delete(id uuid.UUID) *wrapper.ErrorCode
	DeleteImage(id uuid.UUID) *wrapper.ErrorCode
	GetAll(request requests.GetPlaylistsRequest, token string) (wrapper.WithTotalCount[model.EnhancedPlaylist], *wrapper.ErrorCode)
	Get(request requests.GetPlaylistRequest) (model.Playlist, *wrapper.ErrorCode)
	GetFiltersMetadata(
		request requests.GetPlaylistFiltersMetadataRequest,
		token string,
	) (model.PlaylistFiltersMetadata, *wrapper.ErrorCode)
	MoveSong(request requests.MoveSongFromPlaylistRequest) *wrapper.ErrorCode
	RemoveSongs(request requests.RemoveSongsFromPlaylistRequest) *wrapper.ErrorCode
	SaveImage(file *multipart.FileHeader, id uuid.UUID) *wrapper.ErrorCode
	Update(request requests.UpdatePlaylistRequest) *wrapper.ErrorCode
}

type playlistService struct {
	addAlbumsToPlaylist        playlist.AddAlbumsToPlaylist
	addArtistsToPlaylist       playlist.AddArtistsToPlaylist
	addSongsToPlaylist         playlist.AddSongsToPlaylist
	createPlaylist             playlist.CreatePlaylist
	deletePlaylist             playlist.DeletePlaylist
	deleteImageFromPlaylist    playlist.DeleteImageFromPlaylist
	getAllPlaylists            playlist.GetAllPlaylists
	getPlaylist                playlist.GetPlaylist
	getPlaylistFiltersMetadata playlist.GetPlaylistFiltersMetadata
	moveSongFromPlaylist       playlist.MoveSongFromPlaylist
	removeSongsFromPlaylist    playlist.RemoveSongsFromPlaylist
	saveImageToPlaylist        playlist.SaveImageToPlaylist
	updatePlaylist             playlist.UpdatePlaylist
}

func NewPlaylistService(
	addAlbumsToPlaylist playlist.AddAlbumsToPlaylist,
	addArtistsToPlaylist playlist.AddArtistsToPlaylist,
	addSongsToPlaylist playlist.AddSongsToPlaylist,
	createPlaylist playlist.CreatePlaylist,
	deletePlaylist playlist.DeletePlaylist,
	deleteImageFromPlaylist playlist.DeleteImageFromPlaylist,
	getAllPlaylists playlist.GetAllPlaylists,
	getPlaylist playlist.GetPlaylist,
	getPlaylistFiltersMetadata playlist.GetPlaylistFiltersMetadata,
	moveSongFromPlaylist playlist.MoveSongFromPlaylist,
	removeSongFromPlaylist playlist.RemoveSongsFromPlaylist,
	saveImageToPlaylist playlist.SaveImageToPlaylist,
	updatePlaylist playlist.UpdatePlaylist,
) PlaylistService {
	return &playlistService{
		addAlbumsToPlaylist:        addAlbumsToPlaylist,
		addArtistsToPlaylist:       addArtistsToPlaylist,
		addSongsToPlaylist:         addSongsToPlaylist,
		createPlaylist:             createPlaylist,
		deletePlaylist:             deletePlaylist,
		deleteImageFromPlaylist:    deleteImageFromPlaylist,
		getAllPlaylists:            getAllPlaylists,
		getPlaylist:                getPlaylist,
		getPlaylistFiltersMetadata: getPlaylistFiltersMetadata,
		moveSongFromPlaylist:       moveSongFromPlaylist,
		removeSongsFromPlaylist:    removeSongFromPlaylist,
		saveImageToPlaylist:        saveImageToPlaylist,
		updatePlaylist:             updatePlaylist,
	}
}

func (p *playlistService) AddAlbums(request requests.AddAlbumsToPlaylistRequest) *wrapper.ErrorCode {
	return p.addAlbumsToPlaylist.Handle(request)
}

func (p *playlistService) AddArtists(request requests.AddArtistsToPlaylistRequest) *wrapper.ErrorCode {
	return p.addArtistsToPlaylist.Handle(request)
}

func (p *playlistService) AddSongs(request requests.AddSongsToPlaylistRequest) *wrapper.ErrorCode {
	return p.addSongsToPlaylist.Handle(request)
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

func (p *playlistService) GetAll(request requests.GetPlaylistsRequest, token string) (wrapper.WithTotalCount[model.EnhancedPlaylist], *wrapper.ErrorCode) {
	return p.getAllPlaylists.Handle(request, token)
}

func (p *playlistService) Get(request requests.GetPlaylistRequest) (model.Playlist, *wrapper.ErrorCode) {
	return p.getPlaylist.Handle(request)
}

func (p *playlistService) GetFiltersMetadata(
	request requests.GetPlaylistFiltersMetadataRequest,
	token string,
) (model.PlaylistFiltersMetadata, *wrapper.ErrorCode) {
	return p.getPlaylistFiltersMetadata.Handle(request, token)
}

func (p *playlistService) MoveSong(request requests.MoveSongFromPlaylistRequest) *wrapper.ErrorCode {
	return p.moveSongFromPlaylist.Handle(request)
}

func (p *playlistService) RemoveSongs(request requests.RemoveSongsFromPlaylistRequest) *wrapper.ErrorCode {
	return p.removeSongsFromPlaylist.Handle(request)
}

func (p *playlistService) SaveImage(file *multipart.FileHeader, id uuid.UUID) *wrapper.ErrorCode {
	return p.saveImageToPlaylist.Handle(file, id)
}

func (p *playlistService) Update(request requests.UpdatePlaylistRequest) *wrapper.ErrorCode {
	return p.updatePlaylist.Handle(request)
}
