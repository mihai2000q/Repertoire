import WarningModal from '../WarningModal.tsx'
import { useBulkDeleteSongsMutation } from '../../../../state/api/songsApi.ts'
import { toast } from 'react-toastify'
import plural from '../../../../utils/plural.ts'

interface DeleteSongsModalProps {
  ids: string[]
  opened: boolean
  onClose: () => void
  onDelete?: () => void
}

function DeleteSongsModal({ ids, opened, onClose, onDelete }: DeleteSongsModalProps) {
  const [bulkDeleteSongsMutation, { isLoading }] = useBulkDeleteSongsMutation()

  async function handleDelete() {
    await bulkDeleteSongsMutation({ ids: ids }).unwrap()
    toast.success(`${ids.length} song${plural(ids.length)} deleted!`)
    onDelete?.()
  }

  return (
    <WarningModal
      opened={opened}
      onClose={onClose}
      title={`Delete Song${plural(ids)}`}
      description={`Are you sure you want to delete ${ids.length} song${plural(ids)}?`}
      onYes={handleDelete}
      isLoading={isLoading}
    />
  )
}

export default DeleteSongsModal
