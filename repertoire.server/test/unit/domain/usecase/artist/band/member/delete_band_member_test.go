package member

import (
	"errors"
	"net/http"
	"repertoire/server/domain/usecase/artist/band/member"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"slices"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteBandMember_WhenGetArtistFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := member.NewDeleteBandMember(artistRepository)

	id := uuid.New()
	artistID := uuid.New()

	internalError := errors.New("internal error")
	artistRepository.On("GetWithBandMembers", new(model.Artist), artistID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id, artistID)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	artistRepository.AssertExpectations(t)
}

func TestDeleteBandMember_WhenArtistIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := member.NewDeleteBandMember(artistRepository)

	id := uuid.New()
	artistID := uuid.New()

	artistRepository.On("GetWithBandMembers", new(model.Artist), artistID).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(id, artistID)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "artist not found", errCode.Error.Error())

	artistRepository.AssertExpectations(t)
}

func TestDeleteBandMember_WhenBandMemberIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := member.NewDeleteBandMember(artistRepository)

	id := uuid.New()
	artistID := uuid.New()

	// given - mocking
	artist := &model.Artist{
		ID: artistID,
		BandMembers: []model.BandMember{
			{ID: uuid.New(), Order: 0},
		},
	}
	artistRepository.On("GetWithBandMembers", new(model.Artist), artistID).
		Return(nil, artist).
		Once()

	// when
	errCode := _uut.Handle(id, artistID)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "band member not found", errCode.Error.Error())

	artistRepository.AssertExpectations(t)
}

func TestDeleteBandMember_WhenUpdateArtistFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := member.NewDeleteBandMember(artistRepository)

	id := uuid.New()
	artistID := uuid.New()

	// given - mocking
	artist := &model.Artist{
		ID: artistID,
		BandMembers: []model.BandMember{
			{ID: id, Order: 0},
		},
	}
	artistRepository.On("GetWithBandMembers", new(model.Artist), artistID).
		Return(nil, artist).
		Once()

	internalError := errors.New("internal error")
	artistRepository.On("UpdateWithAssociations", mock.IsType(artist)).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(id, artistID)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	artistRepository.AssertExpectations(t)
}

func TestDeleteBandMember_WhenDeleteBandMemberFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := member.NewDeleteBandMember(artistRepository)

	id := uuid.New()
	artistID := uuid.New()

	// given - mocking
	artist := &model.Artist{
		ID: artistID,
		BandMembers: []model.BandMember{
			{ID: id, Order: 0},
		},
	}
	artistRepository.On("GetWithBandMembers", new(model.Artist), artistID).
		Return(nil, artist).
		Once()

	artistRepository.On("UpdateWithAssociations", mock.IsType(artist)).
		Return(nil).
		Once()

	internalError := errors.New("internal error")
	artistRepository.On("DeleteBandMember", id).Return(internalError).Once()

	// when
	errCode := _uut.Handle(id, artistID)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	artistRepository.AssertExpectations(t)
}

func TestDeleteBandMember_WhenSuccessful_ShouldNotReturnAnyError(t *testing.T) {
	tests := []struct {
		name        string
		artist      model.Artist
		memberIndex uint
	}{
		{
			"1 - When it was the only member",
			model.Artist{
				ID: uuid.New(),
				BandMembers: []model.BandMember{
					{ID: uuid.New(), Order: 0},
				},
			},
			0,
		},
		{
			"2 - When there are more members",
			model.Artist{
				ID: uuid.New(),
				BandMembers: []model.BandMember{
					{ID: uuid.New(), Order: 0},
					{ID: uuid.New(), Order: 1},
					{ID: uuid.New(), Order: 2},
				},
			},
			0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			artistRepository := new(repository.ArtistRepositoryMock)
			_uut := member.NewDeleteBandMember(artistRepository)

			id := tt.artist.BandMembers[tt.memberIndex].ID
			artistID := tt.artist.ID

			// given - mocking
			artistRepository.On("GetWithBandMembers", new(model.Artist), artistID).
				Return(nil, &tt.artist).
				Once()

			artistRepository.On("UpdateWithAssociations", mock.IsType(&tt.artist)).
				Run(func(args mock.Arguments) {
					newArtist := args.Get(0).(*model.Artist)

					// members ordered
					members := slices.Clone(newArtist.BandMembers)
					members = slices.DeleteFunc(members, func(a model.BandMember) bool {
						return a.ID == id
					})
					for i, s := range members {
						assert.Equal(t, uint(i), s.Order)
					}
				}).
				Return(nil).
				Once()

			artistRepository.On("DeleteBandMember", id).Return(nil).Once()

			// when
			errCode := _uut.Handle(id, artistID)

			// then
			assert.Nil(t, errCode)

			artistRepository.AssertExpectations(t)
		})
	}
}
