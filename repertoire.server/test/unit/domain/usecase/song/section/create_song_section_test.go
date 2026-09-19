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

	partID := uuid.New()
	request := requests.CreateSongSectionRequest{
		SongID: uuid.New(),
		Name:   "Some Section",
		TypeID: uuid.New(),
		Parts: []requests.CreateSongSectionPartRequest{
			{PartID: &partID},
		},
	}

	internalError := errors.New("internal error")
	songPartRepository.On("GetAllByIDs", new([]model.SongPart), []uuid.UUID{partID}).
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

	partID := uuid.New()
	request := requests.CreateSongSectionRequest{
		SongID: uuid.New(),
		Name:   "Some Section",
		TypeID: uuid.New(),
		Parts: []requests.CreateSongSectionPartRequest{
			{PartID: &partID},
		},
	}

	songPartRepository.On("GetAllByIDs", new([]model.SongPart), []uuid.UUID{partID}).
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

	partID := uuid.New()
	request := requests.CreateSongSectionRequest{
		SongID: uuid.New(),
		Name:   "Some Section",
		TypeID: uuid.New(),
		Parts: []requests.CreateSongSectionPartRequest{
			{PartID: &partID},
		},
	}

	mockParts := &[]model.SongPart{
		{ID: partID, SongID: uuid.New()},
	}
	songPartRepository.On("GetAllByIDs", new([]model.SongPart), []uuid.UUID{partID}).
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

func TestCreateSongSection_WhenCountPartsBySongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	songPartRepository := new(repository.SongPartRepositoryMock)
	_uut := section.NewCreateSongSection(songSectionRepository, songPartRepository)

	request := requests.CreateSongSectionRequest{
		SongID: uuid.New(),
		Name:   "Some Section",
		TypeID: uuid.New(),
		Parts: []requests.CreateSongSectionPartRequest{
			{NewPart: &requests.CreateNewSongPartRequest{Name: "Verse Riff"}},
		},
	}

	sectionsCount := int64(2)
	songSectionRepository.On("CountAllBySong", new(int64), request.SongID).
		Return(nil, &sectionsCount).
		Once()

	internalError := errors.New("internal error")
	songPartRepository.On("CountAllBySong", new(int64), request.SongID).
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
	existingPartID1 := uuid.New()
	existingPartID2 := uuid.New()
	bandMemberID := uuid.New()
	instrumentID := uuid.New()

	tests := []struct {
		name       string
		request    requests.CreateSongSectionRequest
		partsCount int64
	}{
		{
			name: "Without Parts",
			request: requests.CreateSongSectionRequest{
				SongID: uuid.New(),
				Name:   "Chorus",
				TypeID: uuid.New(),
			},
		},
		{
			name: "With Existing Parts",
			request: requests.CreateSongSectionRequest{
				SongID: uuid.New(),
				Name:   "Chorus",
				TypeID: uuid.New(),
				Parts: []requests.CreateSongSectionPartRequest{
					{PartID: &existingPartID1},
					{PartID: &existingPartID2, BandMemberID: &bandMemberID},
				},
			},
		},
		{
			name: "With New Parts",
			request: requests.CreateSongSectionRequest{
				SongID: uuid.New(),
				Name:   "Chorus",
				TypeID: uuid.New(),
				Parts: []requests.CreateSongSectionPartRequest{
					{NewPart: &requests.CreateNewSongPartRequest{Name: "Verse Riff", InstrumentID: &instrumentID}},
					{NewPart: &requests.CreateNewSongPartRequest{Name: "Chorus Riff"}, BandMemberID: &bandMemberID},
				},
			},
			partsCount: 3,
		},
		{
			name: "With Mixed Parts",
			request: requests.CreateSongSectionRequest{
				SongID: uuid.New(),
				Name:   "Chorus",
				TypeID: uuid.New(),
				Parts: []requests.CreateSongSectionPartRequest{
					{PartID: &existingPartID1},
					{NewPart: &requests.CreateNewSongPartRequest{Name: "Verse Riff"}},
				},
			},
			partsCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			songSectionRepository := new(repository.SongSectionRepositoryMock)
			songPartRepository := new(repository.SongPartRepositoryMock)
			_uut := section.NewCreateSongSection(songSectionRepository, songPartRepository)

			var partIDs []uuid.UUID
			hasNewParts := false
			for _, p := range tt.request.Parts {
				if p.PartID != nil {
					partIDs = append(partIDs, *p.PartID)
				} else {
					hasNewParts = true
				}
			}

			if len(partIDs) > 0 {
				var mockParts []model.SongPart
				for _, pid := range partIDs {
					mockParts = append(mockParts, model.SongPart{ID: pid, SongID: tt.request.SongID})
				}
				songPartRepository.On("GetAllByIDs", new([]model.SongPart), partIDs).
					Return(nil, &mockParts).
					Once()
			}

			sectionsCount := int64(3)
			songSectionRepository.On("CountAllBySong", new(int64), tt.request.SongID).
				Return(nil, &sectionsCount).
				Once()

			if hasNewParts {
				partsCount := tt.partsCount
				songPartRepository.On("CountAllBySong", new(int64), tt.request.SongID).
					Return(nil, &partsCount).
					Once()
			}

			songSectionRepository.On("Create", mock.IsType(new(model.SongSection))).
				Run(func(args mock.Arguments) {
					newSection := args.Get(0).(*model.SongSection)
					assertCreatedSongSection(t, tt.request, *newSection, sectionsCount, tt.partsCount)
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
	partsCount int64,
) {
	assert.NotEmpty(t, section.ID)
	assert.Equal(t, request.Name, section.Name)
	assert.Equal(t, request.TypeID, section.SongSectionTypeID)
	assert.Equal(t, request.SongID, section.SongID)
	assert.Equal(t, uint(sectionsCount), section.Order)

	require.Len(t, section.SectionParts, len(request.Parts))

	nextSongOrder := uint(partsCount)
	for i, sp := range section.SectionParts {
		req := request.Parts[i]

		assert.Equal(t, uint(i), sp.Order)
		assert.Equal(t, req.BandMemberID, sp.BandMemberID)

		if req.PartID != nil {
			assert.Equal(t, *req.PartID, sp.PartID)
		} else {
			assert.Equal(t, sp.PartID, sp.Part.ID)
			assert.Equal(t, req.NewPart.Name, sp.Part.Name)
			assert.Equal(t, req.NewPart.InstrumentID, sp.Part.InstrumentID)
			assert.Equal(t, request.SongID, sp.Part.SongID)
			assert.Equal(t, nextSongOrder, sp.Part.SongOrder)
			nextSongOrder++
		}
	}
}
