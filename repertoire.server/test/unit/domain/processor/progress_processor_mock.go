package processor

import (
	"github.com/stretchr/testify/mock"
	"repertoire/server/model"
)

type ProgressProcessorMock struct {
	mock.Mock
}

func (p *ProgressProcessorMock) ComputeRehearsalsScore(history []model.SongSectionHistory) uint64 {
	args := p.Called(history)
	return args.Get(0).(uint64)
}

func (p *ProgressProcessorMock) ComputeConfidenceScore(history []model.SongSectionHistory) uint {
	args := p.Called(history)
	return args.Get(0).(uint)
}
