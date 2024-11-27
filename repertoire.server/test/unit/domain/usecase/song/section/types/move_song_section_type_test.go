package types

import (
	"cmp"
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	types2 "repertoire/server/domain/usecase/song/section/types"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/data/service"
	"slices"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMoveSongSectionType_WhenGetUserIdFromJwtFails_ShouldReturnError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	_uut := &types2.MoveSongSectionType{jwtService: jwtService}

	request := requests.MoveSongSectionTypeRequest{
		ID:     uuid.New(),
		OverID: uuid.New(),
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

func TestMoveSongSectionType_WhenGetSectionTypesFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := &types2.MoveSongSectionType{
		repository: songRepository,
		jwtService: jwtService,
	}

	request := requests.MoveSongSectionTypeRequest{
		ID:     uuid.New(),
		OverID: uuid.New(),
	}
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	internalError := errors.New("internal error")
	songRepository.On("GetSectionTypes", new([]model.SongSectionType), userID).
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

func TestMoveSongSectionType_WhenSectionTypeIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := &types2.MoveSongSectionType{
		repository: songRepository,
		jwtService: jwtService,
	}

	request := requests.MoveSongSectionTypeRequest{
		ID:     uuid.New(),
		OverID: uuid.New(),
	}
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	tunings := &[]model.SongSectionType{}
	songRepository.On("GetSectionTypes", new([]model.SongSectionType), userID).
		Return(nil, tunings).
		Once()

	// when
	errCode := _uut.Handle(request, token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "type not found", errCode.Error.Error())

	jwtService.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}

func TestMoveSongSectionType_WhenOverSectionTypeIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := &types2.MoveSongSectionType{
		repository: songRepository,
		jwtService: jwtService,
	}

	request := requests.MoveSongSectionTypeRequest{
		ID:     uuid.New(),
		OverID: uuid.New(),
	}
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	tunings := &[]model.SongSectionType{
		{ID: request.ID},
	}
	songRepository.On("GetSectionTypes", new([]model.SongSectionType), userID).
		Return(nil, tunings).
		Once()

	// when
	errCode := _uut.Handle(request, token)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "over type not found", errCode.Error.Error())

	jwtService.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}

func TestMoveSongSectionType_WhenUpdateAllSectionTypesFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	jwtService := new(service.JwtServiceMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := &types2.MoveSongSectionType{
		repository: songRepository,
		jwtService: jwtService,
	}

	request := requests.MoveSongSectionTypeRequest{
		ID:     uuid.New(),
		OverID: uuid.New(),
	}
	token := "this is a token"

	userID := uuid.New()
	jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

	tunings := &[]model.SongSectionType{
		{ID: request.ID},
		{ID: request.OverID},
	}
	songRepository.On("GetSectionTypes", new([]model.SongSectionType), userID).
		Return(nil, tunings).
		Once()

	internalError := errors.New("internal error")
	songRepository.On("UpdateAllSectionTypes", mock.IsType(tunings)).
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

func TestMoveSongSectionType_WhenSuccessful_ShouldReturnSongSectionTypes(t *testing.T) {
	tests := []struct {
		name      string
		types     *[]model.SongSectionType
		index     uint
		overIndex uint
	}{
		{
			"Use case 1",
			&[]model.SongSectionType{
				{ID: uuid.New(), Order: 0},
				{ID: uuid.New(), Order: 1},
				{ID: uuid.New(), Order: 2},
				{ID: uuid.New(), Order: 3},
				{ID: uuid.New(), Order: 4},
			},
			1,
			3,
		},
		{
			"Use case 2",
			&[]model.SongSectionType{
				{ID: uuid.New(), Order: 0},
				{ID: uuid.New(), Order: 1},
				{ID: uuid.New(), Order: 2},
				{ID: uuid.New(), Order: 3},
				{ID: uuid.New(), Order: 4},
			},
			3,
			1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			jwtService := new(service.JwtServiceMock)
			songRepository := new(repository.SongRepositoryMock)
			_uut := &types2.MoveSongSectionType{
				repository: songRepository,
				jwtService: jwtService,
			}

			request := requests.MoveSongSectionTypeRequest{
				ID:     (*tt.types)[tt.index].ID,
				OverID: (*tt.types)[tt.overIndex].ID,
			}
			token := "this is a token"

			userID := uuid.New()
			jwtService.On("GetUserIdFromJwt", token).Return(userID, nil).Once()

			songRepository.On("GetSectionTypes", new([]model.SongSectionType), userID).
				Return(nil, tt.types).
				Once()

			songRepository.On("UpdateAllSectionTypes", mock.IsType(tt.types)).
				Run(func(args mock.Arguments) {
					newSongSectionTypes := args.Get(0).(*[]model.SongSectionType)
					sectionTypes := slices.Clone(*newSongSectionTypes)
					slices.SortFunc(sectionTypes, func(a, b model.SongSectionType) int {
						return cmp.Compare(a.Order, b.Order)
					})
					if tt.index < tt.overIndex {
						assert.Equal(t, sectionTypes[tt.overIndex-1].ID, request.OverID)
					} else if tt.index > tt.overIndex {
						assert.Equal(t, sectionTypes[tt.overIndex+1].ID, request.OverID)
					}
					for i, sectionType := range sectionTypes {
						assert.Equal(t, uint(i), sectionType.Order)
					}
				}).
				Return(nil).
				Once()

			// when
			errCode := _uut.Handle(request, token)

			// then
			assert.Nil(t, errCode)

			jwtService.AssertExpectations(t)
			songRepository.AssertExpectations(t)
		})
	}
}
