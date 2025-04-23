package assertion

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"repertoire/server/internal/enums"
	"repertoire/server/model"
	"testing"
	"time"
)

func AlbumFiltersMetadata(t *testing.T, metadata model.AlbumFiltersMetadata, albums []model.Album) {
	artistIDsMap := make(map[uuid.UUID]bool)

	var minReleaseDate *time.Time
	var maxReleaseDate *time.Time

	var minSongsCount *int64
	var maxSongsCount int64 = 0

	var minRehearsals *float64
	var maxRehearsals float64 = 0

	var minConfidence *float64
	var maxConfidence float64 = 0

	var minProgress *float64
	var maxProgress float64 = 0

	var minLastTimePlayed *time.Time
	var maxLastTimePlayed *time.Time

	for _, album := range albums {
		if album.ArtistID != nil {
			artistIDsMap[*album.ArtistID] = true
		}

		if album.ReleaseDate != nil && minReleaseDate != nil && album.ReleaseDate.Before(*minReleaseDate) ||
			album.ReleaseDate != nil && minReleaseDate == nil {
			minReleaseDate = album.ReleaseDate
		}
		if album.ReleaseDate != nil && maxReleaseDate != nil && album.ReleaseDate.After(*maxReleaseDate) ||
			album.ReleaseDate != nil && maxReleaseDate == nil {
			maxReleaseDate = album.ReleaseDate
		}

		songsCount := int64(len(album.Songs))
		if minSongsCount == nil || songsCount < *minSongsCount {
			minSongsCount = &songsCount
		}
		if songsCount > maxSongsCount {
			maxSongsCount = songsCount
		}

		var rehearsals float64 = 0
		var confidence float64 = 0
		var progress float64 = 0
		for _, song := range album.Songs {
			if song.LastTimePlayed != nil && minLastTimePlayed != nil && song.LastTimePlayed.Before(*minLastTimePlayed) ||
				song.LastTimePlayed != nil && minLastTimePlayed == nil {
				minLastTimePlayed = song.LastTimePlayed
			}
			if song.LastTimePlayed != nil && maxLastTimePlayed != nil && song.LastTimePlayed.After(*maxLastTimePlayed) ||
				song.LastTimePlayed != nil && maxLastTimePlayed == nil {
				maxLastTimePlayed = song.LastTimePlayed
			}

			rehearsals += song.Rehearsals
			confidence += song.Confidence
			progress += song.Progress
		}
		if songsCount != 0 {
			if rehearsals != 0 {
				rehearsals = rehearsals / float64(songsCount)
				if minRehearsals == nil || rehearsals < *minRehearsals {
					minRehearsals = &rehearsals
				}
				if rehearsals > maxRehearsals {
					maxRehearsals = rehearsals
				}
			}
			if confidence != 0 {
				confidence = confidence / float64(songsCount)
				if minConfidence == nil || confidence < *minConfidence {
					minConfidence = &confidence
				}
				if confidence > maxConfidence {
					maxConfidence = confidence
				}
			}
			if progress != 0 {
				progress = progress / float64(songsCount)
				if minProgress == nil || progress < *minProgress {
					minProgress = &progress
				}
				if progress > maxProgress {
					maxProgress = progress
				}
			}
		}
	}

	var artistIDs []uuid.UUID
	for key := range artistIDsMap {
		artistIDs = append(artistIDs, key)
	}

	assert.ElementsMatch(t, artistIDs, metadata.ArtistIDs)

	Time(t, minReleaseDate, metadata.MinReleaseDate)
	Time(t, maxReleaseDate, metadata.MaxReleaseDate)

	if minSongsCount == nil {
		assert.Zero(t, metadata.MinSongsCount)
	} else {
		assert.Equal(t, *minSongsCount, metadata.MinSongsCount)
	}
	assert.Equal(t, maxSongsCount, metadata.MaxSongsCount)

	if minRehearsals == nil {
		assert.Zero(t, metadata.MinRehearsals)
	} else {
		assert.Equal(t, *minRehearsals, metadata.MinRehearsals)
	}
	assert.Equal(t, maxRehearsals, metadata.MaxRehearsals)

	if minConfidence == nil {
		assert.Zero(t, metadata.MinConfidence)
	} else {
		assert.Equal(t, *minConfidence, metadata.MinConfidence)
	}
	assert.Equal(t, maxConfidence, metadata.MaxConfidence)

	if minProgress == nil {
		assert.Zero(t, metadata.MinProgress)
	} else {
		assert.Equal(t, *minProgress, metadata.MinProgress)
	}
	assert.Equal(t, maxProgress, metadata.MaxProgress)

	Time(t, minLastTimePlayed, metadata.MinLastTimePlayed)
	Time(t, maxLastTimePlayed, metadata.MaxLastTimePlayed)
}

func ArtistFiltersMetadata(t *testing.T, metadata model.ArtistFiltersMetadata, artists []model.Artist) {
	var minBandMembersCount *int64
	var maxBandMembersCount int64 = 0

	var minAlbumsCount *int64
	var maxAlbumsCount int64 = 0

	var minSongsCount *int64
	var maxSongsCount int64 = 0

	var minRehearsals *float64
	var maxRehearsals float64 = 0

	var minConfidence *float64
	var maxConfidence float64 = 0

	var minProgress *float64
	var maxProgress float64 = 0

	var minLastTimePlayed *time.Time
	var maxLastTimePlayed *time.Time

	for _, artist := range artists {
		bandMembersCount := int64(len(artist.BandMembers))
		if minBandMembersCount == nil || bandMembersCount < *minBandMembersCount {
			minBandMembersCount = &bandMembersCount
		}
		if bandMembersCount > maxBandMembersCount {
			maxBandMembersCount = bandMembersCount
		}

		albumsCount := int64(len(artist.Albums))
		if minAlbumsCount == nil || albumsCount < *minAlbumsCount {
			minAlbumsCount = &albumsCount
		}
		if albumsCount > maxAlbumsCount {
			maxAlbumsCount = albumsCount
		}

		songsCount := int64(len(artist.Songs))
		if minSongsCount == nil || songsCount < *minSongsCount {
			minSongsCount = &songsCount
		}
		if songsCount > maxSongsCount {
			maxSongsCount = songsCount
		}

		var rehearsals float64 = 0
		var confidence float64 = 0
		var progress float64 = 0
		for _, song := range artist.Songs {
			if song.LastTimePlayed != nil && minLastTimePlayed != nil && song.LastTimePlayed.Before(*minLastTimePlayed) ||
				song.LastTimePlayed != nil && minLastTimePlayed == nil {
				minLastTimePlayed = song.LastTimePlayed
			}
			if song.LastTimePlayed != nil && maxLastTimePlayed != nil && song.LastTimePlayed.After(*maxLastTimePlayed) ||
				song.LastTimePlayed != nil && maxLastTimePlayed == nil {
				maxLastTimePlayed = song.LastTimePlayed
			}

			rehearsals += song.Rehearsals
			confidence += song.Confidence
			progress += song.Progress
		}
		if songsCount != 0 {
			if rehearsals != 0 {
				rehearsals = rehearsals / float64(songsCount)
				if minRehearsals == nil || rehearsals < *minRehearsals {
					minRehearsals = &rehearsals
				}
				if rehearsals > maxRehearsals {
					maxRehearsals = rehearsals
				}
			}
			if confidence != 0 {
				confidence = confidence / float64(songsCount)
				if minConfidence == nil || confidence < *minConfidence {
					minConfidence = &confidence
				}
				if confidence > maxConfidence {
					maxConfidence = confidence
				}
			}
			if progress != 0 {
				progress = progress / float64(songsCount)
				if minProgress == nil || progress < *minProgress {
					minProgress = &progress
				}
				if progress > maxProgress {
					maxProgress = progress
				}
			}
		}
	}

	if minBandMembersCount == nil {
		assert.Zero(t, metadata.MinBandMembersCount)
	} else {
		assert.Equal(t, *minBandMembersCount, metadata.MinBandMembersCount)
	}
	assert.Equal(t, maxBandMembersCount, metadata.MaxBandMembersCount)

	if minAlbumsCount == nil {
		assert.Zero(t, metadata.MinAlbumsCount)
	} else {
		assert.Equal(t, *minAlbumsCount, metadata.MinAlbumsCount)
	}
	assert.Equal(t, maxAlbumsCount, metadata.MaxAlbumsCount)

	if minSongsCount == nil {
		assert.Zero(t, metadata.MinSongsCount)
	} else {
		assert.Equal(t, *minSongsCount, metadata.MinSongsCount)
	}
	assert.Equal(t, maxSongsCount, metadata.MaxSongsCount)

	if minRehearsals == nil {
		assert.Zero(t, metadata.MinRehearsals)
	} else {
		assert.Equal(t, *minRehearsals, metadata.MinRehearsals)
	}
	assert.Equal(t, maxRehearsals, metadata.MaxRehearsals)

	if minConfidence == nil {
		assert.Zero(t, metadata.MinConfidence)
	} else {
		assert.Equal(t, *minConfidence, metadata.MinConfidence)
	}
	assert.Equal(t, maxConfidence, metadata.MaxConfidence)

	if minProgress == nil {
		assert.Zero(t, metadata.MinProgress)
	} else {
		assert.Equal(t, *minProgress, metadata.MinProgress)
	}
	assert.Equal(t, maxProgress, metadata.MaxProgress)

	Time(t, minLastTimePlayed, metadata.MinLastTimePlayed)
	Time(t, maxLastTimePlayed, metadata.MaxLastTimePlayed)
}

func PlaylistFiltersMetadata(t *testing.T, metadata model.PlaylistFiltersMetadata, playlists []model.Playlist) {
	var minSongsCount *int64
	var maxSongsCount int64 = 0

	for _, playlist := range playlists {
		songsCount := int64(len(playlist.Songs))
		if minSongsCount == nil || songsCount < *minSongsCount {
			minSongsCount = &songsCount
		}
		if songsCount > maxSongsCount {
			maxSongsCount = songsCount
		}
	}

	if minSongsCount == nil {
		assert.Zero(t, metadata.MinSongsCount)
	} else {
		assert.Equal(t, *minSongsCount, metadata.MinSongsCount)
	}
	assert.Equal(t, maxSongsCount, metadata.MaxSongsCount)
}

func SongFiltersMetadata(t *testing.T, metadata model.SongFiltersMetadata, songs []model.Song) {
	artistIDsMap := make(map[uuid.UUID]bool)
	albumIDsMap := make(map[uuid.UUID]bool)
	guitarTuningIDsMap := make(map[uuid.UUID]bool)
	difficultiesMap := make(map[enums.Difficulty]bool)

	var minReleaseDate *time.Time
	var maxReleaseDate *time.Time

	var minBpm *uint
	var maxBpm *uint

	var minSectionsCount *int64
	var maxSectionsCount int64 = 0

	var minSolosCount *int64
	var maxSolosCount int64 = 0

	var minRiffsCount *int64
	var maxRiffsCount int64 = 0

	var minRehearsals *float64
	var maxRehearsals float64 = 0

	var minConfidence *float64
	var maxConfidence float64 = 0

	var minProgress *float64
	var maxProgress float64 = 0

	var minLastTimePlayed *time.Time
	var maxLastTimePlayed *time.Time

	for _, song := range songs {
		if song.ArtistID != nil {
			artistIDsMap[*song.ArtistID] = true
		}
		if song.AlbumID != nil {
			albumIDsMap[*song.AlbumID] = true
		}
		if song.GuitarTuningID != nil {
			guitarTuningIDsMap[*song.GuitarTuningID] = true
		}
		if song.Difficulty != nil {
			difficultiesMap[*song.Difficulty] = true
		}

		if song.ReleaseDate != nil && minReleaseDate != nil && song.ReleaseDate.Before(*minReleaseDate) ||
			song.ReleaseDate != nil && minReleaseDate == nil {
			minReleaseDate = song.ReleaseDate
		}
		if song.ReleaseDate != nil && maxReleaseDate != nil && song.ReleaseDate.After(*maxReleaseDate) ||
			song.ReleaseDate != nil && maxReleaseDate == nil {
			maxReleaseDate = song.ReleaseDate
		}

		if song.Bpm != nil && minBpm != nil && *song.Bpm < *minBpm ||
			song.Bpm != nil && minBpm == nil {
			minBpm = song.Bpm
		}
		if song.Bpm != nil && maxBpm != nil && *song.Bpm > *maxBpm ||
			song.Bpm != nil && maxBpm == nil {
			maxBpm = song.Bpm
		}

		if minRehearsals == nil || *minRehearsals > song.Rehearsals {
			minRehearsals = &song.Rehearsals
		}
		if maxRehearsals < song.Rehearsals {
			maxRehearsals = song.Rehearsals
		}

		if minConfidence == nil || *minConfidence > song.Confidence {
			minConfidence = &song.Confidence
		}
		if maxConfidence < song.Confidence {
			maxConfidence = song.Confidence
		}

		if minProgress == nil || *minProgress > song.Progress {
			minProgress = &song.Progress
		}
		if maxProgress < song.Progress {
			maxProgress = song.Progress
		}

		if song.LastTimePlayed != nil && minLastTimePlayed != nil && song.LastTimePlayed.Before(*minLastTimePlayed) ||
			song.LastTimePlayed != nil && minLastTimePlayed == nil {
			minLastTimePlayed = song.LastTimePlayed
		}
		if song.LastTimePlayed != nil && maxLastTimePlayed != nil && song.LastTimePlayed.After(*maxLastTimePlayed) ||
			song.LastTimePlayed != nil && maxLastTimePlayed == nil {
			maxLastTimePlayed = song.LastTimePlayed
		}

		var sectionsCount int64 = 0
		var solosCount int64 = 0
		var riffsCount int64 = 0
		for _, section := range song.Sections {
			sectionsCount++
			if section.SongSectionType.Name == "Solo" {
				solosCount++
			}
			if section.SongSectionType.Name == "Riff" {
				riffsCount++
			}
		}

		if minSectionsCount == nil || *minSectionsCount > sectionsCount {
			minSectionsCount = &sectionsCount
		}
		if sectionsCount > maxSectionsCount {
			maxSectionsCount = sectionsCount
		}
		if minSolosCount == nil || *minSolosCount > solosCount {
			minSolosCount = &solosCount
		}
		if solosCount > maxSolosCount {
			maxSolosCount = solosCount
		}
		if minRiffsCount == nil || *minRiffsCount > riffsCount {
			minRiffsCount = &riffsCount
		}
		if riffsCount > maxRiffsCount {
			maxRiffsCount = riffsCount
		}
	}

	var artistIDs []uuid.UUID
	for key := range artistIDsMap {
		artistIDs = append(artistIDs, key)
	}
	var albumIDs []uuid.UUID
	for key := range albumIDsMap {
		albumIDs = append(albumIDs, key)
	}
	var guitarTuningIDs []uuid.UUID
	for key := range guitarTuningIDsMap {
		guitarTuningIDs = append(guitarTuningIDs, key)
	}
	var difficulties []enums.Difficulty
	for key := range difficultiesMap {
		difficulties = append(difficulties, key)
	}

	assert.ElementsMatch(t, artistIDs, metadata.ArtistIDs)
	assert.ElementsMatch(t, albumIDs, metadata.AlbumIDs)
	assert.ElementsMatch(t, guitarTuningIDs, metadata.GuitarTuningIDs)
	assert.ElementsMatch(t, difficulties, metadata.Difficulties)

	Time(t, minReleaseDate, metadata.MinReleaseDate)
	Time(t, maxReleaseDate, metadata.MaxReleaseDate)

	assert.Equal(t, minBpm, metadata.MinBpm)
	assert.Equal(t, maxBpm, metadata.MaxBpm)

	if minSectionsCount == nil {
		assert.Zero(t, metadata.MinSectionsCount)
	} else {
		assert.Equal(t, *minSectionsCount, metadata.MinSectionsCount)
	}
	assert.Equal(t, maxSectionsCount, metadata.MaxSectionsCount)

	if minSolosCount == nil {
		assert.Zero(t, metadata.MinSolosCount)
	} else {
		assert.Equal(t, *minSolosCount, metadata.MinSolosCount)
	}
	assert.Equal(t, maxSolosCount, metadata.MaxSolosCount)

	if minRiffsCount == nil {
		assert.Zero(t, metadata.MinRiffsCount)
	} else {
		assert.Equal(t, *minRiffsCount, metadata.MinRiffsCount)
	}
	assert.Equal(t, maxRiffsCount, metadata.MaxRiffsCount)

	if minRehearsals == nil {
		assert.Zero(t, metadata.MinRehearsals)
	} else {
		assert.Equal(t, *minRehearsals, metadata.MinRehearsals)
	}
	assert.Equal(t, maxRehearsals, metadata.MaxRehearsals)

	if minConfidence == nil {
		assert.Zero(t, metadata.MinConfidence)
	} else {
		assert.Equal(t, *minConfidence, metadata.MinConfidence)
	}
	assert.Equal(t, maxConfidence, metadata.MaxConfidence)

	if minProgress == nil {
		assert.Zero(t, metadata.MinProgress)
	} else {
		assert.Equal(t, *minProgress, metadata.MinProgress)
	}
	assert.Equal(t, maxProgress, metadata.MaxProgress)

	Time(t, minLastTimePlayed, metadata.MinLastTimePlayed)
	Time(t, maxLastTimePlayed, metadata.MaxLastTimePlayed)
}
