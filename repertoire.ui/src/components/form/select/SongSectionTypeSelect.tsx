import { ComboboxItem, Select } from '@mantine/core'
import { useGetSongSectionTypesQuery } from '../../../state/songsApi.ts'

interface SongSectionTypeSelectProps {
  option: ComboboxItem
  onChange: (comboboxItem: ComboboxItem) => void
  error?: boolean
  label?: string
  placeholder?: string
  flex?: number
}

function SongSectionTypeSelect({ option, onChange, error, label, placeholder, flex }: SongSectionTypeSelectProps) {
  const { data: songSectionTypesData } = useGetSongSectionTypesQuery()
  const songSectionTypes = songSectionTypesData?.map((type) => ({
    value: type.id,
    label: type.name
  }))

  return (
    <Select
      w={!flex && 95}
      flex={flex}
      label={label}
      placeholder={placeholder ?? 'Type'}
      data={songSectionTypes}
      value={option ? option.value : null}
      onChange={(_, option) => onChange(option)}
      error={error}
      maxDropdownHeight={150}
      clearable={false}
      searchable
    />
  )
}

export default SongSectionTypeSelect
