import { IconChecks } from '@tabler/icons-react'
import { toast } from 'react-toastify'
import { useAddPerfectSongRehearsalMutation } from '../../../../../state/api/songsApi.ts'
import MenuItemConfirmation from '../MenuItemConfirmation.tsx'

interface PerfectRehearsalMenuItemProps {
  songId: string
  closeMenu: () => void
}

function PerfectRehearsalMenuItem({ songId, closeMenu }: PerfectRehearsalMenuItemProps) {
  const [addPerfectRehearsal, { isLoading }] = useAddPerfectSongRehearsalMutation()

  async function handleAddPerfectRehearsal() {
    await addPerfectRehearsal({ id: songId }).unwrap()
    toast.success(`Perfect rehearsal added!`)
    closeMenu()
  }

  return (
    <MenuItemConfirmation
      isLoading={isLoading}
      onConfirm={handleAddPerfectRehearsal}
      leftSection={<IconChecks size={14} />}
    >
      Perfect Rehearsal
    </MenuItemConfirmation>
  )
}

export default PerfectRehearsalMenuItem
