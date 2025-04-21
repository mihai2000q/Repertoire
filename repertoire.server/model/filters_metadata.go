package model

import (
	"github.com/google/uuid"
	"time"
)

type AlbumFiltersMetadata struct {
	ArtistIDsAgg string      `gorm:"->; column:artist_ids" json:"-"`
	ArtistIDs    []uuid.UUID `gorm:"-" json:"artistIds"`

	MinReleaseDate *time.Time `gorm:"->" json:"minReleaseDate"`
	MaxReleaseDate *time.Time `gorm:"->" json:"maxReleaseDate"`

	MinSongsCount int64 `gorm:"->" json:"minSongsCount"`
	MaxSongsCount int64 `gorm:"->" json:"maxSongsCount"`

	MinRehearsals float64 `gorm:"->" json:"minRehearsals"`
	MaxRehearsals float64 `gorm:"->" json:"maxRehearsals"`

	MinConfidence float64 `gorm:"->" json:"minConfidence"`
	MaxConfidence float64 `gorm:"->" json:"maxConfidence"`

	MinProgress float64 `gorm:"->" json:"minProgress"`
	MaxProgress float64 `gorm:"->" json:"maxProgress"`

	MinLastTimePlayed *time.Time `gorm:"->" json:"minLastTimePlayed"`
	MaxLastTimePlayed *time.Time `gorm:"->" json:"maxLastTimePlayed"`
}

type ArtistFiltersMetadata struct {
	MinBandMembersCount int64 `gorm:"->" json:"minBandMembersCount"`
	MaxBandMembersCount int64 `gorm:"->" json:"maxBandMembersCount"`

	MinAlbumsCount int64 `gorm:"->" json:"minAlbumsCount"`
	MaxAlbumsCount int64 `gorm:"->" json:"maxAlbumsCount"`

	MinSongsCount int64 `gorm:"->" json:"minSongsCount"`
	MaxSongsCount int64 `gorm:"->" json:"maxSongsCount"`

	MinRehearsals float64 `gorm:"->" json:"minRehearsals"`
	MaxRehearsals float64 `gorm:"->" json:"maxRehearsals"`

	MinConfidence float64 `gorm:"->" json:"minConfidence"`
	MaxConfidence float64 `gorm:"->" json:"maxConfidence"`

	MinProgress float64 `gorm:"->" json:"minProgress"`
	MaxProgress float64 `gorm:"->" json:"maxProgress"`

	MinLastTimePlayed *time.Time `gorm:"->" json:"minLastTimePlayed"`
	MaxLastTimePlayed *time.Time `gorm:"->" json:"maxLastTimePlayed"`
}

type PlaylistFiltersMetadata struct {
	MinSongsCount int64 `gorm:"->" json:"minSongsCount"`
	MaxSongsCount int64 `gorm:"->" json:"maxSongsCount"`
}

type SongFiltersMetadata struct {
	ArtistIDsAgg string      `gorm:"->; column:artist_ids" json:"-"`
	ArtistIDs    []uuid.UUID `gorm:"-" json:"artistIds"`

	AlbumIDsAgg string      `gorm:"->; column:album_ids" json:"-"`
	AlbumIDs    []uuid.UUID `gorm:"-" json:"albumIds"`

	MinReleaseDate *time.Time `gorm:"->" json:"minReleaseDate"`
	MaxReleaseDate *time.Time `gorm:"->" json:"maxReleaseDate"`

	MinBpm uint `gorm:"->" json:"minBpm"`
	MaxBpm uint `gorm:"->" json:"maxBpm"`

	DifficultiesAgg string   `gorm:"->; column:difficulties" json:"-"`
	Difficulties    []string `gorm:"-" json:"difficulties"`

	GuitarTuningIDsAgg string      `gorm:"->; column:guitar_tuning_ids" json:"-"`
	GuitarTuningIDs    []uuid.UUID `gorm:"-" json:"guitarTuningIds"`

	MinSectionsCount int64 `gorm:"->" json:"minSectionsCount"`
	MaxSectionsCount int64 `gorm:"->" json:"maxSectionsCount"`

	MinSolosCount int64 `gorm:"->" json:"minSolosCount"`
	MaxSolosCount int64 `gorm:"->" json:"maxSolosCount"`

	MinRiffsCount int64 `gorm:"->" json:"minRiffsCount"`
	MaxRiffsCount int64 `gorm:"->" json:"maxRiffsCount"`

	MinRehearsals float64 `gorm:"->" json:"minRehearsals"`
	MaxRehearsals float64 `gorm:"->" json:"maxRehearsals"`

	MinConfidence float64 `gorm:"->" json:"minConfidence"`
	MaxConfidence float64 `gorm:"->" json:"maxConfidence"`

	MinProgress float64 `gorm:"->" json:"minProgress"`
	MaxProgress float64 `gorm:"->" json:"maxProgress"`

	MinLastTimePlayed *time.Time `gorm:"->" json:"minLastTimePlayed"`
	MaxLastTimePlayed *time.Time `gorm:"->" json:"maxLastTimePlayed"`
}
