import { useState } from 'react'
import { Checkbox, Group, Stack, Text } from '@mantine/core'
import WarningModal from '../../../../../../../../components/modal/WarningModal.tsx'
import { SongSection } from '../../../../../../../../types/models/Song.ts'
import { toast } from 'react-toastify'
import { useDeleteSongSectionMutation } from '../../state/api/songSectionsApi.ts'

interface DeleteSongSectionModalProps {
  opened: boolean
  onClose: () => void
  section: SongSection
  songId: string
}

function DeleteSongSectionModal({ opened, onClose, section, songId }: DeleteSongSectionModalProps) {
  const [deleteSongSectionMutation, { isLoading: isDeleteLoading }] = useDeleteSongSectionMutation()
  const [deleteWithParts, setDeleteWithParts] = useState(false)

  async function handleDelete() {
    await deleteSongSectionMutation({
      id: section.id,
      songId: songId,
      withParts: deleteWithParts
    }).unwrap()
    toast.success(`${section.name} deleted!`)
  }

  return (
    <WarningModal
      opened={opened}
      onClose={onClose}
      title={`Delete Section`}
      description={
        <Stack gap={5}>
          <Group gap={'xxs'}>
            <Text>Are you sure you want to delete</Text>
            <Text fw={600}>{section.name}</Text>
            <Text>?</Text>
          </Group>
          <Checkbox
            checked={deleteWithParts}
            onChange={(event) => setDeleteWithParts(event.currentTarget.checked)}
            label={'Delete all associated parts'}
            c={'dimmed'}
            styles={{ label: { paddingLeft: 8 } }}
          />
        </Stack>
      }
      onYes={handleDelete}
      isLoading={isDeleteLoading}
    />
  )
}

export default DeleteSongSectionModal
