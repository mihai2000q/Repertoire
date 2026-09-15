package processor

import (
	"repertoire/server/data/repository"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type SongProcessorMock struct {
	mock.Mock
}

func (s *SongProcessorMock) AddCustomRehearsal(
	song *model.Song,
	songPartRepository repository.SongPartRepository,
	arrangementID *uuid.UUID,
) (*httperror.ErrorCode, bool) {
	args := s.Called(song, songPartRepository, arrangementID)

	var errCode *httperror.ErrorCode
	if e := args.Get(0); e != nil {
		errCode = e.(*httperror.ErrorCode)
	}

	return errCode, args.Bool(1)
}

func (s *SongProcessorMock) AddPerfectRehearsal(
	song *model.Song,
	songPartRepository repository.SongPartRepository,
) (*httperror.ErrorCode, bool) {
	args := s.Called(song, songPartRepository)

	var errCode *httperror.ErrorCode
	if e := args.Get(0); e != nil {
		errCode = e.(*httperror.ErrorCode)
	}

	return errCode, args.Bool(1)
}

func (s *SongProcessorMock) UpdateSongAfterPartsDeletion(
	songRepository repository.SongRepository,
	songID uuid.UUID,
	partIDs []uuid.UUID,
) *httperror.ErrorCode {
	args := s.Called(songRepository, songID, partIDs)

	var errCode *httperror.ErrorCode
	if e := args.Get(0); e != nil {
		errCode = e.(*httperror.ErrorCode)
	}

	return errCode
}
