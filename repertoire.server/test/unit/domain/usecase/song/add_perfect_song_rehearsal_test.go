package song

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/song"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/domain/processor"
	"slices"
	"testing"
	"time"
)

func TestAddPerfectSongRehearsal_WhenGetSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewAddPerfectSongRehearsal(songRepository, nil)

	request := requests.AddPerfectSongRehearsalRequest{
		ID: uuid.New(),
	}

	internalError := errors.New("internal error")
	songRepository.On("GetWithSections", new(model.Song), request.ID).Return(internalError).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestAddPerfectSongRehearsal_WhenSongIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewAddPerfectSongRehearsal(songRepository, nil)

	request := requests.AddPerfectSongRehearsalRequest{
		ID: uuid.New(),
	}

	songRepository.On("GetWithSections", new(model.Song), request.ID).Return(nil).Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song not found", errCode.Error.Error())

	songRepository.AssertExpectations(t)
}

func TestAddPerfectSongRehearsal_WhenCreateHistoryFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewAddPerfectSongRehearsal(songRepository, nil)

	request := requests.AddPerfectSongRehearsalRequest{
		ID: uuid.New(),
	}

	mockSong := &model.Song{
		ID:       uuid.New(),
		Sections: []model.SongSection{{ID: uuid.New(), Occurrences: 2}},
	}
	songRepository.On("GetWithSections", new(model.Song), request.ID).
		Return(nil, mockSong).
		Once()

	sectionsCount := len(mockSong.Sections)

	internalError := errors.New("internal error")
	songRepository.On("CreateSongSectionHistory", mock.IsType(new(model.SongSectionHistory))).
		Return(internalError).
		Times(sectionsCount)

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestAddPerfectSongRehearsal_WhenGetHistoryFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := song.NewAddPerfectSongRehearsal(songRepository, nil)

	request := requests.AddPerfectSongRehearsalRequest{
		ID: uuid.New(),
	}

	mockSong := &model.Song{
		ID:       uuid.New(),
		Sections: []model.SongSection{{ID: uuid.New(), Occurrences: 2}},
	}
	songRepository.On("GetWithSections", new(model.Song), request.ID).
		Return(nil, mockSong).
		Once()

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
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestAddPerfectSongRehearsal_WhenUpdateFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	progressProcessor := new(processor.ProgressProcessorMock)
	_uut := song.NewAddPerfectSongRehearsal(songRepository, progressProcessor)

	request := requests.AddPerfectSongRehearsalRequest{
		ID: uuid.New(),
	}

	mockSong := &model.Song{
		ID:       uuid.New(),
		Sections: []model.SongSection{{ID: uuid.New(), Occurrences: 2}},
	}
	songRepository.On("GetWithSections", new(model.Song), request.ID).
		Return(nil, mockSong).
		Once()

	sectionsCount := len(mockSong.Sections)

	songRepository.On("CreateSongSectionHistory", mock.IsType(new(model.SongSectionHistory))).
		Return(nil).
		Times(sectionsCount)

	history := &[]model.SongSectionHistory{}
	songRepository.
		On(
			"GetSongSectionHistory",
			new([]model.SongSectionHistory),
			mock.IsType(uuid.UUID{}),
			model.RehearsalsProperty,
		).
		Return(nil, history).
		Times(sectionsCount)

	var newRehearsalScore uint64 = 23
	progressProcessor.On("ComputeRehearsalsScore", *history).Return(newRehearsalScore).Times(sectionsCount)
	var newProgress uint64 = 123
	progressProcessor.On("ComputeProgress", mock.IsType(model.SongSection{})).
		Return(newProgress).
		Times(sectionsCount)

	internalError := errors.New("internal error")
	songRepository.On("UpdateWithAssociations", mock.IsType(new(model.Song))).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
	progressProcessor.AssertExpectations(t)
}

func TestAddPerfectSongRehearsal_WhenSectionsHaveZeroOccurrences_ShouldNotUpdateTheSong(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	progressProcessor := new(processor.ProgressProcessorMock)
	_uut := song.NewAddPerfectSongRehearsal(songRepository, nil)

	request := requests.AddPerfectSongRehearsalRequest{
		ID: uuid.New(),
	}

	mockSong := model.Song{
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
	songRepository.On("GetWithSections", new(model.Song), request.ID).
		Return(nil, &mockSong).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)

	songRepository.AssertExpectations(t)
	progressProcessor.AssertExpectations(t)
}

func TestAddPerfectSongRehearsal_WhenSuccessful_ShouldUpdateSongAndSections(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	progressProcessor := new(processor.ProgressProcessorMock)
	_uut := song.NewAddPerfectSongRehearsal(songRepository, progressProcessor)

	request := requests.AddPerfectSongRehearsalRequest{
		ID: uuid.New(),
	}

	mockSong := model.Song{
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
	songRepository.On("GetWithSections", new(model.Song), request.ID).
		Return(nil, &mockSong).
		Once()

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

	songRepository.On("UpdateWithAssociations", mock.IsType(new(model.Song))).
		Run(func(args mock.Arguments) {
			newSong := args.Get(0).(*model.Song)

			var newSongRehearsals uint = 0
			var newSongProgress uint64 = 0
			for i, section := range newSong.Sections {
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
			assert.Equal(t, float64(newSongProgress)/float64(sectionsCount), newSong.Progress)
			assert.Equal(t, float64(newSongRehearsals)/float64(sectionsCount), newSong.Rehearsals)
			assert.WithinDuration(t, time.Now(), *newSong.LastTimePlayed, 1*time.Minute)
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)

	songRepository.AssertExpectations(t)
	progressProcessor.AssertExpectations(t)
}
