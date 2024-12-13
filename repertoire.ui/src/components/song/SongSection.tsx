import { SongSection as SongSectionType } from '../../types/models/Song.ts'
import {
  ActionIcon,
  alpha,
  Collapse,
  Group,
  Menu,
  NumberFormatter,
  Progress,
  Stack,
  Text,
  Tooltip
} from '@mantine/core'
import {
  IconDots,
  IconEdit,
  IconGripVertical,
  IconLocationPlus,
  IconTrash
} from '@tabler/icons-react'
import { DraggableProvided } from '@hello-pangea/dnd'
import { useDisclosure } from '@mantine/hooks'
import { toast } from 'react-toastify'
import { useDeleteSongSectionMutation, useUpdateSongSectionMutation } from '../../state/songsApi.ts'
import EditSongSectionModal from './modal/EditSongSectionModal.tsx'

interface SongSectionProps {
  section: SongSectionType
  songId: string
  isDragging: boolean
  draggableProvided: DraggableProvided
  showDetails: boolean
  maxSectionProgress: number
}

function SongSection({
  section,
  songId,
  isDragging,
  draggableProvided,
  showDetails,
  maxSectionProgress
}: SongSectionProps) {
  const [updateSongSectionMutation] = useUpdateSongSectionMutation()
  const [deleteSongSectionMutation] = useDeleteSongSectionMutation()

  const [openedEditSongSection, { open: openEditSongSection, close: closeEditSongSection }] =
    useDisclosure(false)

  function handleAddRehearsal() {
    updateSongSectionMutation({
      ...section,
      typeId: section.songSectionType.id,
      rehearsals: section.rehearsals + 1
    })
    toast.info(`${section.name} rehearsals' have been increased by 1`)
  }

  function handleDelete() {
    deleteSongSectionMutation({ id: section.id, songId: songId })
    toast.success(`${section.name} deleted!`)
  }

  return (
    <Stack
      py={'xs'}
      px={'md'}
      sx={(theme) => ({
        cursor: 'default',
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
      <Group align={'center'} gap={'xs'}>
        <ActionIcon variant={'subtle'} size={'lg'} {...draggableProvided.dragHandleProps}>
          <IconGripVertical size={20} />
        </ActionIcon>

        <Text inline fw={600}>
          {section.songSectionType.name}
        </Text>
        <Text flex={1} inline truncate={'end'}>
          {section.name}
        </Text>

        <Group gap={2} align={'center'}>
          <Tooltip label={'Add Rehearsal'} openDelay={200}>
            <ActionIcon variant={'subtle'} size={'md'} onClick={handleAddRehearsal}>
              <IconLocationPlus size={15} />
            </ActionIcon>
          </Tooltip>

          <Menu>
            <Menu.Target>
              <ActionIcon variant={'subtle'} size={'lg'}>
                <IconDots size={20} />
              </ActionIcon>
            </Menu.Target>

            <Menu.Dropdown>
              <Menu.Item leftSection={<IconEdit size={14} />} onClick={openEditSongSection}>
                Edit
              </Menu.Item>
              <Menu.Item leftSection={<IconTrash size={14} />} c={'red.5'} onClick={handleDelete}>
                Delete
              </Menu.Item>
            </Menu.Dropdown>
          </Menu>
        </Group>
      </Group>

      <Collapse in={showDetails}>
        <Group align={'center'} gap={'xl'} px={'md'}>
          <Tooltip.Floating label={`Rehearsals: ${section.rehearsals}`}>
            <Text fw={500} c={'dimmed'} fz={'md'} inline>
              <NumberFormatter value={section.rehearsals} />
            </Text>
          </Tooltip.Floating>

          <Tooltip.Floating label={`Confidence: ${section.confidence}%`}>
            <Progress flex={1} size={'sm'} value={section.confidence} />
          </Tooltip.Floating>

          <Tooltip.Floating
            label={
              <>
                Progress: <NumberFormatter value={section.progress} />
              </>
            }
          >
            <Progress
              flex={1}
              size={'sm'}
              value={section.progress === 0 ? 0 : (section.progress / maxSectionProgress) * 100}
              color={'green'}
            />
          </Tooltip.Floating>
        </Group>
      </Collapse>

      <EditSongSectionModal
        opened={openedEditSongSection}
        onClose={closeEditSongSection}
        section={section}
      />
    </Stack>
  )
}

export default SongSection
