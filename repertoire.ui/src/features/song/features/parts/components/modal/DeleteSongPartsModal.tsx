import WarningModal from '../../../../../../components/modal/WarningModal.tsx'
import { useBulkDeleteSongPartsMutation } from '../../state/api/songPartsApi.ts'
import { toast } from 'react-toastify'
import plural from '../../../../../../utils/plural.ts'
import { useMemo } from 'react'

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

  const [uniqueIds, duplicateIds] = useMemo(() => {
    const seen = new Set<string>()
    const unique: string[] = []
    const duplicates: string[] = []

    ids.forEach((id) => {
      if (seen.has(id)) duplicates.push(id)
      else {
        seen.add(id)
        unique.push(id)
      }
    })

    return [unique, duplicates]
  }, [ids])

  async function handleDelete() {
    await bulkDeleteSongPartsMutation({ ids: uniqueIds, songId: songId }).unwrap()
    toast.success(`${uniqueIds.length} part${plural(uniqueIds)} deleted!`)
    onDelete?.()
  }

  return (
    <WarningModal
      opened={opened}
      onClose={onClose}
      title={`Delete Part${plural(uniqueIds)}`}
      description={
        `Are you sure you want to delete ${uniqueIds.length} part${plural(uniqueIds)}? ` +
        `${
          duplicateIds.length > 0
            ? `(${duplicateIds.length} duplicate${plural(duplicateIds.length)} will be ignored)`
            : ''
        }`
      }
      onYes={handleDelete}
      isLoading={isLoading}
    />
  )
}

export default DeleteSongPartsModal
