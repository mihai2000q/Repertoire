import { toast } from 'react-toastify'
import { useAddPartialSongRehearsalMutation } from '../../../../../state/api/songsApi.ts'
import MenuItemConfirmation from '../MenuItemConfirmation.tsx'
import { IconCheck } from '@tabler/icons-react'

interface PartialRehearsalMenuItemProps {
  songId: string
}

function PartialRehearsalMenuItem({ songId }: PartialRehearsalMenuItemProps) {
  const [addPartialRehearsal, { isLoading }] = useAddPartialSongRehearsalMutation()

  async function handleAddPartialRehearsal() {
    await addPartialRehearsal({ id: songId }).unwrap()
    toast.success(`Partial rehearsal added!`)
  }

  return (
    <MenuItemConfirmation
      isLoading={isLoading}
      onConfirm={handleAddPartialRehearsal}
      leftSection={<IconCheck size={14} />}
    >
      Partial Rehearsal
    </MenuItemConfirmation>
  )
}

export default PartialRehearsalMenuItem
