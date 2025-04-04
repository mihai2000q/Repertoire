package playlist

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/playlist"
	"repertoire/server/internal/message/topics"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreatePlaylist_WhenGetUserIdFromJwtFails_ShouldReturnForbiddenError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	_uut := playlist.NewCreatePlaylist(jwtService, nil, nil)

	request := requests.CreatePlaylistRequest{
		Title: "Some Playlist",
	}
	token := "this is a token"

	forbiddenError := wrapper.ForbiddenError(errors.New("forbidden"))
	jwtService.On("GetUserIdFromJwt", token).Return(uuid.Nil, forbiddenError).Once()

	// when
	id, errCode := _uut.Handle(request, token)

	// then
	assert.Empty(t, id)
	assert.NotNil(t, errCode)
	assert.Equal(t, forbiddenError, errCode)

	jwtService.AssertExpectations(t)
}

func TestCreatePlaylist_WhenCreateFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	jwtService := new(service.JwtServiceMock)
	_uut := playlist.NewCreatePlaylist(jwtService, playlistRepository, nil)

	request := requests.CreatePlaylistRequest{
		Title: "Some Playlist",
	}
	token := "this is a token"
	userID := uuid.New()

	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()
	internalError := errors.New("internal error")
	playlistRepository.On("Create", mock.IsType(new(model.Playlist))).
		Return(internalError).
		Once()

	// when
	id, errCode := _uut.Handle(request, token)

	// then
	assert.Empty(t, id)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	jwtService.AssertExpectations(t)
	playlistRepository.AssertExpectations(t)
}

func TestCreatePlaylist_WhenPublishFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	playlistRepository := new(repository.PlaylistRepositoryMock)
	jwtService := new(service.JwtServiceMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := playlist.NewCreatePlaylist(jwtService, playlistRepository, messagePublisherService)

	request := requests.CreatePlaylistRequest{
		Title: "Some Playlist",
	}
	token := "this is a token"
	userID := uuid.New()

	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	playlistRepository.On("Create", mock.IsType(new(model.Playlist))).
		Return(nil).
		Once()

	internalError := errors.New("internal error")
	messagePublisherService.On("Publish", topics.PlaylistCreatedTopic, mock.IsType(model.Playlist{})).
		Return(internalError).
		Once()

	// when
	id, errCode := _uut.Handle(request, token)

	// then
	assert.Empty(t, id)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	jwtService.AssertExpectations(t)
	playlistRepository.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
}

func TestCreatePlaylist_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	playlistRepository := new(repository.PlaylistRepositoryMock)
	messagePublisherService := new(service.MessagePublisherServiceMock)
	_uut := playlist.NewCreatePlaylist(jwtService, playlistRepository, messagePublisherService)

	request := requests.CreatePlaylistRequest{
		Title: "Some Playlist",
	}
	token := "this is a token"
	userID := uuid.New()

	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	var createdPlaylist model.Playlist
	playlistRepository.On("Create", mock.IsType(new(model.Playlist))).
		Run(func(args mock.Arguments) {
			newPlaylist := args.Get(0).(*model.Playlist)
			assertCreatedPlaylist(t, *newPlaylist, request, userID)
			createdPlaylist = *newPlaylist
		}).
		Return(nil).
		Once()

	messagePublisherService.On("Publish", topics.PlaylistCreatedTopic, mock.IsType(model.Playlist{})).
		Run(func(args mock.Arguments) {
			assert.Equal(t, createdPlaylist, args.Get(1).(model.Playlist))
		}).
		Return(nil).
		Once()

	// when
	id, errCode := _uut.Handle(request, token)

	// then
	assert.Equal(t, createdPlaylist.ID, id)
	assert.Nil(t, errCode)

	jwtService.AssertExpectations(t)
	playlistRepository.AssertExpectations(t)
	messagePublisherService.AssertExpectations(t)
}

func assertCreatedPlaylist(
	t *testing.T,
	playlist model.Playlist,
	request requests.CreatePlaylistRequest,
	userID uuid.UUID,
) {
	assert.Equal(t, request.Title, playlist.Title)
	assert.Equal(t, request.Description, playlist.Description)
	assert.Equal(t, userID, playlist.UserID)
	assert.Nil(t, playlist.ImageURL)
}
