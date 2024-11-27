package section

import (
	"errors"
	"net/http"
	section2 "repertoire/server/domain/usecase/song/section"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"slices"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteSongSection_WhenGetSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := section2.DeleteSongSection{
		songRepository: songRepository,
	}
	id := uuid.New()
	songID := uuid.New()

	internalError := errors.New("internal error")
	songRepository.On("GetWithSections", new(model.Song), songID).
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

func TestDeleteSongSection_WhenSongIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := section2.DeleteSongSection{
		songRepository: songRepository,
	}
	id := uuid.New()
	songID := uuid.New()

	songRepository.On("GetWithSections", new(model.Song), songID).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(id, songID)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song not found", errCode.Error.Error())

	songRepository.AssertExpectations(t)
}

func TestDeleteSongSection_WhenUpdateSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := section2.DeleteSongSection{
		songRepository: songRepository,
	}
	id := uuid.New()
	songID := uuid.New()

	// given - mocking
	song := &model.Song{
		ID: songID,
		Sections: []model.SongSection{
			{ID: id, Order: 0},
		},
	}
	songRepository.On("GetWithSections", new(model.Song), songID).
		Return(nil, song).
		Once()

	internalError := errors.New("internal error")
	songRepository.On("UpdateWithAssociations", mock.IsType(song)).
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

func TestDeleteSongSection_WhenDeleteSectionFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := section2.DeleteSongSection{
		songRepository: songRepository,
	}
	id := uuid.New()
	songID := uuid.New()

	// given - mocking
	song := &model.Song{
		ID: songID,
		Sections: []model.SongSection{
			{ID: id, Order: 0},
		},
	}
	songRepository.On("GetWithSections", new(model.Song), songID).
		Return(nil, song).
		Once()

	songRepository.On("UpdateWithAssociations", mock.IsType(song)).
		Return(nil).
		Once()

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

func TestDeleteSongSection_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := section2.DeleteSongSection{
		songRepository: songRepository,
	}
	id := uuid.New()
	songID := uuid.New()

	// given - mocking
	song := &model.Song{
		ID: songID,
		Sections: []model.SongSection{
			{ID: uuid.New(), Order: 0},
			{ID: id, Order: 1},
			{ID: uuid.New(), Order: 2},
			{ID: uuid.New(), Order: 3},
			{ID: uuid.New(), Order: 4},
		},
	}
	songRepository.On("GetWithSections", new(model.Song), songID).
		Return(nil, song).
		Once()

	songRepository.On("UpdateWithAssociations", mock.IsType(song)).
		Run(func(args mock.Arguments) {
			song := args.Get(0).(*model.Song)
			sections := slices.Clone(song.Sections)

			sections = slices.DeleteFunc(sections, func(a model.SongSection) bool {
				return a.ID == id
			})
			for i, section := range sections {
				assert.Equal(t, uint(i), section.Order)
			}
		}).
		Return(nil).
		Once()

	songRepository.On("DeleteSection", id).Return(nil).Once()

	// when
	errCode := _uut.Handle(id, songID)

	// then
	assert.Nil(t, errCode)

	songRepository.AssertExpectations(t)
}
