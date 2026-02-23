package arrangement

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/song/arrangement"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateSongArrangement_WhenGetFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songArrangementRepository := new(repository.SongArrangementRepositoryMock)
	_uut := arrangement.NewUpdateSongArrangement(songArrangementRepository)

	request := requests.UpdateSongArrangementRequest{
		ID:   uuid.New(),
		Name: "Some Arrangement",
	}

	internalError := errors.New("internal error")
	songArrangementRepository.On("GetWithAssociations", new(model.SongArrangement), request.ID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songArrangementRepository.AssertExpectations(t)
}

func TestUpdateSongArrangement_WhenSongArrangementIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songArrangementRepository := new(repository.SongArrangementRepositoryMock)
	_uut := arrangement.NewUpdateSongArrangement(songArrangementRepository)

	request := requests.UpdateSongArrangementRequest{
		ID:   uuid.New(),
		Name: "Some Arrangement",
	}

	songArrangementRepository.On("GetWithAssociations", new(model.SongArrangement), request.ID).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song arrangement not found", errCode.Error.Error())

	songArrangementRepository.AssertExpectations(t)
}

func TestUpdateSongArrangement_WhenUpdateFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songArrangementRepository := new(repository.SongArrangementRepositoryMock)
	_uut := arrangement.NewUpdateSongArrangement(songArrangementRepository)

	request := requests.UpdateSongArrangementRequest{
		ID:   uuid.New(),
		Name: "Some Arrangement",
	}

	mockArrangement := model.SongArrangement{
		ID: request.ID,
	}
	songArrangementRepository.On("GetWithAssociations", new(model.SongArrangement), request.ID).
		Return(nil, &mockArrangement).
		Once()

	internalError := errors.New("internal error")
	songArrangementRepository.On("UpdateWithAssociations", mock.IsType(new(model.SongArrangement))).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songArrangementRepository.AssertExpectations(t)
}

func TestUpdateSongArrangement_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	// given
	songArrangementRepository := new(repository.SongArrangementRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := arrangement.NewUpdateSongArrangement(songArrangementRepository)

	request := requests.UpdateSongArrangementRequest{
		ID:   uuid.New(),
		Name: "Some Arrangement",
		Occurrences: []requests.UpdateSongSectionOccurrencesRequest{
			{SectionID: uuid.New(), Occurrences: 2},
			{SectionID: uuid.New(), Occurrences: 0},
			{SectionID: uuid.New(), Occurrences: 3},
		},
	}

	mockArrangement := model.SongArrangement{
		ID: request.ID,
		SectionOccurrences: []model.SongSectionOccurrences{
			{SectionID: request.Occurrences[1].SectionID},
			{SectionID: request.Occurrences[0].SectionID},
			{SectionID: request.Occurrences[2].SectionID},
		},
	}
	songArrangementRepository.On("GetWithAssociations", new(model.SongArrangement), request.ID).
		Return(nil, &mockArrangement).
		Once()

	songArrangementRepository.On("UpdateWithAssociations", mock.IsType(new(model.SongArrangement))).
		Run(func(args mock.Arguments) {
			newArrangement := *args.Get(0).(*model.SongArrangement)
			assertUpdatedSongArrangement(t, request, newArrangement)
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)

	songArrangementRepository.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}

func assertUpdatedSongArrangement(
	t *testing.T,
	request requests.UpdateSongArrangementRequest,
	arrangement model.SongArrangement,
) {
	assert.Equal(t, request.Name, arrangement.Name)

	sectionsOccurrencesMap := make(map[uuid.UUID]uint)
	for _, s := range request.Occurrences {
		sectionsOccurrencesMap[s.SectionID] = s.Occurrences
	}
	for i := range arrangement.SectionOccurrences {
		occurrences := sectionsOccurrencesMap[arrangement.SectionOccurrences[i].SectionID]
		assert.Equal(t, arrangement.SectionOccurrences[i].Occurrences, occurrences)
	}
}
