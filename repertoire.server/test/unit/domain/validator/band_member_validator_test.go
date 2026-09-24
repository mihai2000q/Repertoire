package validator

import (
	"errors"
	"net/http"
	"repertoire/server/domain/validator"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBandMemberValidator_WhenNoIds_ShouldReturnNothing(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := validator.NewBandMemberValidator(artistRepository)

	// when
	bandMembers, errCode := _uut.Validate(nil, model.Song{})

	// then
	assert.Nil(t, errCode)
	assert.Empty(t, bandMembers)
	artistRepository.AssertExpectations(t)
}

func TestBandMemberValidator_WhenGetBandMembersFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := validator.NewBandMemberValidator(artistRepository)

	ids := []uuid.UUID{uuid.New(), uuid.New()}
	song := model.Song{ID: uuid.New(), ArtistID: &[]uuid.UUID{uuid.New()}[0]}

	internalError := errors.New("get band members error")
	artistRepository.On("GetBandMembersByIDs", new([]model.BandMember), ids).
		Return(internalError).
		Once()

	// when
	bandMembers, errCode := _uut.Validate(ids, song)

	// then
	assert.Empty(t, bandMembers)
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	artistRepository.AssertExpectations(t)
}

func TestBandMemberValidator_WhenSomeBandMembersNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := validator.NewBandMemberValidator(artistRepository)

	ids := []uuid.UUID{uuid.New(), uuid.New()}
	artistID := uuid.New()
	song := model.Song{ID: uuid.New(), ArtistID: &artistID}

	// only one of the two requested band members exists
	found := []model.BandMember{{ID: ids[0], ArtistID: artistID}}
	artistRepository.On("GetBandMembersByIDs", new([]model.BandMember), ids).
		Return(nil, &found).
		Once()

	// when
	bandMembers, errCode := _uut.Validate(ids, song)

	// then
	assert.Empty(t, bandMembers)
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "band members not found", errCode.Error.Error())

	artistRepository.AssertExpectations(t)
}

func TestBandMemberValidator_WhenAnyBandMemberNotAssociated_ShouldReturnConflictError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := validator.NewBandMemberValidator(artistRepository)

	ids := []uuid.UUID{uuid.New(), uuid.New()}
	artistID := uuid.New()
	song := model.Song{ID: uuid.New(), ArtistID: &artistID}

	found := []model.BandMember{
		{ID: ids[0], ArtistID: artistID},
		{ID: ids[1], ArtistID: uuid.New()}, // different artist
	}
	artistRepository.On("GetBandMembersByIDs", new([]model.BandMember), ids).
		Return(nil, &found).
		Once()

	// when
	bandMembers, errCode := _uut.Validate(ids, song)

	// then
	assert.Empty(t, bandMembers)
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusConflict, errCode.Code)
	assert.Equal(t, "band member is not part of the artist associated with this song", errCode.Error.Error())

	artistRepository.AssertExpectations(t)
}

func TestBandMemberValidator_WhenSongHasNoArtist_ShouldReturnConflictError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := validator.NewBandMemberValidator(artistRepository)

	ids := []uuid.UUID{uuid.New()}
	song := model.Song{ID: uuid.New()} // no artist

	found := []model.BandMember{{ID: ids[0], ArtistID: uuid.New()}}
	artistRepository.On("GetBandMembersByIDs", new([]model.BandMember), ids).
		Return(nil, &found).
		Once()

	// when
	bandMembers, errCode := _uut.Validate(ids, song)

	// then
	assert.Empty(t, bandMembers)
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusConflict, errCode.Code)
	assert.Equal(t, "band member is not part of the artist associated with this song", errCode.Error.Error())

	artistRepository.AssertExpectations(t)
}

func TestBandMemberValidator_WhenSuccessful_ShouldReturnBandMembers(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := validator.NewBandMemberValidator(artistRepository)

	ids := []uuid.UUID{uuid.New(), uuid.New(), uuid.New()}
	artistID := uuid.New()
	song := model.Song{ID: uuid.New(), ArtistID: &artistID}

	found := make([]model.BandMember, len(ids))
	for i, id := range ids {
		found[i] = model.BandMember{ID: id, ArtistID: artistID}
	}
	artistRepository.On("GetBandMembersByIDs", new([]model.BandMember), ids).
		Return(nil, &found).
		Once()

	// when
	bandMembers, errCode := _uut.Validate(ids, song)

	// then
	assert.Nil(t, errCode)
	assert.Equal(t, found, bandMembers)

	artistRepository.AssertExpectations(t)
}
