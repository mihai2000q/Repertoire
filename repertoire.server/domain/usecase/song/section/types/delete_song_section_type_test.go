package types

import (
	"errors"
	"net/http"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"slices"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteSongSectionType_WhenGetUserIdFromJwtFails_ShouldReturnError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	_uut := &DeleteSongSectionType{jwtService: jwtService}

	id := uuid.New()
	token := "this is a token"

	forbiddenError := wrapper.ForbiddenError(errors.New("forbidden error"))
	jwtService.On("GetUserIdFromJwt", token).Return(uuid.Nil, forbiddenError).Once()

	// when
	errCode := _uut.Handle(id, token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, forbiddenError, errCode)

	jwtService.AssertExpectations(t)
}

func TestDeleteSongSectionType_WhenGetSectionTypesFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := &DeleteSongSectionType{
		repository: songRepository,
		jwtService: jwtService,
	}

	id := uuid.New()
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	internalError := errors.New("internal error")
	songRepository.On("GetSectionTypes", new([]model.SongSectionType), userID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id, token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	jwtService.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}

func TestDeleteSongSectionType_WhenUpdateAllSectionTypesFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := &DeleteSongSectionType{
		repository: songRepository,
		jwtService: jwtService,
	}

	id := uuid.New()
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	types := &[]model.SongSectionType{
		{ID: id},
	}
	songRepository.On("GetSectionTypes", new([]model.SongSectionType), userID).
		Return(nil, types).
		Once()

	internalError := errors.New("internal error")
	songRepository.On("UpdateSectionTypes", mock.IsType(types)).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id, token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	jwtService.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}

func TestDeleteSongSectionType_WhenDeleteSectionTypeFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := &DeleteSongSectionType{
		repository: songRepository,
		jwtService: jwtService,
	}

	id := uuid.New()
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	types := &[]model.SongSectionType{
		{ID: id},
	}
	songRepository.On("GetSectionTypes", new([]model.SongSectionType), userID).
		Return(nil, types).
		Once()

	songRepository.On("UpdateAllSectionTypes", mock.IsType(types)).
		Return(nil).
		Once()

	internalError := errors.New("internal error")
	songRepository.On("DeleteSectionType", id).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id, token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	jwtService.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}

func TestDeleteSongSectionType_WhenSuccessful_ShouldReturnGuitarTunings(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := &DeleteSongSectionType{
		repository: songRepository,
		jwtService: jwtService,
	}

	id := uuid.New()
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	types := &[]model.SongSectionType{
		{ID: id},
	}
	songRepository.On("GetSectionTypes", new([]model.SongSectionType), userID).
		Return(nil, types).
		Once()

	songRepository.On("UpdateAllSectionTypes", mock.IsType(types)).
		Run(func(args mock.Arguments) {
			newSectionTypes := args.Get(0).(*[]model.SongSectionType)
			sectionTypes := slices.DeleteFunc(*newSectionTypes, func(s model.SongSectionType) bool {
				return s.ID == id
			})
			for i, sectionType := range sectionTypes {
				assert.Equal(t, i, sectionType.Order)
			}
		}).
		Return(nil).
		Once()

	songRepository.On("DeleteSectionType", id).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(id, token)

	// then
	assert.Nil(t, errCode)

	jwtService.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}
