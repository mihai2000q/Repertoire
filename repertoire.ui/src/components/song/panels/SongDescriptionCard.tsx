import { Stack, Text } from '@mantine/core'
import Song from '../../../types/models/Song.ts'
import EditSongDescriptionModal from '../modal/EditSongDescriptionModal.tsx'
import { useDisclosure } from '@mantine/hooks'
import EditPanelCard from "../../@ui/card/EditPanelCard.tsx";

interface SongDescriptionCardProps {
  song: Song
}

function SongDescriptionCard({ song }: SongDescriptionCardProps) {
  const [openedEdit, { open: openEdit, close: closeEdit }] = useDisclosure(false)

  return (
    <EditPanelCard p={'md'} onEditClick={openEdit}>
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
    </EditPanelCard>
  )
}

export default SongDescriptionCard
