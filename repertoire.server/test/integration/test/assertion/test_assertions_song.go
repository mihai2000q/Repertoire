package assertion

import (
	"repertoire/server/model"
	"slices"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func PerfectSongRehearsal(t *testing.T, song model.Song, newSong model.Song) {
	if newSong.DefaultArrangementID == nil { // nothing changed overall on the song
		assert.Equal(t, song, newSong)
		return
	}

	totalOccurrences := uint(0)
	defaultArrangementOccurrencesMap := make(map[uuid.UUID]uint)
	for _, section := range song.Sections {
		defaultArrangementOccurrences := slices.DeleteFunc(
			slices.Clone(section.ArrangementOccurrences),
			func(occ model.SongSectionOccurrences) bool {
				return occ.ArrangementID != *song.DefaultArrangementID
			})
		totalOccurrences += defaultArrangementOccurrences[0].Occurrences
		defaultArrangementOccurrencesMap[section.ID] = defaultArrangementOccurrences[0].Occurrences
	}

	if totalOccurrences == 0 { // also nothing changed
		assert.Equal(t, song, newSong)
		return
	}

	for i, newSection := range newSong.Sections {
		oldSection := song.Sections[i]

		if defaultArrangementOccurrencesMap[newSection.ID] == 0 { // nothing changed on this section
			assert.Equal(t, oldSection, newSection)
			continue
		}

		assert.Equal(t, newSection.Rehearsals, oldSection.Rehearsals+defaultArrangementOccurrencesMap[newSection.ID])
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
