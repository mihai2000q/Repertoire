import { toast } from 'react-toastify'
import { useRemoveSongsFromArtistMutation } from '../../../state/api/artistsApi.ts'
import WarningModal from '../../@ui/modal/WarningModal.tsx'
import plural from '../../../utils/plural.ts'

interface RemoveSongsFromArtistProps {
  artistId: string
  ids: string[]
  opened: boolean
  onClose: () => void
  onRemove?: () => void
}

function RemoveSongsFromArtist({
  artistId,
  ids,
  opened,
  onClose,
  onRemove
}: RemoveSongsFromArtistProps) {
  const [removeSongsFromArtist, { isLoading }] = useRemoveSongsFromArtistMutation()

  async function handleRemove() {
    await removeSongsFromArtist({ id: artistId, songIds: ids }).unwrap()
    toast.success(`${ids.length} song${plural(ids.length)} removed from artist!`)
    onRemove?.()
  }

  return (
    <WarningModal
      opened={opened}
      onClose={onClose}
      title={`Remove song${plural(ids)} from artist`}
      description={`Are you sure you want to remove ${ids.length} song${plural(ids)} from this artist?`}
      onYes={handleRemove}
      isLoading={isLoading}
    />
  )
}

export default RemoveSongsFromArtist
