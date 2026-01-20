import WarningModal from '../../@ui/modal/WarningModal.tsx'
import { useBulkDeleteSongSectionsMutation } from '../../../state/api/songsApi.ts'
import { toast } from 'react-toastify'
import plural from '../../../utils/plural.ts'

interface DeleteSongSectionsModalProps {
  ids: string[]
  songId: string,
  opened: boolean
  onClose: () => void
  onDelete?: () => void
}

function DeleteSongSectionsModal({ ids, songId, opened, onClose, onDelete }: DeleteSongSectionsModalProps) {
  const [bulkDeleteSongSectionsMutation, { isLoading }] = useBulkDeleteSongSectionsMutation()

  async function handleDelete() {
    await bulkDeleteSongSectionsMutation({ ids: ids, songId: songId }).unwrap()
    toast.success(`${ids.length} section${plural(ids.length)} deleted!`)
    onDelete?.()
  }

  return (
    <WarningModal
      opened={opened}
      onClose={onClose}
      title={`Delete Section${plural(ids)}`}
      description={`Are you sure you want to delete ${ids.length} section${plural(ids)}?`}
      onYes={handleDelete}
      isLoading={isLoading}
    />
  )
}

export default DeleteSongSectionsModal
