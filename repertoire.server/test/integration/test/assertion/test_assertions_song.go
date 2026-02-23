package assertion

import (
	"repertoire/server/model"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func PerfectSongRehearsal(t *testing.T, song model.Song, newSong model.Song) {
	for i, newSection := range newSong.Sections {
		oldSection := song.Sections[i]
		defaultArrangementOccurrences := slices.DeleteFunc(
			slices.Clone(oldSection.ArrangementOccurrences),
			func(occ model.SongSectionOccurrences) bool {
				return occ.ArrangementID != *song.DefaultArrangementID
			})[0]
		if defaultArrangementOccurrences.Occurrences == 0 { // nothing changed
			assert.Equal(t, oldSection, newSection)
			continue
		}

		assert.Equal(t, newSection.Rehearsals, oldSection.Rehearsals+defaultArrangementOccurrences.Occurrences)
		assert.Greater(t, newSection.RehearsalsScore, oldSection.RehearsalsScore)
		assert.GreaterOrEqual(t, newSection.Progress, oldSection.Progress)

		assert.NotEmpty(t, newSection.History[0].ID)
		assert.Equal(t, oldSection.Rehearsals, newSection.History[0].From)
		assert.Equal(t, newSection.Rehearsals, newSection.History[0].To)
		assert.Equal(t, model.RehearsalsProperty, newSection.History[0].Property)
	}

	assert.Greater(t, newSong.Rehearsals, song.Rehearsals)
	assert.Greater(t, newSong.Progress, song.Progress)

	assert.NotNil(t, newSong.LastTimePlayed)
	assert.WithinDuration(t, time.Now(), *newSong.LastTimePlayed, 1*time.Minute)
}
