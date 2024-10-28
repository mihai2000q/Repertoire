package section

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"repertoire/api/requests"
	"repertoire/data/repository"
	"repertoire/model"
	"testing"
)

func TestUpdateSongSection_WhenGetSectionsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := &UpdateSongSection{
		songRepository: songRepository,
	}
	request := requests.UpdateSongSectionRequest{
		ID:         uuid.New(),
		Name:       "Some Artist",
		Rehearsals: 23,
		TypeID:     uuid.New(),
	}

	internalError := errors.New("internal error")
	songRepository.On("GetSection", new(model.SongSection), request.ID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestUpdateSongSection_WhenSectionsIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := &UpdateSongSection{
		songRepository: songRepository,
	}
	request := requests.UpdateSongSectionRequest{
		ID:         uuid.New(),
		Name:       "Some Artist",
		Rehearsals: 23,
		TypeID:     uuid.New(),
	}

	songRepository.On("GetSection", new(model.SongSection), request.ID).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song section not found", errCode.Error.Error())

	songRepository.AssertExpectations(t)
}

func TestUpdateSongSection_WhenUpdateSectionFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := &UpdateSongSection{
		songRepository: songRepository,
	}
	request := requests.UpdateSongSectionRequest{
		ID:         uuid.New(),
		Name:       "Some Artist",
		Rehearsals: 120,
		TypeID:     uuid.New(),
	}

	section := &model.SongSection{
		ID:   request.ID,
		Name: "Old name",
	}
	songRepository.On("GetSection", new(model.SongSection), request.ID).
		Return(nil, section).
		Once()

	internalError := errors.New("internal error")
	songRepository.On("UpdateSection", mock.IsType(new(model.SongSection))).
		Run(func(args mock.Arguments) {
			newSection := args.Get(0).(*model.SongSection)
			assertUpdatedSongSection(t, request, *newSection)
		}).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestUpdateSongSection_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := &UpdateSongSection{
		songRepository: songRepository,
	}
	request := requests.UpdateSongSectionRequest{
		ID:         uuid.New(),
		Name:       "Some Artist",
		Rehearsals: 120,
		TypeID:     uuid.New(),
	}

	section := &model.SongSection{
		ID:   request.ID,
		Name: "Old name",
	}

	songRepository.On("GetSection", new(model.SongSection), request.ID).
		Return(nil, section).
		Once()
	songRepository.On("UpdateSection", mock.IsType(new(model.SongSection))).
		Run(func(args mock.Arguments) {
			newSection := args.Get(0).(*model.SongSection)
			assertUpdatedSongSection(t, request, *newSection)
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)

	songRepository.AssertExpectations(t)
}

func assertUpdatedSongSection(t *testing.T, request requests.UpdateSongSectionRequest, section model.SongSection) {
	assert.Equal(t, request.Name, section.Name)
	assert.Equal(t, request.Rehearsals, section.Rehearsals)
	assert.Equal(t, request.TypeID, section.SongSectionTypeID)
}
