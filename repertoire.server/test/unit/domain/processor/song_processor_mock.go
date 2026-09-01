package processor

import (
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
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
) (*wrapper.ErrorCode, bool) {
	args := s.Called(song, songPartRepository, arrangementID)

	var errCode *wrapper.ErrorCode
	if e := args.Get(0); e != nil {
		errCode = e.(*wrapper.ErrorCode)
	}

	return errCode, args.Bool(1)
}

func (s *SongProcessorMock) AddPerfectRehearsal(
	song *model.Song,
	songPartRepository repository.SongPartRepository,
) (*wrapper.ErrorCode, bool) {
	args := s.Called(song, songPartRepository)

	var errCode *wrapper.ErrorCode
	if e := args.Get(0); e != nil {
		errCode = e.(*wrapper.ErrorCode)
	}

	return errCode, args.Bool(1)
}

func (s *SongProcessorMock) UpdateSongAfterPartsDeletion(
	songRepository repository.SongRepository,
	songID uuid.UUID,
	partIDs []uuid.UUID,
) *wrapper.ErrorCode {
	args := s.Called(songRepository, songID, partIDs)

	var errCode *wrapper.ErrorCode
	if e := args.Get(0); e != nil {
		errCode = e.(*wrapper.ErrorCode)
	}

	return errCode
}
