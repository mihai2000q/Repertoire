import { ComboboxItem, Select, SelectProps } from '@mantine/core'
import { useGetSongSectionTypesQuery } from '../../../../../state/api/songsApi.ts'

interface SongSectionTypeSelectProps extends SelectProps {
  option: ComboboxItem | null
  onOptionChange: (comboboxItem: ComboboxItem | null) => void
}

function SongSectionTypeSelect({
  option,
  onOptionChange,
  label,
  placeholder = 'Type',
  ...others
}: SongSectionTypeSelectProps) {
  const { data: songSectionTypesData, isLoading } = useGetSongSectionTypesQuery()
  const songSectionTypes = songSectionTypesData?.map((type) => ({
    value: type.id,
    label: type.name
  }))

  return (
    <Select
      label={label}
      disabled={isLoading}
      placeholder={placeholder}
      data={songSectionTypes}
      value={option?.value ?? null}
      onChange={(_, option) => onOptionChange(option)}
      maxDropdownHeight={150}
      clearable={false}
      allowDeselect={false}
      searchable
      aria-label={typeof label === 'string' ? label : 'song-section-type'}
      {...others}
    />
  )
}

export default SongSectionTypeSelect
