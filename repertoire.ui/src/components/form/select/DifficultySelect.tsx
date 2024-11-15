import { ComboboxItem, Select } from '@mantine/core'
import Difficulty from '../../../utils/enums/Difficulty.ts'

interface DifficultySelectProps {
  option: ComboboxItem
  onChange: (option: ComboboxItem) => void
}

function DifficultySelect({ option, onChange }: DifficultySelectProps) {
  const difficulties = Object.entries(Difficulty).map(([key, value]) => ({
    value: value,
    label: key
  }))

  return (
    <Select
      flex={1}
      label={'Difficulty'}
      placeholder={'Select Difficulty'}
      data={difficulties}
      value={option ? option.value : null}
      onChange={(_, option) => onChange(option)}
      clearable
    />
  )
}

export default DifficultySelect
