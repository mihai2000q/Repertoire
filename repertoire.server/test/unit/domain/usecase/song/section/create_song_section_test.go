package section

import (
	"errors"
	"math"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/song/section"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"slices"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateSongSection_WhenCountSectionsBySongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	_uut := section.NewCreateSongSection(songSectionRepository, nil, nil)

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

func TestCreateSongSection_WhenGetSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := section.NewCreateSongSection(songSectionRepository, songRepository, nil)

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
	songRepository.On("Get", new(model.Song), request.SongID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songSectionRepository.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}

func TestCreateSongSection_WhenSongIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := section.NewCreateSongSection(songSectionRepository, songRepository, nil)

	request := requests.CreateSongSectionRequest{
		SongID: uuid.New(),
		Name:   "Some Artist",
		TypeID: uuid.New(),
	}

	expectedCount := &[]int64{20}[0]
	songSectionRepository.On("CountAllBySong", mock.IsType(expectedCount), request.SongID).
		Return(nil, expectedCount).
		Once()

	songRepository.On("Get", new(model.Song), request.SongID).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song not found", errCode.Error.Error())

	songSectionRepository.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}

func TestCreateSongSection_WhenIsBandMemberAssociatedWithSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := section.NewCreateSongSection(songSectionRepository, songRepository, nil)

	request := requests.CreateSongSectionRequest{
		SongID:       uuid.New(),
		Name:         "Some Artist",
		TypeID:       uuid.New(),
		BandMemberID: &[]uuid.UUID{uuid.New()}[0],
	}

	expectedCount := &[]int64{20}[0]
	songSectionRepository.On("CountAllBySong", mock.IsType(expectedCount), request.SongID).
		Return(nil, expectedCount).
		Once()

	mockSong := &model.Song{ID: request.SongID}
	songRepository.On("Get", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()

	internalError := errors.New("internal error")
	songRepository.On("IsBandMemberAssociatedWithSong", request.SongID, *request.BandMemberID).
		Return(false, internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songSectionRepository.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}

func TestCreateSongSection_WhenBandMemberIsNotAssociatedWithTheSong_ShouldReturnConflictError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := section.NewCreateSongSection(songSectionRepository, songRepository, nil)

	request := requests.CreateSongSectionRequest{
		SongID:       uuid.New(),
		Name:         "Some Artist",
		TypeID:       uuid.New(),
		BandMemberID: &[]uuid.UUID{uuid.New()}[0],
	}

	expectedCount := &[]int64{20}[0]
	songSectionRepository.On("CountAllBySong", mock.IsType(expectedCount), request.SongID).
		Return(nil, expectedCount).
		Once()

	mockSong := &model.Song{ID: request.SongID}
	songRepository.On("Get", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()

	songRepository.On("IsBandMemberAssociatedWithSong", request.SongID, *request.BandMemberID).
		Return(false, nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusConflict, errCode.Code)
	assert.Equal(t, "band member is not part of the artist associated with this song", errCode.Error.Error())

	songSectionRepository.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}

func TestCreateSongSection_WhenCreateSectionFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := section.NewCreateSongSection(songSectionRepository, songRepository, nil)

	request := requests.CreateSongSectionRequest{
		SongID: uuid.New(),
		Name:   "Some Artist",
		TypeID: uuid.New(),
	}

	expectedCount := &[]int64{20}[0]
	songSectionRepository.On("CountAllBySong", mock.IsType(expectedCount), request.SongID).
		Return(nil, expectedCount).
		Once()

	mockSong := &model.Song{ID: uuid.New()}
	songRepository.On("Get", new(model.Song), request.SongID).
		Return(nil, mockSong).
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
	songRepository.AssertExpectations(t)
}

func TestCreateSongSection_WhenUpdateSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	_uut := section.NewCreateSongSection(songSectionRepository, songRepository, nil)

	request := requests.CreateSongSectionRequest{
		SongID: uuid.New(),
		Name:   "Some Artist",
		TypeID: uuid.New(),
	}

	expectedCount := &[]int64{20}[0]
	songSectionRepository.On("CountAllBySong", mock.IsType(expectedCount), request.SongID).
		Return(nil, expectedCount).
		Once()

	mockSong := &model.Song{ID: uuid.New()}
	songRepository.On("Get", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()

	songSectionRepository.On("Create", mock.IsType(new(model.SongSection))).
		Return(nil).
		Once()

	internalError := errors.New("internal error")
	songRepository.On("Update", mock.IsType(new(model.Song))).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songSectionRepository.AssertExpectations(t)
	songRepository.AssertExpectations(t)
}

func TestCreateSongSection_WhenGetArrangementsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	songArrangementRepository := new(repository.SongArrangementRepositoryMock)
	_uut := section.NewCreateSongSection(songSectionRepository, songRepository, songArrangementRepository)

	request := requests.CreateSongSectionRequest{
		SongID: uuid.New(),
		Name:   "Some Artist",
		TypeID: uuid.New(),
	}

	expectedCount := &[]int64{20}[0]
	songSectionRepository.On("CountAllBySong", mock.IsType(expectedCount), request.SongID).
		Return(nil, expectedCount).
		Once()

	mockSong := &model.Song{ID: uuid.New()}
	songRepository.On("Get", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()

	songSectionRepository.On("Create", mock.IsType(new(model.SongSection))).
		Return(nil).
		Once()

	songRepository.On("Update", mock.IsType(new(model.Song))).
		Return(nil).
		Once()

	internalError := errors.New("internal error")
	songArrangementRepository.On("GetAllBySong", mock.IsType(new([]model.SongArrangement)), request.SongID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songSectionRepository.AssertExpectations(t)
	songRepository.AssertExpectations(t)
	songArrangementRepository.AssertExpectations(t)
}

func TestCreateSongSection_WhenUpdateArrangementsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	songArrangementRepository := new(repository.SongArrangementRepositoryMock)
	_uut := section.NewCreateSongSection(songSectionRepository, songRepository, songArrangementRepository)

	request := requests.CreateSongSectionRequest{
		SongID: uuid.New(),
		Name:   "Some Artist",
		TypeID: uuid.New(),
	}

	expectedCount := &[]int64{20}[0]
	songSectionRepository.On("CountAllBySong", mock.IsType(expectedCount), request.SongID).
		Return(nil, expectedCount).
		Once()

	mockSong := &model.Song{ID: uuid.New()}
	songRepository.On("Get", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()

	songSectionRepository.On("Create", mock.IsType(new(model.SongSection))).
		Return(nil).
		Once()

	songRepository.On("Update", mock.IsType(new(model.Song))).
		Return(nil).
		Once()

	songArrangementRepository.On("GetAllBySong", mock.IsType(new([]model.SongArrangement)), request.SongID).
		Return(nil).
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

	songSectionRepository.AssertExpectations(t)
	songRepository.AssertExpectations(t)
	songArrangementRepository.AssertExpectations(t)
}

func TestCreateSongSection_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	tests := []struct {
		name                   string
		song                   model.Song
		expectedSectionsCount  int64
		expectedSongConfidence float64
		expectedSongRehearsals float64
		expectedSongProgress   float64
		bandMemberID           *uuid.UUID
	}{
		{
			"1 - When there are no precedent sections",
			model.Song{ID: uuid.New()},
			0,
			0,
			0,
			0,
			nil,
		},
		{
			"2 - When there are precedent sections, but with stats 0",
			model.Song{ID: uuid.New()},
			2,
			0,
			0,
			0,
			nil,
		},
		{
			"3 - When there are precedent sections with stats",
			model.Song{
				Confidence: 50,
				Rehearsals: 10,
				Progress:   55,
			},
			2,
			33,
			7,
			37,
			nil,
		},
		{
			"4 - With band member",
			model.Song{
				ID:         uuid.New(),
				Confidence: 0,
				Rehearsals: 0,
				Progress:   0,
			},
			2,
			0,
			0,
			0,
			&[]uuid.UUID{uuid.New()}[0],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			songSectionRepository := new(repository.SongSectionRepositoryMock)
			songRepository := new(repository.SongRepositoryMock)
			songArrangementRepository := new(repository.SongArrangementRepositoryMock)
			_uut := section.NewCreateSongSection(songSectionRepository, songRepository, songArrangementRepository)

			request := requests.CreateSongSectionRequest{
				SongID:       uuid.New(),
				Name:         "Some Artist",
				TypeID:       uuid.New(),
				BandMemberID: tt.bandMemberID,
			}

			songSectionRepository.On("CountAllBySong", mock.IsType(&tt.expectedSectionsCount), request.SongID).
				Return(nil, &tt.expectedSectionsCount).
				Once()

			songRepository.On("Get", new(model.Song), request.SongID).
				Return(nil, &tt.song).
				Once()

			if request.BandMemberID != nil {
				songRepository.On("IsBandMemberAssociatedWithSong", request.SongID, *request.BandMemberID).
					Return(true, nil).
					Once()
			}

			var newSectionID uuid.UUID
			songSectionRepository.On("Create", mock.IsType(new(model.SongSection))).
				Run(func(args mock.Arguments) {
					newSection := args.Get(0).(*model.SongSection)
					newSectionID = newSection.ID
					assertCreatedSongSection(t, request, *newSection, tt.expectedSectionsCount)
				}).
				Return(nil).
				Once()

			songRepository.On("Update", mock.IsType(&tt.song)).
				Run(func(args mock.Arguments) {
					newSong := args.Get(0).(*model.Song)
					assert.Equal(t, tt.expectedSongConfidence, math.Round(newSong.Confidence))
					assert.Equal(t, tt.expectedSongRehearsals, math.Round(newSong.Rehearsals))
					assert.Equal(t, tt.expectedSongProgress, math.Round(newSong.Progress))
				}).
				Return(nil).
				Once()

			arrangements := []model.SongArrangement{
				{
					ID:                 uuid.New(),
					Name:               "Perfect Rehearsal",
					Order:              0,
					SectionOccurrences: []model.SongSectionOccurrences{{SectionID: uuid.New(), Occurrences: 1}},
				},
				{ID: uuid.New(), Name: "Partial Rehearsal", Order: 1},
			}
			oldArrangements := slices.Clone(arrangements)
			songArrangementRepository.On("GetAllBySong", mock.IsType(new([]model.SongArrangement)), request.SongID).
				Return(nil, &arrangements).
				Once()

			songArrangementRepository.On("UpdateAllWithAssociations", mock.IsType(new([]model.SongArrangement))).
				Run(func(args mock.Arguments) {
					newArrangements := args.Get(0).(*[]model.SongArrangement)
					for i, arr := range *newArrangements {
						assert.Len(t, arr.SectionOccurrences, len(oldArrangements[i].SectionOccurrences)+1)
						newOccurrence := arr.SectionOccurrences[len(arr.SectionOccurrences)-1]
						assert.Equal(t, newSectionID, newOccurrence.SectionID)
						assert.Equal(t, arr.ID, newOccurrence.ArrangementID)
						assert.Zero(t, newOccurrence.Occurrences)
					}
				}).
				Return(nil).
				Once()

			// when
			errCode := _uut.Handle(request)

			// then
			assert.Nil(t, errCode)

			songSectionRepository.AssertExpectations(t)
			songRepository.AssertExpectations(t)
			songArrangementRepository.AssertExpectations(t)
		})
	}
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
	assert.Equal(t, model.DefaultSongSectionConfidence, section.Confidence)
	assert.Zero(t, section.RehearsalsScore)
	assert.Zero(t, section.ConfidenceScore)
	assert.Zero(t, section.Progress)
	assert.Equal(t, uint(count), section.Order)
	assert.Equal(t, request.TypeID, section.SongSectionTypeID)
	assert.Equal(t, request.BandMemberID, section.BandMemberID)
	assert.Equal(t, request.InstrumentID, section.InstrumentID)
	assert.Equal(t, request.SongID, section.SongID)
}
