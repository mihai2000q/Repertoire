package album

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/album"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAll_WhenGetUserIdFromJwtFails_ShouldReturnForbiddenError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	_uut := album.NewGetAllAlbums(nil, jwtService)

	request := requests.GetAlbumsRequest{}
	token := "This is a token"

	forbiddenError := wrapper.ForbiddenError(errors.New("forbidden"))
	jwtService.On("GetUserIdFromJwt", token).Return(uuid.Nil, forbiddenError).Once()

	// when
	res, errCode := _uut.Handle(request, token)

	// then
	assert.Empty(t, res)
	assert.NotNil(t, errCode)
	assert.Equal(t, forbiddenError, errCode)

	jwtService.AssertExpectations(t)
}

func TestGetAll_WhenGetAlbumsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	jwtService := new(service.JwtServiceMock)
	_uut := album.NewGetAllAlbums(albumRepository, jwtService)

	request := requests.GetAlbumsRequest{}
	token := "This is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	internalError := errors.New("internal error")
	albumRepository.
		On(
			"GetAllByUser",
			mock.Anything,
			userID,
			request.CurrentPage,
			request.PageSize,
			request.OrderBy,
			request.SearchBy,
		).
		Return(internalError).
		Once()

	// when
	res, errCode := _uut.Handle(request, token)

	// then
	assert.Empty(t, res)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	albumRepository.AssertExpectations(t)
	jwtService.AssertExpectations(t)
}

func TestGetAll_WhenGetAlbumsCountFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	jwtService := new(service.JwtServiceMock)
	_uut := album.NewGetAllAlbums(albumRepository, jwtService)

	request := requests.GetAlbumsRequest{}
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	expectedAlbums := &[]model.EnhancedAlbum{
		{Album: model.Album{Title: "Some Album"}},
	}

	albumRepository.
		On(
			"GetAllByUser",
			mock.IsType(expectedAlbums),
			userID,
			request.CurrentPage,
			request.PageSize,
			request.OrderBy,
			request.SearchBy,
		).
		Return(nil, expectedAlbums).
		Once()

	internalError := errors.New("internal error")
	albumRepository.
		On(
			"GetAllByUserCount",
			mock.Anything,
			userID,
			request.SearchBy,
		).
		Return(internalError).
		Once()

	// when
	res, errCode := _uut.Handle(request, token)

	// then
	assert.Equal(t, expectedAlbums, &res.Models)
	assert.Empty(t, res.TotalCount)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	albumRepository.AssertExpectations(t)
	jwtService.AssertExpectations(t)
}

func TestGetAll_WhenSuccessful_ShouldReturnAlbumsWithTotalCount(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	jwtService := new(service.JwtServiceMock)
	_uut := album.NewGetAllAlbums(albumRepository, jwtService)

	request := requests.GetAlbumsRequest{}
	token := "this is a token"

	expectedAlbums := &[]model.EnhancedAlbum{
		{Album: model.Album{Title: "Some Album"}},
		{Album: model.Album{Title: "Some Other Album"}},
	}
	expectedTotalCount := &[]int64{20}[0]

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	albumRepository.
		On(
			"GetAllByUser",
			mock.IsType(expectedAlbums),
			userID,
			request.CurrentPage,
			request.PageSize,
			request.OrderBy,
			request.SearchBy,
		).
		Return(nil, expectedAlbums).
		Once()

	albumRepository.
		On(
			"GetAllByUserCount",
			mock.IsType(expectedTotalCount),
			userID,
			request.SearchBy,
		).
		Return(nil, expectedTotalCount).
		Once()

	// when
	result, errCode := _uut.Handle(request, token)

	// then
	assert.Equal(t, expectedAlbums, &result.Models)
	assert.Equal(t, expectedTotalCount, &result.TotalCount)
	assert.Nil(t, errCode)

	albumRepository.AssertExpectations(t)
	jwtService.AssertExpectations(t)
}
