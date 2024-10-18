package artist

import (
	"errors"
	"net/http"
	"repertoire/api/requests"
	"repertoire/data/repository"
	"repertoire/data/service"
	"repertoire/models"
	"repertoire/utils/wrapper"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAll_WhenGetUserIdFromJwtFails_ShouldReturnForbiddenError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	_uut := &GetAllArtists{
		jwtService: jwtService,
	}
	request := requests.GetArtistsRequest{}
	token := "This is a token"

	forbiddenError := wrapper.ForbiddenError(errors.New("forbidden error"))
	jwtService.On("GetUserIdFromJwt", token).Return(uuid.Nil, forbiddenError).Once()

	// when
	result, errCode := _uut.Handle(request, token)

	// then
	assert.Empty(t, result)
	assert.NotNil(t, errCode)
	assert.Equal(t, forbiddenError, errCode)

	jwtService.AssertExpectations(t)
}

func TestGetAll_WhenGetArtistsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	jwtService := new(service.JwtServiceMock)
	_uut := &GetAllArtists{
		repository: artistRepository,
		jwtService: jwtService,
	}
	request := requests.GetArtistsRequest{}
	token := "This is a token"

	userId := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userId, nil).Once()

	internalError := errors.New("internal error")
	artistRepository.
		On(
			"GetAllByUser",
			mock.Anything,
			userId,
			request.CurrentPage,
			request.PageSize,
		).
		Return(internalError).
		Once()

	// when
	result, errCode := _uut.Handle(request, token)

	// then
	assert.Empty(t, result)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	artistRepository.AssertExpectations(t)
	jwtService.AssertExpectations(t)
}

func TestGetAll_WhenGetArtistsCountFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	jwtService := new(service.JwtServiceMock)
	_uut := &GetAllArtists{
		repository: artistRepository,
		jwtService: jwtService,
	}
	request := requests.GetArtistsRequest{}
	token := "This is a token"

	expectedArtists := &[]models.Artist{
		{Name: "Some Artist"},
		{Name: "Some other Artist"},
	}

	userId := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userId, nil).Once()

	artistRepository.
		On(
			"GetAllByUser",
			mock.IsType(expectedArtists),
			userId,
			request.CurrentPage,
			request.PageSize,
		).
		Return(nil, expectedArtists).
		Once()

	internalError := errors.New("internal error")
	artistRepository.
		On(
			"GetAllByUserCount",
			mock.Anything,
			userId,
		).
		Return(internalError).
		Once()

	// when
	result, errCode := _uut.Handle(request, token)

	// then
	assert.Equal(t, expectedArtists, &result.Data)
	assert.Empty(t, result.TotalCount)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	artistRepository.AssertExpectations(t)
	jwtService.AssertExpectations(t)
}

func TestGetAll_WhenSuccessful_ShouldReturnArtistsWithTotalCount(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	jwtService := new(service.JwtServiceMock)
	_uut := &GetAllArtists{
		repository: artistRepository,
		jwtService: jwtService,
	}
	request := requests.GetArtistsRequest{}
	token := "This is a token"

	expectedArtists := &[]models.Artist{
		{Name: "Some Artist"},
		{Name: "Some other Artist"},
	}
	expectedTotalCount := &[]int64{20}[0]

	userId := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userId, nil).Once()

	artistRepository.
		On(
			"GetAllByUser",
			mock.IsType(expectedArtists),
			userId,
			request.CurrentPage,
			request.PageSize,
		).
		Return(nil, expectedArtists).
		Once()

	artistRepository.
		On(
			"GetAllByUserCount",
			mock.IsType(expectedTotalCount),
			userId,
		).
		Return(nil, expectedTotalCount).
		Once()

	// when
	result, errCode := _uut.Handle(request, token)

	// then
	assert.Equal(t, expectedArtists, &result.Data)
	assert.Equal(t, expectedTotalCount, &result.TotalCount)
	assert.Nil(t, errCode)

	artistRepository.AssertExpectations(t)
	jwtService.AssertExpectations(t)
}
