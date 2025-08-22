import ArtistProperty from '../../types/enums/properties/ArtistProperty.ts'
import Filter, { FilterValue } from '../../types/Filter.ts'
import FilterOperator from '../../types/enums/FilterOperator.ts'
import { ArtistFiltersMetadata } from '../../types/models/FiltersMetadata.ts'

const artistsFilters: Filter[] = [
  {
    property: ArtistProperty.BandMembers,
    operator: FilterOperator.GreaterThanOrEqual,
    isSet: false
  },
  { property: ArtistProperty.BandMembers, operator: FilterOperator.LessThanOrEqual, isSet: false },

  { property: ArtistProperty.Band, operator: FilterOperator.Equal, isSet: false, value: true },
  { property: ArtistProperty.Band, operator: FilterOperator.NotEqual, isSet: false, value: true },

  { property: ArtistProperty.Albums, operator: FilterOperator.GreaterThanOrEqual, isSet: false },
  { property: ArtistProperty.Albums, operator: FilterOperator.LessThanOrEqual, isSet: false },

  { property: ArtistProperty.Songs, operator: FilterOperator.GreaterThanOrEqual, isSet: false },
  { property: ArtistProperty.Songs, operator: FilterOperator.LessThanOrEqual, isSet: false },

  {
    property: ArtistProperty.Rehearsals,
    operator: FilterOperator.GreaterThanOrEqual,
    isSet: false
  },
  { property: ArtistProperty.Rehearsals, operator: FilterOperator.LessThanOrEqual, isSet: false },

  {
    property: ArtistProperty.Confidence,
    operator: FilterOperator.GreaterThanOrEqual,
    isSet: false
  },
  { property: ArtistProperty.Confidence, operator: FilterOperator.LessThanOrEqual, isSet: false },

  { property: ArtistProperty.Progress, operator: FilterOperator.GreaterThanOrEqual, isSet: false },
  { property: ArtistProperty.Progress, operator: FilterOperator.LessThanOrEqual, isSet: false },

  {
    property: ArtistProperty.LastPlayed,
    operator: FilterOperator.GreaterThanOrEqual,
    isSet: false
  },
  { property: ArtistProperty.LastPlayed, operator: FilterOperator.LessThanOrEqual, isSet: false },
  { property: ArtistProperty.LastPlayed, operator: FilterOperator.IsNull, isSet: false },
  { property: ArtistProperty.LastPlayed, operator: FilterOperator.IsNotNull, isSet: false }
]

export const artistsFiltersMetadataMap: (
  metadata: ArtistFiltersMetadata
) => [string, FilterValue][] = (metadata) => [
  [ArtistProperty.BandMembers + FilterOperator.GreaterThanOrEqual, metadata.minBandMembersCount],
  [ArtistProperty.BandMembers + FilterOperator.LessThanOrEqual, metadata.maxBandMembersCount],

  [ArtistProperty.Albums + FilterOperator.GreaterThanOrEqual, metadata.minAlbumsCount],
  [ArtistProperty.Albums + FilterOperator.LessThanOrEqual, metadata.maxAlbumsCount],

  [ArtistProperty.Songs + FilterOperator.GreaterThanOrEqual, metadata.minSongsCount],
  [ArtistProperty.Songs + FilterOperator.LessThanOrEqual, metadata.maxSongsCount],

  [ArtistProperty.Rehearsals + FilterOperator.GreaterThanOrEqual, metadata.minRehearsals],
  [ArtistProperty.Rehearsals + FilterOperator.LessThanOrEqual, metadata.maxRehearsals],

  [ArtistProperty.Confidence + FilterOperator.GreaterThanOrEqual, metadata.minConfidence],
  [ArtistProperty.Confidence + FilterOperator.LessThanOrEqual, metadata.maxConfidence],

  [ArtistProperty.Progress + FilterOperator.GreaterThanOrEqual, metadata.minProgress],
  [ArtistProperty.Progress + FilterOperator.LessThanOrEqual, metadata.maxProgress],

  [ArtistProperty.LastPlayed + FilterOperator.GreaterThanOrEqual, metadata.minLastTimePlayed],
  [ArtistProperty.LastPlayed + FilterOperator.LessThanOrEqual, metadata.maxLastTimePlayed]
]

export default artistsFilters
