import PlaylistProperty from '../../types/enums/PlaylistProperty.ts'
import Filter, { FilterValue } from '../../types/Filter.ts'
import FilterOperator from '../../types/enums/FilterOperator.ts'
import { PlaylistFiltersMetadata } from '../../types/models/FiltersMetadata.ts'

const playlistsFilters: Filter[] = [
  { property: PlaylistProperty.Songs, operator: FilterOperator.GreaterThanOrEqual, isSet: false },
  { property: PlaylistProperty.Songs, operator: FilterOperator.LessThanOrEqual, isSet: false }
]

export const playlistsFiltersMetadataMap: (
  metadata: PlaylistFiltersMetadata
) => [string, FilterValue][] = (metadata) => [
  [PlaylistProperty.Songs + FilterOperator.GreaterThanOrEqual, metadata.minSongsCount],
  [PlaylistProperty.Songs + FilterOperator.LessThanOrEqual, metadata.maxSongsCount]
]

export default playlistsFilters
