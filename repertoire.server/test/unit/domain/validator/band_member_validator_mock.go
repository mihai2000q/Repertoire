package validator

import (
	"repertoire/server/internal/httperror"
	"repertoire/server/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type BandMemberValidatorMock struct {
	mock.Mock
}

func (b *BandMemberValidatorMock) Validate(ids []uuid.UUID, song model.Song) ([]model.BandMember, *httperror.ErrorCode) {
	args := b.Called(ids, song)
	return args.Get(0).([]model.BandMember), args.Get(1).(*httperror.ErrorCode)
}
