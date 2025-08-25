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

func TestAddPerfectRehearsal_WhenCreateHistoryFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := processor.NewSongProcessor(nil)

	mockSong := &model.Song{
		ID:       uuid.New(),
		Sections: []model.SongSection{{ID: uuid.New(), Occurrences: 2}},
	}
	sectionsCount := len(mockSong.Sections)

	internalError := errors.New("internal error")
	songRepository.On("CreateSongSectionHistory", mock.IsType(new(model.SongSectionHistory))).
		Return(internalError).
		Times(sectionsCount)

	// when
	errCode, updated := _uut.AddPerfectRehearsal(mockSong, songRepository)

	// then
	assert.False(t, updated)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestAddPerfectRehearsal_WhenGetHistoryFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := processor.NewSongProcessor(nil)

	mockSong := &model.Song{
		ID:       uuid.New(),
		Sections: []model.SongSection{{ID: uuid.New(), Occurrences: 2}},
	}
	sectionsCount := len(mockSong.Sections)

	songRepository.On("CreateSongSectionHistory", mock.IsType(new(model.SongSectionHistory))).
		Return(nil).
		Times(sectionsCount)

	internalError := errors.New("internal error")
	songRepository.
		On(
			"GetSongSectionHistory",
			new([]model.SongSectionHistory),
			mock.IsType(uuid.UUID{}),
			model.RehearsalsProperty,
		).
		Return(internalError).
		Times(sectionsCount)

	// when
	errCode, updated := _uut.AddPerfectRehearsal(mockSong, songRepository)

	// then
	assert.False(t, updated)
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestAddPerfectRehearsal_WhenSectionsHaveZeroOccurrences_ShouldNotUpdateTheSong(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	progressProcessor := new(ProgressProcessorMock)
	_uut := processor.NewSongProcessor(progressProcessor)

	mockSong := &model.Song{
		Sections: []model.SongSection{
			{
				ID:         uuid.New(),
				Rehearsals: 23,
			},
			{
				ID: uuid.New(),
			},
		},
	}

	// when
	errCode, updated := _uut.AddPerfectRehearsal(mockSong, songRepository)

	// then
	assert.Nil(t, errCode)
	assert.False(t, updated)

	songRepository.AssertExpectations(t)
	progressProcessor.AssertExpectations(t)
}

func TestAddPerfectRehearsal_WhenSuccessful_ShouldUpdateSongAndSections(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	progressProcessor := new(ProgressProcessorMock)
	_uut := processor.NewSongProcessor(progressProcessor)

	mockSong := &model.Song{
		Sections: []model.SongSection{
			{
				ID:          uuid.New(),
				Rehearsals:  23,
				Occurrences: 2,
			},
			{
				ID:         uuid.New(),
				Rehearsals: 10,
			},
			{
				ID:          uuid.New(),
				Occurrences: 4,
			},
		},
	}

	oldSections := slices.Clone(mockSong.Sections)
	sectionsCount := len(mockSong.Sections)

	sectionsCountWithOcc := len(slices.DeleteFunc(slices.Clone(mockSong.Sections), func(section model.SongSection) bool {
		return section.Occurrences == 0
	}))

	songRepository.On("CreateSongSectionHistory", mock.IsType(new(model.SongSectionHistory))).
		Run(func(args mock.Arguments) {
			newHistory := args.Get(0).(*model.SongSectionHistory)
			assert.NotEmpty(t, newHistory.ID)
			assert.Equal(t, model.RehearsalsProperty, newHistory.Property)

			sections := slices.Clone(mockSong.Sections)
			sections = slices.DeleteFunc(sections, func(section model.SongSection) bool {
				return section.ID != newHistory.SongSectionID
			})

			assert.Equal(t, sections[0].Rehearsals, newHistory.From)
			assert.Equal(t, sections[0].Rehearsals+sections[0].Occurrences, newHistory.To)
		}).
		Return(nil).
		Times(sectionsCountWithOcc)

	history := &[]model.SongSectionHistory{}
	songRepository.
		On(
			"GetSongSectionHistory",
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
	errCode, updated := _uut.AddPerfectRehearsal(mockSong, songRepository)

	// then
	assert.Nil(t, errCode)
	assert.True(t, updated)

	var newSongRehearsals uint = 0
	var newSongProgress uint64 = 0
	for i, section := range mockSong.Sections {
		if section.Occurrences == 0 {
			assert.Equal(t, oldSections[i], section)
			continue
		}
		assert.Equal(t, oldSections[i].Rehearsals+section.Occurrences, section.Rehearsals)
		assert.Equal(t, newRehearsalScore, section.RehearsalsScore)
		assert.Equal(t, newProgress, section.Progress)
		newSongRehearsals += section.Rehearsals
		newSongProgress += section.Progress
	}
	assert.Equal(t, float64(newSongProgress)/float64(sectionsCount), mockSong.Progress)
	assert.Equal(t, float64(newSongRehearsals)/float64(sectionsCount), mockSong.Rehearsals)
	assert.WithinDuration(t, time.Now(), *mockSong.LastTimePlayed, 1*time.Minute)

	songRepository.AssertExpectations(t)
	progressProcessor.AssertExpectations(t)
}
