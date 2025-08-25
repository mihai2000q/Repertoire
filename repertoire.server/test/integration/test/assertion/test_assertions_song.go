package assertion

import (
	"repertoire/server/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func PerfectSongRehearsal(t *testing.T, song model.Song, newSong model.Song) {
	for i, section := range newSong.Sections {
		if section.Occurrences == 0 { // nothing changed
			newSong.Sections[i].History = nil
			assert.Equal(t, song.Sections[i], newSong.Sections[i])
			continue
		}

		assert.Equal(t, section.Rehearsals, song.Sections[i].Rehearsals+song.Sections[i].Occurrences)
		assert.Greater(t, section.RehearsalsScore, song.Sections[i].RehearsalsScore)
		assert.Greater(t, section.Progress, song.Sections[i].Progress)

		assert.NotEmpty(t, section.History[len(section.History)-1].ID)
		assert.Equal(t, song.Sections[i].Rehearsals, section.History[len(section.History)-1].From)
		assert.Equal(t, section.Rehearsals, section.History[len(section.History)-1].To)
		assert.Equal(t, model.RehearsalsProperty, section.History[len(section.History)-1].Property)
	}

	assert.Greater(t, newSong.Rehearsals, song.Rehearsals)
	assert.Greater(t, newSong.Progress, song.Progress)

	assert.NotNil(t, newSong.LastTimePlayed)
	assert.WithinDuration(t, time.Now(), *newSong.LastTimePlayed, 1*time.Minute)
}
