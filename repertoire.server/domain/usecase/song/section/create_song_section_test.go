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

func TestCreateSongSection_WhenCountSectionsBySongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := &CreateSongSection{
		songRepository: songRepository,
	}
	request := requests.CreateSongSectionRequest{
		SongID: uuid.New(),
		Name:   "Some Artist",
		TypeID: uuid.New(),
	}

	internalError := errors.New("internal error")
	songRepository.On("CountSectionsBySong", new(int64), request.SongID).
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

func TestCreateSongSection_WhenCreateSectionFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := &CreateSongSection{
		songRepository: songRepository,
	}
	request := requests.CreateSongSectionRequest{
		SongID: uuid.New(),
		Name:   "Some Artist",
		TypeID: uuid.New(),
	}

	expectedCount := &[]int64{20}[0]
	songRepository.On("CountSectionsBySong", mock.IsType(expectedCount), request.SongID).
		Return(nil, expectedCount).
		Once()

	internalError := errors.New("internal error")
	songRepository.On("CreateSection", mock.IsType(new(model.SongSection))).
		Run(func(args mock.Arguments) {
			newSection := args.Get(0).(*model.SongSection)
			assertCreatedSongSection(t, request, *newSection, *expectedCount)
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

func TestCreateSongSection_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := &CreateSongSection{
		songRepository: songRepository,
	}
	request := requests.CreateSongSectionRequest{
		SongID: uuid.New(),
		Name:   "Some Artist",
		TypeID: uuid.New(),
	}

	expectedCount := &[]int64{20}[0]
	songRepository.On("CountSectionsBySong", mock.IsType(expectedCount), request.SongID).
		Return(nil, expectedCount).
		Once()

	songRepository.On("CreateSection", mock.IsType(new(model.SongSection))).
		Run(func(args mock.Arguments) {
			newSection := args.Get(0).(*model.SongSection)
			assertCreatedSongSection(t, request, *newSection, *expectedCount)
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)

	songRepository.AssertExpectations(t)
}

func assertCreatedSongSection(
	t *testing.T,
	request requests.CreateSongSectionRequest,
	section model.SongSection,
	count int64,
) {
	assert.NotEmpty(t, section.ID)
	assert.Equal(t, request.Name, section.Name)
	assert.Zero(t, section.Rehearsals)
	assert.Equal(t, uint(count), section.Order)
	assert.Equal(t, request.TypeID, section.SongSectionTypeID)
	assert.Equal(t, request.SongID, section.SongID)
}
