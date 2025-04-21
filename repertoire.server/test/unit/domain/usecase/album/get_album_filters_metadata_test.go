package album

import (
	"errors"
	"net/http"
	"repertoire/server/domain/usecase/album"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetAlbumFiltersMetadata_WhenGetUserIdFromJwtFails_ShouldReturnInternalServerErrorError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := album.NewGetAlbumFiltersMetadata(jwtService, albumRepository)

	token := "some token"

	internalError := wrapper.InternalServerError(errors.New("some internal error"))
	jwtService.On("GetUserIdFromJwt", token).Return(uuid.Nil, internalError).Once()

	// when
	result, errCode := _uut.Handle(token)

	// then
	assert.Empty(t, result)
	assert.NotNil(t, errCode)
	assert.Equal(t, internalError, errCode)

	jwtService.AssertExpectations(t)
	albumRepository.AssertExpectations(t)
}

func TestGetAlbumFiltersMetadata_WhenGetFails_ShouldReturnInternalServerErrorError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := album.NewGetAlbumFiltersMetadata(jwtService, albumRepository)

	token := "some token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	internalError := errors.New("some internal error")
	albumRepository.On("GetFiltersMetadata", new(model.AlbumFiltersMetadata), userID).
		Return(internalError).
		Once()

	// when
	result, errCode := _uut.Handle(token)

	// then
	assert.Empty(t, result)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	jwtService.AssertExpectations(t)
	albumRepository.AssertExpectations(t)
}

func TestGetAlbumFiltersMetadata_WhenSuccessful_ShouldReturnAlbumFiltersMetadata(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	albumRepository := new(repository.AlbumRepositoryMock)
	_uut := album.NewGetAlbumFiltersMetadata(jwtService, albumRepository)

	token := "some token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	metadata := model.AlbumFiltersMetadata{
		MinConfidence: float64(0),
	}
	albumRepository.On("GetFiltersMetadata", new(model.AlbumFiltersMetadata), userID).
		Return(nil, &metadata).
		Once()

	// when
	result, errCode := _uut.Handle(token)

	// then
	assert.Nil(t, errCode)
	assert.Equal(t, metadata, result)

	jwtService.AssertExpectations(t)
	albumRepository.AssertExpectations(t)
}
