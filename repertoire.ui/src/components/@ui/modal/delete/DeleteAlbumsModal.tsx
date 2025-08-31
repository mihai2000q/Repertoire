import { Checkbox, Stack, Text } from '@mantine/core'
import WarningModal from '../WarningModal.tsx'
import { useState } from 'react'
import { toast } from 'react-toastify'
import plural from '../../../../utils/plural.ts'
import { useBulkDeleteAlbumsMutation } from '../../../../state/api/albumsApi.ts'

interface DeleteAlbumsModalProps {
  ids: string[]
  opened: boolean
  onClose: () => void
  onDelete?: () => void
}

function DeleteAlbumsModal({ ids, opened, onClose, onDelete }: DeleteAlbumsModalProps) {
  const [bulkDeleteAlbumsMutation, { isLoading: isDeleteLoading }] = useBulkDeleteAlbumsMutation()
  const [deleteWithSongs, setDeleteWithSongs] = useState(false)

  async function handleDelete() {
    await bulkDeleteAlbumsMutation({
      ids: ids,
      withSongs: deleteWithSongs
    }).unwrap()
    toast.success(`${ids.length} album${plural(ids.length)} deleted!`)
    onDelete?.()
  }

  return (
    <WarningModal
      opened={opened}
      onClose={onClose}
      title={`Delete Album${plural(ids)}`}
      description={
        <Stack gap={5}>
          <Text fw={500}>Are you sure you want to delete {ids.length} album{plural(ids)}?</Text>
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

export default DeleteAlbumsModal
