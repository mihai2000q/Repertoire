package section

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/song/section"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/domain/processor"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateSongSection_WhenGetSectionFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	_uut := section.NewUpdateSongSection(songSectionRepository, nil, nil)

	request := requests.UpdateSongSectionRequest{
		ID:     uuid.New(),
		Name:   "Some Artist",
		TypeID: uuid.New(),
	}

	internalError := errors.New("internal error")
	songSectionRepository.On("Get", new(model.SongSection), request.ID).
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

func TestUpdateSongSection_WhenSectionsIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	_uut := section.NewUpdateSongSection(songSectionRepository, nil, nil)

	request := requests.UpdateSongSectionRequest{
		ID:     uuid.New(),
		Name:   "Some Artist",
		TypeID: uuid.New(),
	}

	songSectionRepository.On("Get", new(model.SongSection), request.ID).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song section not found", errCode.Error.Error())

	songSectionRepository.AssertExpectations(t)
}

func TestUpdateSongSection_WhenUpdateSectionFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	_uut := section.NewUpdateSongSection(songSectionRepository, nil, nil)

	request := requests.UpdateSongSectionRequest{
		ID:     uuid.New(),
		Name:   "Some Artist",
		TypeID: uuid.New(),
	}

	// given - mocking
	mockSection := &model.SongSection{
		ID:   request.ID,
		Name: "Old name",
	}
	songSectionRepository.On("Get", new(model.SongSection), request.ID).
		Return(nil, mockSection).
		Once()

	internalError := errors.New("internal error")
	songSectionRepository.On("Update", mock.IsType(new(model.SongSection))).
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

func TestUpdateSongSection_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	progressProcessor := new(processor.ProgressProcessorMock)
	_uut := section.NewUpdateSongSection(songSectionRepository, songRepository, progressProcessor)

	request := requests.UpdateSongSectionRequest{
		ID:     uuid.New(),
		Name:   "Some Artist",
		TypeID: uuid.New(),
	}

	mockSection := &model.SongSection{
		ID:   request.ID,
		Name: "Old name",
	}
	songSectionRepository.On("Get", new(model.SongSection), request.ID).
		Return(nil, mockSection).
		Once()

	songSectionRepository.On("Update", mock.IsType(new(model.SongSection))).
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

	songSectionRepository.AssertExpectations(t)
	songRepository.AssertExpectations(t)
	progressProcessor.AssertExpectations(t)
}

func assertUpdatedSongSection(
	t *testing.T,
	request requests.UpdateSongSectionRequest,
	section model.SongSection,
) {
	assert.Equal(t, request.Name, section.Name)
	assert.Equal(t, request.TypeID, section.SongSectionTypeID)
}
