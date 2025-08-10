package processor

import (
	"repertoire/server/model"

	"github.com/stretchr/testify/mock"
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

func (p *ProgressProcessorMock) ComputeProgress(section model.SongSection) uint64 {
	args := p.Called(section)
	return args.Get(0).(uint64)
}
