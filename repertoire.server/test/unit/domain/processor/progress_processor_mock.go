package processor

import (
	"repertoire/server/model"

	"github.com/stretchr/testify/mock"
)

type ProgressProcessorMock struct {
	mock.Mock
}

func (p *ProgressProcessorMock) ComputeRehearsalsScore(history []model.SongPartHistory) uint64 {
	args := p.Called(history)
	return args.Get(0).(uint64)
}

func (p *ProgressProcessorMock) ComputeConfidenceScore(history []model.SongPartHistory) uint {
	args := p.Called(history)
	return args.Get(0).(uint)
}

func (p *ProgressProcessorMock) ComputeProgress(part model.SongPart) uint64 {
	args := p.Called(part)
	return args.Get(0).(uint64)
}
