package member

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/artist/band/member"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"testing"
)

func TestCreateBandMember_WhenGetArtistFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := member.NewCreateBandMember(artistRepository)

	request := requests.CreateBandMemberRequest{
		ArtistID: uuid.New(),
		Name:     "Some Artist",
		RoleIDs:  []uuid.UUID{uuid.New()},
	}

	internalError := errors.New("internal error")
	artistRepository.On("GetWithBandMembers", new(model.Artist), request.ArtistID).
		Return(internalError).
		Once()

	// when
	id, errCode := _uut.Handle(request)

	// then
	assert.Empty(t, id)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	artistRepository.AssertExpectations(t)
}

func TestCreateBandMember_WhenArtistIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := member.NewCreateBandMember(artistRepository)

	request := requests.CreateBandMemberRequest{
		ArtistID: uuid.New(),
		Name:     "Some Artist",
		RoleIDs:  []uuid.UUID{uuid.New()},
	}

	artistRepository.On("GetWithBandMembers", new(model.Artist), request.ArtistID).
		Return(nil).
		Once()

	// when
	id, errCode := _uut.Handle(request)

	// then
	assert.Empty(t, id)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "artist not found", errCode.Error.Error())

	artistRepository.AssertExpectations(t)
}

func TestCreateBandMember_WhenArtistIsNotBand_ShouldReturnBadRequestError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := member.NewCreateBandMember(artistRepository)

	request := requests.CreateBandMemberRequest{
		ArtistID: uuid.New(),
		Name:     "Some Artist",
		RoleIDs:  []uuid.UUID{uuid.New()},
	}

	artist := &model.Artist{ID: request.ArtistID, IsBand: false}
	artistRepository.On("GetWithBandMembers", mock.IsType(artist), request.ArtistID).
		Return(nil, artist).
		Once()

	// when
	id, errCode := _uut.Handle(request)

	// then
	assert.Empty(t, id)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusBadRequest, errCode.Code)
	assert.Equal(t, "artist is not band", errCode.Error.Error())

	artistRepository.AssertExpectations(t)
}

func TestCreateBandMember_WhenGetBandMemberRolesFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := member.NewCreateBandMember(artistRepository)

	request := requests.CreateBandMemberRequest{
		ArtistID: uuid.New(),
		Name:     "Some Artist",
		RoleIDs:  []uuid.UUID{uuid.New()},
	}

	artist := &model.Artist{ID: request.ArtistID, IsBand: true}
	artistRepository.On("GetWithBandMembers", mock.IsType(artist), request.ArtistID).
		Return(nil, artist).
		Once()

	internalError := errors.New("internal error")
	artistRepository.On("GetBandMemberRolesByIDs", new([]model.BandMemberRole), request.RoleIDs).
		Return(internalError).
		Once()

	// when
	id, errCode := _uut.Handle(request)

	// then
	assert.Empty(t, id)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	artistRepository.AssertExpectations(t)
}

func TestCreateBandMember_WhenCreateBandMemberFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := member.NewCreateBandMember(artistRepository)

	request := requests.CreateBandMemberRequest{
		ArtistID: uuid.New(),
		Name:     "Some Artist",
		RoleIDs:  []uuid.UUID{uuid.New()},
	}

	artist := &model.Artist{ID: request.ArtistID, IsBand: true}
	artistRepository.On("GetWithBandMembers", mock.IsType(artist), request.ArtistID).
		Return(nil, artist).
		Once()

	roles := &[]model.BandMemberRole{
		{ID: request.RoleIDs[0]},
	}
	artistRepository.On("GetBandMemberRolesByIDs", mock.IsType(roles), request.RoleIDs).
		Return(nil, roles).
		Once()

	internalError := errors.New("internal error")
	artistRepository.On("CreateBandMember", mock.IsType(new(model.BandMember))).
		Return(internalError).
		Once()

	// when
	id, errCode := _uut.Handle(request)

	// then
	assert.Empty(t, id)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	artistRepository.AssertExpectations(t)
}

func TestCreateBandMember_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := member.NewCreateBandMember(artistRepository)

	request := requests.CreateBandMemberRequest{
		ArtistID: uuid.New(),
		Name:     "Some Artist",
		RoleIDs:  []uuid.UUID{uuid.New(), uuid.New()},
	}

	artist := &model.Artist{ID: request.ArtistID, IsBand: true}
	artistRepository.On("GetWithBandMembers", mock.IsType(artist), request.ArtistID).
		Return(nil, artist).
		Once()

	roles := &[]model.BandMemberRole{
		{ID: request.RoleIDs[0]},
		{ID: request.RoleIDs[1]},
	}
	artistRepository.On("GetBandMemberRolesByIDs", mock.IsType(roles), request.RoleIDs).
		Return(nil, roles).
		Once()

	var bandMemberID uuid.UUID
	artistRepository.On("CreateBandMember", mock.IsType(new(model.BandMember))).
		Run(func(args mock.Arguments) {
			newBandMember := args.Get(0).(*model.BandMember)
			assertCreatedBandMember(t, request, *newBandMember, len(artist.BandMembers))
			bandMemberID = newBandMember.ID
		}).
		Return(nil).
		Once()

	// when
	id, errCode := _uut.Handle(request)

	// then
	assert.Equal(t, bandMemberID, id)
	assert.Nil(t, errCode)

	artistRepository.AssertExpectations(t)
}

func assertCreatedBandMember(
	t *testing.T,
	request requests.CreateBandMemberRequest,
	bandMember model.BandMember,
	bandMembersLen int,
) {
	assert.NotEmpty(t, bandMember.ID)
	assert.Equal(t, request.Name, bandMember.Name)
	assert.Equal(t, request.ArtistID, bandMember.ArtistID)
	assert.Equal(t, uint(bandMembersLen), bandMember.Order)
	assert.Empty(t, bandMember.ImageURL)
	assert.Len(t, bandMember.Roles, len(request.RoleIDs))
	for i, roleID := range request.RoleIDs {
		assert.Equal(t, roleID, bandMember.Roles[i].ID)
	}
}
