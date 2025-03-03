package artist

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/artist"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateArtist_WhenGetUserIdFromJwtFails_ShouldReturnForbiddenError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	_uut := artist.NewCreateArtist(jwtService, nil, nil)

	request := requests.CreateArtistRequest{
		Name: "Some Artist",
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

func TestCreateArtist_WhenGetArtistFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	jwtService := new(service.JwtServiceMock)
	_uut := artist.NewCreateArtist(jwtService, artistRepository, nil)

	request := requests.CreateArtistRequest{
		Name: "Some Artist",
	}
	token := "this is a token"
	userID := uuid.New()

	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()
	internalError := errors.New("internal error")
	artistRepository.On("Create", mock.IsType(new(model.Artist))).
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
	artistRepository.AssertExpectations(t)
}

func TestCreateArtist_WhenAddToSearchEngineFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	jwtService := new(service.JwtServiceMock)
	searchEngineService := new(service.SearchEngineServiceMock)
	_uut := artist.NewCreateArtist(jwtService, artistRepository, searchEngineService)

	request := requests.CreateArtistRequest{
		Name: "Some Artist",
	}
	token := "this is a token"
	userID := uuid.New()

	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	artistRepository.On("Create", mock.IsType(new(model.Artist))).
		Return(nil).
		Once()

	internalError := wrapper.InternalServerError(errors.New("internal error"))
	searchEngineService.On("Add", mock.IsType([]any{})).
		Return(internalError).
		Once()

	// when
	id, errCode := _uut.Handle(request, token)

	// then
	assert.Empty(t, id)
	assert.NotNil(t, errCode)
	assert.Equal(t, internalError, errCode)

	jwtService.AssertExpectations(t)
	artistRepository.AssertExpectations(t)
	searchEngineService.AssertExpectations(t)
}

func TestCreateArtist_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	jwtService := new(service.JwtServiceMock)
	searchEngineService := new(service.SearchEngineServiceMock)
	_uut := artist.NewCreateArtist(jwtService, artistRepository, searchEngineService)

	request := requests.CreateArtistRequest{
		Name:   "Some Artist",
		IsBand: true,
	}
	token := "this is a token"
	userID := uuid.New()

	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	var artistID uuid.UUID
	artistRepository.On("Create", mock.IsType(new(model.Artist))).
		Run(func(args mock.Arguments) {
			newArtist := args.Get(0).(*model.Artist)
			assertCreatedArtist(t, *newArtist, request, userID)
			artistID = newArtist.ID
		}).
		Return(nil).
		Once()

	searchEngineService.On("Add", mock.IsType([]any{})).
		Run(func(args mock.Arguments) {
			searches := args.Get(0).([]any)
			assert.Len(t, searches, 1)
			assert.Contains(t, searches[0].(model.ArtistSearch).ID, artistID.String())
		}).
		Return(nil).
		Once()

	// when
	id, errCode := _uut.Handle(request, token)

	// then
	assert.Equal(t, artistID, id)
	assert.Nil(t, errCode)

	jwtService.AssertExpectations(t)
	artistRepository.AssertExpectations(t)
	searchEngineService.AssertExpectations(t)
}

func assertCreatedArtist(
	t *testing.T,
	artist model.Artist,
	request requests.CreateArtistRequest,
	userID uuid.UUID,
) {
	assert.Equal(t, request.Name, artist.Name)
	assert.Equal(t, request.IsBand, artist.IsBand)
	assert.Nil(t, artist.ImageURL)
	assert.Equal(t, userID, artist.UserID)
}
