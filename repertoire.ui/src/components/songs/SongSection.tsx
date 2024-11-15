import { SongSection as SongSectionType } from '../../types/models/Song.ts'
import { ActionIcon, alpha, Group, Text } from '@mantine/core'
import { IconDots, IconGripVertical } from '@tabler/icons-react'
import { DraggableProvided } from '@hello-pangea/dnd'

interface SongSectionProps {
  section: SongSectionType
  isDragging: boolean
  draggableProvided: DraggableProvided
}

function SongSection({ section, isDragging, draggableProvided }: SongSectionProps) {
  return (
    <Group
      key={section.id}
      align={'center'}
      gap={'xs'}
      py={'xs'}
      px={'md'}
      sx={(theme) => ({
        transition: '0.25s',
        borderRadius: isDragging ? '16px' : '0px',
        border: isDragging
          ? `1px solid ${alpha(theme.colors.cyan[9], 0.33)}`
          : '1px solid transparent',

        '&:hover': {
          boxShadow: theme.shadows.xl,
          backgroundColor: alpha(theme.colors.cyan[0], 0.15)
        }
      })}
      ref={draggableProvided.innerRef}
      {...draggableProvided.draggableProps}
    >
      <ActionIcon variant={'subtle'} size={'lg'} {...draggableProvided.dragHandleProps}>
        <IconGripVertical size={20} />
      </ActionIcon>

      <Text inline fw={600}>
        {section.songSectionType.name}
      </Text>
      <Text flex={1} inline>
        {section.name}
      </Text>

      <ActionIcon variant={'subtle'} size={'lg'}>
        <IconDots size={20} />
      </ActionIcon>
    </Group>
  )
}

export default SongSection
