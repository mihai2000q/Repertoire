import AlbumProperty from '../../types/enums/AlbumProperty.ts'
import Filter, { FilterValue } from '../../types/Filter.ts'
import FilterOperator from '../../types/enums/FilterOperator.ts'
import { AlbumFiltersMetadata } from '../../types/models/FiltersMetadata.ts'

const albumsFilters: Filter[] = [
  {
    property: AlbumProperty.ArtistId,
    operator: FilterOperator.Equal,
    isSet: false,
    value: undefined
  },

  {
    property: AlbumProperty.ReleaseDate,
    operator: FilterOperator.GreaterThanOrEqual,
    isSet: false
  },
  { property: AlbumProperty.ReleaseDate, operator: FilterOperator.LessThanOrEqual, isSet: false },
  { property: AlbumProperty.ReleaseDate, operator: FilterOperator.IsNull, isSet: false },
  { property: AlbumProperty.ReleaseDate, operator: FilterOperator.IsNotNull, isSet: false },

  { property: AlbumProperty.Songs, operator: FilterOperator.GreaterThanOrEqual, isSet: false },
  { property: AlbumProperty.Songs, operator: FilterOperator.LessThanOrEqual, isSet: false },

  { property: AlbumProperty.Rehearsals, operator: FilterOperator.GreaterThanOrEqual, isSet: false },
  { property: AlbumProperty.Rehearsals, operator: FilterOperator.LessThanOrEqual, isSet: false },

  { property: AlbumProperty.Confidence, operator: FilterOperator.GreaterThanOrEqual, isSet: false },
  { property: AlbumProperty.Confidence, operator: FilterOperator.LessThanOrEqual, isSet: false },

  { property: AlbumProperty.Progress, operator: FilterOperator.GreaterThanOrEqual, isSet: false },
  { property: AlbumProperty.Progress, operator: FilterOperator.LessThanOrEqual, isSet: false },

  { property: AlbumProperty.LastPlayed, operator: FilterOperator.GreaterThanOrEqual, isSet: false },
  { property: AlbumProperty.LastPlayed, operator: FilterOperator.LessThanOrEqual, isSet: false },
  { property: AlbumProperty.LastPlayed, operator: FilterOperator.IsNull, isSet: false },
  { property: AlbumProperty.LastPlayed, operator: FilterOperator.IsNotNull, isSet: false }
]

export const albumsFiltersMetadataMap: (
  metadata: AlbumFiltersMetadata
) => [string, FilterValue][] = (metadata) => [
  [AlbumProperty.ReleaseDate + FilterOperator.GreaterThanOrEqual, metadata.minReleaseDate],
  [AlbumProperty.ReleaseDate + FilterOperator.LessThanOrEqual, metadata.maxReleaseDate],

  [AlbumProperty.Songs + FilterOperator.GreaterThanOrEqual, metadata.minSongsCount],
  [AlbumProperty.Songs + FilterOperator.LessThanOrEqual, metadata.maxSongsCount],

  [AlbumProperty.Rehearsals + FilterOperator.GreaterThanOrEqual, metadata.minRehearsals],
  [AlbumProperty.Rehearsals + FilterOperator.LessThanOrEqual, metadata.maxRehearsals],

  [AlbumProperty.Confidence + FilterOperator.GreaterThanOrEqual, metadata.minConfidence],
  [AlbumProperty.Confidence + FilterOperator.LessThanOrEqual, metadata.maxConfidence],

  [AlbumProperty.Progress + FilterOperator.GreaterThanOrEqual, metadata.minProgress],
  [AlbumProperty.Progress + FilterOperator.LessThanOrEqual, metadata.maxProgress],

  [AlbumProperty.LastPlayed + FilterOperator.GreaterThanOrEqual, metadata.minLastTimePlayed],
  [AlbumProperty.LastPlayed + FilterOperator.LessThanOrEqual, metadata.maxLastTimePlayed]
]

export default albumsFilters
