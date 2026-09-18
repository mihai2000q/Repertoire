import { IconRefresh, IconTrash } from '@tabler/icons-react'
import { ContextMenu } from '../../../../../../../components/menu/ContextMenu.tsx'
import { useDisclosure } from '@mantine/hooks'
import { ReactNode, useEffect, useMemo } from 'react'
import DeleteSongSectionsModal from './modal/DeleteSongSectionsModal.tsx'
import { toast } from 'react-toastify'
import plural from '../../../../../../../utils/plural.ts'
import MenuItemConfirmation from '../../../../../../../components/menu/item/MenuItemConfirmation.tsx'
import { useClickSelect } from '../../../../../../../context/ClickSelectContext.tsx'
import { SongSection } from '../../../../../../../types/models/Song.ts'
import { useBulkUpdateSongPartsMutation } from '../../parts/state/api/songPartsApi.ts'
import DeleteSongPartsModal from '../../parts/components/modal/DeleteSongPartsModal.tsx'

interface SongSectionsContextMenuProps {
  children: ReactNode
  sections: SongSection[]
  songId: string
}

function SongSectionsContextMenu({ children, sections, songId }: SongSectionsContextMenuProps) {
  const { selectedIds, clearSelection } = useClickSelect()
  const selectedSections = useMemo(
    () => sections.filter((s) => selectedIds.includes(`section-${s.id}`)),
    [sections, selectedIds]
  )
  const selectedSectionParts = useMemo(
    () =>
      sections.flatMap((s) => s.parts.filter((p) => selectedIds.includes(`part-${p.id}:${s.id}`))),
    [sections, selectedIds]
  )

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
      `Rehearsals added to ${selectedSectionParts.length} ` + `part${plural(selectedSectionParts)}!`
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
            leftSection={<IconRefresh size={14} />}
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

      {selectedSections.length > 0 ? (
        <DeleteSongSectionsModal
          sections={selectedSections}
          sectionParts={selectedSectionParts}
          songId={songId}
          opened={openedDeleteWarning}
          onClose={closeDeleteWarning}
          onDelete={clearSelection}
        />
      ) : (
        <DeleteSongPartsModal
          ids={selectedSectionParts.map((sp) => sp.id)}
          songId={songId}
          opened={openedDeleteWarning}
          onClose={closeDeleteWarning}
          onDelete={clearSelection}
        />
      )}
    </>
  )
}

export default SongSectionsContextMenu
