package processor

import (
	"errors"
	"net/http"
	"repertoire/server/domain/processor"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"slices"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Add Custom Rehearsal

func TestAddCustomRehearsal_WhenSongHasNoParts_ShouldReturnFalse(t *testing.T) {
	// given
	songPartRepository := new(repository.SongPartRepositoryMock)
	_uut := processor.NewSongProcessor(nil)

	mockSong := &model.Song{ID: uuid.New()}

	// when
	errCode, updated := _uut.AddCustomRehearsal(mockSong, songPartRepository, nil)

	// then
	assert.False(t, updated)
	assert.Nil(t, errCode)

	songPartRepository.AssertExpectations(t)
}

func TestAddCustomRehearsal_WhenSongHasNoArrangements_ShouldReturnFalse(t *testing.T) {
	// given
	songPartRepository := new(repository.SongPartRepositoryMock)
	_uut := processor.NewSongProcessor(nil)

	mockSong := &model.Song{
		ID:    uuid.New(),
		Parts: []model.SongPart{{ID: uuid.New()}},
	}

	// when
	errCode, updated := _uut.AddCustomRehearsal(mockSong, songPartRepository, nil)

	// then
	assert.False(t, updated)
	assert.Nil(t, errCode)

	songPartRepository.AssertExpectations(t)
}

func TestAddCustomRehearsal_WhenArrangementIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songPartRepository := new(repository.SongPartRepositoryMock)
	_uut := processor.NewSongProcessor(nil)

	mockSong := &model.Song{
		ID: uuid.New(),
		Parts: []model.SongPart{
			{
				ID:                     uuid.New(),
				ArrangementOccurrences: []model.SongPartOccurrences{{ArrangementID: uuid.New(), Occurrences: 0}},
			},
		},
	}

	// when
	errCode, updated := _uut.AddCustomRehearsal(mockSong, songPartRepository, &[]uuid.UUID{uuid.New()}[0])

	// then
	assert.False(t, updated)
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song arrangement not found", errCode.Error.Error())

	songPartRepository.AssertExpectations(t)
}

func TestAddCustomRehearsal_WhenCreateHistoryFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songPartRepository := new(repository.SongPartRepositoryMock)
	_uut := processor.NewSongProcessor(nil)

	mockSong := &model.Song{
		ID: uuid.New(),
		Parts: []model.SongPart{
			{
				ID:                     uuid.New(),
				ArrangementOccurrences: []model.SongPartOccurrences{{PartID: uuid.New(), Occurrences: 1}},
			},
		},
	}
	partsCount := len(mockSong.Parts)

	internalError := errors.New("internal error")
	songPartRepository.On("CreateHistory", mock.IsType(new(model.SongPartHistory))).
		Return(internalError).
		Times(partsCount)

	// when
	errCode, updated := _uut.AddCustomRehearsal(mockSong, songPartRepository, nil)

	// then
	assert.False(t, updated)
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songPartRepository.AssertExpectations(t)
}

func TestAddCustomRehearsal_WhenGetHistoryFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songPartRepository := new(repository.SongPartRepositoryMock)
	_uut := processor.NewSongProcessor(nil)

	mockSong := &model.Song{
		ID: uuid.New(),
		Parts: []model.SongPart{
			{
				ID:                     uuid.New(),
				ArrangementOccurrences: []model.SongPartOccurrences{{Occurrences: 1}},
			},
		},
	}
	partsCount := len(mockSong.Parts)

	songPartRepository.On("CreateHistory", mock.IsType(new(model.SongPartHistory))).
		Return(nil).
		Times(partsCount)

	internalError := errors.New("internal error")
	songPartRepository.
		On(
			"GetHistory",
			new([]model.SongPartHistory),
			mock.IsType(uuid.UUID{}),
			model.RehearsalsProperty,
		).
		Return(internalError).
		Times(partsCount)

	// when
	errCode, updated := _uut.AddCustomRehearsal(mockSong, songPartRepository, nil)

	// then
	assert.False(t, updated)
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songPartRepository.AssertExpectations(t)
}

func TestAddCustomRehearsal_WhenPartsHaveZeroOccurrences_ShouldNotUpdateTheSong(t *testing.T) {
	// given
	songPartRepository := new(repository.SongPartRepositoryMock)
	progressProcessor := new(ProgressProcessorMock)
	_uut := processor.NewSongProcessor(progressProcessor)

	mockSong := &model.Song{
		ID: uuid.New(),
		Parts: []model.SongPart{
			{
				ID:                     uuid.New(),
				ArrangementOccurrences: []model.SongPartOccurrences{{Occurrences: 0}},
			},
			{
				ID:                     uuid.New(),
				ArrangementOccurrences: []model.SongPartOccurrences{{Occurrences: 0}},
			},
		},
	}

	// when
	errCode, updated := _uut.AddCustomRehearsal(mockSong, songPartRepository, nil)

	// then
	assert.Nil(t, errCode)
	assert.False(t, updated)

	songPartRepository.AssertExpectations(t)
	progressProcessor.AssertExpectations(t)
}

func TestAddCustomRehearsal_WhenSuccessfulWithoutArrangementID_ShouldUpdateSongAndParts(t *testing.T) {
	// given
	songPartRepository := new(repository.SongPartRepositoryMock)
	progressProcessor := new(ProgressProcessorMock)
	_uut := processor.NewSongProcessor(progressProcessor)

	mockSong := &model.Song{
		ID: uuid.New(),
		Parts: []model.SongPart{
			{
				ID:                     uuid.New(),
				Rehearsals:             23,
				ArrangementOccurrences: []model.SongPartOccurrences{{Occurrences: 23}},
			},
			{
				ID:                     uuid.New(),
				ArrangementOccurrences: []model.SongPartOccurrences{{Occurrences: 0}},
			},
			{
				ID:                     uuid.New(),
				ArrangementOccurrences: []model.SongPartOccurrences{{Occurrences: 4}},
			},
		},
	}

	oldParts := slices.Clone(mockSong.Parts)
	partsCount := len(mockSong.Parts)

	partsCountWithOcc := len(slices.DeleteFunc(slices.Clone(mockSong.Parts), func(part model.SongPart) bool {
		return part.ArrangementOccurrences[0].Occurrences == 0
	}))

	songPartRepository.On("CreateHistory", mock.IsType(new(model.SongPartHistory))).
		Run(func(args mock.Arguments) {
			newHistory := args.Get(0).(*model.SongPartHistory)
			assert.NotEmpty(t, newHistory.ID)
			assert.Equal(t, model.RehearsalsProperty, newHistory.Property)

			parts := slices.Clone(mockSong.Parts)
			parts = slices.DeleteFunc(parts, func(part model.SongPart) bool {
				return part.ID != newHistory.PartID
			})

			assert.Equal(t, parts[0].Rehearsals, newHistory.From)
			assert.Equal(t, parts[0].Rehearsals+parts[0].ArrangementOccurrences[0].Occurrences, newHistory.To)
			assert.WithinDuration(t, newHistory.CreatedAt, time.Now().UTC(), time.Minute)
		}).
		Return(nil).
		Times(partsCountWithOcc)

	history := &[]model.SongPartHistory{}
	songPartRepository.
		On(
			"GetHistory",
			new([]model.SongPartHistory),
			mock.IsType(uuid.UUID{}),
			model.RehearsalsProperty,
		).
		Run(func(args mock.Arguments) {
			partID := args.Get(1).(uuid.UUID)

			parts := slices.Clone(mockSong.Parts)
			parts = slices.DeleteFunc(parts, func(s model.SongPart) bool {
				return s.ID != partID
			})
			assert.Len(t, parts, 1, "ID is not part of the song parts")
		}).
		Return(nil, history).
		Times(partsCountWithOcc)

	var newRehearsalScore uint64 = 23
	progressProcessor.On("ComputeRehearsalsScore", *history).
		Return(newRehearsalScore).
		Times(partsCountWithOcc)

	var newProgress uint64 = 123
	progressProcessor.On("ComputeProgress", mock.IsType(model.SongPart{})).
		Run(func(args mock.Arguments) {
			sec := args.Get(0).(model.SongPart)
			assert.Contains(t, mockSong.Parts, sec)
		}).
		Return(newProgress).
		Times(partsCountWithOcc)

	// when
	errCode, updated := _uut.AddCustomRehearsal(mockSong, songPartRepository, nil)

	// then
	assert.Nil(t, errCode)
	assert.True(t, updated)

	var newSongRehearsals uint = 0
	var newSongProgress uint64 = 0
	for i, part := range mockSong.Parts {
		if part.ArrangementOccurrences[0].Occurrences == 0 {
			assert.Equal(t, oldParts[i], part)
			continue
		}
		assert.Equal(t, oldParts[i].Rehearsals+part.ArrangementOccurrences[0].Occurrences, part.Rehearsals)
		assert.Equal(t, newRehearsalScore, part.RehearsalsScore)
		assert.Equal(t, newProgress, part.Progress)
		newSongRehearsals += part.Rehearsals
		newSongProgress += part.Progress
	}
	assert.Equal(t, float64(newSongProgress)/float64(partsCount), mockSong.Progress)
	assert.Equal(t, float64(newSongRehearsals)/float64(partsCount), mockSong.Rehearsals)
	assert.WithinDuration(t, time.Now(), *mockSong.LastTimePlayed, 1*time.Minute)

	songPartRepository.AssertExpectations(t)
	progressProcessor.AssertExpectations(t)
}

func TestAddCustomRehearsal_WhenSuccessfulWithArrangementID_ShouldUpdateSongAndParts(t *testing.T) {
	// given
	songPartRepository := new(repository.SongPartRepositoryMock)
	progressProcessor := new(ProgressProcessorMock)
	_uut := processor.NewSongProcessor(progressProcessor)

	arrangementID := uuid.New()
	mockSong := &model.Song{
		ID: uuid.New(),
		Parts: []model.SongPart{
			{
				ID:         uuid.New(),
				Rehearsals: 23,
				ArrangementOccurrences: []model.SongPartOccurrences{
					{ArrangementID: arrangementID, Occurrences: 23},
					{ArrangementID: uuid.New(), Occurrences: 4},
				},
			},
			{
				ID: uuid.New(),
				ArrangementOccurrences: []model.SongPartOccurrences{
					{ArrangementID: uuid.New(), Occurrences: 1},
					{ArrangementID: arrangementID, Occurrences: 0},
				},
			},
			{
				ID: uuid.New(),
				ArrangementOccurrences: []model.SongPartOccurrences{
					{ArrangementID: uuid.New(), Occurrences: 0},
					{ArrangementID: arrangementID, Occurrences: 4},
				},
			},
		},
	}

	oldParts := slices.Clone(mockSong.Parts)
	partsCount := len(mockSong.Parts)

	partsCountWithOcc := len(slices.DeleteFunc(slices.Clone(mockSong.Parts), func(part model.SongPart) bool {
		arrangementIndex := slices.IndexFunc(part.ArrangementOccurrences, func(o model.SongPartOccurrences) bool {
			return o.ArrangementID == arrangementID
		})
		return part.ArrangementOccurrences[arrangementIndex].Occurrences == 0
	}))

	songPartRepository.On("CreateHistory", mock.IsType(new(model.SongPartHistory))).
		Run(func(args mock.Arguments) {
			newHistory := args.Get(0).(*model.SongPartHistory)
			assert.NotEmpty(t, newHistory.ID)
			assert.Equal(t, model.RehearsalsProperty, newHistory.Property)

			parts := slices.Clone(mockSong.Parts)
			parts = slices.DeleteFunc(parts, func(part model.SongPart) bool {
				return part.ID != newHistory.PartID
			})
			part := parts[0]

			arrangementIndex := slices.IndexFunc(part.ArrangementOccurrences, func(o model.SongPartOccurrences) bool {
				return o.ArrangementID == arrangementID
			})

			assert.Equal(t, part.Rehearsals, newHistory.From)
			assert.Equal(t, part.Rehearsals+part.ArrangementOccurrences[arrangementIndex].Occurrences, newHistory.To)
			assert.WithinDuration(t, newHistory.CreatedAt, time.Now().UTC(), time.Minute)
		}).
		Return(nil).
		Times(partsCountWithOcc)

	history := &[]model.SongPartHistory{}
	songPartRepository.
		On(
			"GetHistory",
			new([]model.SongPartHistory),
			mock.IsType(uuid.UUID{}),
			model.RehearsalsProperty,
		).
		Run(func(args mock.Arguments) {
			partID := args.Get(1).(uuid.UUID)

			parts := slices.Clone(mockSong.Parts)
			parts = slices.DeleteFunc(parts, func(s model.SongPart) bool {
				return s.ID != partID
			})
			assert.Len(t, parts, 1, "ID is not part of the song parts")
		}).
		Return(nil, history).
		Times(partsCountWithOcc)

	var newRehearsalScore uint64 = 23
	progressProcessor.On("ComputeRehearsalsScore", *history).
		Return(newRehearsalScore).
		Times(partsCountWithOcc)

	var newProgress uint64 = 123
	progressProcessor.On("ComputeProgress", mock.IsType(model.SongPart{})).
		Run(func(args mock.Arguments) {
			sec := args.Get(0).(model.SongPart)
			assert.Contains(t, mockSong.Parts, sec)
		}).
		Return(newProgress).
		Times(partsCountWithOcc)

	// when
	errCode, updated := _uut.AddCustomRehearsal(mockSong, songPartRepository, &arrangementID)

	// then
	assert.Nil(t, errCode)
	assert.True(t, updated)

	var newSongRehearsals uint = 0
	var newSongProgress uint64 = 0
	for i, part := range mockSong.Parts {
		arrangementIndex := slices.IndexFunc(part.ArrangementOccurrences, func(o model.SongPartOccurrences) bool {
			return o.ArrangementID == arrangementID
		})
		if part.ArrangementOccurrences[arrangementIndex].Occurrences == 0 {
			assert.Equal(t, oldParts[i], part)
			continue
		}
		assert.Equal(t, oldParts[i].Rehearsals+part.ArrangementOccurrences[arrangementIndex].Occurrences, part.Rehearsals)
		assert.Equal(t, newRehearsalScore, part.RehearsalsScore)
		assert.Equal(t, newProgress, part.Progress)
		newSongRehearsals += part.Rehearsals
		newSongProgress += part.Progress
	}
	assert.Equal(t, float64(newSongProgress)/float64(partsCount), mockSong.Progress)
	assert.Equal(t, float64(newSongRehearsals)/float64(partsCount), mockSong.Rehearsals)
	assert.WithinDuration(t, time.Now(), *mockSong.LastTimePlayed, 1*time.Minute)

	songPartRepository.AssertExpectations(t)
	progressProcessor.AssertExpectations(t)
}

// Add Perfect Rehearsal

func TestAddPerfectRehearsal_WhenSongHasNoDefaultArrangement_ShouldReturnFalse(t *testing.T) {
	// given
	songPartRepository := new(repository.SongPartRepositoryMock)
	_uut := processor.NewSongProcessor(nil)

	mockSong := &model.Song{ID: uuid.New()}

	// when
	errCode, updated := _uut.AddPerfectRehearsal(mockSong, songPartRepository)

	// then
	assert.False(t, updated)
	assert.Nil(t, errCode)

	songPartRepository.AssertExpectations(t)
}

func TestAddPerfectRehearsal_WhenCreateHistoryFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songPartRepository := new(repository.SongPartRepositoryMock)
	_uut := processor.NewSongProcessor(nil)

	mockSong := &model.Song{
		ID:                   uuid.New(),
		DefaultArrangementID: &[]uuid.UUID{uuid.New()}[0],
		Parts: []model.SongPart{
			{
				ID:                     uuid.New(),
				ArrangementOccurrences: []model.SongPartOccurrences{{PartID: uuid.New(), Occurrences: 1}},
			},
		},
	}
	partsCount := len(mockSong.Parts)

	internalError := errors.New("internal error")
	songPartRepository.On("CreateHistory", mock.IsType(new(model.SongPartHistory))).
		Return(internalError).
		Times(partsCount)

	// when
	errCode, updated := _uut.AddPerfectRehearsal(mockSong, songPartRepository)

	// then
	assert.False(t, updated)
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songPartRepository.AssertExpectations(t)
}

func TestAddPerfectRehearsal_WhenGetHistoryFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songPartRepository := new(repository.SongPartRepositoryMock)
	_uut := processor.NewSongProcessor(nil)

	mockSong := &model.Song{
		ID:                   uuid.New(),
		DefaultArrangementID: &[]uuid.UUID{uuid.New()}[0],
		Parts: []model.SongPart{
			{
				ID:                     uuid.New(),
				ArrangementOccurrences: []model.SongPartOccurrences{{Occurrences: 1}},
			},
		},
	}
	partsCount := len(mockSong.Parts)

	songPartRepository.On("CreateHistory", mock.IsType(new(model.SongPartHistory))).
		Return(nil).
		Times(partsCount)

	internalError := errors.New("internal error")
	songPartRepository.
		On(
			"GetHistory",
			new([]model.SongPartHistory),
			mock.IsType(uuid.UUID{}),
			model.RehearsalsProperty,
		).
		Return(internalError).
		Times(partsCount)

	// when
	errCode, updated := _uut.AddPerfectRehearsal(mockSong, songPartRepository)

	// then
	assert.False(t, updated)
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songPartRepository.AssertExpectations(t)
}

func TestAddPerfectRehearsal_WhenPartsHaveZeroOccurrences_ShouldNotUpdateTheSong(t *testing.T) {
	// given
	songPartRepository := new(repository.SongPartRepositoryMock)
	progressProcessor := new(ProgressProcessorMock)
	_uut := processor.NewSongProcessor(progressProcessor)

	mockSong := &model.Song{
		ID:                   uuid.New(),
		DefaultArrangementID: &[]uuid.UUID{uuid.New()}[0],
		Parts: []model.SongPart{
			{
				ID:                     uuid.New(),
				ArrangementOccurrences: []model.SongPartOccurrences{{Occurrences: 0}},
			},
			{
				ID:                     uuid.New(),
				ArrangementOccurrences: []model.SongPartOccurrences{{Occurrences: 0}},
			},
		},
	}

	// when
	errCode, updated := _uut.AddPerfectRehearsal(mockSong, songPartRepository)

	// then
	assert.Nil(t, errCode)
	assert.False(t, updated)

	songPartRepository.AssertExpectations(t)
	progressProcessor.AssertExpectations(t)
}

func TestAddPerfectRehearsal_WhenSuccessful_ShouldUpdateSongAndParts(t *testing.T) {
	// given
	songPartRepository := new(repository.SongPartRepositoryMock)
	progressProcessor := new(ProgressProcessorMock)
	_uut := processor.NewSongProcessor(progressProcessor)

	mockSong := &model.Song{
		ID:                   uuid.New(),
		DefaultArrangementID: &[]uuid.UUID{uuid.New()}[0],
		Parts: []model.SongPart{
			{
				ID:                     uuid.New(),
				Rehearsals:             23,
				ArrangementOccurrences: []model.SongPartOccurrences{{Occurrences: 23}},
			},
			{
				ID:                     uuid.New(),
				ArrangementOccurrences: []model.SongPartOccurrences{{Occurrences: 0}},
			},
			{
				ID:                     uuid.New(),
				ArrangementOccurrences: []model.SongPartOccurrences{{Occurrences: 4}},
			},
		},
	}

	oldParts := slices.Clone(mockSong.Parts)
	partsCount := len(mockSong.Parts)

	partsCountWithOcc := len(slices.DeleteFunc(slices.Clone(mockSong.Parts), func(part model.SongPart) bool {
		return part.ArrangementOccurrences[0].Occurrences == 0
	}))

	songPartRepository.On("CreateHistory", mock.IsType(new(model.SongPartHistory))).
		Run(func(args mock.Arguments) {
			newHistory := args.Get(0).(*model.SongPartHistory)
			assert.NotEmpty(t, newHistory.ID)
			assert.Equal(t, model.RehearsalsProperty, newHistory.Property)

			parts := slices.Clone(mockSong.Parts)
			parts = slices.DeleteFunc(parts, func(part model.SongPart) bool {
				return part.ID != newHistory.PartID
			})

			assert.Equal(t, parts[0].Rehearsals, newHistory.From)
			assert.Equal(t, parts[0].Rehearsals+parts[0].ArrangementOccurrences[0].Occurrences, newHistory.To)
			assert.WithinDuration(t, newHistory.CreatedAt, time.Now().UTC(), time.Minute)
		}).
		Return(nil).
		Times(partsCountWithOcc)

	history := &[]model.SongPartHistory{}
	songPartRepository.
		On(
			"GetHistory",
			new([]model.SongPartHistory),
			mock.IsType(uuid.UUID{}),
			model.RehearsalsProperty,
		).
		Run(func(args mock.Arguments) {
			partID := args.Get(1).(uuid.UUID)

			parts := slices.Clone(mockSong.Parts)
			parts = slices.DeleteFunc(parts, func(s model.SongPart) bool {
				return s.ID != partID
			})
			assert.Len(t, parts, 1, "ID is not part of the song parts")
		}).
		Return(nil, history).
		Times(partsCountWithOcc)

	var newRehearsalScore uint64 = 23
	progressProcessor.On("ComputeRehearsalsScore", *history).
		Return(newRehearsalScore).
		Times(partsCountWithOcc)

	var newProgress uint64 = 123
	progressProcessor.On("ComputeProgress", mock.IsType(model.SongPart{})).
		Run(func(args mock.Arguments) {
			sec := args.Get(0).(model.SongPart)
			assert.Contains(t, mockSong.Parts, sec)
		}).
		Return(newProgress).
		Times(partsCountWithOcc)

	// when
	errCode, updated := _uut.AddPerfectRehearsal(mockSong, songPartRepository)

	// then
	assert.Nil(t, errCode)
	assert.True(t, updated)

	var newSongRehearsals uint = 0
	var newSongProgress uint64 = 0
	for i, part := range mockSong.Parts {
		if part.ArrangementOccurrences[0].Occurrences == 0 {
			assert.Equal(t, oldParts[i], part)
			continue
		}
		assert.Equal(t, oldParts[i].Rehearsals+part.ArrangementOccurrences[0].Occurrences, part.Rehearsals)
		assert.Equal(t, newRehearsalScore, part.RehearsalsScore)
		assert.Equal(t, newProgress, part.Progress)
		newSongRehearsals += part.Rehearsals
		newSongProgress += part.Progress
	}
	assert.Equal(t, float64(newSongProgress)/float64(partsCount), mockSong.Progress)
	assert.Equal(t, float64(newSongRehearsals)/float64(partsCount), mockSong.Rehearsals)
	assert.WithinDuration(t, time.Now(), *mockSong.LastTimePlayed, 1*time.Minute)

	songPartRepository.AssertExpectations(t)
	progressProcessor.AssertExpectations(t)
}
