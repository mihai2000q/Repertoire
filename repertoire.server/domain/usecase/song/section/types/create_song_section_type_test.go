package types

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateSongSectionType_WhenGetUserIdFromJwtFails_ShouldReturnError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	_uut := &CreateSongSectionType{
		jwtService: jwtService,
	}

	request := requests.CreateSongSectionTypeRequest{
		Name: "New Type",
	}
	token := "this is a token"

	forbiddenError := wrapper.ForbiddenError(errors.New("forbidden error"))
	jwtService.On("GetUserIdFromJwt", token).Return(uuid.Nil, forbiddenError).Once()

	// when
	errCode := _uut.Handle(request, token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, forbiddenError, errCode)

	jwtService.AssertExpectations(t)
}

func TestCreateSongSectionType_WhenCountSectionTypesFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := &CreateSongSectionType{
		repository: songRepository,
		jwtService: jwtService,
	}

	request := requests.CreateSongSectionTypeRequest{
		Name: "New Type",
	}
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	internalError := errors.New("internal error")
	songRepository.On("CountSectionTypes", new(int64), userID).Return(internalError).Once()

	// when
	errCode := _uut.Handle(request, token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	jwtService.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}

func TestCreateSongSectionType_WhenCreateSectionTypeFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := &CreateSongSectionType{
		repository: songRepository,
		jwtService: jwtService,
	}

	request := requests.CreateSongSectionTypeRequest{
		Name: "New Type",
	}
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	songRepository.On("CountSectionTypes", new(int64), userID).Return(nil).Once()

	internalError := errors.New("internal error")
	songRepository.On("CreateSectionType", mock.IsType(new(model.SongSectionType))).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request, token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	jwtService.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}

func TestCreateSongSectionType_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := &CreateSongSectionType{
		repository: songRepository,
		jwtService: jwtService,
	}

	request := requests.CreateSongSectionTypeRequest{
		Name: "New Type",
	}
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	count := &[]int64{20}[0]
	songRepository.On("CountSectionTypes", mock.IsType(count), userID).
		Return(nil, count).
		Once()

	songRepository.On("CreateSectionType", mock.IsType(new(model.SongSectionType))).
		Run(func(args mock.Arguments) {
			newType := args.Get(0).(*model.SongSectionType)
			assertCreatedSongSectionType(t, *newType, request, userID, count)
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request, token)

	// then
	assert.Nil(t, errCode)

	jwtService.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}

func assertCreatedSongSectionType(
	t *testing.T,
	sectionType model.SongSectionType,
	request requests.CreateSongSectionTypeRequest,
	userID uuid.UUID,
	count *int64,
) {
	assert.NotEmpty(t, sectionType.ID)
	assert.Equal(t, request.Name, sectionType.Name)
	assert.Equal(t, uint(*count), sectionType.Order)
	assert.Equal(t, userID, sectionType.UserID)
}
