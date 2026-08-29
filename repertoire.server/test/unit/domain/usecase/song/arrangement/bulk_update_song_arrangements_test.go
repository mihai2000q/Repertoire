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

func TestBulkUpdateSongArrangements_WhenGetFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songArrangementRepository := new(repository.SongArrangementRepositoryMock)
	_uut := arrangement.NewBulkUpdateSongArrangements(songArrangementRepository)

	request := requests.BulkUpdateSongArrangementsRequest{
		Requests: []requests.UpdateSongArrangementRequest{
			{
				ID:   uuid.New(),
				Name: "Some Arrangement",
			},
		},
	}

	internalError := errors.New("internal error")
	songArrangementRepository.
		On("GetAllBySongWithPartOccurrences",
			new([]model.SongArrangement),
			getIds(request),
			request.SongID,
		).
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

func TestBulkUpdateSongArrangements_WhenSongArrangementIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songArrangementRepository := new(repository.SongArrangementRepositoryMock)
	_uut := arrangement.NewBulkUpdateSongArrangements(songArrangementRepository)

	request := requests.BulkUpdateSongArrangementsRequest{
		Requests: []requests.UpdateSongArrangementRequest{
			{
				ID:   uuid.New(),
				Name: "Some Arrangement",
			},
		},
	}

	songArrangementRepository.
		On("GetAllBySongWithPartOccurrences",
			new([]model.SongArrangement),
			getIds(request),
			request.SongID,
		).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song arrangements not found", errCode.Error.Error())

	songArrangementRepository.AssertExpectations(t)
}

func TestBulkUpdateSongArrangements_WhenUpdateFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songArrangementRepository := new(repository.SongArrangementRepositoryMock)
	_uut := arrangement.NewBulkUpdateSongArrangements(songArrangementRepository)

	request := requests.BulkUpdateSongArrangementsRequest{
		Requests: []requests.UpdateSongArrangementRequest{
			{
				ID:   uuid.New(),
				Name: "Some Arrangement",
			},
		},
	}

	mockArrangements := []model.SongArrangement{
		{ID: request.Requests[0].ID},
	}
	songArrangementRepository.
		On("GetAllBySongWithPartOccurrences",
			new([]model.SongArrangement),
			getIds(request),
			request.SongID,
		).
		Return(nil, &mockArrangements).
		Once()

	internalError := errors.New("internal error")
	songArrangementRepository.On("UpdateAllWithAssociations", mock.IsType(new([]model.SongArrangement))).
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

func TestBulkUpdateSongArrangements_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	// given
	songArrangementRepository := new(repository.SongArrangementRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := arrangement.NewBulkUpdateSongArrangements(songArrangementRepository)

	request := requests.BulkUpdateSongArrangementsRequest{
		Requests: []requests.UpdateSongArrangementRequest{
			{
				ID:   uuid.New(),
				Name: "Some Arrangement 1",
				Occurrences: []requests.UpdateSongPartOccurrencesRequest{
					{PartID: uuid.New(), Occurrences: 5},
					{PartID: uuid.New(), Occurrences: 0},
					{PartID: uuid.New(), Occurrences: 1},
				},
			},
			{
				ID:   uuid.New(),
				Name: "Some Arrangement 2",
				Occurrences: []requests.UpdateSongPartOccurrencesRequest{
					{PartID: uuid.New(), Occurrences: 2},
					{PartID: uuid.New(), Occurrences: 1},
					{PartID: uuid.New(), Occurrences: 3},
				},
			},
		},
	}

	mockArrangements := []model.SongArrangement{
		{
			ID: request.Requests[0].ID,
			PartOccurrences: []model.SongPartOccurrences{
				{PartID: request.Requests[0].Occurrences[1].PartID},
				{PartID: request.Requests[0].Occurrences[0].PartID},
				{PartID: request.Requests[0].Occurrences[2].PartID},
			},
		},
		{
			ID: request.Requests[1].ID,
			PartOccurrences: []model.SongPartOccurrences{
				{PartID: request.Requests[1].Occurrences[0].PartID},
				{PartID: request.Requests[1].Occurrences[1].PartID},
				{PartID: request.Requests[1].Occurrences[2].PartID},
			},
		},
	}
	songArrangementRepository.
		On("GetAllBySongWithPartOccurrences",
			new([]model.SongArrangement),
			getIds(request),
			request.SongID,
		).
		Return(nil, &mockArrangements).
		Once()

	songArrangementRepository.On("UpdateAllWithAssociations", mock.IsType(new([]model.SongArrangement))).
		Run(func(args mock.Arguments) {
			newArrangements := *args.Get(0).(*[]model.SongArrangement)
			for i := range newArrangements {
				assertUpdatedSongArrangement(t, request.Requests[i], newArrangements[i])
			}
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

func getIds(request requests.BulkUpdateSongArrangementsRequest) []uuid.UUID {
	var ids []uuid.UUID
	for _, r := range request.Requests {
		ids = append(ids, r.ID)
	}
	return ids
}

func assertUpdatedSongArrangement(
	t *testing.T,
	request requests.UpdateSongArrangementRequest,
	arrangement model.SongArrangement,
) {
	assert.Equal(t, request.Name, arrangement.Name)

	sectionsOccurrencesMap := make(map[uuid.UUID]uint)
	for _, s := range request.Occurrences {
		sectionsOccurrencesMap[s.PartID] = s.Occurrences
	}
	for i := range arrangement.PartOccurrences {
		occurrences := sectionsOccurrencesMap[arrangement.PartOccurrences[i].PartID]
		assert.Equal(t, arrangement.PartOccurrences[i].Occurrences, occurrences)
	}
}
