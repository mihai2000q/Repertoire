import Filter from '../../types/Filter.ts'
import { Checkbox, Group, NumberInput, RangeSlider, Stack, Text } from '@mantine/core'
import { useEffect, useState } from 'react'
import { ArtistSearch } from '../../types/models/Search.ts'
import ArtistSelect from '../@ui/form/select/ArtistSelect.tsx'
import { DatePickerInput } from '@mantine/dates'
import { IconCalendarCheck, IconCalendarRepeat } from '@tabler/icons-react'
import FiltersDrawer from '../@ui/drawer/FiltersDrawer.tsx'
import { useLazyGetAlbumFiltersMetadataQuery } from '../../state/api/albumsApi.ts'
import AlbumProperty from '../../types/enums/AlbumProperty.ts'
import { albumsFiltersMetadataMap } from '../../data/albums/albumsFilters.ts'
import FilterOperator from '../../types/enums/FilterOperator.ts'
import { useDidUpdate } from '@mantine/hooks'
import useFiltersMetadata from '../../hooks/filter/useFiltersMetadata.ts'
import useFiltersHandlers from '../../hooks/filter/useFiltersHandlers.ts'

interface AlbumFiltersProps {
  opened: boolean
  onClose: () => void
  filters: Map<string, Filter>
  setFilters: (filters: Map<string, Filter>) => void
  isAlbumsLoading?: boolean
}

function AlbumFilters({
  opened,
  onClose,
  filters,
  setFilters,
  isAlbumsLoading
}: AlbumFiltersProps) {
  const [getFiltersMetadata, { data: filtersMetadata, isLoading }] =
    useLazyGetAlbumFiltersMetadataQuery()
  useEffect(() => {
    getFiltersMetadata(undefined, true)
  }, [])

  const [internalFilters, setInternalFilters] = useState(filters)
  const initialFilters = useFiltersMetadata(
    filtersMetadata,
    filters,
    setFilters,
    albumsFiltersMetadataMap
  )

  const { handleIsSetChange, handleValueChange, handleDoubleValueChange, getDateRangeValues } =
    useFiltersHandlers(internalFilters, setInternalFilters, initialFilters)

  const [artist, setArtist] = useState<ArtistSearch>(null)
  useDidUpdate(() => {
    handleValueChange(AlbumProperty.ArtistId + FilterOperator.Equal, artist?.id)
  }, [artist])

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
      additionalReset={() => setArtist(null)}
    >
      <Stack px={'lg'} pb={'lg'}>
        <ArtistSelect
          artist={artist}
          setArtist={setArtist}
          ids={filtersMetadata?.artistIds}
          disabled={isLoading}
        />

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
              values[0]?.toISOString() ?? null,
              AlbumProperty.ReleaseDate + FilterOperator.LessThanOrEqual,
              values[1]?.toISOString() ?? null
            )
          }
        />

        <Stack gap={'xxs'}>
          <Text fw={500} fz={'sm'}>
            Has Release Date?
          </Text>
          <Group>
            <Checkbox
              label={'Yes'}
              styles={{ label: { paddingLeft: 8 } }}
              disabled={isLoading}
              checked={
                internalFilters.get(AlbumProperty.ReleaseDate + FilterOperator.IsNotNull).isSet
              }
              onChange={(value) =>
                handleIsSetChange(
                  AlbumProperty.ReleaseDate + FilterOperator.IsNotNull,
                  value.currentTarget.checked
                )
              }
            />
            <Checkbox
              label={'No'}
              styles={{ label: { paddingLeft: 8 } }}
              disabled={isLoading}
              checked={internalFilters.get(AlbumProperty.ReleaseDate + FilterOperator.IsNull).isSet}
              onChange={(value) =>
                handleIsSetChange(
                  AlbumProperty.ReleaseDate + FilterOperator.IsNull,
                  value.currentTarget.checked
                )
              }
            />
          </Group>
        </Stack>

        <Stack gap={2}>
          <Text fw={500} fz={'sm'}>
            Songs
          </Text>
          <Group gap={'xs'}>
            <NumberInput
              flex={1}
              allowNegative={false}
              allowDecimal={false}
              disabled={isLoading}
              value={
                internalFilters.get(AlbumProperty.Songs + FilterOperator.GreaterThanOrEqual)
                  .value as string | number | null | undefined
              }
              onChange={(value) =>
                handleValueChange(AlbumProperty.Songs + FilterOperator.GreaterThanOrEqual, value)
              }
            />
            <Text>-</Text>
            <NumberInput
              flex={1}
              allowNegative={false}
              allowDecimal={false}
              disabled={isLoading}
              value={
                internalFilters.get(AlbumProperty.Songs + FilterOperator.LessThanOrEqual).value as
                  | string
                  | number
                  | null
                  | undefined
              }
              onChange={(value) =>
                handleValueChange(AlbumProperty.Songs + FilterOperator.LessThanOrEqual, value)
              }
            />
          </Group>
        </Stack>

        <Stack gap={2}>
          <Text fw={500} fz={'sm'}>
            Rehearsals
          </Text>
          <Group gap={'xs'}>
            <NumberInput
              flex={1}
              allowNegative={false}
              allowDecimal={false}
              disabled={isLoading}
              value={
                internalFilters.get(AlbumProperty.Rehearsals + FilterOperator.GreaterThanOrEqual)
                  .value as string | number | null | undefined
              }
              onChange={(value) =>
                handleValueChange(
                  AlbumProperty.Rehearsals + FilterOperator.GreaterThanOrEqual,
                  value
                )
              }
            />
            <Text>-</Text>
            <NumberInput
              flex={1}
              allowNegative={false}
              allowDecimal={false}
              disabled={isLoading}
              value={
                internalFilters.get(AlbumProperty.Rehearsals + FilterOperator.LessThanOrEqual)
                  .value as string | number | null | undefined
              }
              onChange={(value) =>
                handleValueChange(AlbumProperty.Rehearsals + FilterOperator.LessThanOrEqual, value)
              }
            />
          </Group>
        </Stack>

        <Stack gap={'xxs'}>
          <Text fw={500} fz={'sm'}>
            Confidence
          </Text>

          <RangeSlider
            disabled={isLoading}
            thumbFromLabel={'confidence'}
            label={(value) => `${value}%`}
            value={[
              internalFilters.get(AlbumProperty.Confidence + FilterOperator.GreaterThanOrEqual)
                .value as number,
              internalFilters.get(AlbumProperty.Confidence + FilterOperator.LessThanOrEqual)
                .value as number
            ]}
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

        <Stack gap={2}>
          <Text fw={500} fz={'sm'}>
            Progress
          </Text>
          <Group gap={'xs'}>
            <NumberInput
              flex={1}
              allowNegative={false}
              allowDecimal={false}
              disabled={isLoading}
              value={
                internalFilters.get(AlbumProperty.Progress + FilterOperator.GreaterThanOrEqual)
                  .value as string | number | null | undefined
              }
              onChange={(value) =>
                handleValueChange(AlbumProperty.Progress + FilterOperator.GreaterThanOrEqual, value)
              }
            />
            <Text>-</Text>
            <NumberInput
              flex={1}
              allowNegative={false}
              allowDecimal={false}
              disabled={isLoading}
              value={
                internalFilters.get(AlbumProperty.Progress + FilterOperator.LessThanOrEqual)
                  .value as string | number | null | undefined
              }
              onChange={(value) =>
                handleValueChange(AlbumProperty.Progress + FilterOperator.LessThanOrEqual, value)
              }
            />
          </Group>
        </Stack>

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
              values[0]?.toISOString() ?? null,
              AlbumProperty.LastPlayed + FilterOperator.LessThanOrEqual,
              values[1]?.toISOString() ?? null
            )
          }
        />
        <Stack gap={'xxs'}>
          <Text fw={500} fz={'sm'}>
            Has Been Played Before?
          </Text>
          <Group>
            <Checkbox
              label={'Yes'}
              styles={{ label: { paddingLeft: 8 } }}
              disabled={isLoading}
              checked={
                internalFilters.get(AlbumProperty.LastPlayed + FilterOperator.IsNotNull).isSet
              }
              onChange={(value) =>
                handleIsSetChange(
                  AlbumProperty.LastPlayed + FilterOperator.IsNotNull,
                  value.currentTarget.checked
                )
              }
            />
            <Checkbox
              label={'Never'}
              styles={{ label: { paddingLeft: 8 } }}
              disabled={isLoading}
              checked={internalFilters.get(AlbumProperty.LastPlayed + FilterOperator.IsNull).isSet}
              onChange={(value) =>
                handleIsSetChange(
                  AlbumProperty.LastPlayed + FilterOperator.IsNull,
                  value.currentTarget.checked
                )
              }
            />
          </Group>
        </Stack>
      </Stack>
    </FiltersDrawer>
  )
}

export default AlbumFilters
