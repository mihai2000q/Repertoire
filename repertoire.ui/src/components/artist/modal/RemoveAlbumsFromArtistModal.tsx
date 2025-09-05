import { toast } from 'react-toastify'
import { useRemoveAlbumsFromArtistMutation } from '../../../state/api/artistsApi.ts'
import WarningModal from '../../@ui/modal/WarningModal.tsx'
import plural from '../../../utils/plural.ts'

interface RemoveAlbumsFromArtistProps {
  artistId: string
  ids: string[]
  opened: boolean
  onClose: () => void
  onRemove?: () => void
}

function RemoveAlbumsFromArtist({
  artistId,
  ids,
  opened,
  onClose,
  onRemove
}: RemoveAlbumsFromArtistProps) {
  const [removeAlbumsFromArtist, { isLoading }] = useRemoveAlbumsFromArtistMutation()

  async function handleRemove() {
    await removeAlbumsFromArtist({ id: artistId, albumIds: ids }).unwrap()
    toast.success(`${ids.length} album${plural(ids.length)} removed from artist!`)
    onRemove?.()
  }

  return (
    <WarningModal
      opened={opened}
      onClose={onClose}
      title={`Remove album${plural(ids)} from artist`}
      description={`Are you sure you want to remove ${ids.length} album${plural(ids)} from this artist?`}
      onYes={handleRemove}
      isLoading={isLoading}
    />
  )
}

export default RemoveAlbumsFromArtist
