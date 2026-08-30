package part

import (
	"cmp"
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/song/part"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"slices"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMoveSongPartInSong_WhenGetSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := part.NewMoveSongPartInSong(songRepository)

	request := requests.MoveSongPartInSongRequest{
		ID:     uuid.New(),
		OverID: uuid.New(),
		SongID: uuid.New(),
	}

	// given - mocking
	internalError := errors.New("internal error")
	songRepository.On("GetWithParts", new(model.Song), request.SongID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestMoveSongPartInSong_WhenSongIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := part.NewMoveSongPartInSong(songRepository)

	request := requests.MoveSongPartInSongRequest{
		ID:     uuid.New(),
		OverID: uuid.New(),
		SongID: uuid.New(),
	}

	// given - mocking
	songRepository.On("GetWithParts", new(model.Song), request.SongID).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "song not found", errCode.Error.Error())

	songRepository.AssertExpectations(t)
}

func TestMoveSongPartInSong_WhenPartIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := part.NewMoveSongPartInSong(songRepository)

	song := &model.Song{ID: uuid.New()}

	request := requests.MoveSongPartInSongRequest{
		ID:     uuid.New(),
		OverID: uuid.New(),
		SongID: song.ID,
	}

	// given - mocking
	songRepository.On("GetWithParts", new(model.Song), request.SongID).
		Return(nil, song).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "part not found", errCode.Error.Error())

	songRepository.AssertExpectations(t)
}

func TestMoveSongPartInSong_WhenOverPartIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := part.NewMoveSongPartInSong(songRepository)

	song := &model.Song{
		ID: uuid.New(),
		Parts: []model.SongPart{
			{ID: uuid.New(), SongOrder: 0},
		},
	}

	request := requests.MoveSongPartInSongRequest{
		ID:     song.Parts[0].ID,
		OverID: uuid.New(),
		SongID: song.ID,
	}

	// given - mocking
	songRepository.On("GetWithParts", new(model.Song), request.SongID).
		Return(nil, song).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "over part not found", errCode.Error.Error())

	songRepository.AssertExpectations(t)
}

func TestMoveSongPartInSong_WhenUpdateFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := part.NewMoveSongPartInSong(songRepository)

	song := &model.Song{
		ID: uuid.New(),
		Parts: []model.SongPart{
			{ID: uuid.New(), SongOrder: 0},
			{ID: uuid.New(), SongOrder: 1},
		},
	}

	request := requests.MoveSongPartInSongRequest{
		ID:     song.Parts[0].ID,
		OverID: song.Parts[1].ID,
		SongID: song.ID,
	}

	// given - mocking
	songRepository.On("GetWithParts", new(model.Song), request.SongID).
		Return(nil, song).
		Once()

	internalError := errors.New("internal error")
	songRepository.On("UpdateWithAssociations", mock.IsType(new(model.Song))).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songRepository.AssertExpectations(t)
}

func TestMoveSongPartInSong_WhenIsValid_ShouldNotReturnAnyError(t *testing.T) {
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
				Parts: []model.SongPart{
					{ID: uuid.New(), SongOrder: 0},
					{ID: uuid.New(), SongOrder: 1},
					{ID: uuid.New(), SongOrder: 2},
					{ID: uuid.New(), SongOrder: 3},
					{ID: uuid.New(), SongOrder: 4},
				},
			},
			1,
			3,
		},
		{
			"Use case 2",
			&model.Song{
				ID: uuid.New(),
				Parts: []model.SongPart{
					{ID: uuid.New(), SongOrder: 0},
					{ID: uuid.New(), SongOrder: 1},
					{ID: uuid.New(), SongOrder: 2},
					{ID: uuid.New(), SongOrder: 3},
					{ID: uuid.New(), SongOrder: 4},
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
			_uut := part.NewMoveSongPartInSong(songRepository)

			request := requests.MoveSongPartInSongRequest{
				ID:     tt.song.Parts[tt.index].ID,
				OverID: tt.song.Parts[tt.overIndex].ID,
				SongID: tt.song.ID,
			}

			// given - mocking
			songRepository.On("GetWithParts", new(model.Song), request.SongID).
				Return(nil, tt.song).
				Once()

			songRepository.On("UpdateWithAssociations", mock.IsType(new(model.Song))).
				Run(func(args mock.Arguments) {
					song := args.Get(0).(*model.Song)
					parts := slices.Clone(song.Parts)
					slices.SortFunc(parts, func(a, b model.SongPart) int {
						return cmp.Compare(a.SongOrder, b.SongOrder)
					})
					if tt.index < tt.overIndex {
						assert.Equal(t, parts[tt.overIndex-1].ID, request.OverID)
					} else if tt.index > tt.overIndex {
						assert.Equal(t, parts[tt.overIndex+1].ID, request.OverID)
					}
					assert.Equal(t, parts[tt.overIndex].ID, request.ID)
					for i, s := range parts {
						assert.Equal(t, uint(i), s.SongOrder)
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
