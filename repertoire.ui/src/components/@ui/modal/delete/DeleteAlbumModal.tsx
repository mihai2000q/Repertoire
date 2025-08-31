import { useState } from 'react'
import { Checkbox, Group, Stack, Text } from '@mantine/core'
import WarningModal from '../WarningModal.tsx'
import Album from '../../../../types/models/Album.ts'
import { useDeleteAlbumMutation } from '../../../../state/api/albumsApi.ts'
import { toast } from 'react-toastify'

interface DeleteArtistModalProps {
  opened: boolean
  onClose: () => void
  album: Album
  onDelete?: () => void
  withName?: boolean
}

function DeleteAlbumModal({ opened, onClose, album, onDelete, withName }: DeleteArtistModalProps) {
  const [deleteAlbumMutation, { isLoading: isDeleteLoading }] = useDeleteAlbumMutation()
  const [deleteWithSongs, setDeleteWithSongs] = useState(false)

  async function handleDelete() {
    await deleteAlbumMutation({ id: album.id, withSongs: deleteWithSongs }).unwrap()
    toast.success(`${album.title} deleted!`)
    onDelete?.()
  }

  return (
    <WarningModal
      opened={opened}
      onClose={onClose}
      title={'Delete Album'}
      description={
        <Stack gap={5}>
          {withName === true ? (
            <Group gap={'xxs'}>
              <Text>Are you sure you want to delete</Text>
              <Text fw={600}>{album.title}</Text>
              <Text>?</Text>
            </Group>
          ) : (
            <Text>Are you sure you want to delete this album?</Text>
          )}
          <Checkbox
            checked={deleteWithSongs}
            onChange={(event) => setDeleteWithSongs(event.currentTarget.checked)}
            label={'Delete all associated songs'}
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

export default DeleteAlbumModal
