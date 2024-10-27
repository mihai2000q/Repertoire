package song

import (
	"errors"
	"net/http"
	"repertoire/api/requests"
	"repertoire/data/repository"
	"repertoire/data/service"
	"repertoire/model"
	"repertoire/utils/wrapper"
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

	userId := uuid.New()

	jwtService.On("GetUserIdFromJwt", token).Return(userId, nil).Once()

	internalError := errors.New("internal error")
	songRepository.
		On(
			"GetAllByUser",
			mock.Anything,
			userId,
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

	userId := uuid.New()
	expectedSongs := &[]model.Song{
		{Title: "Some Song"},
		{Title: "Some other Song"},
	}

	jwtService.On("GetUserIdFromJwt", token).Return(userId, nil).Once()

	songRepository.
		On(
			"GetAllByUser",
			mock.IsType(expectedSongs),
			userId,
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
			userId,
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

	userId := uuid.New()
	expectedSongs := &[]model.Song{
		{Title: "Some Song"},
		{Title: "Some other Song"},
	}
	expectedTotalCount := &[]int64{20}[0]

	jwtService.On("GetUserIdFromJwt", token).Return(userId, nil).Once()

	songRepository.
		On(
			"GetAllByUser",
			mock.IsType(expectedSongs),
			userId,
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
			userId,
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
