package playlist

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

func TestCreatePlaylist_WhenGetUserIdFromJwtFails_ShouldReturnForbiddenError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	_uut := &CreatePlaylist{
		jwtService: jwtService,
	}
	request := requests.CreatePlaylistRequest{
		Title: "Some Playlist",
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

func TestCreatePlaylist_WhenGetPlaylistFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	jwtService := new(service.JwtServiceMock)
	_uut := &CreatePlaylist{
		repository: playlistRepository,
		jwtService: jwtService,
	}
	request := requests.CreatePlaylistRequest{
		Title: "Some Playlist",
	}
	token := "this is a token"
	userID := uuid.New()

	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()
	internalError := errors.New("internal error")
	playlistRepository.On("Create", mock.IsType(new(model.Playlist))).
		Run(func(args mock.Arguments) {
			newPlaylist := args.Get(0).(*model.Playlist)
			assert.Equal(t, request.Title, newPlaylist.Title)
			assert.Equal(t, request.Description, newPlaylist.Description)
			assert.Equal(t, userID, newPlaylist.UserID)
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
	playlistRepository.AssertExpectations(t)
}

func TestCreatePlaylist_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	jwtService := new(service.JwtServiceMock)
	_uut := &CreatePlaylist{
		repository: playlistRepository,
		jwtService: jwtService,
	}
	request := requests.CreatePlaylistRequest{
		Title: "Some Playlist",
	}
	token := "this is a token"
	userID := uuid.New()

	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()
	playlistRepository.On("Create", mock.IsType(new(model.Playlist))).
		Run(func(args mock.Arguments) {
			newPlaylist := args.Get(0).(*model.Playlist)
			assert.Equal(t, request.Title, newPlaylist.Title)
			assert.Equal(t, request.Description, newPlaylist.Description)
			assert.Equal(t, userID, newPlaylist.UserID)
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request, token)

	// then
	assert.Nil(t, errCode)

	jwtService.AssertExpectations(t)
	playlistRepository.AssertExpectations(t)
}
