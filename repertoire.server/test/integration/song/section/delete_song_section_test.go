package section

import (
	"net/http"
	"net/http/httptest"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/core"
	songData "repertoire/server/test/integration/test/data/song"
	"repertoire/server/test/integration/test/utils"
	"slices"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestDeleteSongSection_WhenSongIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().DELETE(w, "/api/songs/sections/"+uuid.New().String()+"/from/"+uuid.New().String())

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteSongSection_WhenSectionIsNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	song := songData.Songs[0]

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().DELETE(w, "/api/songs/sections/"+uuid.New().String()+"/from/"+song.ID.String())

	// then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteSongSection_WhenSuccessful_ShouldDeleteSection(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	// song with sections and previous stats
	song := songData.Songs[0]
	section := songData.SongSections[0]

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().DELETE(w, "/api/songs/sections/"+section.ID.String()+"/from/"+section.SongID.String())

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)

	var newSong model.Song
	db.Preload("Sections", func(db *gorm.DB) *gorm.DB {
		return db.Order("\"order\"")
	}).
		Find(&newSong, song.ID)

	assert.True(t,
		slices.IndexFunc(newSong.Sections, func(t model.SongSection) bool {
			return t.ID == section.ID
		}) == -1,
		"Song Section has not been deleted",
	)

	for i, s := range newSong.Sections {
		assert.Equal(t, uint(i), s.Order)
	}
}

func TestDeleteSongSection_WhenSuccessfulWithParts_ShouldDeleteSectionAndParts(t *testing.T) {
	// given
	utils.SeedAndCleanupData(t, songData.Users, songData.SeedData)

	song := songData.Songs[0]
	section := songData.SongSections[0]

	oldSectionParts := slices.Clone(slices.DeleteFunc(slices.Clone(songData.SongSectionParts), func(sp model.SongSectionPart) bool {
		return sp.SectionID != section.ID
	}))

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().DELETE(w, "/api/songs/sections/"+section.ID.String()+"/from/"+section.SongID.String()+"?withParts=true")

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	db := utils.GetDatabase(t)
	var newSong model.Song
	db.Preload("Parts", func(db *gorm.DB) *gorm.DB {
		return db.Order("song_order")
	}).
		Preload("Sections", func(db *gorm.DB) *gorm.DB {
			return db.Order("\"order\"")
		}).
		Find(&newSong, song.ID)

	assert.False(t,
		slices.ContainsFunc(newSong.Sections, func(t model.SongSection) bool {
			return t.ID == section.ID
		}),
		"Song Section has not been deleted",
	)

	for _, sp := range oldSectionParts {
		assert.False(t,
			slices.ContainsFunc(newSong.Parts, func(p model.SongPart) bool {
				return p.ID == sp.PartID
			}),
			"Part should have been deleted",
		)
	}

	// reordered sections
	for i, s := range newSong.Sections {
		assert.Equal(t, uint(i), s.Order)
	}

	// reordered parts
	for i, s := range newSong.Parts {
		assert.Equal(t, uint(i), s.SongOrder)
	}

	// stats decreased due to parts removal
	assert.Less(t, newSong.Confidence, song.Confidence)
	assert.Less(t, newSong.Rehearsals, song.Rehearsals)
	assert.Less(t, newSong.Progress, song.Progress)
}
