package section

import (
	"errors"
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

func TestUpdateSongSectionsPartialOccurrences_WhenGetSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := section.NewUpdateSongSectionsPartialOccurrences(songRepository)

	request := requests.UpdateSongSectionsPartialOccurrencesRequest{
		SongID:   uuid.New(),
		Sections: []requests.UpdateSectionPartialOccurrencesRequest{{ID: uuid.New()}},
	}

	internalError := errors.New("internal error")
	songRepository.On("GetWithSections", new(model.Song), request.SongID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestUpdateSongSectionsPartialOccurrences_WhenSongIsEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := section.NewUpdateSongSectionsPartialOccurrences(songRepository)

	request := requests.UpdateSongSectionsPartialOccurrencesRequest{
		SongID:   uuid.New(),
		Sections: []requests.UpdateSectionPartialOccurrencesRequest{{ID: uuid.New()}},
	}

	songRepository.On("GetWithSections", new(model.Song), request.SongID).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song not found", errCode.Error.Error())

	songRepository.AssertExpectations(t)
}

func TestUpdateSongSectionsPartialOccurrences_WhenUpdateSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := section.NewUpdateSongSectionsPartialOccurrences(songRepository)

	request := requests.UpdateSongSectionsPartialOccurrencesRequest{
		SongID:   uuid.New(),
		Sections: []requests.UpdateSectionPartialOccurrencesRequest{{ID: uuid.New()}},
	}

	mockSong := &model.Song{
		ID:       request.SongID,
		Title:    "Some Song",
		Sections: []model.SongSection{{ID: request.Sections[0].ID}},
	}
	songRepository.On("GetWithSections", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()

	internalError := errors.New("internal error")
	songRepository.On("UpdateWithAssociations", mock.IsType(mockSong)).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestUpdateSongSectionsPartialOccurrences_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := section.NewUpdateSongSectionsPartialOccurrences(songRepository)

	request := requests.UpdateSongSectionsPartialOccurrencesRequest{
		SongID: uuid.New(),
		Sections: []requests.UpdateSectionPartialOccurrencesRequest{
			{
				ID:                 uuid.New(),
				PartialOccurrences: uint(1),
			},
			{
				ID:                 uuid.New(),
				PartialOccurrences: uint(2),
			},
			{
				ID:                 uuid.New(),
				PartialOccurrences: uint(3),
			},
			{
				ID:                 uuid.New(),
				PartialOccurrences: uint(23), // ignored
			},
		},
	}

	mockSong := &model.Song{
		ID:    request.SongID,
		Title: "Some Song",
		Sections: []model.SongSection{
			{ID: request.Sections[2].ID},
			{ID: request.Sections[0].ID},
			{ID: request.Sections[1].ID},
			{ID: uuid.New()}, // also ignored
		},
	}
	songRepository.On("GetWithSections", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()

	songRepository.On("UpdateWithAssociations", mock.IsType(mockSong)).
		Run(func(args mock.Arguments) {
			newSong := args.Get(0).(*model.Song)
			for _, newSection := range newSong.Sections {
				requestSections := slices.Clone(request.Sections)
				sectionOccurrence := slices.DeleteFunc(requestSections, func(r requests.UpdateSectionPartialOccurrencesRequest) bool {
					return r.ID != newSection.ID
				})
				if len(sectionOccurrence) == 0 {
					continue
				}
				assert.Equal(t, sectionOccurrence[0].ID, newSection.ID)
				assert.Equal(t, sectionOccurrence[0].PartialOccurrences, newSection.PartialOccurrences)
			}
		}).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)

	songRepository.AssertExpectations(t)
}
