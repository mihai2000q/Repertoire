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

func TestCreateSongArrangement_WhenCountSectionsBySongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songArrangementRepository := new(repository.SongArrangementRepositoryMock)
	_uut := arrangement.NewCreateSongArrangement(songArrangementRepository, nil)

	request := requests.CreateSongArrangementRequest{
		SongID: uuid.New(),
		Name:   "Some Artist",
	}

	internalError := errors.New("internal error")
	songArrangementRepository.On("CountBySong", new(int64), request.SongID).
		Return(internalError).
		Once()

	// when
	id, errCode := _uut.Handle(request)

	// then
	assert.Empty(t, id)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songArrangementRepository.AssertExpectations(t)
}

func TestCreateSongArrangement_WhenGetSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songArrangementRepository := new(repository.SongArrangementRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := arrangement.NewCreateSongArrangement(songArrangementRepository, songRepository)

	request := requests.CreateSongArrangementRequest{
		SongID: uuid.New(),
		Name:   "Some Artist",
	}

	songArrangementRepository.On("CountBySong", new(int64), request.SongID).
		Return(nil).
		Once()

	internalError := errors.New("internal error")
	songRepository.On("GetWithParts", new(model.Song), request.SongID).
		Return(internalError).
		Once()

	// when
	id, errCode := _uut.Handle(request)

	// then
	assert.Empty(t, id)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songArrangementRepository.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}

func TestCreateSongArrangement_WhenSongIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songArrangementRepository := new(repository.SongArrangementRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := arrangement.NewCreateSongArrangement(songArrangementRepository, songRepository)

	request := requests.CreateSongArrangementRequest{
		SongID: uuid.New(),
		Name:   "Some Artist",
	}

	songArrangementRepository.On("CountBySong", new(int64), request.SongID).
		Return(nil).
		Once()

	songRepository.On("GetWithParts", new(model.Song), request.SongID).
		Return(nil).
		Once()

	// when
	id, errCode := _uut.Handle(request)

	// then
	assert.Empty(t, id)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song not found", errCode.Error.Error())

	songArrangementRepository.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}

func TestCreateSongArrangement_WhenCreateFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songArrangementRepository := new(repository.SongArrangementRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := arrangement.NewCreateSongArrangement(songArrangementRepository, songRepository)

	request := requests.CreateSongArrangementRequest{
		SongID: uuid.New(),
		Name:   "Some Artist",
	}

	songArrangementRepository.On("CountBySong", new(int64), request.SongID).
		Return(nil).
		Once()

	mockSong := model.Song{
		ID: request.SongID,
	}
	songRepository.On("GetWithParts", new(model.Song), request.SongID).
		Return(nil, &mockSong).
		Once()

	internalError := errors.New("internal error")
	songArrangementRepository.On("Create", mock.IsType(new(model.SongArrangement))).
		Return(internalError).
		Once()

	// when
	id, errCode := _uut.Handle(request)

	// then
	assert.Empty(t, id)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songArrangementRepository.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}

func TestCreateSongArrangement_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	// given
	songArrangementRepository := new(repository.SongArrangementRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := arrangement.NewCreateSongArrangement(songArrangementRepository, songRepository)

	request := requests.CreateSongArrangementRequest{
		SongID: uuid.New(),
		Name:   "Some Artist",
	}

	var arrangementsCount int64 = 10
	songArrangementRepository.On("CountBySong", new(int64), request.SongID).
		Return(nil, &arrangementsCount).
		Once()

	mockSong := model.Song{
		ID: request.SongID,
		Parts: []model.SongPart{
			{ID: uuid.New()},
			{ID: uuid.New()},
		},
	}
	songRepository.On("GetWithParts", new(model.Song), request.SongID).
		Return(nil, &mockSong).
		Once()

	var newId uuid.UUID
	songArrangementRepository.On("Create", mock.IsType(new(model.SongArrangement))).
		Run(func(args mock.Arguments) {
			newArrangement := *args.Get(0).(*model.SongArrangement)
			newId = newArrangement.ID
			assertCreatedSongArrangement(t, request, newArrangement, arrangementsCount, mockSong.Parts)
		}).
		Return(nil).
		Once()

	// when
	id, errCode := _uut.Handle(request)

	// then
	assert.Equal(t, id, newId)
	assert.Nil(t, errCode)

	songArrangementRepository.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}

func assertCreatedSongArrangement(
	t *testing.T,
	request requests.CreateSongArrangementRequest,
	arrangement model.SongArrangement,
	arrangementsCount int64,
	songParts []model.SongPart,
) {
	assert.NotEmpty(t, arrangement.ID)
	assert.Equal(t, request.Name, arrangement.Name)
	assert.Equal(t, request.SongID, arrangement.SongID)
	assert.Equal(t, uint(arrangementsCount), arrangement.Order)
	for i, part := range songParts {
		assert.Equal(t, part.ID, arrangement.PartOccurrences[i].PartID)
		assert.Equal(t, arrangement.ID, arrangement.PartOccurrences[i].ArrangementID)
		assert.Zero(t, arrangement.PartOccurrences[i].Occurrences)
	}
}
