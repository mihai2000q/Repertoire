package album

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

func TestCreateAlbum_WhenGetUserIdFromJwtFails_ShouldReturnForbiddenError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	_uut := &CreateAlbum{
		jwtService: jwtService,
	}
	request := requests.CreateAlbumRequest{
		Title: "Some Album",
	}
	token := "this is a token"

	forbiddenError := wrapper.ForbiddenError(errors.New("forbidden"))
	jwtService.On("GetUserIdFromJwt", token).Return(uuid.Nil, forbiddenError).Once()

	// when
	errCode := _uut.Handle(request, token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, forbiddenError, errCode)

	jwtService.AssertExpectations(t)
}

func TestCreateAlbum_WhenGetAlbumFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	jwtService := new(service.JwtServiceMock)
	_uut := &CreateAlbum{
		repository: albumRepository,
		jwtService: jwtService,
	}
	request := requests.CreateAlbumRequest{
		Title: "Some Album",
	}
	token := "this is a token"
	userID := uuid.New()

	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()
	internalError := errors.New("internal error")
	albumRepository.On("Create", mock.IsType(new(models.Album))).
		Run(func(args mock.Arguments) {
			newAlbum := args.Get(0).(*models.Album)
			assert.Equal(t, request.Title, newAlbum.Title)
			assert.Equal(t, userID, newAlbum.UserID)
		}).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request, token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	jwtService.AssertExpectations(t)
	albumRepository.AssertExpectations(t)
}

func TestCreateAlbum_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	// given
	albumRepository := new(repository.AlbumRepositoryMock)
	jwtService := new(service.JwtServiceMock)
	_uut := &CreateAlbum{
		repository: albumRepository,
		jwtService: jwtService,
	}
	request := requests.CreateAlbumRequest{
		Title: "Some Album",
	}
	token := "this is a token"
	userID := uuid.New()

	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()
	albumRepository.On("Create", mock.IsType(new(models.Album))).
		Run(func(args mock.Arguments) {
			newAlbum := args.Get(0).(*models.Album)
			assert.Equal(t, request.Title, newAlbum.Title)
			assert.Equal(t, userID, newAlbum.UserID)
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request, token)

	// then
	assert.Nil(t, errCode)

	jwtService.AssertExpectations(t)
	albumRepository.AssertExpectations(t)
}
