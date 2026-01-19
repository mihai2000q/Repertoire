import { ActionIcon, Tooltip } from '@mantine/core'
import SelectionDrawer from '../@ui/drawer/SelectionDrawer.tsx'
import { useDisclosure } from '@mantine/hooks'
import { IconLocationPlus, IconTrash } from '@tabler/icons-react'
import plural from '../../utils/plural.ts'
import DeleteSongSectionsModal from './modal/DeleteSongSectionsModal.tsx'
import { useBulkRehearsalsSongSectionsMutation } from '../../state/api/songsApi.ts'
import { toast } from 'react-toastify'
import { useClickSelect } from '../../context/ClickSelectContext.tsx'

function SongSectionsSelectionDrawer({ songId }: { songId: string }) {
  const { selectedIds, clearSelection, isSelectionActive } = useClickSelect()

  const [openedDeleteWarning, { open: openDeleteWarning, close: closeDeleteWarning }] =
    useDisclosure(false)

  const [bulkRehearsals, { isLoading: bulkRehearsalsIsLoading }] =
    useBulkRehearsalsSongSectionsMutation()

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
      <SelectionDrawer
        aria-label={'song-sections-selection-drawer'}
        opened={isSelectionActive}
        onClose={clearSelection}
        text={`${selectedIds.length} section${plural(selectedIds)} selected`}
        actionIcons={
          <Tooltip.Group openDelay={200}>
            <Tooltip label={'Add Rehearsals'} openDelay={200}>
              <ActionIcon
                aria-label={'add-rehearsals'}
                variant={'grey-primary'}
                loading={bulkRehearsalsIsLoading}
                onClick={handleAddRehearsals}
              >
                <IconLocationPlus size={15} />
              </ActionIcon>
            </Tooltip>
            <Tooltip label={'Delete sections'}>
              <ActionIcon
                aria-label={'delete'}
                variant={'grey-primary'}
                onClick={openDeleteWarning}
              >
                <IconTrash size={18} />
              </ActionIcon>
            </Tooltip>
          </Tooltip.Group>
        }
      />

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

export default SongSectionsSelectionDrawer
