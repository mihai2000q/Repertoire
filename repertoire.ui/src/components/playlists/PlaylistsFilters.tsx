import Filter from '../../types/Filter.ts'
import { Stack } from '@mantine/core'
import { useEffect, useState } from 'react'
import FiltersDrawer from '../@ui/drawer/FiltersDrawer.tsx'
import {
  useGetPlaylistFiltersMetadataQuery,
  useLazyGetPlaylistFiltersMetadataQuery
} from '../../state/api/playlistsApi.ts'
import PlaylistProperty from '../../types/enums/properties/PlaylistProperty.ts'
import { playlistsFiltersMetadataMap } from '../../data/playlists/playlistsFilters.ts'
import FilterOperator from '../../types/enums/FilterOperator.ts'
import useFiltersMetadata from '../../hooks/filter/useFiltersMetadata.ts'
import useFiltersHandlers from '../../hooks/filter/useFiltersHandlers.ts'
import NumberInputRange from '../@ui/form/input/NumberInputRange.tsx'
import useSearchBy from '../../hooks/api/useSearchBy.ts'

interface PlaylistFiltersProps {
  opened: boolean
  onClose: () => void
  filters: Map<string, Filter>
  setFilters: (filters: Map<string, Filter>, withSearchParams?: boolean) => void
  isPlaylistsLoading?: boolean
}

function PlaylistsFilters({
  opened,
  onClose,
  filters,
  setFilters,
  isPlaylistsLoading
}: PlaylistFiltersProps) {
  const [getFiltersMetadata, { data: initialFiltersMetadata, isLoading }] =
    useLazyGetPlaylistFiltersMetadataQuery()
  useEffect(() => {
    getFiltersMetadata({}, true)
  }, [])

  const searchBy = useSearchBy(filters)
  const { data: filtersMetadata } = useGetPlaylistFiltersMetadataQuery({ searchBy: searchBy })

  const [internalFilters, setInternalFilters] = useState(filters)
  const initialFilters = useFiltersMetadata(
    initialFiltersMetadata,
    filtersMetadata,
    filters,
    setFilters,
    setInternalFilters,
    playlistsFiltersMetadataMap
  )

  const { handleValueChange } = useFiltersHandlers(
    internalFilters,
    setInternalFilters,
    initialFilters
  )

  return (
    <FiltersDrawer
      opened={opened}
      onClose={onClose}
      filters={filters}
      setFilters={setFilters}
      internalFilters={internalFilters}
      setInternalFilters={setInternalFilters}
      initialFilters={initialFilters}
      isLoading={isPlaylistsLoading}
    >
      <Stack px={'lg'} pb={'lg'}>
        <NumberInputRange
          label={'Songs'}
          isLoading={isLoading}
          value1={
            internalFilters.get(PlaylistProperty.Songs + FilterOperator.GreaterThanOrEqual)
              .value as string | number
          }
          onChange1={(value) =>
            handleValueChange(PlaylistProperty.Songs + FilterOperator.GreaterThanOrEqual, value)
          }
          value2={
            internalFilters.get(PlaylistProperty.Songs + FilterOperator.LessThanOrEqual).value as
              | string
              | number
          }
          onChange2={(value) =>
            handleValueChange(PlaylistProperty.Songs + FilterOperator.LessThanOrEqual, value)
          }
          max={
            initialFilters.get(PlaylistProperty.Songs + FilterOperator.LessThanOrEqual)
              .value as number
          }
        />
      </Stack>
    </FiltersDrawer>
  )
}

export default PlaylistsFilters
