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
	for _, section := range newSong.Sections {
		defaultArrangementOccurrences := slices.DeleteFunc(
			slices.Clone(section.ArrangementOccurrences),
			func(occ model.SongSectionOccurrences) bool {
				return occ.ArrangementID != *newSong.DefaultArrangementID
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

		newRehearsals := oldSection.Rehearsals + defaultArrangementOccurrencesMap[newSection.ID]
		assert.Equal(t, newRehearsals, newSection.Rehearsals)
		assert.Greater(t, newSection.RehearsalsScore, oldSection.RehearsalsScore)
		assert.GreaterOrEqual(t, newSection.Progress, oldSection.Progress)

		assert.NotEmpty(t, newSection.History[i].ID)
		assert.Equal(t, oldSection.Rehearsals, newSection.History[i].From)
		assert.Equal(t, newSection.Rehearsals, newSection.History[i].To)
		assert.Equal(t, model.RehearsalsProperty, newSection.History[0].Property)
	}

	assert.Greater(t, newSong.Rehearsals, song.Rehearsals)
	assert.Greater(t, newSong.Progress, song.Progress)

	assert.NotNil(t, newSong.LastTimePlayed)
	assert.WithinDuration(t, time.Now(), *newSong.LastTimePlayed, 1*time.Minute)
}

func PerfectSongRehearsalWithDuplicates(t *testing.T, song model.Song, newSong model.Song, duplicates int) {
	if newSong.DefaultArrangementID == nil { // nothing changed overall on the song
		assert.Equal(t, song, newSong)
		return
	}

	totalOccurrences := uint(0)
	defaultArrangementOccurrencesMap := make(map[uuid.UUID]uint)
	for _, section := range newSong.Sections {
		defaultArrangementOccurrences := slices.DeleteFunc(
			slices.Clone(section.ArrangementOccurrences),
			func(occ model.SongSectionOccurrences) bool {
				return occ.ArrangementID != *newSong.DefaultArrangementID
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

		newRehearsals := oldSection.Rehearsals + defaultArrangementOccurrencesMap[newSection.ID]*uint(duplicates+1)
		assert.Equal(t, newRehearsals, newSection.Rehearsals)
		assert.Greater(t, newSection.RehearsalsScore, oldSection.RehearsalsScore)
		assert.GreaterOrEqual(t, newSection.Progress, oldSection.Progress)

		for i := 0; i <= duplicates; i++ {
			fromDiff := defaultArrangementOccurrencesMap[newSection.ID] * uint(duplicates-i)
			toDiff := defaultArrangementOccurrencesMap[newSection.ID] * uint(i)
			assert.NotEmpty(t, newSection.History[i].ID)
			assert.Equal(t, oldSection.Rehearsals+fromDiff, newSection.History[i].From)
			assert.Equal(t, newSection.Rehearsals-toDiff, newSection.History[i].To)
			assert.Equal(t, model.RehearsalsProperty, newSection.History[0].Property)
		}
	}

	assert.Greater(t, newSong.Rehearsals, song.Rehearsals)
	assert.Greater(t, newSong.Progress, song.Progress)

	assert.NotNil(t, newSong.LastTimePlayed)
	assert.WithinDuration(t, time.Now(), *newSong.LastTimePlayed, 1*time.Minute)
}
