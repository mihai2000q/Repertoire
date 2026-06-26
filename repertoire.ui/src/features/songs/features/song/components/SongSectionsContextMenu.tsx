import { Menu } from '@mantine/core'
import { IconLocationPlus, IconTrash } from '@tabler/icons-react'
import { ContextMenu } from '../@ui/menu/ContextMenu.tsx'
import { useDisclosure } from '@mantine/hooks'
import { ReactNode, useEffect } from 'react'
import DeleteSongSectionsModal from './modal/DeleteSongSectionsModal.tsx'
import { useBulkRehearsalsSongSectionsMutation } from '../../state/api/songsApi.ts'
import { toast } from 'react-toastify'
import plural from '../../utils/plural.ts'
import MenuItemConfirmation from '../@ui/menu/item/MenuItemConfirmation.tsx'
import { useClickSelect } from '../../context/ClickSelectContext.tsx'

function SongSectionsContextMenu({ children, songId }: { children: ReactNode; songId: string }) {
  const { selectedIds, clearSelection } = useClickSelect()

  const [openedMenu, { open: openMenu, close: closeMenu }] = useDisclosure(false)

  const [openedDeleteWarning, { open: openDeleteWarning, close: closeDeleteWarning }] =
    useDisclosure(false)

  const [bulkRehearsals, { isLoading: bulkRehearsalsIsLoading }] =
    useBulkRehearsalsSongSectionsMutation()

  useEffect(() => {
    if (selectedIds.length === 0) closeMenu()
  }, [selectedIds])

  async function handleAddRehearsals() {
    await bulkRehearsals({
      sections: selectedIds.map((id) => ({ id: id, rehearsals: 1 })),
      songId: songId
    }).unwrap()
    toast.success(`Rehearsals added to ${selectedIds.length} section${plural(selectedIds)}!`)
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
            isLoading={bulkRehearsalsIsLoading}
            onConfirm={handleAddRehearsals}
            leftSection={<IconLocationPlus size={14} />}
          >
            Add Rehearsals
          </MenuItemConfirmation>

          <Menu.Item c={'red'} leftSection={<IconTrash size={14} />} onClick={openDeleteWarning}>
            Delete
          </Menu.Item>
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
