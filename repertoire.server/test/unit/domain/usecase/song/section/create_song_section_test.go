package section

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/song/section"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateSongSection_WhenCountSectionsBySongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	_uut := section.NewCreateSongSection(songSectionRepository)

	request := requests.CreateSongSectionRequest{
		SongID: uuid.New(),
		Name:   "Some Artist",
		TypeID: uuid.New(),
	}

	internalError := errors.New("internal error")
	songSectionRepository.On("CountAllBySong", new(int64), request.SongID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songSectionRepository.AssertExpectations(t)
}

func TestCreateSongSection_WhenCreateSectionFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	_uut := section.NewCreateSongSection(songSectionRepository)

	request := requests.CreateSongSectionRequest{
		SongID: uuid.New(),
		Name:   "Some Artist",
		TypeID: uuid.New(),
	}

	expectedCount := &[]int64{20}[0]
	songSectionRepository.On("CountAllBySong", mock.IsType(expectedCount), request.SongID).
		Return(nil, expectedCount).
		Once()

	internalError := errors.New("internal error")
	songSectionRepository.On("Create", mock.IsType(new(model.SongSection))).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songSectionRepository.AssertExpectations(t)
}

func TestCreateSongSection_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	_uut := section.NewCreateSongSection(songSectionRepository)

	request := requests.CreateSongSectionRequest{
		SongID: uuid.New(),
		Name:   "Some Artist",
		TypeID: uuid.New(),
	}

	sectionsCount := int64(12)
	songSectionRepository.On("CountAllBySong", mock.IsType(&sectionsCount), request.SongID).
		Return(nil, &sectionsCount).
		Once()

	songSectionRepository.On("Create", mock.IsType(new(model.SongSection))).
		Run(func(args mock.Arguments) {
			newSection := args.Get(0).(*model.SongSection)
			assertCreatedSongSection(t, request, *newSection, sectionsCount)
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)

	songSectionRepository.AssertExpectations(t)
}

func assertCreatedSongSection(
	t *testing.T,
	request requests.CreateSongSectionRequest,
	section model.SongSection,
	sectionsCount int64,
) {
	assert.NotEmpty(t, section.ID)
	assert.Equal(t, request.Name, section.Name)
	assert.Equal(t, request.TypeID, section.SongSectionTypeID)
	assert.Equal(t, request.SongID, section.SongID)
	assert.Equal(t, uint(sectionsCount), section.Order)
}
