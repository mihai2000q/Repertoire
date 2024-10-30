package song

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAll_WhenGetUserIdFromJwtFails_ShouldReturnForbiddenError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	_uut := &GetAllSongs{
		jwtService: jwtService,
	}
	request := requests.GetSongsRequest{}
	token := "this is the token"

	forbiddenError := wrapper.UnauthorizedError(errors.New("forbidden error"))
	jwtService.On("GetUserIdFromJwt", token).Return(uuid.Nil, forbiddenError).Once()

	// when
	result, errCode := _uut.Handle(request, token)

	// then
	assert.Empty(t, result)
	assert.NotNil(t, errCode)
	assert.Equal(t, forbiddenError, errCode)

	jwtService.AssertExpectations(t)
}

func TestGetAll_WhenGetSongsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	jwtService := new(service.JwtServiceMock)
	_uut := &GetAllSongs{
		repository: songRepository,
		jwtService: jwtService,
	}
	request := requests.GetSongsRequest{}
	token := "this is the token"

	userID := uuid.New()

	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	internalError := errors.New("internal error")
	songRepository.
		On(
			"GetAllByUser",
			mock.Anything,
			userID,
			request.CurrentPage,
			request.PageSize,
			request.OrderBy,
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

	songRepository.AssertExpectations(t)
	jwtService.AssertExpectations(t)
}

func TestGetAll_WhenGetSongsCountFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	jwtService := new(service.JwtServiceMock)
	_uut := &GetAllSongs{
		repository: songRepository,
		jwtService: jwtService,
	}
	request := requests.GetSongsRequest{}
	token := "this is the token"

	userID := uuid.New()
	expectedSongs := &[]model.Song{
		{Title: "Some Song"},
		{Title: "Some other Song"},
	}

	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	songRepository.
		On(
			"GetAllByUser",
			mock.IsType(expectedSongs),
			userID,
			request.CurrentPage,
			request.PageSize,
			request.OrderBy,
		).
		Return(nil, expectedSongs).
		Once()

	internalError := errors.New("internal error")
	songRepository.
		On(
			"GetAllByUserCount",
			mock.Anything,
			userID,
		).
		Return(internalError).
		Once()

	// when
	result, errCode := _uut.Handle(request, token)

	// then
	assert.Equal(t, expectedSongs, &result.Models)
	assert.Empty(t, result.TotalCount)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
	jwtService.AssertExpectations(t)
}

func TestGetAll_WhenSuccessful_ShouldReturnSongsWithTotalCount(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	jwtService := new(service.JwtServiceMock)
	_uut := &GetAllSongs{
		repository: songRepository,
		jwtService: jwtService,
	}
	request := requests.GetSongsRequest{}
	token := "this is the token"

	userID := uuid.New()
	expectedSongs := &[]model.Song{
		{Title: "Some Song"},
		{Title: "Some other Song"},
	}
	expectedTotalCount := &[]int64{20}[0]

	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	songRepository.
		On(
			"GetAllByUser",
			mock.IsType(expectedSongs),
			userID,
			request.CurrentPage,
			request.PageSize,
			request.OrderBy,
		).
		Return(nil, expectedSongs).
		Once()

	songRepository.
		On(
			"GetAllByUserCount",
			mock.IsType(expectedTotalCount),
			userID,
		).
		Return(nil, expectedTotalCount).
		Once()

	// when
	result, errCode := _uut.Handle(request, token)

	// then
	assert.Equal(t, expectedSongs, &result.Models)
	assert.Equal(t, expectedTotalCount, &result.TotalCount)
	assert.Nil(t, errCode)

	songRepository.AssertExpectations(t)
	jwtService.AssertExpectations(t)
}
