import { Stack, Text } from '@mantine/core'
import Song from '../../../types/models/Song.ts'
import EditSongDescriptionModal from '../modal/EditSongDescriptionModal.tsx'
import { useDisclosure } from '@mantine/hooks'
import EditWidget from '../../@ui/widget/EditWidget.tsx'

interface SongDescriptionWidgetProps {
  song: Song
}

function SongDescriptionWidget({ song }: SongDescriptionWidgetProps) {
  const [openedEdit, { open: openEdit, close: closeEdit }] = useDisclosure(false)

  return (
    <EditWidget p={'md'} ariaLabel={'description-widget'} onEditClick={openEdit}>
      <Stack gap={'xs'}>
        <Text fw={600}>Description</Text>
        {song.description ? (
          <Text>{song.description}</Text>
        ) : (
          <Text fw={500} c={'dimmed'} fs={'italic'}>
            No Description
          </Text>
        )}
      </Stack>

      <EditSongDescriptionModal song={song} opened={openedEdit} onClose={closeEdit} />
    </EditWidget>
  )
}

export default SongDescriptionWidget
