package member

import (
	"cmp"
	"errors"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/artist/band/member"
	"repertoire/server/model"
	"repertoire/server/test/unit/data/repository"
	"slices"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMoveBandMember_WhenGetArtistFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := member.NewMoveBandMember(artistRepository)

	request := requests.MoveBandMemberRequest{
		ID:       uuid.New(),
		OverID:   uuid.New(),
		ArtistID: uuid.New(),
	}

	// given - mocking
	internalError := errors.New("internal error")
	artistRepository.On("GetWithBandMembers", new(model.Artist), request.ArtistID).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	artistRepository.AssertExpectations(t)
}

func TestMoveBandMember_WhenArtistIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := member.NewMoveBandMember(artistRepository)

	request := requests.MoveBandMemberRequest{
		ID:       uuid.New(),
		OverID:   uuid.New(),
		ArtistID: uuid.New(),
	}

	// given - mocking
	artistRepository.On("GetWithBandMembers", new(model.Artist), request.ArtistID).
		Return(nil).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "artist not found", errCode.Error.Error())

	artistRepository.AssertExpectations(t)
}

func TestMoveBandMember_WhenBandMemberIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := member.NewMoveBandMember(artistRepository)

	artist := &model.Artist{ID: uuid.New()}

	request := requests.MoveBandMemberRequest{
		ID:       uuid.New(),
		OverID:   uuid.New(),
		ArtistID: artist.ID,
	}

	// given - mocking
	artistRepository.On("GetWithBandMembers", new(model.Artist), request.ArtistID).
		Return(nil, artist).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "band member not found", errCode.Error.Error())

	artistRepository.AssertExpectations(t)
}

func TestMoveBandMember_WhenOverBandMemberIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := member.NewMoveBandMember(artistRepository)

	artist := &model.Artist{
		ID: uuid.New(),
		BandMembers: []model.BandMember{
			{ID: uuid.New(), Order: 0},
		},
	}

	request := requests.MoveBandMemberRequest{
		ID:       artist.BandMembers[0].ID,
		OverID:   uuid.New(),
		ArtistID: artist.ID,
	}

	// given - mocking
	artistRepository.On("GetWithBandMembers", new(model.Artist), request.ArtistID).
		Return(nil, artist).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "over band member not found", errCode.Error.Error())

	artistRepository.AssertExpectations(t)
}

func TestMoveBandMember_WhenUpdateFails_ShouldReturnInternalServerError(t *testing.T) {
	// given
	artistRepository := new(repository.ArtistRepositoryMock)
	_uut := member.NewMoveBandMember(artistRepository)

	artist := &model.Artist{
		ID: uuid.New(),
		BandMembers: []model.BandMember{
			{ID: uuid.New(), Order: 0},
			{ID: uuid.New(), Order: 1},
		},
	}

	request := requests.MoveBandMemberRequest{
		ID:       artist.BandMembers[0].ID,
		OverID:   artist.BandMembers[1].ID,
		ArtistID: artist.ID,
	}

	// given - mocking
	artistRepository.On("GetWithBandMembers", new(model.Artist), request.ArtistID).
		Return(nil, artist).
		Once()

	internalError := errors.New("internal error")
	artistRepository.On("UpdateWithAssociations", mock.IsType(new(model.Artist))).
		Return(internalError).
		Once()

	// when
	errCode := _uut.Handle(request)

	// then
	assert.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Equal(t, internalError, errCode.Error)

	artistRepository.AssertExpectations(t)
}

func TestMoveBandMember_WhenIsValid_ShouldNotReturnAnyError(t *testing.T) {
	tests := []struct {
		name      string
		artist    *model.Artist
		index     uint
		overIndex uint
	}{
		{
			"Use case 1",
			&model.Artist{
				ID: uuid.New(),
				BandMembers: []model.BandMember{
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
			&model.Artist{
				ID: uuid.New(),
				BandMembers: []model.BandMember{
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
			artistRepository := new(repository.ArtistRepositoryMock)
			_uut := member.NewMoveBandMember(artistRepository)

			request := requests.MoveBandMemberRequest{
				ID:       tt.artist.BandMembers[tt.index].ID,
				OverID:   tt.artist.BandMembers[tt.overIndex].ID,
				ArtistID: tt.artist.ID,
			}

			// given - mocking
			artistRepository.On("GetWithBandMembers", new(model.Artist), request.ArtistID).
				Return(nil, tt.artist).
				Once()

			artistRepository.On("UpdateWithAssociations", mock.IsType(new(model.Artist))).
				Run(func(args mock.Arguments) {
					artist := args.Get(0).(*model.Artist)
					members := slices.Clone(artist.BandMembers)
					slices.SortFunc(members, func(a, b model.BandMember) int {
						return cmp.Compare(a.Order, b.Order)
					})
					if tt.index < tt.overIndex {
						assert.Equal(t, members[tt.overIndex-1].ID, request.OverID)
					} else if tt.index > tt.overIndex {
						assert.Equal(t, members[tt.overIndex+1].ID, request.OverID)
					}
					assert.Equal(t, members[tt.overIndex].ID, request.ID)
					for i, s := range members {
						assert.Equal(t, uint(i), s.Order)
					}
				}).
				Return(nil).
				Once()

			// when
			errCode := _uut.Handle(request)

			// then
			assert.Nil(t, errCode)

			artistRepository.AssertExpectations(t)
		})
	}
}
