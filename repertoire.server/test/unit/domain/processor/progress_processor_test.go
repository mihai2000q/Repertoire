package processor

import (
	"github.com/stretchr/testify/assert"
	"repertoire/server/domain/processor"
	"repertoire/server/model"
	"testing"
	"time"
)

func TestComputeRehearsalsScore_WhenHistoryLengthIs0_ShouldReturn0(t *testing.T) {
	// given
	_uut := processor.NewProgressProcessor()

	// when
	result := _uut.ComputeRehearsalsScore([]model.SongSectionHistory{})

	// then
	assert.Zero(t, result)
}

func TestComputeRehearsalsScore_WhenHistoryLengthIs1_ShouldReturnLatestRehearsalsValue(t *testing.T) {
	// given
	_uut := processor.NewProgressProcessor()
	history := []model.SongSectionHistory{
		{To: 1},
	}

	// when
	result := _uut.ComputeRehearsalsScore(history)

	// then
	assert.Equal(t, uint64(history[0].To), result)
}

func TestComputeRehearsalsScore_WhenSuccessful_ShouldComputeAndReturnRehearsalsScore(t *testing.T) {
	tests := []struct {
		name           string
		history        []model.SongSectionHistory
		expectedResult uint64
	}{
		{
			"Use Case 1",
			[]model.SongSectionHistory{
				{
					From:      0,
					To:        1,
					Property:  model.RehearsalsProperty,
					CreatedAt: time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					From:      1,
					To:        2,
					Property:  model.RehearsalsProperty,
					CreatedAt: time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					From:      2,
					To:        3,
					Property:  model.RehearsalsProperty,
					CreatedAt: time.Date(2024, time.January, 2, 0, 0, 0, 0, time.UTC),
				},
			},
			45,
		},
		{
			"Use Case 2",
			[]model.SongSectionHistory{
				{
					From:      0,
					To:        1,
					Property:  model.RehearsalsProperty,
					CreatedAt: time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					From:      1,
					To:        2,
					Property:  model.RehearsalsProperty,
					CreatedAt: time.Date(2024, time.January, 5, 0, 0, 0, 0, time.UTC),
				},
				{
					From:      2,
					To:        3,
					Property:  model.RehearsalsProperty,
					CreatedAt: time.Date(2024, time.January, 15, 0, 0, 0, 0, time.UTC),
				},
			},
			9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			_uut := processor.NewProgressProcessor()

			// when
			rehearsalsScore := _uut.ComputeRehearsalsScore(tt.history)

			// then
			assert.Equal(t, tt.expectedResult, rehearsalsScore)
		})
	}
}

func TestComputeConfidenceScore_WhenHistoryLengthIs0_ShouldReturnDefaultConfidence(t *testing.T) {
	// given
	_uut := processor.NewProgressProcessor()

	// when
	result := _uut.ComputeConfidenceScore([]model.SongSectionHistory{})

	// then
	assert.Equal(t, model.DefaultSongSectionConfidence, result)
}

func TestComputeConfidenceScore_WhenSuccessful_ShouldComputeAndReturnConfidenceScore(t *testing.T) {
	tests := []struct {
		name           string
		history        []model.SongSectionHistory
		expectedResult uint
	}{
		{
			"1 - Smooth Increase",
			[]model.SongSectionHistory{
				{
					From:     0,
					To:       10,
					Property: model.ConfidenceProperty,
				},
				{
					From:     10,
					To:       20,
					Property: model.ConfidenceProperty,
				},
				{
					From:     20,
					To:       30,
					Property: model.ConfidenceProperty,
				},
				{
					From:     30,
					To:       40,
					Property: model.ConfidenceProperty,
				},
			},
			47,
		},
		{
			"2 - Fluctuations with final high",
			[]model.SongSectionHistory{
				{
					From:     0,
					To:       10,
					Property: model.ConfidenceProperty,
				},
				{
					From:     10,
					To:       20,
					Property: model.ConfidenceProperty,
				},
				{
					From:     20,
					To:       30,
					Property: model.ConfidenceProperty,
				},
				{
					From:     30,
					To:       40,
					Property: model.ConfidenceProperty,
				},
				{
					From:     40,
					To:       30,
					Property: model.ConfidenceProperty,
				},
				{
					From:     30,
					To:       35,
					Property: model.ConfidenceProperty,
				},
				{
					From:     35,
					To:       50,
					Property: model.ConfidenceProperty,
				},
			},
			48,
		},
		{
			"3 - Fluctuations with low final result",
			[]model.SongSectionHistory{
				{
					From:     0,
					To:       10,
					Property: model.ConfidenceProperty,
				},
				{
					From:     10,
					To:       20,
					Property: model.ConfidenceProperty,
				},
				{
					From:     20,
					To:       30,
					Property: model.ConfidenceProperty,
				},
				{
					From:     30,
					To:       40,
					Property: model.ConfidenceProperty,
				},
				{
					From:     40,
					To:       30,
					Property: model.ConfidenceProperty,
				},
				{
					From:     30,
					To:       35,
					Property: model.ConfidenceProperty,
				},
				{
					From:     35,
					To:       25,
					Property: model.ConfidenceProperty,
				},
			},
			22,
		},
		{
			"4 - Extreme Increasing Jumps",
			[]model.SongSectionHistory{
				{
					From:     0,
					To:       50,
					Property: model.ConfidenceProperty,
				},
				{
					From:     50,
					To:       80,
					Property: model.ConfidenceProperty,
				},
			},
			76,
		},
		{
			"5 - Extreme Decreasing Jumps",
			[]model.SongSectionHistory{
				{
					From:     0,
					To:       50,
					Property: model.ConfidenceProperty,
				},
				{
					From:     50,
					To:       80,
					Property: model.ConfidenceProperty,
				},
				{
					From:     80,
					To:       50,
					Property: model.ConfidenceProperty,
				},
			},
			32,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			_uut := processor.NewProgressProcessor()

			// when
			confidenceScore := _uut.ComputeConfidenceScore(tt.history)

			// then
			assert.Equal(t, tt.expectedResult, confidenceScore)
		})
	}
}
