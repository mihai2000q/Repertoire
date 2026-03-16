package assertion

import (
	"repertoire/server/model"
	"slices"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func CustomSongRehearsal(t *testing.T, song model.Song, newSong model.Song, arrangementID uuid.UUID) {
	if len(newSong.Sections[0].ArrangementOccurrences) == 0 { // nothing changed overall on the song
		assert.Equal(t, song, newSong)
		return
	}

	assertSongRehearsal(t, song, newSong, 0, &arrangementID)
}

func PerfectSongRehearsal(t *testing.T, song model.Song, newSong model.Song) {
	if newSong.DefaultArrangementID == nil { // nothing changed overall on the song
		assert.Equal(t, song, newSong)
		return
	}

	assertSongRehearsal(t, song, newSong, 0, nil)
}

func PerfectSongRehearsalWithDuplicates(t *testing.T, song model.Song, newSong model.Song, duplicates int) {
	if newSong.DefaultArrangementID == nil { // nothing changed overall on the song
		assert.Equal(t, song, newSong)
		return
	}

	assertSongRehearsal(t, song, newSong, duplicates, nil)
}

func assertSongRehearsal(t *testing.T, song model.Song, newSong model.Song, duplicates int, arrangementID *uuid.UUID) {
	if arrangementID == nil {
		arrangementID = newSong.DefaultArrangementID
	}

	totalOccurrences := uint(0)
	arrangementOccurrencesMap := make(map[uuid.UUID]uint)
	for _, section := range newSong.Sections {
		arrangementOccurrences := slices.DeleteFunc(
			slices.Clone(section.ArrangementOccurrences),
			func(occ model.SongSectionOccurrences) bool {
				return occ.ArrangementID != *arrangementID
			})
		totalOccurrences += arrangementOccurrences[0].Occurrences
		arrangementOccurrencesMap[section.ID] = arrangementOccurrences[0].Occurrences
	}

	if totalOccurrences == 0 { // also nothing changed
		assert.Equal(t, song, newSong)
		return
	}

	for i, newSection := range newSong.Sections {
		oldSection := song.Sections[i]

		if arrangementOccurrencesMap[newSection.ID] == 0 { // nothing changed on this section
			assert.Equal(t, oldSection, newSection)
			continue
		}

		newRehearsals := oldSection.Rehearsals + arrangementOccurrencesMap[newSection.ID]*uint(duplicates+1)
		assert.Equal(t, newRehearsals, newSection.Rehearsals)
		assert.Greater(t, newSection.RehearsalsScore, oldSection.RehearsalsScore)
		assert.GreaterOrEqual(t, newSection.Progress, oldSection.Progress)

		for j := 0; j <= duplicates; j++ {
			fromDiff := arrangementOccurrencesMap[newSection.ID] * uint(duplicates-j)
			toDiff := arrangementOccurrencesMap[newSection.ID] * uint(j)
			assert.NotEmpty(t, newSection.History[j].ID)
			assert.Equal(t, oldSection.Rehearsals+fromDiff, newSection.History[j].From)
			assert.Equal(t, newSection.Rehearsals-toDiff, newSection.History[j].To)
			assert.Equal(t, model.RehearsalsProperty, newSection.History[0].Property)
		}
	}

	assert.Greater(t, newSong.Rehearsals, song.Rehearsals)
	assert.Greater(t, newSong.Progress, song.Progress)

	assert.NotNil(t, newSong.LastTimePlayed)
	assert.WithinDuration(t, time.Now(), *newSong.LastTimePlayed, 1*time.Minute)
}
