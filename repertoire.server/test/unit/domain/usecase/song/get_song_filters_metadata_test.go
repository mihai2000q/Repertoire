package song

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/song"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetSongFiltersMetadata_WhenGetUserIdFromJwtFails_ShouldReturnInternalServerErrorError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewGetSongFiltersMetadata(jwtService, songRepository)

	request := requests.GetSongFiltersMetadataRequest{}
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
	songRepository.AssertExpectations(t)
}

func TestGetSongFiltersMetadata_WhenGetFails_ShouldReturnInternalServerErrorError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewGetSongFiltersMetadata(jwtService, songRepository)

	request := requests.GetSongFiltersMetadataRequest{}
	token := "some token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	internalError := errors.New("some internal error")
	songRepository.On("GetFiltersMetadata", new(model.SongFiltersMetadata), userID, request.SearchBy).
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
	songRepository.AssertExpectations(t)
}

func TestGetSongFiltersMetadata_WhenSuccessful_ShouldReturnSongFiltersMetadata(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewGetSongFiltersMetadata(jwtService, songRepository)

	request := requests.GetSongFiltersMetadataRequest{}
	token := "some token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	metadata := model.SongFiltersMetadata{
		MinConfidence: float64(0),
	}
	songRepository.On("GetFiltersMetadata", new(model.SongFiltersMetadata), userID, request.SearchBy).
		Return(nil, &metadata).
		Once()

	// when
	result, errCode := _uut.Handle(request, token)

	// then
	assert.Nil(t, errCode)
	assert.Equal(t, metadata, result)

	jwtService.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}
