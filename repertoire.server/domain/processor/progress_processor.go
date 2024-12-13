package processor

import (
	"math"
	"repertoire/server/model"
	"time"
)

type ProgressProcessor interface {
	ComputeRehearsalsScore(history []model.SongSectionHistory) uint64
	ComputeConfidenceScore(history []model.SongSectionHistory) uint
	ComputeProgress(section model.SongSection) uint64
}

type progressProcessor struct{}

func NewProgressProcessor() ProgressProcessor {
	return &progressProcessor{}
}

func (progressProcessor) ComputeRehearsalsScore(history []model.SongSectionHistory) uint64 {
	if len(history) == 0 {
		return 0
	}
	if len(history) == 1 {
		return uint64(history[0].To)
	}

	var rehearsalsScore uint64 = 0
	for i := 1; i < len(history); i++ {
		previousRehearsal := history[i-1]
		currentRehearsal := history[i]

		daysDifference := currentRehearsal.CreatedAt.Add(24 * time.Hour).Sub(previousRehearsal.CreatedAt)
		daysValue := 30 / (daysDifference.Hours() / 24)
		rehearsalsDifference := currentRehearsal.To - currentRehearsal.From
		rehearsalsScore += uint64(math.Round(float64(rehearsalsDifference) * daysValue))
	}
	return rehearsalsScore
}

func (progressProcessor) ComputeConfidenceScore(history []model.SongSectionHistory) uint {
	historyLength := len(history)
	if historyLength == 0 {
		return model.DefaultSongSectionConfidence
	}

	overallConfidence := int(history[historyLength-1].To - history[0].From)

	// calculate volatility of fluctuations
	totalFluctuation := 0 // mean
	positiveFluctuationsLen := 0
	positiveFluctuation := 0
	extremeFluctuations := 0
	for _, h := range history {
		currentFluctuation := int(h.To - h.From)
		totalFluctuation += currentFluctuation
		if math.Abs(float64(currentFluctuation)) > 30 {
			extremeFluctuations++
		}

		if currentFluctuation > 0 {
			positiveFluctuationsLen++
			positiveFluctuation += currentFluctuation
		}
	}
	averageFluctuation := float64(totalFluctuation) / float64(historyLength)
	averagePositiveFluctuation := float64(positiveFluctuation) / float64(positiveFluctuationsLen)

	negativeFluctuation := totalFluctuation - positiveFluctuation
	negativeFluctuationsLen := historyLength - positiveFluctuationsLen
	var averageNegativeFluctuation float64 = 0
	if negativeFluctuationsLen != 0 {
		averageNegativeFluctuation = float64(negativeFluctuation) / float64(negativeFluctuationsLen)
	}

	// apply penalty for fluctuations and calculate standard deviation
	standardDeviation := 0.0
	for _, h := range history {
		currentFluctuation := int(h.To - h.From)
		standardDeviation += math.Pow(float64(currentFluctuation)-averageFluctuation, 2)
	}
	standardDeviation = math.Sqrt(standardDeviation / float64(historyLength))

	penaltyDeviation := math.Pow(standardDeviation, 1.25) * 0.25
	penaltyNegative := math.Abs(averageNegativeFluctuation) * 0.55
	penaltyExtreme := float64(extremeFluctuations) / float64(historyLength) * 55 / float64(negativeFluctuationsLen+1)
	penalty := penaltyDeviation + penaltyNegative + penaltyExtreme

	// reward for smooth increase
	reward := averagePositiveFluctuation * 0.7

	finalScore := float64(overallConfidence) - penalty + reward
	confidenceScore := uint(math.Min(math.Max(math.Round(finalScore), 0), 100))

	return confidenceScore
}

func (progressProcessor) ComputeProgress(section model.SongSection) uint64 {
	progress := float64(section.ConfidenceScore) * float64(section.RehearsalsScore) / 100
	if progress < 0.5 {
		return 1
	}
	return uint64(math.Round(progress))
}
