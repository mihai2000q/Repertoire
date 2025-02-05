package role

import (
	"cmp"
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/udata/band/member/role"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"slices"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMoveBandMemberRole_WhenGetUserIdFromJwtFails_ShouldReturnError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	_uut := role.NewMoveBandMemberRole(nil, jwtService)

	request := requests.MoveBandMemberRoleRequest{
		ID:     uuid.New(),
		OverID: uuid.New(),
	}
	token := "this is a token"

	forbiddenError := wrapper.ForbiddenError(errors.New("forbidden error"))
	jwtService.On("GetUserIdFromJwt", token).Return(uuid.Nil, forbiddenError).Once()

	// when
	errCode := _uut.Handle(request, token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, forbiddenError, errCode)

	jwtService.AssertExpectations(t)
}

func TestMoveBandMemberRole_WhenGetBandMemberRolesFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := role.NewMoveBandMemberRole(artistRepository, jwtService)

	request := requests.MoveBandMemberRoleRequest{
		ID:     uuid.New(),
		OverID: uuid.New(),
	}
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	internalError := errors.New("internal error")
	artistRepository.On("GetBandMemberRoles", new([]model.BandMemberRole), userID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request, token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	jwtService.AssertExpectations(t)
	artistRepository.AssertExpectations(t)
}

func TestMoveBandMemberRole_WhenBandMemberRoleIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := role.NewMoveBandMemberRole(artistRepository, jwtService)

	request := requests.MoveBandMemberRoleRequest{
		ID:     uuid.New(),
		OverID: uuid.New(),
	}
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	expectedRoles := &[]model.BandMemberRole{}
	artistRepository.On("GetBandMemberRoles", new([]model.BandMemberRole), userID).
		Return(nil, expectedRoles).
		Once()

	// when
	errCode := _uut.Handle(request, token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "role not found", errCode.Error.Error())

	jwtService.AssertExpectations(t)
	artistRepository.AssertExpectations(t)
}

func TestMoveBandMemberRole_WhenOverBandMemberRoleIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := role.NewMoveBandMemberRole(artistRepository, jwtService)

	request := requests.MoveBandMemberRoleRequest{
		ID:     uuid.New(),
		OverID: uuid.New(),
	}
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	expectedRoles := &[]model.BandMemberRole{
		{ID: request.ID},
	}
	artistRepository.On("GetBandMemberRoles", new([]model.BandMemberRole), userID).
		Return(nil, expectedRoles).
		Once()

	// when
	errCode := _uut.Handle(request, token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "over role not found", errCode.Error.Error())

	jwtService.AssertExpectations(t)
	artistRepository.AssertExpectations(t)
}

func TestMoveBandMemberRole_WhenUpdateAllBandMemberRolesFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := role.NewMoveBandMemberRole(artistRepository, jwtService)

	request := requests.MoveBandMemberRoleRequest{
		ID:     uuid.New(),
		OverID: uuid.New(),
	}
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	expectedRoles := &[]model.BandMemberRole{
		{ID: request.ID},
		{ID: request.OverID},
	}
	artistRepository.On("GetBandMemberRoles", new([]model.BandMemberRole), userID).
		Return(nil, expectedRoles).
		Once()

	internalError := errors.New("internal error")
	artistRepository.On("UpdateAllBandMemberRoles", mock.IsType(expectedRoles)).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request, token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	jwtService.AssertExpectations(t)
	artistRepository.AssertExpectations(t)
}

func TestMoveBandMemberRole_WhenSuccessful_ShouldReturnBandMemberRoles(t *testing.T) {
	tests := []struct {
		name      string
		role      *[]model.BandMemberRole
		index     uint
		overIndex uint
	}{
		{
			"Use case 1",
			&[]model.BandMemberRole{
				{ID: uuid.New(), Order: 0},
				{ID: uuid.New(), Order: 1},
				{ID: uuid.New(), Order: 2},
				{ID: uuid.New(), Order: 3},
				{ID: uuid.New(), Order: 4},
			},
			1,
			3,
		},
		{
			"Use case 2",
			&[]model.BandMemberRole{
				{ID: uuid.New(), Order: 0},
				{ID: uuid.New(), Order: 1},
				{ID: uuid.New(), Order: 2},
				{ID: uuid.New(), Order: 3},
				{ID: uuid.New(), Order: 4},
			},
			3,
			1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			jwtService := new(service.JwtServiceMock)
			artistRepository := new(repository.ArtistRepositoryMock)
			_uut := role.NewMoveBandMemberRole(artistRepository, jwtService)

			request := requests.MoveBandMemberRoleRequest{
				ID:     (*tt.role)[tt.index].ID,
				OverID: (*tt.role)[tt.overIndex].ID,
			}
			token := "this is a token"

			// given - mocking
			userID := uuid.New()
			jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

			artistRepository.On("GetBandMemberRoles", new([]model.BandMemberRole), userID).
				Return(nil, tt.role).
				Once()

			artistRepository.On("UpdateAllBandMemberRoles", mock.IsType(tt.role)).
				Run(func(args mock.Arguments) {
					newBandMemberRoles := args.Get(0).(*[]model.BandMemberRole)
					sortedBandMemberRoles := slices.Clone(*newBandMemberRoles)
					slices.SortFunc(sortedBandMemberRoles, func(a, b model.BandMemberRole) int {
						return cmp.Compare(a.Order, b.Order)
					})
					if tt.index < tt.overIndex {
						assert.Equal(t, sortedBandMemberRoles[tt.overIndex-1].ID, request.OverID)
					} else if tt.index > tt.overIndex {
						assert.Equal(t, sortedBandMemberRoles[tt.overIndex+1].ID, request.OverID)
					}
					for i, sectionRole := range sortedBandMemberRoles {
						assert.Equal(t, uint(i), sectionRole.Order)
					}
				}).
				Return(nil).
				Once()

			// when
			errCode := _uut.Handle(request, token)

			// then
			assert.Nil(t, errCode)

			jwtService.AssertExpectations(t)
			artistRepository.AssertExpectations(t)
		})
	}
}
