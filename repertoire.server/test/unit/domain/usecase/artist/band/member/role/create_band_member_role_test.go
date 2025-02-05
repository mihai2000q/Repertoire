package role

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/udata/band/member/role"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateBandMemberRole_WhenGetUserIdFromJwtFails_ShouldReturnError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	_uut := role.NewCreateBandMemberRole(nil, jwtService)

	request := requests.CreateBandMemberRoleRequest{
		Name: "New Type",
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

func TestCreateBandMemberRole_WhenCountBandMemberRolesFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := role.NewCreateBandMemberRole(artistRepository, jwtService)

	request := requests.CreateBandMemberRoleRequest{
		Name: "New Role",
	}
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	internalError := errors.New("internal error")
	artistRepository.On("CountBandMemberRoles", new(int64), userID).Return(internalError).Once()

	// when
	errCode := _uut.Handle(request, token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	jwtService.AssertExpectations(t)
	artistRepository.AssertExpectations(t)
}

func TestCreateBandMemberRole_WhenCreateBandMemberRoleFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := role.NewCreateBandMemberRole(artistRepository, jwtService)

	request := requests.CreateBandMemberRoleRequest{
		Name: "New Type",
	}
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	artistRepository.On("CountBandMemberRoles", new(int64), userID).Return(nil).Once()

	internalError := errors.New("internal error")
	artistRepository.On("CreateBandMemberRole", mock.IsType(new(model.BandMemberRole))).
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

func TestCreateBandMemberRole_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := role.NewCreateBandMemberRole(artistRepository, jwtService)

	request := requests.CreateBandMemberRoleRequest{
		Name: "New Type",
	}
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	count := &[]int64{20}[0]
	artistRepository.On("CountBandMemberRoles", mock.IsType(count), userID).
		Return(nil, count).
		Once()

	artistRepository.On("CreateBandMemberRole", mock.IsType(new(model.BandMemberRole))).
		Run(func(args mock.Arguments) {
			newType := args.Get(0).(*model.BandMemberRole)
			assertCreatedBandMemberRole(t, *newType, request, userID, count)
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request, token)

	// then
	assert.Nil(t, errCode)

	jwtService.AssertExpectations(t)
	artistRepository.AssertExpectations(t)
}

func assertCreatedBandMemberRole(
	t *testing.T,
	bandMemberRole model.BandMemberRole,
	request requests.CreateBandMemberRoleRequest,
	userID uuid.UUID,
	count *int64,
) {
	assert.NotEmpty(t, bandMemberRole.ID)
	assert.Equal(t, request.Name, bandMemberRole.Name)
	assert.Equal(t, uint(*count), bandMemberRole.Order)
	assert.Equal(t, userID, bandMemberRole.UserID)
}
