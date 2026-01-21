import WarningModal from '../WarningModal.tsx'
import { useBulkDeletePlaylistsMutation } from '../../../../state/api/playlistsApi.ts'
import { toast } from 'react-toastify'
import plural from '../../../../utils/plural.ts'

interface DeletePlaylistsModalProps {
  ids: string[]
  opened: boolean
  onClose: () => void
  onDelete?: () => void
}

function DeletePlaylistsModal({ ids, opened, onClose, onDelete }: DeletePlaylistsModalProps) {
  const [bulkDeletePlaylistsMutation, { isLoading }] = useBulkDeletePlaylistsMutation()

  async function handleDelete() {
    await bulkDeletePlaylistsMutation({ ids: ids }).unwrap()
    toast.success(`${ids.length} playlist${plural(ids.length)} deleted!`)
    onDelete?.()
  }

  return (
    <WarningModal
      opened={opened}
      onClose={onClose}
      title={`Delete Playlist${plural(ids)}`}
      description={`Are you sure you want to delete ${ids.length} playlist${plural(ids)}?`}
      onYes={handleDelete}
      isLoading={isLoading}
    />
  )
}

export default DeletePlaylistsModal
