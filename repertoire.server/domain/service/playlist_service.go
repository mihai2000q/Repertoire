package service

import (
	"mime/multipart"
	"repertoire/server/api/requests"
	"repertoire/server/api/responses"
	"repertoire/server/domain/usecase/playlist"
	"repertoire/server/domain/usecase/playlist/song"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type PlaylistService interface {
	AddAlbums(request requests.AddAlbumsToPlaylistRequest) (*responses.AddAlbumsToPlaylistResponse, *wrapper.ErrorCode)
	AddArtists(request requests.AddArtistsToPlaylistRequest) (*responses.AddArtistsToPlaylistResponse, *wrapper.ErrorCode)
	BulkDelete(request requests.BulkDeletePlaylistsRequest) *wrapper.ErrorCode
	Create(request requests.CreatePlaylistRequest, token string) (uuid.UUID, *wrapper.ErrorCode)
	Delete(id uuid.UUID) *wrapper.ErrorCode
	DeleteImage(id uuid.UUID) *wrapper.ErrorCode
	GetAll(request requests.GetPlaylistsRequest, token string) (wrapper.WithTotalCount[model.EnhancedPlaylist], *wrapper.ErrorCode)
	Get(request requests.GetPlaylistRequest) (model.Playlist, *wrapper.ErrorCode)
	GetFiltersMetadata(
		request requests.GetPlaylistFiltersMetadataRequest,
		token string,
	) (model.PlaylistFiltersMetadata, *wrapper.ErrorCode)
	SaveImage(file *multipart.FileHeader, id uuid.UUID) *wrapper.ErrorCode
	Update(request requests.UpdatePlaylistRequest) *wrapper.ErrorCode

	AddSongs(request requests.AddSongsToPlaylistRequest) (*responses.AddSongsToPlaylistResponse, *wrapper.ErrorCode)
	GetSongs(request requests.GetPlaylistSongsRequest) (wrapper.WithTotalCount[model.Song], *wrapper.ErrorCode)
	MoveSong(request requests.MoveSongFromPlaylistRequest) *wrapper.ErrorCode
	RemoveSongs(request requests.RemoveSongsFromPlaylistRequest) *wrapper.ErrorCode
	ShuffleSongs(request requests.ShufflePlaylistSongsRequest) *wrapper.ErrorCode
}

type playlistService struct {
	addAlbumsToPlaylist        playlist.AddAlbumsToPlaylist
	addArtistsToPlaylist       playlist.AddArtistsToPlaylist
	bulkDeletePlaylists        playlist.BulkDeletePlaylists
	createPlaylist             playlist.CreatePlaylist
	deletePlaylist             playlist.DeletePlaylist
	deleteImageFromPlaylist    playlist.DeleteImageFromPlaylist
	getAllPlaylists            playlist.GetAllPlaylists
	getPlaylist                playlist.GetPlaylist
	getPlaylistFiltersMetadata playlist.GetPlaylistFiltersMetadata
	saveImageToPlaylist        playlist.SaveImageToPlaylist
	updatePlaylist             playlist.UpdatePlaylist

	addSongsToPlaylist      song.AddSongsToPlaylist
	getPlaylistSongs        song.GetPlaylistSongs
	moveSongFromPlaylist    song.MoveSongFromPlaylist
	removeSongsFromPlaylist song.RemoveSongsFromPlaylist
	shufflePlaylistSongs    song.ShufflePlaylistSongs
}

func NewPlaylistService(
	addAlbumsToPlaylist playlist.AddAlbumsToPlaylist,
	addArtistsToPlaylist playlist.AddArtistsToPlaylist,
	bulkDeletePlaylists playlist.BulkDeletePlaylists,
	createPlaylist playlist.CreatePlaylist,
	deletePlaylist playlist.DeletePlaylist,
	deleteImageFromPlaylist playlist.DeleteImageFromPlaylist,
	getAllPlaylists playlist.GetAllPlaylists,
	getPlaylist playlist.GetPlaylist,
	getPlaylistFiltersMetadata playlist.GetPlaylistFiltersMetadata,
	saveImageToPlaylist playlist.SaveImageToPlaylist,
	updatePlaylist playlist.UpdatePlaylist,

	addSongsToPlaylist song.AddSongsToPlaylist,
	getPlaylistSongs song.GetPlaylistSongs,
	moveSongFromPlaylist song.MoveSongFromPlaylist,
	removeSongFromPlaylist song.RemoveSongsFromPlaylist,
	shufflePlaylistSongs song.ShufflePlaylistSongs,
) PlaylistService {
	return &playlistService{
		addAlbumsToPlaylist:        addAlbumsToPlaylist,
		addArtistsToPlaylist:       addArtistsToPlaylist,
		bulkDeletePlaylists:        bulkDeletePlaylists,
		createPlaylist:             createPlaylist,
		deletePlaylist:             deletePlaylist,
		deleteImageFromPlaylist:    deleteImageFromPlaylist,
		getAllPlaylists:            getAllPlaylists,
		getPlaylist:                getPlaylist,
		getPlaylistFiltersMetadata: getPlaylistFiltersMetadata,
		saveImageToPlaylist:        saveImageToPlaylist,
		updatePlaylist:             updatePlaylist,

		addSongsToPlaylist:      addSongsToPlaylist,
		getPlaylistSongs:        getPlaylistSongs,
		moveSongFromPlaylist:    moveSongFromPlaylist,
		removeSongsFromPlaylist: removeSongFromPlaylist,
		shufflePlaylistSongs:    shufflePlaylistSongs,
	}
}

func (p *playlistService) AddAlbums(request requests.AddAlbumsToPlaylistRequest) (*responses.AddAlbumsToPlaylistResponse, *wrapper.ErrorCode) {
	return p.addAlbumsToPlaylist.Handle(request)
}

func (p *playlistService) AddArtists(
	request requests.AddArtistsToPlaylistRequest,
) (*responses.AddArtistsToPlaylistResponse, *wrapper.ErrorCode) {
	return p.addArtistsToPlaylist.Handle(request)
}

func (p *playlistService) BulkDelete(request requests.BulkDeletePlaylistsRequest) *wrapper.ErrorCode {
	return p.bulkDeletePlaylists.Handle(request)
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

func (p *playlistService) SaveImage(file *multipart.FileHeader, id uuid.UUID) *wrapper.ErrorCode {
	return p.saveImageToPlaylist.Handle(file, id)
}

func (p *playlistService) Update(request requests.UpdatePlaylistRequest) *wrapper.ErrorCode {
	return p.updatePlaylist.Handle(request)
}

// songs

func (p *playlistService) AddSongs(
	request requests.AddSongsToPlaylistRequest,
) (*responses.AddSongsToPlaylistResponse, *wrapper.ErrorCode) {
	return p.addSongsToPlaylist.Handle(request)
}

func (p *playlistService) GetSongs(request requests.GetPlaylistSongsRequest) (wrapper.WithTotalCount[model.Song], *wrapper.ErrorCode) {
	return p.getPlaylistSongs.Handle(request)
}

func (p *playlistService) MoveSong(request requests.MoveSongFromPlaylistRequest) *wrapper.ErrorCode {
	return p.moveSongFromPlaylist.Handle(request)
}

func (p *playlistService) RemoveSongs(request requests.RemoveSongsFromPlaylistRequest) *wrapper.ErrorCode {
	return p.removeSongsFromPlaylist.Handle(request)
}

func (p *playlistService) ShuffleSongs(request requests.ShufflePlaylistSongsRequest) *wrapper.ErrorCode {
	return p.shufflePlaylistSongs.Handle(request)
}
