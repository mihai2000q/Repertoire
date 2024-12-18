package section

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"math"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/song/section"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"testing"
)

func TestCreateSongSection_WhenCountSectionsBySongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := section.NewCreateSongSection(songRepository)

	request := requests.CreateSongSectionRequest{
		SongID: uuid.New(),
		Name:   "Some Artist",
		TypeID: uuid.New(),
	}

	internalError := errors.New("internal error")
	songRepository.On("CountSectionsBySong", new(int64), request.SongID).
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

func TestCreateSongSection_WhenCreateSectionFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := section.NewCreateSongSection(songRepository)

	request := requests.CreateSongSectionRequest{
		SongID: uuid.New(),
		Name:   "Some Artist",
		TypeID: uuid.New(),
	}

	expectedCount := &[]int64{20}[0]
	songRepository.On("CountSectionsBySong", mock.IsType(expectedCount), request.SongID).
		Return(nil, expectedCount).
		Once()

	internalError := errors.New("internal error")
	songRepository.On("CreateSection", mock.IsType(new(model.SongSection))).
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

func TestCreateSongSection_WhenGetSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := section.NewCreateSongSection(songRepository)

	request := requests.CreateSongSectionRequest{
		SongID: uuid.New(),
		Name:   "Some Artist",
		TypeID: uuid.New(),
	}

	expectedCount := &[]int64{20}[0]
	songRepository.On("CountSectionsBySong", mock.IsType(expectedCount), request.SongID).
		Return(nil, expectedCount).
		Once()

	songRepository.On("CreateSection", mock.IsType(new(model.SongSection))).
		Return(nil).
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

	songRepository.AssertExpectations(t)
}

func TestCreateSongSection_WhenUpdateSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := section.NewCreateSongSection(songRepository)

	request := requests.CreateSongSectionRequest{
		SongID: uuid.New(),
		Name:   "Some Artist",
		TypeID: uuid.New(),
	}

	expectedCount := &[]int64{20}[0]
	songRepository.On("CountSectionsBySong", mock.IsType(expectedCount), request.SongID).
		Return(nil, expectedCount).
		Once()

	songRepository.On("CreateSection", mock.IsType(new(model.SongSection))).
		Return(nil).
		Once()

	songRepository.On("Get", new(model.Song), request.SongID).
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

	songRepository.AssertExpectations(t)
}

func TestCreateSongSection_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	tests := []struct {
		name                   string
		song                   model.Song
		expectedSectionsCount  int64
		expectedSongConfidence float64
		expectedSongRehearsals float64
		expectedSongProgress   float64
	}{
		{
			"1 - When there are no precedent sections",
			model.Song{},
			0,
			0,
			0,
			0,
		},
		{
			"1 - When there are precedent sections, but with stats 0",
			model.Song{},
			2,
			0,
			0,
			0,
		},
		{
			"1 - When there are precedent sections with stats",
			model.Song{
				Confidence: 50,
				Rehearsals: 10,
				Progress:   55,
			},
			2,
			33,
			7,
			37,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			songRepository := new(repository.SongRepositoryMock)
			_uut := section.NewCreateSongSection(songRepository)

			request := requests.CreateSongSectionRequest{
				SongID: uuid.New(),
				Name:   "Some Artist",
				TypeID: uuid.New(),
			}

			songRepository.On("CountSectionsBySong", mock.IsType(&tt.expectedSectionsCount), request.SongID).
				Return(nil, &tt.expectedSectionsCount).
				Once()

			songRepository.On("CreateSection", mock.IsType(new(model.SongSection))).
				Run(func(args mock.Arguments) {
					newSection := args.Get(0).(*model.SongSection)
					assertCreatedSongSection(t, request, *newSection, tt.expectedSectionsCount)
				}).
				Return(nil).
				Once()

			songRepository.On("Get", new(model.Song), request.SongID).
				Return(nil, &tt.song).
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

			// when
			errCode := _uut.Handle(request)

			// then
			assert.Nil(t, errCode)

			songRepository.AssertExpectations(t)
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
	assert.Equal(t, request.SongID, section.SongID)
}
