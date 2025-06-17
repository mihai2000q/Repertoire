import { Checkbox, Group, Stack, Text } from '@mantine/core'
import WarningModal from './WarningModal.tsx'
import { useDeleteArtistMutation } from '../../../state/api/artistsApi.ts'
import { useState } from 'react'
import { toast } from 'react-toastify'
import Artist from '../../../types/models/Artist.ts'

interface DeleteArtistModalProps {
  opened: boolean
  onClose: () => void
  artist: Artist
  onDelete?: () => void
  withName?: boolean
}

function DeleteArtistModal({
  opened,
  onClose,
  artist,
  onDelete,
  withName
}: DeleteArtistModalProps) {
  const [deleteArtistMutation, { isLoading: isDeleteLoading }] = useDeleteArtistMutation()
  const [deleteWithAssociations, setDeleteWithAssociations] = useState(false)

  async function handleDelete() {
    await deleteArtistMutation({
      id: artist.id,
      withAlbums: deleteWithAssociations,
      withSongs: deleteWithAssociations
    }).unwrap()
    toast.success(`${artist.name} deleted!`)
    onDelete?.()
  }

  return (
    <WarningModal
      opened={opened}
      onClose={onClose}
      title={'Delete Artist'}
      description={
        <Stack gap={'xs'}>
          {withName === true ? (
            <Group gap={'xxs'}>
              <Text>Are you sure you want to delete</Text>
              <Text fw={600}>{artist.name}</Text>
              <Text>?</Text>
            </Group>
          ) : (
            <Text fw={500}>Are you sure you want to delete this artist?</Text>
          )}
          <Checkbox
            checked={deleteWithAssociations}
            onChange={(event) => setDeleteWithAssociations(event.currentTarget.checked)}
            label={
              <Text c={'dimmed'}>
                Delete all associated <b>albums</b> and <b>songs</b>
              </Text>
            }
            styles={{ label: { paddingLeft: 8 } }}
          />
        </Stack>
      }
      onYes={handleDelete}
      isLoading={isDeleteLoading}
    />
  )
}

export default DeleteArtistModal
