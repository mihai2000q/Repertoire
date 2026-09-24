package section

import (
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/song/section"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"repertoire/server/test/unit/domain/validator"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateSongSection_WhenGetSongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := section.NewCreateSongSection(nil, nil, songRepository, nil, nil)

	request := requests.CreateSongSectionRequest{
		SongID: uuid.New(),
		Name:   "Some Section",
		TypeID: uuid.New(),
	}

	internalError := errors.New("internal error")
	songRepository.On("Get", new(model.Song), request.SongID).
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

func TestCreateSongSection_WhenSongNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	_uut := section.NewCreateSongSection(nil, nil, songRepository, nil, nil)

	request := requests.CreateSongSectionRequest{
		SongID: uuid.New(),
		Name:   "Some Section",
		TypeID: uuid.New(),
	}

	songRepository.On("Get", new(model.Song), request.SongID).
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

func TestCreateSongSection_WhenSongPartValidationFails_ShouldReturnItsError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	songPartValidator := new(validator.SongPartValidatorMock)
	_uut := section.NewCreateSongSection(nil, nil, songRepository, nil, songPartValidator)

	partID := uuid.New()
	request := requests.CreateSongSectionRequest{
		SongID: uuid.New(),
		Name:   "Some Section",
		TypeID: uuid.New(),
		Parts: []requests.CreateSongSectionPartRequest{
			{PartID: &partID},
		},
	}

	mockSong := &model.Song{ID: request.SongID}
	songRepository.On("Get", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()

	validationError := httperror.NotFoundError(errors.New("parts not found"))
	songPartValidator.On("Validate", []uuid.UUID{partID}, request.SongID).
		Return(validationError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, validationError, errCode)

	songRepository.AssertExpectations(t)
	songPartValidator.AssertExpectations(t)
}

func TestCreateSongSection_WhenBandMemberValidationFails_ShouldReturnItsError(t *testing.T) {
	// given
	songRepository := new(repository.SongRepositoryMock)
	bandMemberValidator := new(validator.BandMemberValidatorMock)
	_uut := section.NewCreateSongSection(nil, nil, songRepository, bandMemberValidator, nil)

	bandMemberID := uuid.New()
	request := requests.CreateSongSectionRequest{
		SongID: uuid.New(),
		Name:   "Some Section",
		TypeID: uuid.New(),
		Parts: []requests.CreateSongSectionPartRequest{
			{
				NewPart:       &requests.CreateNewSongPartRequest{Name: "Verse Riff"},
				BandMemberIDs: []uuid.UUID{bandMemberID},
			},
		},
	}

	mockSong := &model.Song{ID: request.SongID}
	songRepository.On("Get", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()

	validationError := httperror.ConflictError(
		errors.New("band member is not part of the artist associated with this song"),
	)
	bandMemberValidator.On("Validate", []uuid.UUID{bandMemberID}, *mockSong).
		Return(nil, validationError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, validationError, errCode)

	songRepository.AssertExpectations(t)
	bandMemberValidator.AssertExpectations(t)
}

func TestCreateSongSection_WhenCountSectionsBySongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	bandMemberValidator := new(validator.BandMemberValidatorMock)
	_uut := section.NewCreateSongSection(songSectionRepository, nil, songRepository, bandMemberValidator, nil)

	request := requests.CreateSongSectionRequest{
		SongID: uuid.New(),
		Name:   "Some Section",
		TypeID: uuid.New(),
	}

	mockSong := &model.Song{ID: request.SongID}
	songRepository.On("Get", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()

	bandMemberValidator.On("Validate", []uuid.UUID{}, *mockSong).
		Return(nil, nil).
		Once()

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

	songRepository.AssertExpectations(t)
	bandMemberValidator.AssertExpectations(t)
	songSectionRepository.AssertExpectations(t)
}

func TestCreateSongSection_WhenCountPartsBySongFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	songPartRepository := new(repository.SongPartRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	bandMemberValidator := new(validator.BandMemberValidatorMock)
	_uut := section.NewCreateSongSection(
		songSectionRepository,
		songPartRepository,
		songRepository,
		bandMemberValidator,
		nil,
	)

	request := requests.CreateSongSectionRequest{
		SongID: uuid.New(),
		Name:   "Some Section",
		TypeID: uuid.New(),
		Parts: []requests.CreateSongSectionPartRequest{
			{NewPart: &requests.CreateNewSongPartRequest{Name: "Verse Riff"}},
		},
	}

	mockSong := &model.Song{ID: request.SongID}
	songRepository.On("Get", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()

	bandMemberValidator.On("Validate", []uuid.UUID{}, *mockSong).
		Return(nil, nil).
		Once()

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

	songRepository.AssertExpectations(t)
	bandMemberValidator.AssertExpectations(t)
	songSectionRepository.AssertExpectations(t)
	songPartRepository.AssertExpectations(t)
}

func TestCreateSongSection_WhenCreateSectionFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	songSectionRepository := new(repository.SongSectionRepositoryMock)
	songRepository := new(repository.SongRepositoryMock)
	bandMemberValidator := new(validator.BandMemberValidatorMock)
	_uut := section.NewCreateSongSection(songSectionRepository, nil, songRepository, bandMemberValidator, nil)

	request := requests.CreateSongSectionRequest{
		SongID: uuid.New(),
		Name:   "Some Section",
		TypeID: uuid.New(),
	}

	mockSong := &model.Song{ID: request.SongID}
	songRepository.On("Get", new(model.Song), request.SongID).
		Return(nil, mockSong).
		Once()

	bandMemberValidator.On("Validate", []uuid.UUID{}, *mockSong).
		Return(nil, nil).
		Once()

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

	songRepository.AssertExpectations(t)
	bandMemberValidator.AssertExpectations(t)
	songSectionRepository.AssertExpectations(t)
}

func TestCreateSongSection_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	existingPartID1 := uuid.New()
	existingPartID2 := uuid.New()
	bandMemberID1 := uuid.New()
	bandMemberID2 := uuid.New()
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
					{PartID: &existingPartID2, BandMemberIDs: []uuid.UUID{bandMemberID1, bandMemberID2}},
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
					{
						NewPart:       &requests.CreateNewSongPartRequest{Name: "Chorus Riff"},
						BandMemberIDs: []uuid.UUID{bandMemberID1},
					},
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
			songRepository := new(repository.SongRepositoryMock)
			bandMemberValidator := new(validator.BandMemberValidatorMock)
			songPartValidator := new(validator.SongPartValidatorMock)
			_uut := section.NewCreateSongSection(
				songSectionRepository,
				songPartRepository,
				songRepository,
				bandMemberValidator,
				songPartValidator,
			)

			mockSong := &model.Song{ID: tt.request.SongID}
			songRepository.On("Get", new(model.Song), tt.request.SongID).
				Return(nil, mockSong).
				Once()

			// existing-part validation: only called when the request references at least one PartID
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
				songPartValidator.On("Validate", partIDs, tt.request.SongID).
					Return(nil).
					Once()
			}

			// band member validation: called once with the distinct ids across every part,
			// in first-seen order - always called, even with an empty (non-nil) slice
			var bandMemberIDs []uuid.UUID
			seen := make(map[uuid.UUID]bool)
			for _, p := range tt.request.Parts {
				for _, id := range p.BandMemberIDs {
					if !seen[id] {
						seen[id] = true
						bandMemberIDs = append(bandMemberIDs, id)
					}
				}
			}
			expectedBandMemberIDs := bandMemberIDs
			if expectedBandMemberIDs == nil {
				expectedBandMemberIDs = []uuid.UUID{}
			}
			bandMembers := make([]model.BandMember, len(bandMemberIDs))
			for i, id := range bandMemberIDs {
				bandMembers[i] = model.BandMember{ID: id}
			}
			bandMemberValidator.On("Validate", expectedBandMemberIDs, *mockSong).
				Return(bandMembers, nil).
				Once()

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

			songRepository.AssertExpectations(t)
			bandMemberValidator.AssertExpectations(t)
			songPartValidator.AssertExpectations(t)
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

		bandMemberIDs := make([]uuid.UUID, len(sp.BandMembers))
		for j, bm := range sp.BandMembers {
			bandMemberIDs[j] = bm.ID
		}
		assert.ElementsMatch(t, req.BandMemberIDs, bandMemberIDs)

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
