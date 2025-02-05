package role

import (
	"errors"
	"net/http"
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

func TestDeleteBandMemberRole_WhenGetUserIdFromJwtFails_ShouldReturnError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	_uut := role.NewDeleteBandMemberRole(nil, jwtService)

	id := uuid.New()
	token := "this is a token"

	forbiddenError := wrapper.ForbiddenError(errors.New("forbidden error"))
	jwtService.On("GetUserIdFromJwt", token).Return(uuid.Nil, forbiddenError).Once()

	// when
	errCode := _uut.Handle(id, token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, forbiddenError, errCode)

	jwtService.AssertExpectations(t)
}

func TestDeleteBandMemberRole_WhenGetBandMemberRolesFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := role.NewDeleteBandMemberRole(artistRepository, jwtService)

	id := uuid.New()
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	internalError := errors.New("internal error")
	artistRepository.On("GetBandMemberRoles", new([]model.BandMemberRole), userID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id, token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	jwtService.AssertExpectations(t)
	artistRepository.AssertExpectations(t)
}

func TestDeleteBandMemberRole_WhenRoleIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := role.NewDeleteBandMemberRole(artistRepository, jwtService)

	id := uuid.New()
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	bandMemberRoles := &[]model.BandMemberRole{
		{ID: uuid.New()},
	}
	artistRepository.On("GetBandMemberRoles", new([]model.BandMemberRole), userID).
		Return(nil, bandMemberRoles).
		Once()

	// when
	errCode := _uut.Handle(id, token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "band member role not found", errCode.Error.Error())

	jwtService.AssertExpectations(t)
	artistRepository.AssertExpectations(t)
}

func TestDeleteBandMemberRole_WhenUpdateAllBandMemberRolesFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := role.NewDeleteBandMemberRole(artistRepository, jwtService)

	id := uuid.New()
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	bandMemberRoles := &[]model.BandMemberRole{
		{ID: id},
	}
	artistRepository.On("GetBandMemberRoles", new([]model.BandMemberRole), userID).
		Return(nil, bandMemberRoles).
		Once()

	internalError := errors.New("internal error")
	artistRepository.On("UpdateAllBandMemberRoles", mock.IsType(bandMemberRoles)).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id, token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	jwtService.AssertExpectations(t)
	artistRepository.AssertExpectations(t)
}

func TestDeleteBandMemberRole_WhenDeleteBandMemberRoleFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := role.NewDeleteBandMemberRole(artistRepository, jwtService)

	id := uuid.New()
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	bandMemberRoles := &[]model.BandMemberRole{
		{ID: id},
	}
	artistRepository.On("GetBandMemberRoles", new([]model.BandMemberRole), userID).
		Return(nil, bandMemberRoles).
		Once()

	artistRepository.On("UpdateAllBandMemberRoles", mock.IsType(bandMemberRoles)).
		Return(nil).
		Once()

	internalError := errors.New("internal error")
	artistRepository.On("DeleteBandMemberRole", id).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id, token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	jwtService.AssertExpectations(t)
	artistRepository.AssertExpectations(t)
}

func TestDeleteBandMemberRole_WhenSuccessful_ShouldReturnGuitarTunings(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := role.NewDeleteBandMemberRole(artistRepository, jwtService)

	id := uuid.New()
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	bandMemberRoles := &[]model.BandMemberRole{
		{ID: id},
	}
	artistRepository.On("GetBandMemberRoles", new([]model.BandMemberRole), userID).
		Return(nil, bandMemberRoles).
		Once()

	artistRepository.On("UpdateAllBandMemberRoles", mock.IsType(bandMemberRoles)).
		Run(func(args mock.Arguments) {
			newBandMemberRoles := args.Get(0).(*[]model.BandMemberRole)
			sortedBandMemberRoles := slices.DeleteFunc(*newBandMemberRoles, func(s model.BandMemberRole) bool {
				return s.ID == id
			})
			for i, bandMemberRole := range sortedBandMemberRoles {
				assert.Equal(t, i, bandMemberRole.Order)
			}
		}).
		Return(nil).
		Once()

	artistRepository.On("DeleteBandMemberRole", id).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(id, token)

	// then
	assert.Nil(t, errCode)

	jwtService.AssertExpectations(t)
	artistRepository.AssertExpectations(t)
}
