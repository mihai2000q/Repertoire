import Filter from '../../types/Filter.ts'
import { RangeSlider, Stack, Text } from '@mantine/core'
import { useEffect, useState } from 'react'
import { DatePickerInput } from '@mantine/dates'
import { IconCalendarCheck, IconCalendarRepeat } from '@tabler/icons-react'
import FiltersDrawer from '../@ui/drawer/FiltersDrawer.tsx'
import {
  useGetSongFiltersMetadataQuery,
  useLazyGetSongFiltersMetadataQuery
} from '../../state/api/songsApi.ts'
import SongProperty from '../../types/enums/SongProperty.ts'
import { songsFiltersMetadataMap } from '../../data/songs/songsFilters.ts'
import FilterOperator from '../../types/enums/FilterOperator.ts'
import useFiltersMetadata from '../../hooks/filter/useFiltersMetadata.ts'
import useFiltersHandlers from '../../hooks/filter/useFiltersHandlers.ts'
import NumberInputRange from '../@ui/form/input/NumberInputRange.tsx'
import DoubleCheckbox from '../@ui/filter/DoubleCheckbox.tsx'
import DifficultyMultiSelect from '../@ui/form/select/multi/DifficultyMultiSelect.tsx'
import Difficulty from '../../types/enums/Difficulty.ts'
import GuitarTuningMultiSelect from '../@ui/form/select/multi/GuitarTuningMultiSelect.tsx'
import InstrumentMultiSelect from '../@ui/form/select/multi/InstrumentMultiSelect.tsx'
import useSearchBy from '../../hooks/api/useSearchBy.ts'

interface SongsFiltersProps {
  opened: boolean
  onClose: () => void
  filters: Map<string, Filter>
  setFilters: (filters: Map<string, Filter>, withSearchParams?: boolean) => void
  isSongsLoading?: boolean
}

function SongsFilters({ opened, onClose, filters, setFilters, isSongsLoading }: SongsFiltersProps) {
  const [getFiltersMetadata, { data: initialFiltersMetadata, isLoading }] =
    useLazyGetSongFiltersMetadataQuery()
  useEffect(() => {
    getFiltersMetadata({}, true)
  }, [])

  const searchBy = useSearchBy(filters)
  const { data: filtersMetadata } = useGetSongFiltersMetadataQuery({ searchBy: searchBy })

  const [internalFilters, setInternalFilters] = useState(filters)
  const initialFilters = useFiltersMetadata(
    initialFiltersMetadata,
    filtersMetadata,
    filters,
    setFilters,
    setInternalFilters,
    songsFiltersMetadataMap
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
      isLoading={isSongsLoading}
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
            SongProperty.ReleaseDate + FilterOperator.GreaterThanOrEqual,
            SongProperty.ReleaseDate + FilterOperator.LessThanOrEqual
          )}
          onChange={(values) =>
            handleDoubleValueChange(
              SongProperty.ReleaseDate + FilterOperator.GreaterThanOrEqual,
              values[0]?.toISOString() ?? null,
              SongProperty.ReleaseDate + FilterOperator.LessThanOrEqual,
              values[1]?.toISOString() ?? null
            )
          }
        />
        <DoubleCheckbox
          title={'Has Release Date?'}
          label1={'Yes'}
          checked1={internalFilters.get(SongProperty.ReleaseDate + FilterOperator.IsNotNull).isSet}
          onChange1={(value) =>
            handleIsSetChange(SongProperty.ReleaseDate + FilterOperator.IsNotNull, value)
          }
          label2={'No'}
          checked2={internalFilters.get(SongProperty.ReleaseDate + FilterOperator.IsNull).isSet}
          onChange2={(value) =>
            handleIsSetChange(SongProperty.ReleaseDate + FilterOperator.IsNull, value)
          }
          disabled={isLoading}
        />

        <NumberInputRange
          label={'Bpm'}
          isLoading={isLoading}
          value1={
            internalFilters.get(SongProperty.BPM + FilterOperator.GreaterThanOrEqual).value as
              | string
              | number
          }
          onChange1={(value) =>
            handleValueChange(SongProperty.BPM + FilterOperator.GreaterThanOrEqual, value)
          }
          value2={
            internalFilters.get(SongProperty.BPM + FilterOperator.LessThanOrEqual).value as
              | string
              | number
          }
          onChange2={(value) =>
            handleValueChange(SongProperty.BPM + FilterOperator.LessThanOrEqual, value)
          }
          max={
            initialFilters.get(SongProperty.BPM + FilterOperator.LessThanOrEqual).value as number
          }
        />
        <DoubleCheckbox
          title={'Has BPM?'}
          label1={'Yes'}
          checked1={internalFilters.get(SongProperty.BPM + FilterOperator.IsNotNull).isSet}
          onChange1={(value) =>
            handleIsSetChange(SongProperty.BPM + FilterOperator.IsNotNull, value)
          }
          label2={'No'}
          checked2={internalFilters.get(SongProperty.BPM + FilterOperator.IsNull).isSet}
          onChange2={(value) => handleIsSetChange(SongProperty.BPM + FilterOperator.IsNull, value)}
          disabled={isLoading}
        />

        <DifficultyMultiSelect
          label={'Difficulties'}
          placeholder={'Select difficulties'}
          difficulties={
            (internalFilters.get(SongProperty.Difficulty + FilterOperator.In)
              .value as Difficulty[]) ?? []
          }
          setDifficulties={(ids) => {
            handleValueChange(SongProperty.Difficulty + FilterOperator.In, ids as Difficulty[])
          }}
          availableDifficulties={filtersMetadata?.difficulties}
          disabled={isLoading}
        />
        <DoubleCheckbox
          title={'Has Difficulty?'}
          label1={'Yes'}
          checked1={internalFilters.get(SongProperty.Difficulty + FilterOperator.IsNotNull).isSet}
          onChange1={(value) =>
            handleIsSetChange(SongProperty.Difficulty + FilterOperator.IsNotNull, value)
          }
          label2={'No'}
          checked2={internalFilters.get(SongProperty.Difficulty + FilterOperator.IsNull).isSet}
          onChange2={(value) =>
            handleIsSetChange(SongProperty.Difficulty + FilterOperator.IsNull, value)
          }
          disabled={isLoading}
        />

        <GuitarTuningMultiSelect
          label={'Guitar Tunings'}
          placeholder={'Select tunings'}
          ids={
            (internalFilters.get(SongProperty.GuitarTuningId + FilterOperator.In)
              .value as string[]) ?? []
          }
          setIds={(ids) => {
            handleValueChange(SongProperty.GuitarTuningId + FilterOperator.In, ids as string[])
          }}
          availableIds={filtersMetadata?.guitarTuningIds}
          disabled={isLoading}
        />
        <DoubleCheckbox
          title={'Has Guitar Tuning?'}
          label1={'Yes'}
          checked1={
            internalFilters.get(SongProperty.GuitarTuningId + FilterOperator.IsNotNull).isSet
          }
          onChange1={(value) =>
            handleIsSetChange(SongProperty.GuitarTuningId + FilterOperator.IsNotNull, value)
          }
          label2={'No'}
          checked2={internalFilters.get(SongProperty.GuitarTuningId + FilterOperator.IsNull).isSet}
          onChange2={(value) =>
            handleIsSetChange(SongProperty.GuitarTuningId + FilterOperator.IsNull, value)
          }
          disabled={isLoading}
        />

        <InstrumentMultiSelect
          label={'Instruments'}
          placeholder={'Select instruments'}
          ids={
            (internalFilters.get(SongProperty.InstrumentId + FilterOperator.In)
              .value as string[]) ?? []
          }
          setIds={(ids) => {
            handleValueChange(SongProperty.InstrumentId + FilterOperator.In, ids as string[])
          }}
          availableIds={filtersMetadata?.instrumentIds}
          disabled={isLoading}
        />

        <DoubleCheckbox
          title={'Is Recorded?'}
          label1={'Yes'}
          checked1={internalFilters.get(SongProperty.IsRecorded + FilterOperator.Equal).isSet}
          onChange1={(value) =>
            handleIsSetChange(SongProperty.IsRecorded + FilterOperator.Equal, value)
          }
          label2={'No'}
          checked2={internalFilters.get(SongProperty.IsRecorded + FilterOperator.NotEqual).isSet}
          onChange2={(value) =>
            handleIsSetChange(SongProperty.IsRecorded + FilterOperator.NotEqual, value)
          }
          disabled={isLoading}
        />

        <NumberInputRange
          label={'Sections'}
          isLoading={isLoading}
          value1={
            internalFilters.get(SongProperty.Sections + FilterOperator.GreaterThanOrEqual).value as
              | string
              | number
          }
          onChange1={(value) =>
            handleValueChange(SongProperty.Sections + FilterOperator.GreaterThanOrEqual, value)
          }
          value2={
            internalFilters.get(SongProperty.Sections + FilterOperator.LessThanOrEqual).value as
              | string
              | number
          }
          onChange2={(value) =>
            handleValueChange(SongProperty.Sections + FilterOperator.LessThanOrEqual, value)
          }
          max={
            initialFilters.get(SongProperty.Sections + FilterOperator.LessThanOrEqual)
              .value as number
          }
        />

        <NumberInputRange
          label={'Solos'}
          isLoading={isLoading}
          value1={
            internalFilters.get(SongProperty.Solos + FilterOperator.GreaterThanOrEqual).value as
              | string
              | number
          }
          onChange1={(value) =>
            handleValueChange(SongProperty.Solos + FilterOperator.GreaterThanOrEqual, value)
          }
          value2={
            internalFilters.get(SongProperty.Solos + FilterOperator.LessThanOrEqual).value as
              | string
              | number
          }
          onChange2={(value) =>
            handleValueChange(SongProperty.Solos + FilterOperator.LessThanOrEqual, value)
          }
          max={
            initialFilters.get(SongProperty.Solos + FilterOperator.LessThanOrEqual).value as number
          }
        />

        <NumberInputRange
          label={'Riffs'}
          isLoading={isLoading}
          value1={
            internalFilters.get(SongProperty.Riffs + FilterOperator.GreaterThanOrEqual).value as
              | string
              | number
          }
          onChange1={(value) =>
            handleValueChange(SongProperty.Riffs + FilterOperator.GreaterThanOrEqual, value)
          }
          value2={
            internalFilters.get(SongProperty.Riffs + FilterOperator.LessThanOrEqual).value as
              | string
              | number
          }
          onChange2={(value) =>
            handleValueChange(SongProperty.Riffs + FilterOperator.LessThanOrEqual, value)
          }
          max={
            initialFilters.get(SongProperty.Riffs + FilterOperator.LessThanOrEqual).value as number
          }
        />

        <NumberInputRange
          label={'Rehearsals'}
          isLoading={isLoading}
          value1={
            internalFilters.get(SongProperty.Rehearsals + FilterOperator.GreaterThanOrEqual)
              .value as string | number
          }
          onChange1={(value) =>
            handleValueChange(SongProperty.Rehearsals + FilterOperator.GreaterThanOrEqual, value)
          }
          value2={
            internalFilters.get(SongProperty.Rehearsals + FilterOperator.LessThanOrEqual).value as
              | string
              | number
          }
          onChange2={(value) =>
            handleValueChange(SongProperty.Rehearsals + FilterOperator.LessThanOrEqual, value)
          }
          max={
            initialFilters.get(SongProperty.Rehearsals + FilterOperator.LessThanOrEqual)
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
              SongProperty.Confidence + FilterOperator.GreaterThanOrEqual,
              SongProperty.Confidence + FilterOperator.LessThanOrEqual
            )}
            onChange={(values) =>
              handleDoubleValueChange(
                SongProperty.Confidence + FilterOperator.GreaterThanOrEqual,
                values[0],
                SongProperty.Confidence + FilterOperator.LessThanOrEqual,
                values[1]
              )
            }
          />
        </Stack>

        <NumberInputRange
          label={'Progress'}
          isLoading={isLoading}
          value1={
            internalFilters.get(SongProperty.Progress + FilterOperator.GreaterThanOrEqual).value as
              | string
              | number
          }
          onChange1={(value) =>
            handleValueChange(SongProperty.Progress + FilterOperator.GreaterThanOrEqual, value)
          }
          value2={
            internalFilters.get(SongProperty.Progress + FilterOperator.LessThanOrEqual).value as
              | string
              | number
          }
          onChange2={(value) =>
            handleValueChange(SongProperty.Progress + FilterOperator.LessThanOrEqual, value)
          }
          max={
            initialFilters.get(SongProperty.Progress + FilterOperator.LessThanOrEqual)
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
            SongProperty.LastPlayed + FilterOperator.GreaterThanOrEqual,
            SongProperty.LastPlayed + FilterOperator.LessThanOrEqual
          )}
          onChange={(values) =>
            handleDoubleValueChange(
              SongProperty.LastPlayed + FilterOperator.GreaterThanOrEqual,
              values[0]?.toISOString() ?? null,
              SongProperty.LastPlayed + FilterOperator.LessThanOrEqual,
              values[1]?.toISOString() ?? null
            )
          }
        />
        <DoubleCheckbox
          title={'Has Been Played Before?'}
          label1={'Yes'}
          checked1={internalFilters.get(SongProperty.LastPlayed + FilterOperator.IsNotNull).isSet}
          onChange1={(value) =>
            handleIsSetChange(SongProperty.LastPlayed + FilterOperator.IsNotNull, value)
          }
          label2={'Never'}
          checked2={internalFilters.get(SongProperty.LastPlayed + FilterOperator.IsNull).isSet}
          onChange2={(value) =>
            handleIsSetChange(SongProperty.LastPlayed + FilterOperator.IsNull, value)
          }
          disabled={isLoading}
        />

        <DoubleCheckbox
          title={'Has Songsterr Link?'}
          label1={'Yes'}
          checked1={internalFilters.get(SongProperty.SongsterrLink + FilterOperator.IsNotNull).isSet}
          onChange1={(value) =>
            handleIsSetChange(SongProperty.SongsterrLink + FilterOperator.IsNotNull, value)
          }
          label2={'No'}
          checked2={internalFilters.get(SongProperty.SongsterrLink + FilterOperator.IsNull).isSet}
          onChange2={(value) =>
            handleIsSetChange(SongProperty.SongsterrLink + FilterOperator.IsNull, value)
          }
          disabled={isLoading}
        />

        <DoubleCheckbox
          title={'Has Youtube Link?'}
          label1={'Yes'}
          checked1={internalFilters.get(SongProperty.YoutubeLink + FilterOperator.IsNotNull).isSet}
          onChange1={(value) =>
            handleIsSetChange(SongProperty.YoutubeLink + FilterOperator.IsNotNull, value)
          }
          label2={'No'}
          checked2={internalFilters.get(SongProperty.YoutubeLink + FilterOperator.IsNull).isSet}
          onChange2={(value) =>
            handleIsSetChange(SongProperty.YoutubeLink + FilterOperator.IsNull, value)
          }
          disabled={isLoading}
        />
      </Stack>
    </FiltersDrawer>
  )
}

export default SongsFilters
