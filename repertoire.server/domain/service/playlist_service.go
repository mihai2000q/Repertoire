package service

import (
	"mime/multipart"
	"repertoire/server/api/requests"
	"repertoire/server/api/responses"
	"repertoire/server/domain/usecase/playlist"
	"repertoire/server/domain/usecase/playlist/song"
	"repertoire/server/internal/httperror"
	"repertoire/server/internal/pagination"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type PlaylistService interface {
	AddAlbums(request requests.AddAlbumsToPlaylistRequest) (*responses.AddAlbumsToPlaylistResponse, *httperror.ErrorCode)
	AddArtists(request requests.AddArtistsToPlaylistRequest) (*responses.AddArtistsToPlaylistResponse, *httperror.ErrorCode)
	AddPerfectRehearsals(request requests.AddPerfectRehearsalsToPlaylistsRequest) *httperror.ErrorCode
	BulkDelete(request requests.BulkDeletePlaylistsRequest) *httperror.ErrorCode
	Create(request requests.CreatePlaylistRequest, token string) (uuid.UUID, *httperror.ErrorCode)
	Delete(id uuid.UUID) *httperror.ErrorCode
	DeleteImage(id uuid.UUID) *httperror.ErrorCode
	GetAll(request requests.GetPlaylistsRequest, token string) (pagination.WithTotalCount[model.EnhancedPlaylist], *httperror.ErrorCode)
	Get(request requests.GetPlaylistRequest) (model.Playlist, *httperror.ErrorCode)
	GetFiltersMetadata(
		request requests.GetPlaylistFiltersMetadataRequest,
		token string,
	) (model.PlaylistFiltersMetadata, *httperror.ErrorCode)
	SaveImage(file *multipart.FileHeader, id uuid.UUID) *httperror.ErrorCode
	Update(request requests.UpdatePlaylistRequest) *httperror.ErrorCode

	AddSongs(request requests.AddSongsToPlaylistRequest) (*responses.AddSongsToPlaylistResponse, *httperror.ErrorCode)
	AddPerfectSongRehearsals(request requests.AddPerfectPlaylistSongRehearsalsRequest) *httperror.ErrorCode
	GetSongs(request requests.GetPlaylistSongsRequest) (pagination.WithTotalCount[model.Song], *httperror.ErrorCode)
	MoveSong(request requests.MoveSongFromPlaylistRequest) *httperror.ErrorCode
	RemoveSongs(request requests.RemoveSongsFromPlaylistRequest) *httperror.ErrorCode
	ShuffleSongs(request requests.ShufflePlaylistSongsRequest) *httperror.ErrorCode
}

type playlistService struct {
	addAlbumsToPlaylist             playlist.AddAlbumsToPlaylist
	addPerfectRehearsalsToPlaylists playlist.AddPerfectRehearsalsToPlaylists
	addArtistsToPlaylist            playlist.AddArtistsToPlaylist
	bulkDeletePlaylists             playlist.BulkDeletePlaylists
	createPlaylist                  playlist.CreatePlaylist
	deletePlaylist                  playlist.DeletePlaylist
	deleteImageFromPlaylist         playlist.DeleteImageFromPlaylist
	getAllPlaylists                 playlist.GetAllPlaylists
	getPlaylist                     playlist.GetPlaylist
	getPlaylistFiltersMetadata      playlist.GetPlaylistFiltersMetadata
	saveImageToPlaylist             playlist.SaveImageToPlaylist
	updatePlaylist                  playlist.UpdatePlaylist

	addSongsToPlaylist               song.AddSongsToPlaylist
	addPerfectPlaylistSongRehearsals song.AddPerfectPlaylistSongRehearsals
	getPlaylistSongs                 song.GetPlaylistSongs
	moveSongFromPlaylist             song.MoveSongFromPlaylist
	removeSongsFromPlaylist          song.RemoveSongsFromPlaylist
	shufflePlaylistSongs             song.ShufflePlaylistSongs
}

func NewPlaylistService(
	addAlbumsToPlaylist playlist.AddAlbumsToPlaylist,
	addArtistsToPlaylist playlist.AddArtistsToPlaylist,
	addPerfectRehearsalsToPlaylists playlist.AddPerfectRehearsalsToPlaylists,
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
	addPerfectPlaylistSongRehearsals song.AddPerfectPlaylistSongRehearsals,
	getPlaylistSongs song.GetPlaylistSongs,
	moveSongFromPlaylist song.MoveSongFromPlaylist,
	removeSongFromPlaylist song.RemoveSongsFromPlaylist,
	shufflePlaylistSongs song.ShufflePlaylistSongs,
) PlaylistService {
	return &playlistService{
		addAlbumsToPlaylist:             addAlbumsToPlaylist,
		addArtistsToPlaylist:            addArtistsToPlaylist,
		addPerfectRehearsalsToPlaylists: addPerfectRehearsalsToPlaylists,
		bulkDeletePlaylists:             bulkDeletePlaylists,
		createPlaylist:                  createPlaylist,
		deletePlaylist:                  deletePlaylist,
		deleteImageFromPlaylist:         deleteImageFromPlaylist,
		getAllPlaylists:                 getAllPlaylists,
		getPlaylist:                     getPlaylist,
		getPlaylistFiltersMetadata:      getPlaylistFiltersMetadata,
		saveImageToPlaylist:             saveImageToPlaylist,
		updatePlaylist:                  updatePlaylist,

		addSongsToPlaylist:               addSongsToPlaylist,
		addPerfectPlaylistSongRehearsals: addPerfectPlaylistSongRehearsals,
		getPlaylistSongs:                 getPlaylistSongs,
		moveSongFromPlaylist:             moveSongFromPlaylist,
		removeSongsFromPlaylist:          removeSongFromPlaylist,
		shufflePlaylistSongs:             shufflePlaylistSongs,
	}
}

func (p *playlistService) AddAlbums(request requests.AddAlbumsToPlaylistRequest) (*responses.AddAlbumsToPlaylistResponse, *httperror.ErrorCode) {
	return p.addAlbumsToPlaylist.Handle(request)
}

func (p *playlistService) AddArtists(
	request requests.AddArtistsToPlaylistRequest,
) (*responses.AddArtistsToPlaylistResponse, *httperror.ErrorCode) {
	return p.addArtistsToPlaylist.Handle(request)
}

func (p *playlistService) AddPerfectRehearsals(request requests.AddPerfectRehearsalsToPlaylistsRequest) *httperror.ErrorCode {
	return p.addPerfectRehearsalsToPlaylists.Handle(request)
}

func (p *playlistService) BulkDelete(request requests.BulkDeletePlaylistsRequest) *httperror.ErrorCode {
	return p.bulkDeletePlaylists.Handle(request)
}

func (p *playlistService) Create(request requests.CreatePlaylistRequest, token string) (uuid.UUID, *httperror.ErrorCode) {
	return p.createPlaylist.Handle(request, token)
}

func (p *playlistService) Delete(id uuid.UUID) *httperror.ErrorCode {
	return p.deletePlaylist.Handle(id)
}

func (p *playlistService) DeleteImage(id uuid.UUID) *httperror.ErrorCode {
	return p.deleteImageFromPlaylist.Handle(id)
}

func (p *playlistService) GetAll(request requests.GetPlaylistsRequest, token string) (pagination.WithTotalCount[model.EnhancedPlaylist], *httperror.ErrorCode) {
	return p.getAllPlaylists.Handle(request, token)
}

func (p *playlistService) Get(request requests.GetPlaylistRequest) (model.Playlist, *httperror.ErrorCode) {
	return p.getPlaylist.Handle(request)
}

func (p *playlistService) GetFiltersMetadata(
	request requests.GetPlaylistFiltersMetadataRequest,
	token string,
) (model.PlaylistFiltersMetadata, *httperror.ErrorCode) {
	return p.getPlaylistFiltersMetadata.Handle(request, token)
}

func (p *playlistService) SaveImage(file *multipart.FileHeader, id uuid.UUID) *httperror.ErrorCode {
	return p.saveImageToPlaylist.Handle(file, id)
}

func (p *playlistService) Update(request requests.UpdatePlaylistRequest) *httperror.ErrorCode {
	return p.updatePlaylist.Handle(request)
}

// songs

func (p *playlistService) AddSongs(
	request requests.AddSongsToPlaylistRequest,
) (*responses.AddSongsToPlaylistResponse, *httperror.ErrorCode) {
	return p.addSongsToPlaylist.Handle(request)
}

func (p *playlistService) AddPerfectSongRehearsals(
	request requests.AddPerfectPlaylistSongRehearsalsRequest,
) *httperror.ErrorCode {
	return p.addPerfectPlaylistSongRehearsals.Handle(request)
}

func (p *playlistService) GetSongs(request requests.GetPlaylistSongsRequest) (pagination.WithTotalCount[model.Song], *httperror.ErrorCode) {
	return p.getPlaylistSongs.Handle(request)
}

func (p *playlistService) MoveSong(request requests.MoveSongFromPlaylistRequest) *httperror.ErrorCode {
	return p.moveSongFromPlaylist.Handle(request)
}

func (p *playlistService) RemoveSongs(request requests.RemoveSongsFromPlaylistRequest) *httperror.ErrorCode {
	return p.removeSongsFromPlaylist.Handle(request)
}

func (p *playlistService) ShuffleSongs(request requests.ShufflePlaylistSongsRequest) *httperror.ErrorCode {
	return p.shufflePlaylistSongs.Handle(request)
}
