import { ComboboxItem, Group, Loader, Select, Text } from '@mantine/core'
import { useGetGuitarTuningsQuery } from '../../../../state/songsApi.ts'
import CustomIconGuitarHead from '../../icons/CustomIconGuitarHead.tsx'

interface GuitarTuningSelectProps {
  option: ComboboxItem | null
  onChange: (comboboxItem: ComboboxItem | null) => void
}

function GuitarTuningSelect({ option, onChange }: GuitarTuningSelectProps) {
  const { data: guitarTuningsData, isLoading } = useGetGuitarTuningsQuery()
  const guitarTunings = guitarTuningsData?.map((guitarTuning) => ({
    value: guitarTuning.id,
    label: guitarTuning.name
  }))

  return isLoading ? (
    <Group gap={'xs'} flex={1.25}>
      <Loader size={25} />
      <Text fz={'sm'} c={'dimmed'}>
        Loading Tunings...
      </Text>
    </Group>
  ) : (
    <Select
      flex={1.25}
      leftSection={<CustomIconGuitarHead size={20} />}
      label={'Guitar Tuning'}
      placeholder={'Select Guitar Tuning'}
      data={guitarTunings}
      value={option?.value ?? null}
      onChange={(_, option) => onChange(option)}
      maxDropdownHeight={150}
      clearable
      searchable
    />
  )
}

export default GuitarTuningSelect
