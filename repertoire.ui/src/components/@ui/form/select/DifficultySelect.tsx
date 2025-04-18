import { ComboboxItem, Select } from '@mantine/core'
import Difficulty from '../../../../types/enums/Difficulty.ts'
import { IconStarFilled } from '@tabler/icons-react'

const difficulties = Object.entries(Difficulty).map(([key, value]) => ({
  value: value,
  label: key
}))

interface DifficultySelectProps {
  option: ComboboxItem | null
  onChange: (option: ComboboxItem | null) => void
}

function DifficultySelect({ option, onChange }: DifficultySelectProps) {
  return (
    <Select
      flex={1}
      leftSection={<IconStarFilled size={20} />}
      label={'Difficulty'}
      placeholder={'Select Difficulty'}
      data={difficulties}
      value={option?.value ?? null}
      onChange={(_, option) => onChange(option)}
      clearable
    />
  )
}

export default DifficultySelect
