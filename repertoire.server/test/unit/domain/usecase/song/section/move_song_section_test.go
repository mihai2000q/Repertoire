package section

import (
	"cmp"
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

func TestMoveSongSection_WhenGetSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := section.NewMoveSongSection(songRepository)

	request := requests.MoveSongSectionRequest{
		ID:     uuid.New(),
		OverID: uuid.New(),
		SongID: uuid.New(),
	}

	// given - mocking
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

func TestMoveSongSection_WhenSongIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := section.NewMoveSongSection(songRepository)

	request := requests.MoveSongSectionRequest{
		ID:     uuid.New(),
		OverID: uuid.New(),
		SongID: uuid.New(),
	}

	// given - mocking
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

func TestMoveSongSection_WhenSectionIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := section.NewMoveSongSection(songRepository)

	song := &model.Song{ID: uuid.New()}

	request := requests.MoveSongSectionRequest{
		ID:     uuid.New(),
		OverID: uuid.New(),
		SongID: song.ID,
	}

	// given - mocking
	songRepository.On("GetWithSections", new(model.Song), request.SongID).
		Return(nil, song).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "section not found", errCode.Error.Error())

	songRepository.AssertExpectations(t)
}

func TestMoveSongSection_WhenOverSectionIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := section.NewMoveSongSection(songRepository)

	song := &model.Song{
		ID: uuid.New(),
		Sections: []model.SongSection{
			{ID: uuid.New(), Order: 0},
		},
	}

	request := requests.MoveSongSectionRequest{
		ID:     song.Sections[0].ID,
		OverID: uuid.New(),
		SongID: song.ID,
	}

	// given - mocking
	songRepository.On("GetWithSections", new(model.Song), request.SongID).
		Return(nil, song).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "over section not found", errCode.Error.Error())

	songRepository.AssertExpectations(t)
}

func TestMoveSongSection_WhenUpdateFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := section.NewMoveSongSection(songRepository)

	song := &model.Song{
		ID: uuid.New(),
		Sections: []model.SongSection{
			{ID: uuid.New(), Order: 0},
			{ID: uuid.New(), Order: 1},
		},
	}

	request := requests.MoveSongSectionRequest{
		ID:     song.Sections[0].ID,
		OverID: song.Sections[1].ID,
		SongID: song.ID,
	}

	// given - mocking
	songRepository.On("GetWithSections", new(model.Song), request.SongID).
		Return(nil, song).
		Once()

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
}

func TestMoveSongSection_WhenIsValid_ShouldNotReturnAnyError(t *testing.T) {
	tests := []struct {
		name      string
		song      *model.Song
		index     uint
		overIndex uint
	}{
		{
			"Use case 1",
			&model.Song{
				ID: uuid.New(),
				Sections: []model.SongSection{
					{ID: uuid.New(), Order: 0},
					{ID: uuid.New(), Order: 1},
					{ID: uuid.New(), Order: 2},
					{ID: uuid.New(), Order: 3},
					{ID: uuid.New(), Order: 4},
				},
			},
			1,
			3,
		},
		{
			"Use case 2",
			&model.Song{
				ID: uuid.New(),
				Sections: []model.SongSection{
					{ID: uuid.New(), Order: 0},
					{ID: uuid.New(), Order: 1},
					{ID: uuid.New(), Order: 2},
					{ID: uuid.New(), Order: 3},
					{ID: uuid.New(), Order: 4},
				},
			},
			3,
			1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			songRepository := new(repository.SongRepositoryMock)
			_uut := section.NewMoveSongSection(songRepository)

			request := requests.MoveSongSectionRequest{
				ID:     tt.song.Sections[tt.index].ID,
				OverID: tt.song.Sections[tt.overIndex].ID,
				SongID: tt.song.ID,
			}

			// given - mocking
			songRepository.On("GetWithSections", new(model.Song), request.SongID).
				Return(nil, tt.song).
				Once()

			songRepository.On("UpdateWithAssociations", mock.IsType(new(model.Song))).
				Run(func(args mock.Arguments) {
					song := args.Get(0).(*model.Song)
					sections := slices.Clone(song.Sections)
					slices.SortFunc(sections, func(a, b model.SongSection) int {
						return cmp.Compare(a.Order, b.Order)
					})
					if tt.index < tt.overIndex {
						assert.Equal(t, sections[tt.overIndex-1].ID, request.OverID)
					} else if tt.index > tt.overIndex {
						assert.Equal(t, sections[tt.overIndex+1].ID, request.OverID)
					}
					assert.Equal(t, sections[tt.overIndex].ID, request.ID)
					for i, s := range sections {
						assert.Equal(t, uint(i), s.Order)
					}
				}).
				Return(nil).
				Once()

			// when
			errCode := _uut.Handle(request)

			// then
			assert.Nil(t, errCode)

			songRepository.AssertExpectations(t)
		})
	}
}
