package assertion

import (
	"repertoire/server/model"
	"slices"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func CustomSongRehearsalWithDuplicates(
	t *testing.T,
	song model.Song,
	newSong model.Song,
	arrangementID uuid.UUID,
	duplicates int,
) {
	assertSongRehearsal(t, song, newSong, duplicates, &arrangementID)
}

func CustomSongRehearsal(t *testing.T, song model.Song, newSong model.Song, arrangementID uuid.UUID) {
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
	for _, part := range newSong.Parts {
		arrangementIndex := slices.IndexFunc(
			part.ArrangementOccurrences,
			func(occ model.SongPartOccurrences) bool {
				return occ.ArrangementID == *arrangementID
			})
		totalOccurrences += part.ArrangementOccurrences[arrangementIndex].Occurrences
		arrangementOccurrencesMap[part.ID] = part.ArrangementOccurrences[arrangementIndex].Occurrences
	}

	if totalOccurrences == 0 { // also nothing changed
		assert.Equal(t, song, newSong)
		return
	}

	for i, newPart := range newSong.Parts {
		oldPart := song.Parts[i]

		if arrangementOccurrencesMap[newPart.ID] == 0 { // nothing changed on this section
			assert.Equal(t, oldPart, newPart)
			continue
		}

		newRehearsals := oldPart.Rehearsals + arrangementOccurrencesMap[newPart.ID]*uint(duplicates+1)
		assert.Equal(t, newRehearsals, newPart.Rehearsals)
		assert.Greater(t, newPart.RehearsalsScore, oldPart.RehearsalsScore)
		assert.GreaterOrEqual(t, newPart.Progress, oldPart.Progress)

		for j := 0; j <= duplicates; j++ {
			fromDiff := arrangementOccurrencesMap[newPart.ID] * uint(duplicates-j)
			toDiff := arrangementOccurrencesMap[newPart.ID] * uint(j)
			assert.Equal(t, oldPart.Rehearsals+fromDiff, newPart.History[j].From)
			assert.Equal(t, newPart.Rehearsals-toDiff, newPart.History[j].To)
		}
	}

	assert.Greater(t, newSong.Rehearsals, song.Rehearsals)
	assert.Greater(t, newSong.Progress, song.Progress)

	assert.NotNil(t, newSong.LastTimePlayed)
	assert.WithinDuration(t, time.Now(), *newSong.LastTimePlayed, 1*time.Minute)
}
