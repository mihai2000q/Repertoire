import { SongSectionType } from '../../../../../../../../types/models/Song.ts'
import { Text } from '@mantine/core'

interface SongSectionTypeBadgeProps {
  songSectionType: SongSectionType
}

function SongSectionTypeBadge({ songSectionType }: SongSectionTypeBadgeProps) {
  return (
    <Text
      fw={600}
      fz={'xs'}
      c={'gray.5'}
      bg={'gray.1'}
      px={'6px'}
      py={'3px'}
      inline
      style={(theme) => ({
        borderRadius: theme.radius.md
      })}
    >
      {songSectionType.name}
    </Text>
  )
}

export default SongSectionTypeBadge
