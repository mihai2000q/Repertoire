package repository

import (
	"repertoire/server/model"

	"github.com/stretchr/testify/mock"

	"github.com/google/uuid"
)

type ArtistRepositoryMock struct {
	mock.Mock
}

func (a *ArtistRepositoryMock) Get(artist *model.Artist, id uuid.UUID) error {
	args := a.Called(artist, id)

	if len(args) > 1 {
		*artist = *args.Get(1).(*model.Artist)
	}

	return args.Error(0)
}

func (a *ArtistRepositoryMock) GetWithAssociations(artist *model.Artist, id uuid.UUID) error {
	args := a.Called(artist, id)

	if len(args) > 1 {
		*artist = *args.Get(1).(*model.Artist)
	}

	return args.Error(0)
}

func (a *ArtistRepositoryMock) GetWithBandMembers(artist *model.Artist, id uuid.UUID) error {
	args := a.Called(artist, id)

	if len(args) > 1 {
		*artist = *args.Get(1).(*model.Artist)
	}

	return args.Error(0)
}

func (a *ArtistRepositoryMock) GetAllByIDsWithSongs(artists *[]model.Artist, ids []uuid.UUID) error {
	args := a.Called(artists, ids)

	if len(args) > 1 {
		*artists = *args.Get(1).(*[]model.Artist)
	}

	return args.Error(0)
}

func (a *ArtistRepositoryMock) GetAllByUser(
	artists *[]model.Artist,
	userID uuid.UUID,
	currentPage *int,
	pageSize *int,
	orderBy []string,
	searchBy []string,
) error {
	args := a.Called(artists, userID, currentPage, pageSize, orderBy, searchBy)

	if len(args) > 1 {
		*artists = *args.Get(1).(*[]model.Artist)
	}

	return args.Error(0)
}

func (a *ArtistRepositoryMock) GetAllByUserCount(count *int64, userID uuid.UUID, searchBy []string) error {
	args := a.Called(count, userID, searchBy)

	if len(args) > 1 {
		*count = *args.Get(1).(*int64)
	}

	return args.Error(0)
}

func (a *ArtistRepositoryMock) Create(artist *model.Artist) error {
	args := a.Called(artist)
	return args.Error(0)
}

func (a *ArtistRepositoryMock) Update(artist *model.Artist) error {
	args := a.Called(artist)
	return args.Error(0)
}

func (a *ArtistRepositoryMock) UpdateWithAssociations(artist *model.Artist) error {
	args := a.Called(artist)
	return args.Error(0)
}

func (a *ArtistRepositoryMock) Delete(id uuid.UUID) error {
	args := a.Called(id)
	return args.Error(0)
}

func (a *ArtistRepositoryMock) DeleteAlbums(id uuid.UUID) error {
	args := a.Called(id)
	return args.Error(0)
}

func (a *ArtistRepositoryMock) DeleteSongs(id uuid.UUID) error {
	args := a.Called(id)
	return args.Error(0)
}

// Band Member

func (a *ArtistRepositoryMock) GetBandMember(bandMember *model.BandMember, id uuid.UUID) error {
	args := a.Called(bandMember, id)

	if len(args) > 1 {
		*bandMember = *args.Get(1).(*model.BandMember)
	}

	return args.Error(0)
}

func (a *ArtistRepositoryMock) GetBandMemberWithArtist(bandMember *model.BandMember, id uuid.UUID) error {
	args := a.Called(bandMember, id)

	if len(args) > 1 {
		*bandMember = *args.Get(1).(*model.BandMember)
	}

	return args.Error(0)
}

func (a *ArtistRepositoryMock) CreateBandMember(bandMember *model.BandMember) error {
	args := a.Called(bandMember)
	return args.Error(0)
}

func (a *ArtistRepositoryMock) UpdateBandMember(bandMember *model.BandMember) error {
	args := a.Called(bandMember)
	return args.Error(0)
}

func (a *ArtistRepositoryMock) ReplaceRolesFromBandMember(roles *[]model.BandMemberRole, bandMember *model.BandMember) error {
	args := a.Called(roles, bandMember)
	return args.Error(0)
}

func (a *ArtistRepositoryMock) DeleteBandMember(id uuid.UUID) error {
	args := a.Called(id)
	return args.Error(0)
}

// Band Member - Roles

func (a *ArtistRepositoryMock) GetBandMemberRoles(bandMemberRoles *[]model.BandMemberRole, userID uuid.UUID) error {
	args := a.Called(bandMemberRoles, userID)

	if len(args) > 1 {
		*bandMemberRoles = *args.Get(1).(*[]model.BandMemberRole)
	}

	return args.Error(0)
}

func (a *ArtistRepositoryMock) GetBandMemberRolesByIDs(bandMemberRoles *[]model.BandMemberRole, ids []uuid.UUID) error {
	args := a.Called(bandMemberRoles, ids)

	if len(args) > 1 {
		*bandMemberRoles = *args.Get(1).(*[]model.BandMemberRole)
	}

	return args.Error(0)
}

func (a *ArtistRepositoryMock) CountBandMemberRoles(count *int64, userID uuid.UUID) error {
	args := a.Called(count, userID)

	if len(args) > 1 {
		*count = *args.Get(1).(*int64)
	}

	return args.Error(0)
}

func (a *ArtistRepositoryMock) CreateBandMemberRole(bandMemberRole *model.BandMemberRole) error {
	args := a.Called(bandMemberRole)
	return args.Error(0)
}

func (a *ArtistRepositoryMock) UpdateBandMemberRole(bandMemberRole *model.BandMemberRole) error {
	args := a.Called(bandMemberRole)
	return args.Error(0)
}

func (a *ArtistRepositoryMock) UpdateAllBandMemberRoles(bandMemberRoles *[]model.BandMemberRole) error {
	args := a.Called(bandMemberRoles)
	return args.Error(0)
}

func (a *ArtistRepositoryMock) DeleteBandMemberRole(id uuid.UUID) error {
	args := a.Called(id)
	return args.Error(0)
}
