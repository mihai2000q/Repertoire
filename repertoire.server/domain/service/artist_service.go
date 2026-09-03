package service

import (
	"mime/multipart"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/artist"
	"repertoire/server/domain/usecase/artist/bandmember"
	"repertoire/server/internal/httperror"
	"repertoire/server/internal/pagination"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type ArtistService interface {
	AddAlbums(request requests.AddAlbumsToArtistRequest) *httperror.ErrorCode
	AddPerfectRehearsals(request requests.AddPerfectRehearsalsToArtistsRequest) *httperror.ErrorCode
	AddSongs(request requests.AddSongsToArtistRequest) *httperror.ErrorCode
	BulkDelete(request requests.BulkDeleteArtistsRequest) *httperror.ErrorCode
	Create(request requests.CreateArtistRequest, token string) (uuid.UUID, *httperror.ErrorCode)
	Delete(request requests.DeleteArtistRequest) *httperror.ErrorCode
	DeleteImage(id uuid.UUID) *httperror.ErrorCode
	GetAll(request requests.GetArtistsRequest, token string) (pagination.WithTotalCount[model.EnhancedArtist], *httperror.ErrorCode)
	Get(id uuid.UUID) (model.Artist, *httperror.ErrorCode)
	GetFiltersMetadata(request requests.GetArtistFiltersMetadataRequest, token string) (model.ArtistFiltersMetadata, *httperror.ErrorCode)
	RemoveAlbums(request requests.RemoveAlbumsFromArtistRequest) *httperror.ErrorCode
	RemoveSongs(request requests.RemoveSongsFromArtistRequest) *httperror.ErrorCode
	SaveImage(file *multipart.FileHeader, id uuid.UUID) *httperror.ErrorCode
	Update(request requests.UpdateArtistRequest) *httperror.ErrorCode

	CreateBandMember(request requests.CreateBandMemberRequest) (uuid.UUID, *httperror.ErrorCode)
	DeleteBandMember(id uuid.UUID, artistID uuid.UUID) *httperror.ErrorCode
	DeleteBandMemberImage(id uuid.UUID) *httperror.ErrorCode
	MoveBandMember(request requests.MoveBandMemberRequest) *httperror.ErrorCode
	SaveBandMemberImage(file *multipart.FileHeader, id uuid.UUID) *httperror.ErrorCode
	UpdateBandMember(request requests.UpdateBandMemberRequest) *httperror.ErrorCode

	GetBandMemberRoles(token string) ([]model.BandMemberRole, *httperror.ErrorCode)
}

type artistService struct {
	addAlbumsToArtist             artist.AddAlbumsToArtist
	addPerfectRehearsalsToArtists artist.AddPerfectRehearsalsToArtists
	addSongsToArtist              artist.AddSongsToArtist
	bulkDeleteArtists             artist.BulkDeleteArtists
	createArtist                  artist.CreateArtist
	deleteArtist                  artist.DeleteArtist
	deleteImageFromArtist         artist.DeleteImageFromArtist
	getAllArtists                 artist.GetAllArtists
	getArtist                     artist.GetArtist
	getArtistFiltersMetadata      artist.GetArtistFiltersMetadata
	removeAlbumsFromArtist        artist.RemoveAlbumsFromArtist
	removeSongsFromArtist         artist.RemoveSongsFromArtist
	saveImageToArtist             artist.SaveImageToArtist
	updateArtist                  artist.UpdateArtist

	createBandMember          bandmember.CreateBandMember
	deleteBandMember          bandmember.DeleteBandMember
	deleteImageFromBandMember bandmember.DeleteImageFromBandMember
	moveBandMember            bandmember.MoveBandMember
	updateBandMember          bandmember.UpdateBandMember
	saveImageToBandMember     bandmember.SaveImageToBandMember

	getBandMemberRoles bandmember.GetBandMemberRoles
}

func NewArtistService(
	addAlbumsToArtist artist.AddAlbumsToArtist,
	addPerfectRehearsalsToArtists artist.AddPerfectRehearsalsToArtists,
	addSongsToArtist artist.AddSongsToArtist,
	bulkDeleteArtists artist.BulkDeleteArtists,
	createArtist artist.CreateArtist,
	deleteArtist artist.DeleteArtist,
	deleteImageFromArtist artist.DeleteImageFromArtist,
	getAllArtists artist.GetAllArtists,
	getArtist artist.GetArtist,
	getArtistFiltersMetadata artist.GetArtistFiltersMetadata,
	removeAlbumsFromArtist artist.RemoveAlbumsFromArtist,
	removeSongsFromArtist artist.RemoveSongsFromArtist,
	saveImageToArtist artist.SaveImageToArtist,
	updateArtist artist.UpdateArtist,

	createBandMember bandmember.CreateBandMember,
	deleteBandMember bandmember.DeleteBandMember,
	deleteImageFromBandMember bandmember.DeleteImageFromBandMember,
	moveBandMember bandmember.MoveBandMember,
	saveImageToBandMember bandmember.SaveImageToBandMember,
	updateBandMember bandmember.UpdateBandMember,

	getBandMemberRoles bandmember.GetBandMemberRoles,
) ArtistService {
	return &artistService{
		addAlbumsToArtist:             addAlbumsToArtist,
		addPerfectRehearsalsToArtists: addPerfectRehearsalsToArtists,
		addSongsToArtist:              addSongsToArtist,
		bulkDeleteArtists:             bulkDeleteArtists,
		createArtist:                  createArtist,
		deleteArtist:                  deleteArtist,
		deleteImageFromArtist:         deleteImageFromArtist,
		getAllArtists:                 getAllArtists,
		getArtist:                     getArtist,
		getArtistFiltersMetadata:      getArtistFiltersMetadata,
		removeAlbumsFromArtist:        removeAlbumsFromArtist,
		removeSongsFromArtist:         removeSongsFromArtist,
		saveImageToArtist:             saveImageToArtist,
		updateArtist:                  updateArtist,

		createBandMember:          createBandMember,
		deleteBandMember:          deleteBandMember,
		deleteImageFromBandMember: deleteImageFromBandMember,
		moveBandMember:            moveBandMember,
		saveImageToBandMember:     saveImageToBandMember,
		updateBandMember:          updateBandMember,

		getBandMemberRoles: getBandMemberRoles,
	}
}

func (a *artistService) AddAlbums(request requests.AddAlbumsToArtistRequest) *httperror.ErrorCode {
	return a.addAlbumsToArtist.Handle(request)
}

func (a *artistService) AddPerfectRehearsals(request requests.AddPerfectRehearsalsToArtistsRequest) *httperror.ErrorCode {
	return a.addPerfectRehearsalsToArtists.Handle(request)
}

func (a *artistService) AddSongs(request requests.AddSongsToArtistRequest) *httperror.ErrorCode {
	return a.addSongsToArtist.Handle(request)
}

func (a *artistService) BulkDelete(request requests.BulkDeleteArtistsRequest) *httperror.ErrorCode {
	return a.bulkDeleteArtists.Handle(request)
}

func (a *artistService) Create(request requests.CreateArtistRequest, token string) (uuid.UUID, *httperror.ErrorCode) {
	return a.createArtist.Handle(request, token)
}

func (a *artistService) Delete(request requests.DeleteArtistRequest) *httperror.ErrorCode {
	return a.deleteArtist.Handle(request)
}

func (a *artistService) DeleteImage(id uuid.UUID) *httperror.ErrorCode {
	return a.deleteImageFromArtist.Handle(id)
}

func (a *artistService) GetAll(request requests.GetArtistsRequest, token string) (pagination.WithTotalCount[model.EnhancedArtist], *httperror.ErrorCode) {
	return a.getAllArtists.Handle(request, token)
}

func (a *artistService) Get(id uuid.UUID) (model.Artist, *httperror.ErrorCode) {
	return a.getArtist.Handle(id)
}

func (a *artistService) GetFiltersMetadata(
	request requests.GetArtistFiltersMetadataRequest,
	token string,
) (model.ArtistFiltersMetadata, *httperror.ErrorCode) {
	return a.getArtistFiltersMetadata.Handle(request, token)
}

func (a *artistService) RemoveAlbums(request requests.RemoveAlbumsFromArtistRequest) *httperror.ErrorCode {
	return a.removeAlbumsFromArtist.Handle(request)
}

func (a *artistService) RemoveSongs(request requests.RemoveSongsFromArtistRequest) *httperror.ErrorCode {
	return a.removeSongsFromArtist.Handle(request)
}

func (a *artistService) SaveImage(file *multipart.FileHeader, id uuid.UUID) *httperror.ErrorCode {
	return a.saveImageToArtist.Handle(file, id)
}

func (a *artistService) Update(request requests.UpdateArtistRequest) *httperror.ErrorCode {
	return a.updateArtist.Handle(request)
}

// Band Member

func (a *artistService) CreateBandMember(request requests.CreateBandMemberRequest) (uuid.UUID, *httperror.ErrorCode) {
	return a.createBandMember.Handle(request)
}

func (a *artistService) DeleteBandMember(id uuid.UUID, artistID uuid.UUID) *httperror.ErrorCode {
	return a.deleteBandMember.Handle(id, artistID)
}

func (a *artistService) DeleteBandMemberImage(id uuid.UUID) *httperror.ErrorCode {
	return a.deleteImageFromBandMember.Handle(id)
}

func (a *artistService) MoveBandMember(request requests.MoveBandMemberRequest) *httperror.ErrorCode {
	return a.moveBandMember.Handle(request)
}

func (a *artistService) SaveBandMemberImage(file *multipart.FileHeader, id uuid.UUID) *httperror.ErrorCode {
	return a.saveImageToBandMember.Handle(file, id)
}

func (a *artistService) UpdateBandMember(request requests.UpdateBandMemberRequest) *httperror.ErrorCode {
	return a.updateBandMember.Handle(request)
}

// Band Member - Roles

func (a *artistService) GetBandMemberRoles(token string) ([]model.BandMemberRole, *httperror.ErrorCode) {
	return a.getBandMemberRoles.Handle(token)
}
