package service

import (
	"mime/multipart"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/artist"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type ArtistService interface {
	AddAlbums(request requests.AddAlbumsToArtistRequest) *wrapper.ErrorCode
	AddSongs(request requests.AddSongsToArtistRequest) *wrapper.ErrorCode
	Create(request requests.CreateArtistRequest, token string) (uuid.UUID, *wrapper.ErrorCode)
	Delete(id uuid.UUID) *wrapper.ErrorCode
	DeleteImage(id uuid.UUID) *wrapper.ErrorCode
	GetAll(request requests.GetArtistsRequest, token string) (wrapper.WithTotalCount[model.Artist], *wrapper.ErrorCode)
	Get(id uuid.UUID) (model.Artist, *wrapper.ErrorCode)
	RemoveAlbum(id uuid.UUID, albumID uuid.UUID) *wrapper.ErrorCode
	RemoveSong(id uuid.UUID, songID uuid.UUID) *wrapper.ErrorCode
	SaveImage(file *multipart.FileHeader, id uuid.UUID) *wrapper.ErrorCode
	Update(request requests.UpdateArtistRequest) *wrapper.ErrorCode
}

type artistService struct {
	addAlbumsToArtist     artist.AddAlbumsToArtist
	addSongsToArtist      artist.AddSongsToArtist
	createArtist          artist.CreateArtist
	deleteArtist          artist.DeleteArtist
	deleteImageFromArtist artist.DeleteImageFromArtist
	getAllArtists         artist.GetAllArtists
	getArtist             artist.GetArtist
	removeAlbumFromArtist artist.RemoveAlbumFromArtist
	removeSongFromArtist  artist.RemoveSongFromArtist
	saveImageToArtist     artist.SaveImageToArtist
	updateArtist          artist.UpdateArtist
}

func NewArtistService(
	addAlbumsToArtist artist.AddAlbumsToArtist,
	addSongsToArtist artist.AddSongsToArtist,
	createArtist artist.CreateArtist,
	deleteArtist artist.DeleteArtist,
	deleteImageFromArtist artist.DeleteImageFromArtist,
	getAllArtists artist.GetAllArtists,
	getArtist artist.GetArtist,
	removeAlbumFromArtist artist.RemoveAlbumFromArtist,
	removeSongFromArtist artist.RemoveSongFromArtist,
	saveImageToArtist artist.SaveImageToArtist,
	updateArtist artist.UpdateArtist,
) ArtistService {
	return &artistService{
		addAlbumsToArtist:     addAlbumsToArtist,
		addSongsToArtist:      addSongsToArtist,
		createArtist:          createArtist,
		deleteArtist:          deleteArtist,
		deleteImageFromArtist: deleteImageFromArtist,
		getAllArtists:         getAllArtists,
		getArtist:             getArtist,
		removeAlbumFromArtist: removeAlbumFromArtist,
		removeSongFromArtist:  removeSongFromArtist,
		saveImageToArtist:     saveImageToArtist,
		updateArtist:          updateArtist,
	}
}

func (a *artistService) AddAlbums(request requests.AddAlbumsToArtistRequest) *wrapper.ErrorCode {
	return a.addAlbumsToArtist.Handle(request)
}

func (a *artistService) AddSongs(request requests.AddSongsToArtistRequest) *wrapper.ErrorCode {
	return a.addSongsToArtist.Handle(request)
}

func (a *artistService) Create(request requests.CreateArtistRequest, token string) (uuid.UUID, *wrapper.ErrorCode) {
	return a.createArtist.Handle(request, token)
}

func (a *artistService) Delete(id uuid.UUID) *wrapper.ErrorCode {
	return a.deleteArtist.Handle(id)
}

func (a *artistService) DeleteImage(id uuid.UUID) *wrapper.ErrorCode {
	return a.deleteImageFromArtist.Handle(id)
}

func (a *artistService) GetAll(request requests.GetArtistsRequest, token string) (wrapper.WithTotalCount[model.Artist], *wrapper.ErrorCode) {
	return a.getAllArtists.Handle(request, token)
}

func (a *artistService) Get(id uuid.UUID) (model.Artist, *wrapper.ErrorCode) {
	return a.getArtist.Handle(id)
}

func (a *artistService) RemoveAlbum(id uuid.UUID, albumID uuid.UUID) *wrapper.ErrorCode {
	return a.removeAlbumFromArtist.Handle(id, albumID)
}

func (a *artistService) RemoveSong(id uuid.UUID, songID uuid.UUID) *wrapper.ErrorCode {
	return a.removeSongFromArtist.Handle(id, songID)
}

func (a *artistService) SaveImage(file *multipart.FileHeader, id uuid.UUID) *wrapper.ErrorCode {
	return a.saveImageToArtist.Handle(file, id)
}

func (a *artistService) Update(request requests.UpdateArtistRequest) *wrapper.ErrorCode {
	return a.updateArtist.Handle(request)
}
