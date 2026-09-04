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
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateSongSection_WhenGetPartsFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	songPartRepository := new(repository.SongPartRepositoryMock)
	_uut := section.NewCreateSongSection(songSectionRepository, songPartRepository)

	request := requests.CreateSongSectionRequest{
		SongID:  uuid.New(),
		Name:    "Some Section",
		TypeID:  uuid.New(),
		PartIDs: []uuid.UUID{uuid.New()},
	}

	internalError := errors.New("internal error")
	songPartRepository.On("GetAllByIDs", new([]model.SongPart), request.PartIDs).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songSectionRepository.AssertExpectations(t)
	songPartRepository.AssertExpectations(t)
}

func TestCreateSongSection_WhenPartsLenIsNotTheSameAsRequest_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	songPartRepository := new(repository.SongPartRepositoryMock)
	_uut := section.NewCreateSongSection(songSectionRepository, songPartRepository)

	request := requests.CreateSongSectionRequest{
		SongID:  uuid.New(),
		Name:    "Some Section",
		TypeID:  uuid.New(),
		PartIDs: []uuid.UUID{uuid.New()},
	}

	songPartRepository.On("GetAllByIDs", new([]model.SongPart), request.PartIDs).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "parts not found", errCode.Error.Error())

	songSectionRepository.AssertExpectations(t)
	songPartRepository.AssertExpectations(t)
}

func TestCreateSongSection_WhenPartsDoNotBelongToSameSong_ShouldReturnConflictError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	songPartRepository := new(repository.SongPartRepositoryMock)
	_uut := section.NewCreateSongSection(songSectionRepository, songPartRepository)

	request := requests.CreateSongSectionRequest{
		SongID:  uuid.New(),
		Name:    "Some Section",
		TypeID:  uuid.New(),
		PartIDs: []uuid.UUID{uuid.New()},
	}

	mockParts := &[]model.SongPart{
		{ID: request.PartIDs[0], SongID: uuid.New()},
	}
	songPartRepository.On("GetAllByIDs", new([]model.SongPart), request.PartIDs).
		Return(nil, mockParts).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusConflict, errCode.Code)
	assert.Equal(t, "song part does not belong to the same song as the section", errCode.Error.Error())

	songSectionRepository.AssertExpectations(t)
	songPartRepository.AssertExpectations(t)
}

func TestCreateSongSection_WhenCountSectionsBySongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	_uut := section.NewCreateSongSection(songSectionRepository, nil)

	request := requests.CreateSongSectionRequest{
		SongID: uuid.New(),
		Name:   "Some Section",
		TypeID: uuid.New(),
	}

	internalError := errors.New("internal error")
	songSectionRepository.On("CountAllBySong", new(int64), request.SongID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songSectionRepository.AssertExpectations(t)
}

func TestCreateSongSection_WhenCreateSectionFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	_uut := section.NewCreateSongSection(songSectionRepository, nil)

	request := requests.CreateSongSectionRequest{
		SongID: uuid.New(),
		Name:   "Some Section",
		TypeID: uuid.New(),
	}

	sectionsCount := int64(5)
	songSectionRepository.On("CountAllBySong", new(int64), request.SongID).
		Return(nil, &sectionsCount).
		Once()

	internalError := errors.New("create error")
	songSectionRepository.On("Create", mock.IsType(new(model.SongSection))).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	songSectionRepository.AssertExpectations(t)
}

func TestCreateSongSection_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	tests := []struct {
		name    string
		request requests.CreateSongSectionRequest
	}{
		{
			"Without Parts",
			requests.CreateSongSectionRequest{
				SongID: uuid.New(),
				Name:   "Chorus",
				TypeID: uuid.New(),
			},
		},
		{
			"With Parts",
			requests.CreateSongSectionRequest{
				SongID:  uuid.New(),
				Name:    "Chorus",
				TypeID:  uuid.New(),
				PartIDs: []uuid.UUID{uuid.New(), uuid.New()},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			songSectionRepository := new(repository.SongSectionRepositoryMock)
			songPartRepository := new(repository.SongPartRepositoryMock)
			_uut := section.NewCreateSongSection(songSectionRepository, songPartRepository)

			if len(tt.request.PartIDs) > 0 {
				var mockParts []model.SongPart
				for _, pid := range tt.request.PartIDs {
					mockParts = append(mockParts, model.SongPart{ID: pid, SongID: tt.request.SongID})
				}
				songPartRepository.On("GetAllByIDs", new([]model.SongPart), tt.request.PartIDs).
					Return(nil, &mockParts).
					Once()
			}

			sectionsCount := int64(3)
			songSectionRepository.On("CountAllBySong", new(int64), tt.request.SongID).
				Return(nil, &sectionsCount).
				Once()

			songSectionRepository.On("Create", mock.IsType(new(model.SongSection))).
				Run(func(args mock.Arguments) {
					newSection := args.Get(0).(*model.SongSection)
					assertCreatedSongSection(t, tt.request, *newSection, sectionsCount)
				}).
				Return(nil).
				Once()

			// when
			errCode := _uut.Handle(tt.request)

			// then
			assert.Nil(t, errCode)

			songSectionRepository.AssertExpectations(t)
			songPartRepository.AssertExpectations(t)
		})
	}
}

func assertCreatedSongSection(
	t *testing.T,
	request requests.CreateSongSectionRequest,
	section model.SongSection,
	sectionsCount int64,
) {
	assert.NotEmpty(t, section.ID)
	assert.Equal(t, request.Name, section.Name)
	assert.Equal(t, request.TypeID, section.SongSectionTypeID)
	assert.Equal(t, request.SongID, section.SongID)
	assert.Equal(t, uint(sectionsCount), section.Order)

	assert.Len(t, section.SectionParts, len(request.PartIDs))
	for i, sp := range section.SectionParts {
		assert.Equal(t, request.PartIDs[i], sp.PartID)
		assert.Equal(t, uint(i), sp.Order)
	}
}
