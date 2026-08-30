package bandmember

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/reorder"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type MoveBandMember struct {
	artistRepository repository.ArtistRepository
}

func NewMoveBandMember(repository repository.ArtistRepository) MoveBandMember {
	return MoveBandMember{
		artistRepository: repository,
	}
}

func (m MoveBandMember) Handle(request requests.MoveBandMemberRequest) *wrapper.ErrorCode {
	var artist model.Artist
	err := m.artistRepository.GetWithBandMembers(&artist, request.ArtistID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(artist).IsZero() {
		return wrapper.NotFoundError(errors.New("artist not found"))
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

	err = m.artistRepository.UpdateWithAssociations(&artist)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
