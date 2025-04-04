package section

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/song/section"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"testing"
)

func TestUpdateAllSongSections_WhenGetSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := section.NewUpdateAllSongSections(songRepository)

	request := requests.UpdateAllSongSectionsRequest{
		SongID: uuid.New(),
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

func TestUpdateAllSongSections_WhenSettingsAreEmpty_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := section.NewUpdateAllSongSections(songRepository)

	request := requests.UpdateAllSongSectionsRequest{
		SongID: uuid.New(),
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

func TestUpdateAllSongSections_WhenUpdateFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := section.NewUpdateAllSongSections(songRepository)

	request := requests.UpdateAllSongSectionsRequest{
		SongID: uuid.New(),
	}

	mockSong := &model.Song{
		ID: request.SongID,
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

func TestUpdateAllSongSections_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := section.NewUpdateAllSongSections(songRepository)

	request := requests.UpdateAllSongSectionsRequest{
		SongID:       uuid.New(),
		InstrumentID: &[]uuid.UUID{uuid.New()}[0],
		BandMemberID: &[]uuid.UUID{uuid.New()}[0],
	}

	mockSong := &model.Song{
		ID: request.SongID,
		Sections: []model.SongSection{
			{ID: uuid.New()},
			{ID: uuid.New()},
			{ID: uuid.New()},
		},
	}

	songRepository.On("GetWithSections", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()

	songRepository.On("UpdateWithAssociations", mock.IsType(mockSong)).
		Run(func(args mock.Arguments) {
			newSong := args.Get(0).(*model.Song)
			for _, songSection := range newSong.Sections {
				if request.InstrumentID != nil {
					assert.Equal(t, request.InstrumentID, songSection.InstrumentID)
				}
				if request.BandMemberID != nil {
					assert.Equal(t, request.BandMemberID, songSection.BandMemberID)
				}
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
