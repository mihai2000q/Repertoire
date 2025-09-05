import { toast } from 'react-toastify'
import { useRemoveSongsFromPlaylistMutation } from '../../../state/api/playlistsApi.ts'
import WarningModal from '../../@ui/modal/WarningModal.tsx'
import plural from '../../../utils/plural.ts'

interface RemoveSongsFromPlaylistProps {
  playlistId: string
  ids: string[]
  opened: boolean
  onClose: () => void
  onRemove?: () => void
}

function RemoveSongsFromPlaylist({
  playlistId,
  ids,
  opened,
  onClose,
  onRemove
}: RemoveSongsFromPlaylistProps) {
  const [removeSongsFromPlaylist, { isLoading }] = useRemoveSongsFromPlaylistMutation()

  async function handleRemove() {
    await removeSongsFromPlaylist({ id: playlistId, playlistSongIds: ids }).unwrap()
    toast.success(`${ids.length} song${plural(ids.length)} removed from playlist!`)
    onRemove?.()
  }

  return (
    <WarningModal
      opened={opened}
      onClose={onClose}
      title={`Remove song${plural(ids)} from playlist`}
      description={`Are you sure you want to remove ${ids.length} song${plural(ids)} from this playlist?`}
      onYes={handleRemove}
      isLoading={isLoading}
    />
  )
}

export default RemoveSongsFromPlaylist
