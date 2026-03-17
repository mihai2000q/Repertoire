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
)

// Add Custom Rehearsal

func TestAddCustomRehearsal_WhenSongHasNoSections_ShouldReturnFalse(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	_uut := processor.NewSongProcessor(nil)

	mockSong := &model.Song{ID: uuid.New()}

	// when
	errCode, updated := _uut.AddCustomRehearsal(mockSong, songSectionRepository, nil)

	// then
	assert.False(t, updated)
	assert.Nil(t, errCode)

	songSectionRepository.AssertExpectations(t)
}

func TestAddCustomRehearsal_WhenSongHasNoArrangements_ShouldReturnFalse(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	_uut := processor.NewSongProcessor(nil)

	mockSong := &model.Song{
		ID:       uuid.New(),
		Sections: []model.SongSection{{ID: uuid.New()}},
	}

	// when
	errCode, updated := _uut.AddCustomRehearsal(mockSong, songSectionRepository, nil)

	// then
	assert.False(t, updated)
	assert.Nil(t, errCode)

	songSectionRepository.AssertExpectations(t)
}

func TestAddCustomRehearsal_WhenArrangementIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	_uut := processor.NewSongProcessor(nil)

	mockSong := &model.Song{
		ID: uuid.New(),
		Sections: []model.SongSection{
			{
				ID:                     uuid.New(),
				ArrangementOccurrences: []model.SongSectionOccurrences{{ArrangementID: uuid.New(), Occurrences: 0}},
			},
		},
	}

	// when
	errCode, updated := _uut.AddCustomRehearsal(mockSong, songSectionRepository, &[]uuid.UUID{uuid.New()}[0])

	// then
	assert.False(t, updated)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song arrangement not found", errCode.Error.Error())

	songSectionRepository.AssertExpectations(t)
}

func TestAddCustomRehearsal_WhenCreateHistoryFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	_uut := processor.NewSongProcessor(nil)

	mockSong := &model.Song{
		ID: uuid.New(),
		Sections: []model.SongSection{
			{
				ID:                     uuid.New(),
				ArrangementOccurrences: []model.SongSectionOccurrences{{SectionID: uuid.New(), Occurrences: 1}},
			},
		},
	}
	sectionsCount := len(mockSong.Sections)

	internalError := errors.New("internal error")
	songSectionRepository.On("CreateHistory", mock.IsType(new(model.SongSectionHistory))).
		Return(internalError).
		Times(sectionsCount)

	// when
	errCode, updated := _uut.AddCustomRehearsal(mockSong, songSectionRepository, nil)

	// then
	assert.False(t, updated)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songSectionRepository.AssertExpectations(t)
}

func TestAddCustomRehearsal_WhenGetHistoryFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	_uut := processor.NewSongProcessor(nil)

	mockSong := &model.Song{
		ID: uuid.New(),
		Sections: []model.SongSection{
			{
				ID:                     uuid.New(),
				ArrangementOccurrences: []model.SongSectionOccurrences{{Occurrences: 1}},
			},
		},
	}
	sectionsCount := len(mockSong.Sections)

	songSectionRepository.On("CreateHistory", mock.IsType(new(model.SongSectionHistory))).
		Return(nil).
		Times(sectionsCount)

	internalError := errors.New("internal error")
	songSectionRepository.
		On(
			"GetHistory",
			new([]model.SongSectionHistory),
			mock.IsType(uuid.UUID{}),
			model.RehearsalsProperty,
		).
		Return(internalError).
		Times(sectionsCount)

	// when
	errCode, updated := _uut.AddCustomRehearsal(mockSong, songSectionRepository, nil)

	// then
	assert.False(t, updated)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songSectionRepository.AssertExpectations(t)
}

func TestAddCustomRehearsal_WhenSectionsHaveZeroOccurrences_ShouldNotUpdateTheSong(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	progressProcessor := new(ProgressProcessorMock)
	_uut := processor.NewSongProcessor(progressProcessor)

	mockSong := &model.Song{
		ID: uuid.New(),
		Sections: []model.SongSection{
			{
				ID:                     uuid.New(),
				ArrangementOccurrences: []model.SongSectionOccurrences{{Occurrences: 0}},
			},
			{
				ID:                     uuid.New(),
				ArrangementOccurrences: []model.SongSectionOccurrences{{Occurrences: 0}},
			},
		},
	}

	// when
	errCode, updated := _uut.AddCustomRehearsal(mockSong, songSectionRepository, nil)

	// then
	assert.Nil(t, errCode)
	assert.False(t, updated)

	songSectionRepository.AssertExpectations(t)
	progressProcessor.AssertExpectations(t)
}

func TestAddCustomRehearsal_WhenSuccessfulWithoutArrangementID_ShouldUpdateSongAndSections(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	progressProcessor := new(ProgressProcessorMock)
	_uut := processor.NewSongProcessor(progressProcessor)

	mockSong := &model.Song{
		ID: uuid.New(),
		Sections: []model.SongSection{
			{
				ID:                     uuid.New(),
				Rehearsals:             23,
				ArrangementOccurrences: []model.SongSectionOccurrences{{Occurrences: 23}},
			},
			{
				ID:                     uuid.New(),
				ArrangementOccurrences: []model.SongSectionOccurrences{{Occurrences: 0}},
			},
			{
				ID:                     uuid.New(),
				ArrangementOccurrences: []model.SongSectionOccurrences{{Occurrences: 4}},
			},
		},
	}

	oldSections := slices.Clone(mockSong.Sections)
	sectionsCount := len(mockSong.Sections)

	sectionsCountWithOcc := len(slices.DeleteFunc(slices.Clone(mockSong.Sections), func(section model.SongSection) bool {
		return section.ArrangementOccurrences[0].Occurrences == 0
	}))

	songSectionRepository.On("CreateHistory", mock.IsType(new(model.SongSectionHistory))).
		Run(func(args mock.Arguments) {
			newHistory := args.Get(0).(*model.SongSectionHistory)
			assert.NotEmpty(t, newHistory.ID)
			assert.Equal(t, model.RehearsalsProperty, newHistory.Property)

			sections := slices.Clone(mockSong.Sections)
			sections = slices.DeleteFunc(sections, func(section model.SongSection) bool {
				return section.ID != newHistory.SongSectionID
			})

			assert.Equal(t, sections[0].Rehearsals, newHistory.From)
			assert.Equal(t, sections[0].Rehearsals+sections[0].ArrangementOccurrences[0].Occurrences, newHistory.To)
			assert.WithinDuration(t, newHistory.CreatedAt, time.Now().UTC(), time.Minute)
		}).
		Return(nil).
		Times(sectionsCountWithOcc)

	history := &[]model.SongSectionHistory{}
	songSectionRepository.
		On(
			"GetHistory",
			new([]model.SongSectionHistory),
			mock.IsType(uuid.UUID{}),
			model.RehearsalsProperty,
		).
		Run(func(args mock.Arguments) {
			sectionID := args.Get(1).(uuid.UUID)

			sections := slices.Clone(mockSong.Sections)
			sections = slices.DeleteFunc(sections, func(s model.SongSection) bool {
				return s.ID != sectionID
			})
			assert.Len(t, sections, 1, "ID is not part of the song sections")
		}).
		Return(nil, history).
		Times(sectionsCountWithOcc)

	var newRehearsalScore uint64 = 23
	progressProcessor.On("ComputeRehearsalsScore", *history).
		Return(newRehearsalScore).
		Times(sectionsCountWithOcc)

	var newProgress uint64 = 123
	progressProcessor.On("ComputeProgress", mock.IsType(model.SongSection{})).
		Run(func(args mock.Arguments) {
			sec := args.Get(0).(model.SongSection)
			assert.Contains(t, mockSong.Sections, sec)
		}).
		Return(newProgress).
		Times(sectionsCountWithOcc)

	// when
	errCode, updated := _uut.AddCustomRehearsal(mockSong, songSectionRepository, nil)

	// then
	assert.Nil(t, errCode)
	assert.True(t, updated)

	var newSongRehearsals uint = 0
	var newSongProgress uint64 = 0
	for i, section := range mockSong.Sections {
		if section.ArrangementOccurrences[0].Occurrences == 0 {
			assert.Equal(t, oldSections[i], section)
			continue
		}
		assert.Equal(t, oldSections[i].Rehearsals+section.ArrangementOccurrences[0].Occurrences, section.Rehearsals)
		assert.Equal(t, newRehearsalScore, section.RehearsalsScore)
		assert.Equal(t, newProgress, section.Progress)
		newSongRehearsals += section.Rehearsals
		newSongProgress += section.Progress
	}
	assert.Equal(t, float64(newSongProgress)/float64(sectionsCount), mockSong.Progress)
	assert.Equal(t, float64(newSongRehearsals)/float64(sectionsCount), mockSong.Rehearsals)
	assert.WithinDuration(t, time.Now(), *mockSong.LastTimePlayed, 1*time.Minute)

	songSectionRepository.AssertExpectations(t)
	progressProcessor.AssertExpectations(t)
}

func TestAddCustomRehearsal_WhenSuccessfulWithArrangementID_ShouldUpdateSongAndSections(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	progressProcessor := new(ProgressProcessorMock)
	_uut := processor.NewSongProcessor(progressProcessor)

	arrangementID := uuid.New()
	mockSong := &model.Song{
		ID: uuid.New(),
		Sections: []model.SongSection{
			{
				ID:         uuid.New(),
				Rehearsals: 23,
				ArrangementOccurrences: []model.SongSectionOccurrences{
					{ArrangementID: arrangementID, Occurrences: 23},
					{ArrangementID: uuid.New(), Occurrences: 4},
				},
			},
			{
				ID: uuid.New(),
				ArrangementOccurrences: []model.SongSectionOccurrences{
					{ArrangementID: uuid.New(), Occurrences: 1},
					{ArrangementID: arrangementID, Occurrences: 0},
				},
			},
			{
				ID: uuid.New(),
				ArrangementOccurrences: []model.SongSectionOccurrences{
					{ArrangementID: uuid.New(), Occurrences: 0},
					{ArrangementID: arrangementID, Occurrences: 4},
				},
			},
		},
	}

	oldSections := slices.Clone(mockSong.Sections)
	sectionsCount := len(mockSong.Sections)

	sectionsCountWithOcc := len(slices.DeleteFunc(slices.Clone(mockSong.Sections), func(section model.SongSection) bool {
		arrangementIndex := slices.IndexFunc(section.ArrangementOccurrences, func(o model.SongSectionOccurrences) bool {
			return o.ArrangementID == arrangementID
		})
		return section.ArrangementOccurrences[arrangementIndex].Occurrences == 0
	}))

	songSectionRepository.On("CreateHistory", mock.IsType(new(model.SongSectionHistory))).
		Run(func(args mock.Arguments) {
			newHistory := args.Get(0).(*model.SongSectionHistory)
			assert.NotEmpty(t, newHistory.ID)
			assert.Equal(t, model.RehearsalsProperty, newHistory.Property)

			sections := slices.Clone(mockSong.Sections)
			sections = slices.DeleteFunc(sections, func(section model.SongSection) bool {
				return section.ID != newHistory.SongSectionID
			})
			section := sections[0]

			arrangementIndex := slices.IndexFunc(section.ArrangementOccurrences, func(o model.SongSectionOccurrences) bool {
				return o.ArrangementID == arrangementID
			})

			assert.Equal(t, section.Rehearsals, newHistory.From)
			assert.Equal(t, section.Rehearsals+section.ArrangementOccurrences[arrangementIndex].Occurrences, newHistory.To)
			assert.WithinDuration(t, newHistory.CreatedAt, time.Now().UTC(), time.Minute)
		}).
		Return(nil).
		Times(sectionsCountWithOcc)

	history := &[]model.SongSectionHistory{}
	songSectionRepository.
		On(
			"GetHistory",
			new([]model.SongSectionHistory),
			mock.IsType(uuid.UUID{}),
			model.RehearsalsProperty,
		).
		Run(func(args mock.Arguments) {
			sectionID := args.Get(1).(uuid.UUID)

			sections := slices.Clone(mockSong.Sections)
			sections = slices.DeleteFunc(sections, func(s model.SongSection) bool {
				return s.ID != sectionID
			})
			assert.Len(t, sections, 1, "ID is not part of the song sections")
		}).
		Return(nil, history).
		Times(sectionsCountWithOcc)

	var newRehearsalScore uint64 = 23
	progressProcessor.On("ComputeRehearsalsScore", *history).
		Return(newRehearsalScore).
		Times(sectionsCountWithOcc)

	var newProgress uint64 = 123
	progressProcessor.On("ComputeProgress", mock.IsType(model.SongSection{})).
		Run(func(args mock.Arguments) {
			sec := args.Get(0).(model.SongSection)
			assert.Contains(t, mockSong.Sections, sec)
		}).
		Return(newProgress).
		Times(sectionsCountWithOcc)

	// when
	errCode, updated := _uut.AddCustomRehearsal(mockSong, songSectionRepository, &arrangementID)

	// then
	assert.Nil(t, errCode)
	assert.True(t, updated)

	var newSongRehearsals uint = 0
	var newSongProgress uint64 = 0
	for i, section := range mockSong.Sections {
		arrangementIndex := slices.IndexFunc(section.ArrangementOccurrences, func(o model.SongSectionOccurrences) bool {
			return o.ArrangementID == arrangementID
		})
		if section.ArrangementOccurrences[arrangementIndex].Occurrences == 0 {
			assert.Equal(t, oldSections[i], section)
			continue
		}
		assert.Equal(t, oldSections[i].Rehearsals+section.ArrangementOccurrences[arrangementIndex].Occurrences, section.Rehearsals)
		assert.Equal(t, newRehearsalScore, section.RehearsalsScore)
		assert.Equal(t, newProgress, section.Progress)
		newSongRehearsals += section.Rehearsals
		newSongProgress += section.Progress
	}
	assert.Equal(t, float64(newSongProgress)/float64(sectionsCount), mockSong.Progress)
	assert.Equal(t, float64(newSongRehearsals)/float64(sectionsCount), mockSong.Rehearsals)
	assert.WithinDuration(t, time.Now(), *mockSong.LastTimePlayed, 1*time.Minute)

	songSectionRepository.AssertExpectations(t)
	progressProcessor.AssertExpectations(t)
}

// Add Perfect Rehearsal

func TestAddPerfectRehearsal_WhenSongHasNoDefaultArrangement_ShouldReturnFalse(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	_uut := processor.NewSongProcessor(nil)

	mockSong := &model.Song{ID: uuid.New()}

	// when
	errCode, updated := _uut.AddPerfectRehearsal(mockSong, songSectionRepository)

	// then
	assert.False(t, updated)
	assert.Nil(t, errCode)

	songSectionRepository.AssertExpectations(t)
}

func TestAddPerfectRehearsal_WhenCreateHistoryFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	_uut := processor.NewSongProcessor(nil)

	mockSong := &model.Song{
		ID:                   uuid.New(),
		DefaultArrangementID: &[]uuid.UUID{uuid.New()}[0],
		Sections: []model.SongSection{
			{
				ID:                     uuid.New(),
				ArrangementOccurrences: []model.SongSectionOccurrences{{SectionID: uuid.New(), Occurrences: 1}},
			},
		},
	}
	sectionsCount := len(mockSong.Sections)

	internalError := errors.New("internal error")
	songSectionRepository.On("CreateHistory", mock.IsType(new(model.SongSectionHistory))).
		Return(internalError).
		Times(sectionsCount)

	// when
	errCode, updated := _uut.AddPerfectRehearsal(mockSong, songSectionRepository)

	// then
	assert.False(t, updated)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songSectionRepository.AssertExpectations(t)
}

func TestAddPerfectRehearsal_WhenGetHistoryFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	_uut := processor.NewSongProcessor(nil)

	mockSong := &model.Song{
		ID:                   uuid.New(),
		DefaultArrangementID: &[]uuid.UUID{uuid.New()}[0],
		Sections: []model.SongSection{
			{
				ID:                     uuid.New(),
				ArrangementOccurrences: []model.SongSectionOccurrences{{Occurrences: 1}},
			},
		},
	}
	sectionsCount := len(mockSong.Sections)

	songSectionRepository.On("CreateHistory", mock.IsType(new(model.SongSectionHistory))).
		Return(nil).
		Times(sectionsCount)

	internalError := errors.New("internal error")
	songSectionRepository.
		On(
			"GetHistory",
			new([]model.SongSectionHistory),
			mock.IsType(uuid.UUID{}),
			model.RehearsalsProperty,
		).
		Return(internalError).
		Times(sectionsCount)

	// when
	errCode, updated := _uut.AddPerfectRehearsal(mockSong, songSectionRepository)

	// then
	assert.False(t, updated)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songSectionRepository.AssertExpectations(t)
}

func TestAddPerfectRehearsal_WhenSectionsHaveZeroOccurrences_ShouldNotUpdateTheSong(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	progressProcessor := new(ProgressProcessorMock)
	_uut := processor.NewSongProcessor(progressProcessor)

	mockSong := &model.Song{
		ID:                   uuid.New(),
		DefaultArrangementID: &[]uuid.UUID{uuid.New()}[0],
		Sections: []model.SongSection{
			{
				ID:                     uuid.New(),
				ArrangementOccurrences: []model.SongSectionOccurrences{{Occurrences: 0}},
			},
			{
				ID:                     uuid.New(),
				ArrangementOccurrences: []model.SongSectionOccurrences{{Occurrences: 0}},
			},
		},
	}

	// when
	errCode, updated := _uut.AddPerfectRehearsal(mockSong, songSectionRepository)

	// then
	assert.Nil(t, errCode)
	assert.False(t, updated)

	songSectionRepository.AssertExpectations(t)
	progressProcessor.AssertExpectations(t)
}

func TestAddPerfectRehearsal_WhenSuccessful_ShouldUpdateSongAndSections(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	progressProcessor := new(ProgressProcessorMock)
	_uut := processor.NewSongProcessor(progressProcessor)

	mockSong := &model.Song{
		ID:                   uuid.New(),
		DefaultArrangementID: &[]uuid.UUID{uuid.New()}[0],
		Sections: []model.SongSection{
			{
				ID:                     uuid.New(),
				Rehearsals:             23,
				ArrangementOccurrences: []model.SongSectionOccurrences{{Occurrences: 23}},
			},
			{
				ID:                     uuid.New(),
				ArrangementOccurrences: []model.SongSectionOccurrences{{Occurrences: 0}},
			},
			{
				ID:                     uuid.New(),
				ArrangementOccurrences: []model.SongSectionOccurrences{{Occurrences: 4}},
			},
		},
	}

	oldSections := slices.Clone(mockSong.Sections)
	sectionsCount := len(mockSong.Sections)

	sectionsCountWithOcc := len(slices.DeleteFunc(slices.Clone(mockSong.Sections), func(section model.SongSection) bool {
		return section.ArrangementOccurrences[0].Occurrences == 0
	}))

	songSectionRepository.On("CreateHistory", mock.IsType(new(model.SongSectionHistory))).
		Run(func(args mock.Arguments) {
			newHistory := args.Get(0).(*model.SongSectionHistory)
			assert.NotEmpty(t, newHistory.ID)
			assert.Equal(t, model.RehearsalsProperty, newHistory.Property)

			sections := slices.Clone(mockSong.Sections)
			sections = slices.DeleteFunc(sections, func(section model.SongSection) bool {
				return section.ID != newHistory.SongSectionID
			})

			assert.Equal(t, sections[0].Rehearsals, newHistory.From)
			assert.Equal(t, sections[0].Rehearsals+sections[0].ArrangementOccurrences[0].Occurrences, newHistory.To)
			assert.WithinDuration(t, newHistory.CreatedAt, time.Now().UTC(), time.Minute)
		}).
		Return(nil).
		Times(sectionsCountWithOcc)

	history := &[]model.SongSectionHistory{}
	songSectionRepository.
		On(
			"GetHistory",
			new([]model.SongSectionHistory),
			mock.IsType(uuid.UUID{}),
			model.RehearsalsProperty,
		).
		Run(func(args mock.Arguments) {
			sectionID := args.Get(1).(uuid.UUID)

			sections := slices.Clone(mockSong.Sections)
			sections = slices.DeleteFunc(sections, func(s model.SongSection) bool {
				return s.ID != sectionID
			})
			assert.Len(t, sections, 1, "ID is not part of the song sections")
		}).
		Return(nil, history).
		Times(sectionsCountWithOcc)

	var newRehearsalScore uint64 = 23
	progressProcessor.On("ComputeRehearsalsScore", *history).
		Return(newRehearsalScore).
		Times(sectionsCountWithOcc)

	var newProgress uint64 = 123
	progressProcessor.On("ComputeProgress", mock.IsType(model.SongSection{})).
		Run(func(args mock.Arguments) {
			sec := args.Get(0).(model.SongSection)
			assert.Contains(t, mockSong.Sections, sec)
		}).
		Return(newProgress).
		Times(sectionsCountWithOcc)

	// when
	errCode, updated := _uut.AddPerfectRehearsal(mockSong, songSectionRepository)

	// then
	assert.Nil(t, errCode)
	assert.True(t, updated)

	var newSongRehearsals uint = 0
	var newSongProgress uint64 = 0
	for i, section := range mockSong.Sections {
		if section.ArrangementOccurrences[0].Occurrences == 0 {
			assert.Equal(t, oldSections[i], section)
			continue
		}
		assert.Equal(t, oldSections[i].Rehearsals+section.ArrangementOccurrences[0].Occurrences, section.Rehearsals)
		assert.Equal(t, newRehearsalScore, section.RehearsalsScore)
		assert.Equal(t, newProgress, section.Progress)
		newSongRehearsals += section.Rehearsals
		newSongProgress += section.Progress
	}
	assert.Equal(t, float64(newSongProgress)/float64(sectionsCount), mockSong.Progress)
	assert.Equal(t, float64(newSongRehearsals)/float64(sectionsCount), mockSong.Rehearsals)
	assert.WithinDuration(t, time.Now(), *mockSong.LastTimePlayed, 1*time.Minute)

	songSectionRepository.AssertExpectations(t)
	progressProcessor.AssertExpectations(t)
}
