import { IconLocationPlus, IconTrash } from '@tabler/icons-react'
import { ContextMenu } from '../../../../../../../components/menu/ContextMenu.tsx'
import { useDisclosure } from '@mantine/hooks'
import { ReactNode, useEffect, useRef } from 'react'
import DeleteSongPartsModal from './modal/DeleteSongPartsModal.tsx'
import { useBulkUpdateSongPartsMutation } from '../state/api/songPartsApi.ts'
import { toast } from 'react-toastify'
import plural from '../../../../../../../utils/plural.ts'
import MenuItemConfirmation from '../../../../../../../components/menu/item/MenuItemConfirmation.tsx'
import { useClickSelect } from '../../../../../../../context/ClickSelectContext.tsx'
import { SongPart } from '../../../../../../../types/models/Song.ts'

interface SongPartsContextMenuProps {
  children: ReactNode
  songId: string
  parts: SongPart[]
}

function SongPartsContextMenu({ children, songId, parts }: SongPartsContextMenuProps) {
  const { selectedIds, clearSelection } = useClickSelect()
  const selectedParts = useRef<SongPart[]>([])
  useEffect(() => {
    selectedParts.current = parts.filter((p) => selectedIds.some((pId) => pId === p.id))
  }, [selectedIds])

  const [openedMenu, { open: openMenu, close: closeMenu }] = useDisclosure(false)

  const [openedDeleteWarning, { open: openDeleteWarning, close: closeDeleteWarning }] =
    useDisclosure(false)

  const [bulkUpdate, { isLoading: bulkUpdateIsLoading }] = useBulkUpdateSongPartsMutation()

  useEffect(() => {
    if (selectedIds.length === 0) closeMenu()
  }, [selectedIds])

  async function handleAddRehearsals() {
    await bulkUpdate({
      requests: selectedParts.current.map((p) => ({
        id: p.id,
        confidence: p.confidence,
        rehearsals: p.rehearsals + 1
      })),
      songId: songId
    }).unwrap()
    toast.success(`Rehearsals added to ${selectedIds.length} part${plural(selectedIds)}!`)
    clearSelection()
  }

  return (
    <>
      <ContextMenu
        aria-label={'song-parts-context-menu'}
        opened={openedMenu}
        onClose={closeMenu}
        onOpen={openMenu}
        disabled={selectedIds.length === 0}
      >
        <ContextMenu.Target>{children}</ContextMenu.Target>

        <ContextMenu.Dropdown>
          <ContextMenu.Label></ContextMenu.Label>

          <MenuItemConfirmation
            isLoading={bulkUpdateIsLoading}
            onConfirm={handleAddRehearsals}
            leftSection={<IconLocationPlus size={14} />}
          >
            Add Rehearsals
          </MenuItemConfirmation>

          <ContextMenu.Item
            c={'red'}
            leftSection={<IconTrash size={14} />}
            onClick={openDeleteWarning}
          >
            Delete
          </ContextMenu.Item>
        </ContextMenu.Dropdown>
      </ContextMenu>

      <DeleteSongPartsModal
        ids={selectedIds}
        songId={songId}
        opened={openedDeleteWarning}
        onClose={closeDeleteWarning}
        onDelete={clearSelection}
      />
    </>
  )
}

export default SongPartsContextMenu
