import { ComboboxItem, Select } from '@mantine/core'
import { useGetSongSectionTypesQuery } from '../../../../state/songsApi.ts'

interface SongSectionTypeSelectProps {
  option: ComboboxItem | null
  onChange: (comboboxItem: ComboboxItem | null) => void
  error?: boolean
  label?: string
  placeholder?: string
  flex?: number
}

function SongSectionTypeSelect({
  option,
  onChange,
  error,
  label,
  flex,
  placeholder = 'Type'
}: SongSectionTypeSelectProps) {
  const { data: songSectionTypesData, isLoading } = useGetSongSectionTypesQuery()
  const songSectionTypes = songSectionTypesData?.map((type) => ({
    value: type.id,
    label: type.name
  }))

  return (
    <Select
      w={!flex && 95}
      flex={flex}
      label={label}
      disabled={isLoading}
      placeholder={placeholder}
      data={songSectionTypes}
      value={option?.value}
      onChange={(_, option) => onChange(option)}
      error={error}
      maxDropdownHeight={150}
      clearable={false}
      searchable
    />
  )
}

export default SongSectionTypeSelect
