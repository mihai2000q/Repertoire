import { ComboboxItem, Select, SelectProps } from '@mantine/core'
import { useGetSongSectionsQuery } from '../../../sections/state/api/songSectionsApi.ts'

interface SongSectionSelectProps extends SelectProps {
  option: ComboboxItem | null
  onOptionChange: (comboboxItem: ComboboxItem | null) => void
  songId: string
}

function SongSectionSelect({
  option,
  onOptionChange,
  label,
  songId,
  ...props
}: SongSectionSelectProps) {
  const { data: songSectionsData, isLoading } = useGetSongSectionsQuery({ songId: songId })
  const songSections = songSectionsData?.map((songSection) => ({
    value: songSection.id,
    label: songSection.name
  }))

  return (
    <Select
      disabled={isLoading}
      placeholder={'Section'}
      data={songSections}
      value={option?.value ?? null}
      onChange={(_, option) => onOptionChange(option)}
      maxDropdownHeight={150}
      comboboxProps={{
        width: 'max-content',
        position: 'bottom-start'
      }}
      searchable
      aria-label={typeof label === 'string' ? label : 'song-section'}
      {...props}
    />
  )
}

export default SongSectionSelect
