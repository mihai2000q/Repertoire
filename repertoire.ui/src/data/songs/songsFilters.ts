import SongProperty from '../../types/enums/SongProperty.ts'
import Filter, { FilterValue } from '../../types/Filter.ts'
import FilterOperator from '../../types/enums/FilterOperator.ts'
import { SongFiltersMetadata } from '../../types/models/FiltersMetadata.ts'

const songsFilters: Filter[] = [
  { property: SongProperty.AlbumId, operator: FilterOperator.Equal, isSet: false },

  { property: SongProperty.ArtistId, operator: FilterOperator.Equal, isSet: false },

  { property: SongProperty.IsRecorded, operator: FilterOperator.Equal, isSet: false, value: true },
  {
    property: SongProperty.IsRecorded,
    operator: FilterOperator.NotEqual,
    isSet: false,
    value: true
  },

  { property: SongProperty.ReleaseDate, operator: FilterOperator.GreaterThanOrEqual, isSet: false },
  { property: SongProperty.ReleaseDate, operator: FilterOperator.LessThanOrEqual, isSet: false },
  { property: SongProperty.ReleaseDate, operator: FilterOperator.IsNull, isSet: false },
  { property: SongProperty.ReleaseDate, operator: FilterOperator.IsNotNull, isSet: false },

  { property: SongProperty.BPM, operator: FilterOperator.GreaterThanOrEqual, isSet: false },
  { property: SongProperty.BPM, operator: FilterOperator.LessThanOrEqual, isSet: false },
  { property: SongProperty.BPM, operator: FilterOperator.IsNull, isSet: false },
  { property: SongProperty.BPM, operator: FilterOperator.IsNotNull, isSet: false },

  { property: SongProperty.Difficulty, operator: FilterOperator.In, isSet: false },
  { property: SongProperty.Difficulty, operator: FilterOperator.IsNull, isSet: false },
  { property: SongProperty.Difficulty, operator: FilterOperator.IsNotNull, isSet: false },

  { property: SongProperty.GuitarTuningId, operator: FilterOperator.In, isSet: false },
  { property: SongProperty.GuitarTuningId, operator: FilterOperator.IsNull, isSet: false },
  { property: SongProperty.GuitarTuningId, operator: FilterOperator.IsNotNull, isSet: false },

  { property: SongProperty.InstrumentId, operator: FilterOperator.In, isSet: false },

  { property: SongProperty.Sections, operator: FilterOperator.GreaterThanOrEqual, isSet: false },
  { property: SongProperty.Sections, operator: FilterOperator.LessThanOrEqual, isSet: false },

  { property: SongProperty.Solos, operator: FilterOperator.GreaterThanOrEqual, isSet: false },
  { property: SongProperty.Solos, operator: FilterOperator.LessThanOrEqual, isSet: false },

  { property: SongProperty.Riffs, operator: FilterOperator.GreaterThanOrEqual, isSet: false },
  { property: SongProperty.Riffs, operator: FilterOperator.LessThanOrEqual, isSet: false },

  { property: SongProperty.Rehearsals, operator: FilterOperator.GreaterThanOrEqual, isSet: false },
  { property: SongProperty.Rehearsals, operator: FilterOperator.LessThanOrEqual, isSet: false },

  { property: SongProperty.Confidence, operator: FilterOperator.GreaterThanOrEqual, isSet: false },
  { property: SongProperty.Confidence, operator: FilterOperator.LessThanOrEqual, isSet: false },

  { property: SongProperty.Progress, operator: FilterOperator.GreaterThanOrEqual, isSet: false },
  { property: SongProperty.Progress, operator: FilterOperator.LessThanOrEqual, isSet: false },

  { property: SongProperty.LastPlayed, operator: FilterOperator.GreaterThanOrEqual, isSet: false },
  { property: SongProperty.LastPlayed, operator: FilterOperator.LessThanOrEqual, isSet: false },
  { property: SongProperty.LastPlayed, operator: FilterOperator.IsNull, isSet: false },
  { property: SongProperty.LastPlayed, operator: FilterOperator.IsNotNull, isSet: false }
]

export const songsFiltersMetadataMap: (metadata: SongFiltersMetadata) => [string, FilterValue][] = (
  metadata
) => [
  [SongProperty.ReleaseDate + FilterOperator.GreaterThanOrEqual, metadata.minReleaseDate],
  [SongProperty.ReleaseDate + FilterOperator.LessThanOrEqual, metadata.maxReleaseDate],

  [SongProperty.BPM + FilterOperator.GreaterThanOrEqual, metadata.minBpm],
  [SongProperty.BPM + FilterOperator.LessThanOrEqual, metadata.maxBpm],

  [SongProperty.Sections + FilterOperator.GreaterThanOrEqual, metadata.minSectionsCount],
  [SongProperty.Sections + FilterOperator.LessThanOrEqual, metadata.maxSectionsCount],

  [SongProperty.Solos + FilterOperator.GreaterThanOrEqual, metadata.minSolosCount],
  [SongProperty.Solos + FilterOperator.LessThanOrEqual, metadata.maxSolosCount],

  [SongProperty.Riffs + FilterOperator.GreaterThanOrEqual, metadata.minRiffsCount],
  [SongProperty.Riffs + FilterOperator.LessThanOrEqual, metadata.maxRiffsCount],

  [SongProperty.Rehearsals + FilterOperator.GreaterThanOrEqual, metadata.minRehearsals],
  [SongProperty.Rehearsals + FilterOperator.LessThanOrEqual, metadata.maxRehearsals],

  [SongProperty.Confidence + FilterOperator.GreaterThanOrEqual, metadata.minConfidence],
  [SongProperty.Confidence + FilterOperator.LessThanOrEqual, metadata.maxConfidence],

  [SongProperty.Progress + FilterOperator.GreaterThanOrEqual, metadata.minProgress],
  [SongProperty.Progress + FilterOperator.LessThanOrEqual, metadata.maxProgress],

  [SongProperty.LastPlayed + FilterOperator.GreaterThanOrEqual, metadata.minLastTimePlayed],
  [SongProperty.LastPlayed + FilterOperator.LessThanOrEqual, metadata.maxLastTimePlayed]
]

export default songsFilters
