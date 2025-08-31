import { Checkbox, Stack, Text } from '@mantine/core'
import WarningModal from '../WarningModal.tsx'
import { useBulkDeleteArtistsMutation } from '../../../../state/api/artistsApi.ts'
import { useState } from 'react'
import { toast } from 'react-toastify'
import plural from '../../../../utils/plural.ts'

interface DeleteArtistsModalProps {
  ids: string[]
  opened: boolean
  onClose: () => void
  onDelete?: () => void
}

function DeleteArtistsModal({ ids, opened, onClose, onDelete }: DeleteArtistsModalProps) {
  const [bulkDeleteArtistsMutation, { isLoading: isDeleteLoading }] = useBulkDeleteArtistsMutation()
  const [deleteWithAssociations, setDeleteWithAssociations] = useState(false)

  async function handleDelete() {
    await bulkDeleteArtistsMutation({
      ids: ids,
      withAlbums: deleteWithAssociations,
      withSongs: deleteWithAssociations
    }).unwrap()
    toast.success(`${ids.length} artist${plural(ids.length)} deleted!`)
    onDelete?.()
  }

  return (
    <WarningModal
      opened={opened}
      onClose={onClose}
      title={`Delete Artist${plural(ids)}`}
      description={
        <Stack gap={'xs'}>
          <Text fw={500}>Are you sure you want to delete {ids.length} artist{plural(ids)}?</Text>
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

export default DeleteArtistsModal
