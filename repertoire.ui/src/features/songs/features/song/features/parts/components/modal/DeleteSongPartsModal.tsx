import WarningModal from '../../../../../../../../components/modal/WarningModal.tsx'
import { useBulkDeleteSongPartsMutation } from '../../state/api/songPartsApi.ts'
import { toast } from 'react-toastify'
import plural from '../../../../../../../../utils/plural.ts'

interface DeleteSongPartsModalProps {
  ids: string[]
  songId: string
  opened: boolean
  onClose: () => void
  onDelete?: () => void
}

function DeleteSongPartsModal({
  ids,
  songId,
  opened,
  onClose,
  onDelete
}: DeleteSongPartsModalProps) {
  const [bulkDeleteSongPartsMutation, { isLoading }] = useBulkDeleteSongPartsMutation()

  async function handleDelete() {
    await bulkDeleteSongPartsMutation({ ids: ids, songId: songId }).unwrap()
    toast.success(`${ids.length} part${plural(ids.length)} deleted!`)
    onDelete?.()
  }

  return (
    <WarningModal
      opened={opened}
      onClose={onClose}
      title={`Delete Part${plural(ids)}`}
      description={`Are you sure you want to delete ${ids.length} part${plural(ids)}?`}
      onYes={handleDelete}
      isLoading={isLoading}
    />
  )
}

export default DeleteSongPartsModal
