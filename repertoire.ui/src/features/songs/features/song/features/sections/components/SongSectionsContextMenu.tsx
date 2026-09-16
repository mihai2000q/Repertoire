import { IconLocationPlus, IconTrash } from '@tabler/icons-react'
import { ContextMenu } from '../../../../../../../components/menu/ContextMenu.tsx'
import { useDisclosure } from '@mantine/hooks'
import { ReactNode, useEffect, useRef, useState } from 'react'
import DeleteSongSectionsModal from './modal/DeleteSongSectionsModal.tsx'
import { toast } from 'react-toastify'
import plural from '../../../../../../../utils/plural.ts'
import MenuItemConfirmation from '../../../../../../../components/menu/item/MenuItemConfirmation.tsx'
import { useClickSelect } from '../../../../../../../context/ClickSelectContext.tsx'
import { SongPart, SongSection } from '../../../../../../../types/models/Song.ts'
import { useBulkUpdateSongPartsMutation } from '../../parts/state/api/songPartsApi.ts'

interface SongSectionsContextMenuProps {
  children: ReactNode
  sections: SongSection[]
  songId: string
}

function SongSectionsContextMenu({ children, sections, songId }: SongSectionsContextMenuProps) {
  const { selectedIds, clearSelection } = useClickSelect()
  const selectedSections = useRef<SongSection[]>([])
  const [selectedSectionParts, setSelectedSectionParts] = useState<SongPart[]>([])
  useEffect(() => {
    selectedSections.current = sections.filter((s) => selectedIds.some((sId) => sId === s.id))
    setSelectedSectionParts(selectedSections.current.flatMap((s) => s.parts))
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
      requests: selectedSectionParts.map((p) => ({
        id: p.id,
        rehearsals: p.rehearsals + 1,
        confidence: p.confidence
      })),
      songId: songId
    }).unwrap()
    toast.success(
      `Rehearsals added to
      ${selectedSectionParts.length} part${plural(selectedSectionParts)}
      of the selected section${plural(selectedIds)}!`
    )
    clearSelection()
  }

  return (
    <>
      <ContextMenu
        aria-label={'song-sections-context-menu'}
        opened={openedMenu}
        onClose={closeMenu}
        onOpen={openMenu}
        disabled={selectedIds.length === 0}
      >
        <ContextMenu.Target>{children}</ContextMenu.Target>

        <ContextMenu.Dropdown>
          <MenuItemConfirmation
            isLoading={bulkUpdateIsLoading}
            onConfirm={handleAddRehearsals}
            leftSection={<IconLocationPlus size={14} />}
            disabled={selectedSectionParts.length === 0}
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

      <DeleteSongSectionsModal
        ids={selectedIds}
        songId={songId}
        opened={openedDeleteWarning}
        onClose={closeDeleteWarning}
        onDelete={clearSelection}
      />
    </>
  )
}

export default SongSectionsContextMenu
