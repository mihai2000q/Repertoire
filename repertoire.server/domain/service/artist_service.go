package service

import (
	"mime/multipart"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/artist"
	"repertoire/server/domain/usecase/artist/band/member/role"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type ArtistService interface {
	AddAlbums(request requests.AddAlbumsToArtistRequest) *wrapper.ErrorCode
	AddSongs(request requests.AddSongsToArtistRequest) *wrapper.ErrorCode
	Create(request requests.CreateArtistRequest, token string) (uuid.UUID, *wrapper.ErrorCode)
	Delete(request requests.DeleteArtistRequest) *wrapper.ErrorCode
	DeleteImage(id uuid.UUID) *wrapper.ErrorCode
	GetAll(request requests.GetArtistsRequest, token string) (wrapper.WithTotalCount[model.Artist], *wrapper.ErrorCode)
	Get(id uuid.UUID) (model.Artist, *wrapper.ErrorCode)
	RemoveAlbums(request requests.RemoveAlbumsFromArtistRequest) *wrapper.ErrorCode
	RemoveSongs(request requests.RemoveSongsFromArtistRequest) *wrapper.ErrorCode
	SaveImage(file *multipart.FileHeader, id uuid.UUID) *wrapper.ErrorCode
	Update(request requests.UpdateArtistRequest) *wrapper.ErrorCode
	CreateBandMemberRole(request requests.CreateBandMemberRoleRequest, token string) *wrapper.ErrorCode
	DeleteBandMemberRole(id uuid.UUID, token string) *wrapper.ErrorCode
	GetBandMemberRoles(token string) ([]model.BandMemberRole, *wrapper.ErrorCode)
	MoveBandMemberRole(request requests.MoveBandMemberRoleRequest, token string) *wrapper.ErrorCode
}

type artistService struct {
	addAlbumsToArtist      artist.AddAlbumsToArtist
	addSongsToArtist       artist.AddSongsToArtist
	createArtist           artist.CreateArtist
	deleteArtist           artist.DeleteArtist
	deleteImageFromArtist  artist.DeleteImageFromArtist
	getAllArtists          artist.GetAllArtists
	getArtist              artist.GetArtist
	removeAlbumsFromArtist artist.RemoveAlbumsFromArtist
	removeSongsFromArtist  artist.RemoveSongsFromArtist
	saveImageToArtist      artist.SaveImageToArtist
	updateArtist           artist.UpdateArtist
	createBandMemberRole role.CreateBandMemberRole
	deleteBandMemberRole role.DeleteBandMemberRole
	getBandMemberRoles   role.GetBandMemberRoles
	moveBandMemberRole   role.MoveBandMemberRole
}

func NewArtistService(
	addAlbumsToArtist artist.AddAlbumsToArtist,
	addSongsToArtist artist.AddSongsToArtist,
	createArtist artist.CreateArtist,
	deleteArtist artist.DeleteArtist,
	deleteImageFromArtist artist.DeleteImageFromArtist,
	getAllArtists artist.GetAllArtists,
	getArtist artist.GetArtist,
	removeAlbumsFromArtist artist.RemoveAlbumsFromArtist,
	removeSongsFromArtist artist.RemoveSongsFromArtist,
	saveImageToArtist artist.SaveImageToArtist,
	updateArtist artist.UpdateArtist,
	createBandMemberRole role.CreateBandMemberRole,
	deleteBandMemberRole role.DeleteBandMemberRole,
	getBandMemberRoles role.GetBandMemberRoles,
	moveBandMemberRole role.MoveBandMemberRole,
) ArtistService {
	return &artistService{
		addAlbumsToArtist:      addAlbumsToArtist,
		addSongsToArtist:       addSongsToArtist,
		createArtist:           createArtist,
		deleteArtist:           deleteArtist,
		deleteImageFromArtist:  deleteImageFromArtist,
		getAllArtists:          getAllArtists,
		getArtist:              getArtist,
		removeAlbumsFromArtist: removeAlbumsFromArtist,
		removeSongsFromArtist:  removeSongsFromArtist,
		saveImageToArtist:      saveImageToArtist,
		updateArtist:           updateArtist,
		createBandMemberRole:   createBandMemberRole,
		deleteBandMemberRole:   deleteBandMemberRole,
		getBandMemberRoles:     getBandMemberRoles,
		moveBandMemberRole:     moveBandMemberRole,
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

func (a *artistService) Delete(request requests.DeleteArtistRequest) *wrapper.ErrorCode {
	return a.deleteArtist.Handle(request)
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

func (a *artistService) RemoveAlbums(request requests.RemoveAlbumsFromArtistRequest) *wrapper.ErrorCode {
	return a.removeAlbumsFromArtist.Handle(request)
}

func (a *artistService) RemoveSongs(request requests.RemoveSongsFromArtistRequest) *wrapper.ErrorCode {
	return a.removeSongsFromArtist.Handle(request)
}

func (a *artistService) SaveImage(file *multipart.FileHeader, id uuid.UUID) *wrapper.ErrorCode {
	return a.saveImageToArtist.Handle(file, id)
}

func (a *artistService) Update(request requests.UpdateArtistRequest) *wrapper.ErrorCode {
	return a.updateArtist.Handle(request)
}
// Band Member - Roles

func (a *artistService) CreateBandMemberRole(request requests.CreateBandMemberRoleRequest, token string) *wrapper.ErrorCode {
	return a.createBandMemberRole.Handle(request, token)
}

func (a *artistService) DeleteBandMemberRole(id uuid.UUID, token string) *wrapper.ErrorCode {
	return a.deleteBandMemberRole.Handle(id, token)
}

func (a *artistService) GetBandMemberRoles(token string) ([]model.BandMemberRole, *wrapper.ErrorCode) {
	return a.getBandMemberRoles.Handle(token)
}

func (a *artistService) MoveBandMemberRole(request requests.MoveBandMemberRoleRequest, token string) *wrapper.ErrorCode {
	return a.moveBandMemberRole.Handle(request, token)
}
