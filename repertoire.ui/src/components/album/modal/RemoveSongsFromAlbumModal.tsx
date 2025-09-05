import { toast } from 'react-toastify'
import { useRemoveSongsFromAlbumMutation } from '../../../state/api/albumsApi.ts'
import WarningModal from '../../@ui/modal/WarningModal.tsx'
import plural from '../../../utils/plural.ts'

interface RemoveSongsFromAlbumProps {
  albumId: string
  ids: string[]
  opened: boolean
  onClose: () => void
  onRemove?: () => void
}

function RemoveSongsFromAlbum({
  albumId,
  ids,
  opened,
  onClose,
  onRemove
}: RemoveSongsFromAlbumProps) {
  const [removeSongsFromAlbum, { isLoading }] = useRemoveSongsFromAlbumMutation()

  async function handleRemove() {
    await removeSongsFromAlbum({ id: albumId, songIds: ids }).unwrap()
    toast.success(`${ids.length} song${plural(ids.length)} removed from album!`)
    onRemove?.()
  }

  return (
    <WarningModal
      opened={opened}
      onClose={onClose}
      title={`Remove song${plural(ids)} from album`}
      description={`Are you sure you want to remove ${ids.length} song${plural(ids)} from this album?`}
      onYes={handleRemove}
      isLoading={isLoading}
    />
  )
}

export default RemoveSongsFromAlbum
