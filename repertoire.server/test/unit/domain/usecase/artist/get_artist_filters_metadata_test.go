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
)

func TestGetArtistFiltersMetadata_WhenGetUserIdFromJwtFails_ShouldReturnInternalServerErrorError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := artist.NewGetArtistFiltersMetadata(jwtService, artistRepository)

	request := requests.GetArtistFiltersMetadataRequest{}
	token := "some token"

	internalError := wrapper.InternalServerError(errors.New("some internal error"))
	jwtService.On("GetUserIdFromJwt", token).Return(uuid.Nil, internalError).Once()

	// when
	result, errCode := _uut.Handle(request, token)

	// then
	assert.Empty(t, result)
	assert.NotNil(t, errCode)
	assert.Equal(t, internalError, errCode)

	jwtService.AssertExpectations(t)
	artistRepository.AssertExpectations(t)
}

func TestGetArtistFiltersMetadata_WhenGetFails_ShouldReturnInternalServerErrorError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := artist.NewGetArtistFiltersMetadata(jwtService, artistRepository)

	request := requests.GetArtistFiltersMetadataRequest{}
	token := "some token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	internalError := errors.New("some internal error")
	artistRepository.On("GetFiltersMetadata", new(model.ArtistFiltersMetadata), userID, request.SearchBy).
		Return(internalError).
		Once()

	// when
	result, errCode := _uut.Handle(request, token)

	// then
	assert.Empty(t, result)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	jwtService.AssertExpectations(t)
	artistRepository.AssertExpectations(t)
}

func TestGetArtistFiltersMetadata_WhenSuccessful_ShouldReturnArtistFiltersMetadata(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := artist.NewGetArtistFiltersMetadata(jwtService, artistRepository)

	request := requests.GetArtistFiltersMetadataRequest{}
	token := "some token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	metadata := model.ArtistFiltersMetadata{
		MinConfidence: float64(0),
	}
	artistRepository.On("GetFiltersMetadata", new(model.ArtistFiltersMetadata), userID, request.SearchBy).
		Return(nil, &metadata).
		Once()

	// when
	result, errCode := _uut.Handle(request, token)

	// then
	assert.Nil(t, errCode)
	assert.Equal(t, metadata, result)

	jwtService.AssertExpectations(t)
	artistRepository.AssertExpectations(t)
}
