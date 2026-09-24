package bandmember

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/httperror"
	"repertoire/server/internal/reorder"
	"repertoire/server/model"
)

type MoveBandMember struct {
	artistRepository repository.ArtistRepository
}

func NewMoveBandMember(artistRepository repository.ArtistRepository) MoveBandMember {
	return MoveBandMember{
		artistRepository: artistRepository,
	}
}

func (m MoveBandMember) Handle(request requests.MoveBandMemberRequest) *httperror.ErrorCode {
	var artist model.Artist
	if err := m.artistRepository.GetWithBandMembers(&artist, request.ArtistID); err != nil {
		return httperror.DatabaseError(err)
	}
	if reflect.ValueOf(artist).IsZero() {
		return httperror.NotFoundError(errors.New("artist not found"))
	}

	errCode := reorder.MoveEntity(
		artist.BandMembers,
		request.ID,
		request.OverID,
		&reorder.Config{
			EntityNotFoundMsg:     "band member not found",
			OverEntityNotFoundMsg: "over band member not found",
		},
	)
	if errCode != nil {
		return errCode
	}

	if err := m.artistRepository.UpdateWithAssociations(&artist); err != nil {
		return httperror.DatabaseError(err)
	}

	return nil
}
