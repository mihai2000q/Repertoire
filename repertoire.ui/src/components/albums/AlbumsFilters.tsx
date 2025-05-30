import Filter from '../../types/Filter.ts'
import { RangeSlider, Stack, Text } from '@mantine/core'
import { useEffect, useState } from 'react'
import { DatePickerInput } from '@mantine/dates'
import { IconCalendarCheck, IconCalendarRepeat } from '@tabler/icons-react'
import FiltersDrawer from '../@ui/drawer/FiltersDrawer.tsx'
import {
  useGetAlbumFiltersMetadataQuery,
  useLazyGetAlbumFiltersMetadataQuery
} from '../../state/api/albumsApi.ts'
import AlbumProperty from '../../types/enums/AlbumProperty.ts'
import { albumsFiltersMetadataMap } from '../../data/albums/albumsFilters.ts'
import FilterOperator from '../../types/enums/FilterOperator.ts'
import useFiltersMetadata from '../../hooks/filter/useFiltersMetadata.ts'
import useFiltersHandlers from '../../hooks/filter/useFiltersHandlers.ts'
import NumberInputRange from '../@ui/form/input/NumberInputRange.tsx'
import DoubleCheckbox from '../@ui/filter/DoubleCheckbox.tsx'
import useSearchBy from '../../hooks/api/useSearchBy.ts'

interface AlbumsFiltersProps {
  opened: boolean
  onClose: () => void
  filters: Map<string, Filter>
  setFilters: (filters: Map<string, Filter>, withSearchParams?: boolean) => void
  isAlbumsLoading?: boolean
}

function AlbumsFilters({
  opened,
  onClose,
  filters,
  setFilters,
  isAlbumsLoading
}: AlbumsFiltersProps) {
  const [getFiltersMetadata, { data: initialFiltersMetadata, isLoading }] =
    useLazyGetAlbumFiltersMetadataQuery()
  useEffect(() => {
    getFiltersMetadata({}, true)
  }, [])

  const searchBy = useSearchBy(filters)
  const { data: filtersMetadata } = useGetAlbumFiltersMetadataQuery({ searchBy: searchBy })

  const [internalFilters, setInternalFilters] = useState(filters)
  const initialFilters = useFiltersMetadata(
    initialFiltersMetadata,
    filtersMetadata,
    filters,
    setFilters,
    setInternalFilters,
    albumsFiltersMetadataMap
  )

  const {
    handleIsSetChange,
    handleValueChange,
    handleDoubleValueChange,
    getDateRangeValues,
    getSliderValues
  } = useFiltersHandlers(internalFilters, setInternalFilters, initialFilters)

  return (
    <FiltersDrawer
      opened={opened}
      onClose={onClose}
      filters={filters}
      setFilters={setFilters}
      internalFilters={internalFilters}
      setInternalFilters={setInternalFilters}
      initialFilters={initialFilters}
      isLoading={isAlbumsLoading}
    >
      <Stack px={'lg'} pb={'lg'}>
        <DatePickerInput
          flex={1}
          type={'range'}
          label={'Release Date'}
          placeholder={'Release Date'}
          valueFormat={'DD MMM YYYY'}
          leftSection={<IconCalendarRepeat size={20} />}
          disabled={isLoading}
          value={getDateRangeValues(
            AlbumProperty.ReleaseDate + FilterOperator.GreaterThanOrEqual,
            AlbumProperty.ReleaseDate + FilterOperator.LessThanOrEqual
          )}
          onChange={(values) =>
            handleDoubleValueChange(
              AlbumProperty.ReleaseDate + FilterOperator.GreaterThanOrEqual,
              values[0] ?? null,
              AlbumProperty.ReleaseDate + FilterOperator.LessThanOrEqual,
              values[1] ?? null
            )
          }
        />
        <DoubleCheckbox
          title={'Has Release Date?'}
          label1={'Yes'}
          checked1={internalFilters.get(AlbumProperty.ReleaseDate + FilterOperator.IsNotNull).isSet}
          onChange1={(value) =>
            handleIsSetChange(AlbumProperty.ReleaseDate + FilterOperator.IsNotNull, value)
          }
          label2={'No'}
          checked2={internalFilters.get(AlbumProperty.ReleaseDate + FilterOperator.IsNull).isSet}
          onChange2={(value) =>
            handleIsSetChange(AlbumProperty.ReleaseDate + FilterOperator.IsNull, value)
          }
          disabled={isLoading}
        />

        <NumberInputRange
          label={'Songs'}
          isLoading={isLoading}
          value1={
            internalFilters.get(AlbumProperty.Songs + FilterOperator.GreaterThanOrEqual).value as
              | string
              | number
          }
          onChange1={(value) =>
            handleValueChange(AlbumProperty.Songs + FilterOperator.GreaterThanOrEqual, value)
          }
          value2={
            internalFilters.get(AlbumProperty.Songs + FilterOperator.LessThanOrEqual).value as
              | string
              | number
          }
          onChange2={(value) =>
            handleValueChange(AlbumProperty.Songs + FilterOperator.LessThanOrEqual, value)
          }
          max={
            initialFilters.get(AlbumProperty.Songs + FilterOperator.LessThanOrEqual).value as number
          }
        />

        <NumberInputRange
          label={'Rehearsals'}
          isLoading={isLoading}
          value1={
            internalFilters.get(AlbumProperty.Rehearsals + FilterOperator.GreaterThanOrEqual)
              .value as string | number
          }
          onChange1={(value) =>
            handleValueChange(AlbumProperty.Rehearsals + FilterOperator.GreaterThanOrEqual, value)
          }
          value2={
            internalFilters.get(AlbumProperty.Rehearsals + FilterOperator.LessThanOrEqual).value as
              | string
              | number
          }
          onChange2={(value) =>
            handleValueChange(AlbumProperty.Rehearsals + FilterOperator.LessThanOrEqual, value)
          }
          max={
            initialFilters.get(AlbumProperty.Rehearsals + FilterOperator.LessThanOrEqual)
              .value as number
          }
        />

        <Stack gap={'xxs'}>
          <Text fw={500} fz={'sm'}>
            Confidence
          </Text>

          <RangeSlider
            disabled={isLoading}
            thumbFromLabel={'confidence-from'}
            thumbToLabel={'confidence-to'}
            label={(value) => `${value}%`}
            value={getSliderValues(
              AlbumProperty.Confidence + FilterOperator.GreaterThanOrEqual,
              AlbumProperty.Confidence + FilterOperator.LessThanOrEqual
            )}
            onChange={(values) =>
              handleDoubleValueChange(
                AlbumProperty.Confidence + FilterOperator.GreaterThanOrEqual,
                values[0],
                AlbumProperty.Confidence + FilterOperator.LessThanOrEqual,
                values[1]
              )
            }
          />
        </Stack>

        <NumberInputRange
          label={'Progress'}
          isLoading={isLoading}
          value1={
            internalFilters.get(AlbumProperty.Progress + FilterOperator.GreaterThanOrEqual)
              .value as string | number
          }
          onChange1={(value) =>
            handleValueChange(AlbumProperty.Progress + FilterOperator.GreaterThanOrEqual, value)
          }
          value2={
            internalFilters.get(AlbumProperty.Progress + FilterOperator.LessThanOrEqual).value as
              | string
              | number
          }
          onChange2={(value) =>
            handleValueChange(AlbumProperty.Progress + FilterOperator.LessThanOrEqual, value)
          }
          max={
            initialFilters.get(AlbumProperty.Progress + FilterOperator.LessThanOrEqual)
              .value as number
          }
        />

        <DatePickerInput
          flex={1}
          type={'range'}
          label={'Last Played'}
          placeholder={'Last Played'}
          valueFormat={'DD MMM YYYY'}
          leftSection={<IconCalendarCheck size={20} />}
          disabled={isLoading}
          value={getDateRangeValues(
            AlbumProperty.LastPlayed + FilterOperator.GreaterThanOrEqual,
            AlbumProperty.LastPlayed + FilterOperator.LessThanOrEqual
          )}
          onChange={(values) =>
            handleDoubleValueChange(
              AlbumProperty.LastPlayed + FilterOperator.GreaterThanOrEqual,
              values[0] ?? null,
              AlbumProperty.LastPlayed + FilterOperator.LessThanOrEqual,
              values[1] ?? null
            )
          }
        />
        <DoubleCheckbox
          title={'Has Been Played Before?'}
          label1={'Yes'}
          checked1={internalFilters.get(AlbumProperty.LastPlayed + FilterOperator.IsNotNull).isSet}
          onChange1={(value) =>
            handleIsSetChange(AlbumProperty.LastPlayed + FilterOperator.IsNotNull, value)
          }
          label2={'Never'}
          checked2={internalFilters.get(AlbumProperty.LastPlayed + FilterOperator.IsNull).isSet}
          onChange2={(value) =>
            handleIsSetChange(AlbumProperty.LastPlayed + FilterOperator.IsNull, value)
          }
          disabled={isLoading}
        />
      </Stack>
    </FiltersDrawer>
  )
}

export default AlbumsFilters
