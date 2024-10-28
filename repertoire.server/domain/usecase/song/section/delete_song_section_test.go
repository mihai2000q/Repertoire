package section

import (
	"errors"
	"net/http"
	"repertoire/data/repository"
	"repertoire/model"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteSongSection_WhenGetSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := DeleteSongSection{
		songRepository: songRepository,
	}
	id := uuid.New()
	songID := uuid.New()

	internalError := errors.New("internal error")
	songRepository.On("Get", new(model.Song), songID).Return(internalError).Once()

	// when
	errCode := _uut.Handle(id, songID)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestDeleteSongSection_WhenDeleteSectionFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := DeleteSongSection{
		songRepository: songRepository,
	}
	id := uuid.New()
	songID := uuid.New()

	// given - mocking
	song := &model.Song{
		ID: songID,
		Sections: []model.SongSection{
			{ID: id, Order: 0},
			{ID: uuid.New(), Order: 1},
		},
	}
	songRepository.On("Get", new(model.Song), songID).Return(nil, song).Once()

	internalError := errors.New("internal error")
	songRepository.On("DeleteSection", id).Return(internalError).Once()

	// when
	errCode := _uut.Handle(id, songID)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestDeleteSongSection_WhenUpdateSectionFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := DeleteSongSection{
		songRepository: songRepository,
	}
	id := uuid.New()
	songID := uuid.New()

	// given - mocking
	song := &model.Song{
		ID: songID,
		Sections: []model.SongSection{
			{ID: id, Order: 0},
			{ID: uuid.New(), Order: 1},
		},
	}
	songRepository.On("Get", new(model.Song), songID).Return(nil, song).Once()

	songRepository.On("DeleteSection", id).Return(nil).Once()

	internalError := errors.New("internal error")
	songRepository.On("UpdateSection", mock.IsType(&song.Sections[0])).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id, songID)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestDeleteSongSection_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	initialSections := []model.SongSection{
		{ID: uuid.New(), Order: 0},
		{ID: uuid.New(), Order: 1},
		{ID: uuid.New(), Order: 2},
		{ID: uuid.New(), Order: 3},
		{ID: uuid.New(), Order: 4},
	}

	tests := []struct {
		name                  string
		sections              []model.SongSection
		id                    uuid.UUID
		indexesThatGetUpdated []int
	}{
		{
			"Use case 1",
			initialSections,
			initialSections[0].ID,
			[]int{1, 2, 3, 4},
		},
		{
			"Use case 2",
			initialSections,
			initialSections[2].ID,
			[]int{3, 4},
		},
		{
			"Use case 3",
			initialSections,
			initialSections[4].ID,
			[]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			songRepository := new(repository.SongRepositoryMock)
			_uut := DeleteSongSection{
				songRepository: songRepository,
			}
			songID := uuid.New()

			// given - mocking
			song := &model.Song{ID: songID, Sections: tt.sections}
			songRepository.On("Get", new(model.Song), songID).Return(nil, song).Once()

			songRepository.On("DeleteSection", tt.id).Return(nil).Once()

			for _, i := range tt.indexesThatGetUpdated {
				songRepository.On("UpdateSection", mock.IsType(&song.Sections[i])).
					Run(func(args mock.Arguments) {
						newSection := args.Get(0).(*model.SongSection)
						assert.Equal(t, song.Sections[i].Order-1, newSection.Order)
					}).
					Return(nil).
					Once()
			}

			// when
			errCode := _uut.Handle(tt.id, songID)

			// then
			assert.Nil(t, errCode)

			songRepository.AssertExpectations(t)
		})
	}
}
