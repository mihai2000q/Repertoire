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

func TestUpdateBandMember_WhenGetBandMembersFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := member.NewUpdateBandMember(artistRepository)

	request := requests.UpdateBandMemberRequest{
		ID:      uuid.New(),
		Name:    "Some Artist",
		RoleIDs: []uuid.UUID{uuid.New()},
	}

	internalError := errors.New("internal error")
	artistRepository.On("GetBandMember", new(model.BandMember), request.ID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	artistRepository.AssertExpectations(t)
}

func TestUpdateBandMember_WhenBandMembersIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := member.NewUpdateBandMember(artistRepository)

	request := requests.UpdateBandMemberRequest{
		ID:      uuid.New(),
		Name:    "Some Artist",
		RoleIDs: []uuid.UUID{uuid.New()},
	}

	artistRepository.On("GetBandMember", new(model.BandMember), request.ID).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "band member not found", errCode.Error.Error())

	artistRepository.AssertExpectations(t)
}

func TestUpdateBandMember_WhenGetRolesFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := member.NewUpdateBandMember(artistRepository)

	request := requests.UpdateBandMemberRequest{
		ID:      uuid.New(),
		Name:    "Some Artist",
		RoleIDs: []uuid.UUID{uuid.New()},
	}

	// given - mocking
	mockBandMember := &model.BandMember{ID: request.ID}
	artistRepository.On("GetBandMember", new(model.BandMember), request.ID).
		Return(nil, mockBandMember).
		Once()

	internalError := errors.New("internal error")
	artistRepository.On("GetBandMemberRolesByIDs", new([]model.BandMemberRole), request.RoleIDs).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	artistRepository.AssertExpectations(t)
}

func TestUpdateBandMember_WhenUpdateBandMemberFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := member.NewUpdateBandMember(artistRepository)

	request := requests.UpdateBandMemberRequest{
		ID:      uuid.New(),
		Name:    "Some Artist",
		RoleIDs: []uuid.UUID{uuid.New()},
	}

	// given - mocking
	mockBandMember := &model.BandMember{ID: request.ID}
	artistRepository.On("GetBandMember", new(model.BandMember), request.ID).
		Return(nil, mockBandMember).
		Once()

	artistRepository.On("GetBandMemberRolesByIDs", new([]model.BandMemberRole), request.RoleIDs).
		Return(nil).
		Once()

	internalError := errors.New("internal error")
	artistRepository.On("UpdateBandMember", mock.IsType(new(model.BandMember))).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	artistRepository.AssertExpectations(t)
}

func TestUpdateBandMember_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := member.NewUpdateBandMember(artistRepository)

	request := requests.UpdateBandMemberRequest{
		ID:      uuid.New(),
		Name:    "Some Artist",
		RoleIDs: []uuid.UUID{uuid.New()},
	}

	// given - mocking
	mockBandMember := &model.BandMember{ID: request.ID}
	artistRepository.On("GetBandMember", new(model.BandMember), request.ID).
		Return(nil, mockBandMember).
		Once()

	roles := []model.BandMemberRole{
		{ID: request.RoleIDs[0]},
	}
	artistRepository.On("GetBandMemberRolesByIDs", new([]model.BandMemberRole), request.RoleIDs).
		Return(nil, &roles).
		Once()

	artistRepository.On("UpdateBandMember", mock.IsType(new(model.BandMember))).
		Run(func(args mock.Arguments) {
			newBandMember := args.Get(0).(*model.BandMember)
			assertUpdatedBandMember(t, request, *newBandMember)
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)

	artistRepository.AssertExpectations(t)
}

func assertUpdatedBandMember(
	t *testing.T,
	request requests.UpdateBandMemberRequest,
	member model.BandMember,
) {
	assert.Equal(t, request.Name, member.Name)
	assert.Len(t, member.Roles, len(request.RoleIDs))
	for i, roleID := range request.RoleIDs {
		assert.Equal(t, roleID, member.Roles[i].ID)
	}
}
