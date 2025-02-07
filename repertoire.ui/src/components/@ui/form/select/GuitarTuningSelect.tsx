import { ComboboxItem, Select } from '@mantine/core'
import { useGetGuitarTuningsQuery } from '../../../../state/api/songsApi.ts'
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

  return (
    <Select
      flex={1.25}
      leftSection={<CustomIconGuitarHead size={20} />}
      label={'Guitar Tuning'}
      disabled={isLoading}
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
