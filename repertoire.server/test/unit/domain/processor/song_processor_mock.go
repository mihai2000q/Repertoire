package processor

import (
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/stretchr/testify/mock"
)

type SongProcessorMock struct {
	mock.Mock
}

func (s *SongProcessorMock) AddPerfectRehearsal(
	song *model.Song,
	repository repository.SongRepository,
) (*wrapper.ErrorCode, bool) {
	args := s.Called(song, repository)

	var errCode *wrapper.ErrorCode
	if e := args.Get(0); e != nil {
		errCode = e.(*wrapper.ErrorCode)
	}
	
	return errCode, args.Bool(1)
}
