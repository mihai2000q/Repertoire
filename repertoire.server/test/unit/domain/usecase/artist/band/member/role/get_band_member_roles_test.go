package role

import (
	"errors"
	"net/http"
	"repertoire/server/domain/usecase/artist/band/member/role"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetBandMemberRoles_WhenGetUserIdFromJwtFails_ShouldReturnError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	_uut := role.NewGetBandMemberRoles(nil, jwtService)

	token := "this is a token"

	forbiddenError := wrapper.ForbiddenError(errors.New("forbidden error"))
	jwtService.On("GetUserIdFromJwt", token).Return(uuid.Nil, forbiddenError).Once()

	// when
	BandMemberRoles, errCode := _uut.Handle(token)

	// then
	assert.Empty(t, BandMemberRoles)
	assert.NotNil(t, errCode)
	assert.Equal(t, forbiddenError, errCode)

	jwtService.AssertExpectations(t)
}

func TestGetBandMemberRoles_WhenGetBandMemberRolesFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := role.NewGetBandMemberRoles(artistRepository, jwtService)

	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	internalError := errors.New("internal error")
	artistRepository.On("GetBandMemberRoles", new([]model.BandMemberRole), userID).Return(internalError).Once()

	// when
	BandMemberRoles, errCode := _uut.Handle(token)

	// then
	assert.Empty(t, BandMemberRoles)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	jwtService.AssertExpectations(t)
	artistRepository.AssertExpectations(t)
}

func TestGetBandMemberRoles_WhenSuccessful_ShouldReturnBandMemberRoles(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := role.NewGetBandMemberRoles(artistRepository, jwtService)

	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	expectedRoles := &[]model.BandMemberRole{
		{ID: uuid.New(), Name: "Guitarist"},
		{ID: uuid.New(), Name: "Vocalist"},
	}
	artistRepository.On("GetBandMemberRoles", mock.IsType(expectedRoles), userID).
		Return(nil, expectedRoles).
		Once()

	// when
	BandMemberRoles, errCode := _uut.Handle(token)

	// then
	assert.Equal(t, expectedRoles, &BandMemberRoles)
	assert.Nil(t, errCode)

	jwtService.AssertExpectations(t)
	artistRepository.AssertExpectations(t)
}
