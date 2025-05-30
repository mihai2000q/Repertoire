import Filter from '../../types/Filter.ts'
import { RangeSlider, Stack, Text } from '@mantine/core'
import { useEffect, useState } from 'react'
import { DatePickerInput } from '@mantine/dates'
import { IconCalendarCheck } from '@tabler/icons-react'
import FiltersDrawer from '../@ui/drawer/FiltersDrawer.tsx'
import {
  useGetArtistFiltersMetadataQuery,
  useLazyGetArtistFiltersMetadataQuery
} from '../../state/api/artistsApi.ts'
import ArtistProperty from '../../types/enums/ArtistProperty.ts'
import { artistsFiltersMetadataMap } from '../../data/artists/artistsFilters.ts'
import FilterOperator from '../../types/enums/FilterOperator.ts'
import useFiltersMetadata from '../../hooks/filter/useFiltersMetadata.ts'
import useFiltersHandlers from '../../hooks/filter/useFiltersHandlers.ts'
import NumberInputRange from '../@ui/form/input/NumberInputRange.tsx'
import DoubleCheckbox from '../@ui/filter/DoubleCheckbox.tsx'
import useSearchBy from '../../hooks/api/useSearchBy.ts'

interface ArtistsFiltersProps {
  opened: boolean
  onClose: () => void
  filters: Map<string, Filter>
  setFilters: (filters: Map<string, Filter>, withSearchParams?: boolean) => void
  isArtistsLoading?: boolean
}

function ArtistsFilters({
  opened,
  onClose,
  filters,
  setFilters,
  isArtistsLoading
}: ArtistsFiltersProps) {
  const [getFiltersMetadata, { data: initialFiltersMetadata, isLoading }] =
    useLazyGetArtistFiltersMetadataQuery()
  useEffect(() => {
    getFiltersMetadata({}, true)
  }, [])

  const searchBy = useSearchBy(filters)
  const { data: filtersMetadata } = useGetArtistFiltersMetadataQuery({ searchBy: searchBy })

  const [internalFilters, setInternalFilters] = useState(filters)
  const initialFilters = useFiltersMetadata(
    initialFiltersMetadata,
    filtersMetadata,
    filters,
    setFilters,
    setInternalFilters,
    artistsFiltersMetadataMap
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
      isLoading={isArtistsLoading}
    >
      <Stack px={'lg'} pb={'lg'}>
        <NumberInputRange
          label={'Band Members'}
          isLoading={isLoading}
          value1={
            internalFilters.get(ArtistProperty.BandMembers + FilterOperator.GreaterThanOrEqual)
              .value as string | number
          }
          onChange1={(value) =>
            handleValueChange(ArtistProperty.BandMembers + FilterOperator.GreaterThanOrEqual, value)
          }
          value2={
            internalFilters.get(ArtistProperty.BandMembers + FilterOperator.LessThanOrEqual)
              .value as string | number
          }
          onChange2={(value) =>
            handleValueChange(ArtistProperty.BandMembers + FilterOperator.LessThanOrEqual, value)
          }
          max={
            initialFilters.get(ArtistProperty.BandMembers + FilterOperator.LessThanOrEqual)
              .value as number
          }
        />

        <DoubleCheckbox
          title={'Is a Band?'}
          label1={'Yes'}
          checked1={internalFilters.get(ArtistProperty.Band + FilterOperator.Equal).isSet}
          onChange1={(value) =>
            handleIsSetChange(ArtistProperty.Band + FilterOperator.Equal, value)
          }
          label2={'No'}
          checked2={internalFilters.get(ArtistProperty.Band + FilterOperator.NotEqual).isSet}
          onChange2={(value) =>
            handleIsSetChange(ArtistProperty.Band + FilterOperator.NotEqual, value)
          }
          disabled={isLoading}
        />

        <NumberInputRange
          label={'Albums'}
          isLoading={isLoading}
          value1={
            internalFilters.get(ArtistProperty.Albums + FilterOperator.GreaterThanOrEqual).value as
              | string
              | number
          }
          onChange1={(value) =>
            handleValueChange(ArtistProperty.Albums + FilterOperator.GreaterThanOrEqual, value)
          }
          value2={
            internalFilters.get(ArtistProperty.Albums + FilterOperator.LessThanOrEqual).value as
              | string
              | number
          }
          onChange2={(value) =>
            handleValueChange(ArtistProperty.Albums + FilterOperator.LessThanOrEqual, value)
          }
          max={
            initialFilters.get(ArtistProperty.Albums + FilterOperator.LessThanOrEqual)
              .value as number
          }
        />

        <NumberInputRange
          label={'Songs'}
          isLoading={isLoading}
          value1={
            internalFilters.get(ArtistProperty.Songs + FilterOperator.GreaterThanOrEqual).value as
              | string
              | number
          }
          onChange1={(value) =>
            handleValueChange(ArtistProperty.Songs + FilterOperator.GreaterThanOrEqual, value)
          }
          value2={
            internalFilters.get(ArtistProperty.Songs + FilterOperator.LessThanOrEqual).value as
              | string
              | number
          }
          onChange2={(value) =>
            handleValueChange(ArtistProperty.Songs + FilterOperator.LessThanOrEqual, value)
          }
          max={
            initialFilters.get(ArtistProperty.Songs + FilterOperator.LessThanOrEqual)
              .value as number
          }
        />

        <NumberInputRange
          label={'Rehearsals'}
          isLoading={isLoading}
          value1={
            internalFilters.get(ArtistProperty.Rehearsals + FilterOperator.GreaterThanOrEqual)
              .value as string | number
          }
          onChange1={(value) =>
            handleValueChange(ArtistProperty.Rehearsals + FilterOperator.GreaterThanOrEqual, value)
          }
          value2={
            internalFilters.get(ArtistProperty.Rehearsals + FilterOperator.LessThanOrEqual)
              .value as string | number
          }
          onChange2={(value) =>
            handleValueChange(ArtistProperty.Rehearsals + FilterOperator.LessThanOrEqual, value)
          }
          max={
            initialFilters.get(ArtistProperty.Rehearsals + FilterOperator.LessThanOrEqual)
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
              ArtistProperty.Confidence + FilterOperator.GreaterThanOrEqual,
              ArtistProperty.Confidence + FilterOperator.LessThanOrEqual
            )}
            onChange={(values) =>
              handleDoubleValueChange(
                ArtistProperty.Confidence + FilterOperator.GreaterThanOrEqual,
                values[0],
                ArtistProperty.Confidence + FilterOperator.LessThanOrEqual,
                values[1]
              )
            }
          />
        </Stack>

        <NumberInputRange
          label={'Progress'}
          isLoading={isLoading}
          value1={
            internalFilters.get(ArtistProperty.Progress + FilterOperator.GreaterThanOrEqual)
              .value as string | number
          }
          onChange1={(value) =>
            handleValueChange(ArtistProperty.Progress + FilterOperator.GreaterThanOrEqual, value)
          }
          value2={
            internalFilters.get(ArtistProperty.Progress + FilterOperator.LessThanOrEqual).value as
              | string
              | number
          }
          onChange2={(value) =>
            handleValueChange(ArtistProperty.Progress + FilterOperator.LessThanOrEqual, value)
          }
          max={
            initialFilters.get(ArtistProperty.Progress + FilterOperator.LessThanOrEqual)
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
            ArtistProperty.LastPlayed + FilterOperator.GreaterThanOrEqual,
            ArtistProperty.LastPlayed + FilterOperator.LessThanOrEqual
          )}
          onChange={(values) =>
            handleDoubleValueChange(
              ArtistProperty.LastPlayed + FilterOperator.GreaterThanOrEqual,
              values[0] ?? null,
              ArtistProperty.LastPlayed + FilterOperator.LessThanOrEqual,
              values[1] ?? null
            )
          }
        />
        <DoubleCheckbox
          title={'Has Been Played Before?'}
          label1={'Yes'}
          checked1={internalFilters.get(ArtistProperty.LastPlayed + FilterOperator.IsNotNull).isSet}
          onChange1={(value) =>
            handleIsSetChange(ArtistProperty.LastPlayed + FilterOperator.IsNotNull, value)
          }
          label2={'Never'}
          checked2={internalFilters.get(ArtistProperty.LastPlayed + FilterOperator.IsNull).isSet}
          onChange2={(value) =>
            handleIsSetChange(ArtistProperty.LastPlayed + FilterOperator.IsNull, value)
          }
          disabled={isLoading}
        />
      </Stack>
    </FiltersDrawer>
  )
}

export default ArtistsFilters
