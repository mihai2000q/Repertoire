package section

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/song/section"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAllSongSections_WhenGetAllBySongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	_uut := section.NewGetAllSongSections(songSectionRepository)

	request := requests.GetSongSectionsRequest{
		SongID: uuid.New(),
	}

	internalError := errors.New("internal error")
	songSectionRepository.On("GetAllBySong", new([]model.SongSection), request.SongID).
		Return(internalError).
		Once()

	// when
	sections, errCode := _uut.Handle(request)

	// then
	assert.Empty(t, sections)
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songSectionRepository.AssertExpectations(t)
}

func TestGetAllSongSections_WhenSuccessful_ShouldReturnSections(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	_uut := section.NewGetAllSongSections(songSectionRepository)

	request := requests.GetSongSectionsRequest{
		SongID: uuid.New(),
	}

	expectedSections := []model.SongSection{
		{ID: uuid.New(), Name: "Verse"},
		{ID: uuid.New(), Name: "Chorus"},
		{ID: uuid.New(), Name: "Bridge"},
	}

	songSectionRepository.On("GetAllBySong", new([]model.SongSection), request.SongID).
		Return(nil, &expectedSections).
		Once()

	// when
	sections, errCode := _uut.Handle(request)

	// then
	assert.Nil(t, errCode)
	assert.Equal(t, expectedSections, sections)

	songSectionRepository.AssertExpectations(t)
}
