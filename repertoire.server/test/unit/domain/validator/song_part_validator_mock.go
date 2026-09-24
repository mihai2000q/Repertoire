package validator

import (
	"repertoire/server/internal/httperror"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type SongPartValidatorMock struct {
	mock.Mock
}

func (s *SongPartValidatorMock) Validate(ids []uuid.UUID, songID uuid.UUID) *httperror.ErrorCode {
	args := s.Called(ids, songID)

	var errCode *httperror.ErrorCode
	if args.Get(0) != nil {
		errCode = args.Get(0).(*httperror.ErrorCode)
	}

	return errCode
}
