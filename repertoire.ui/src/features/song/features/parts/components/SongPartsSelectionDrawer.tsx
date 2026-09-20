import { ActionIcon, Tooltip } from '@mantine/core'
import SelectionDrawer from '../../../../../components/drawer/SelectionDrawer.tsx'
import { useDisclosure } from '@mantine/hooks'
import { IconRefresh, IconTrash } from '@tabler/icons-react'
import plural from '../../../../../utils/plural.ts'
import DeleteSongPartsModal from './modal/DeleteSongPartsModal.tsx'
import { useBulkUpdateSongPartsMutation } from '../state/api/songPartsApi.ts'
import { toast } from 'react-toastify'
import { useClickSelect } from '../../../../../context/ClickSelectContext.tsx'
import { SongPart } from '../../../../../types/models/Song.ts'
import { useEffect, useRef } from 'react'

interface SongPartsSelectionDrawerProps {
  songId: string
  parts: SongPart[]
}

function SongPartsSelectionDrawer({ songId, parts }: SongPartsSelectionDrawerProps) {
  const { selectedIds, clearSelection, isClickSelectionActive } = useClickSelect()
  const selectedParts = useRef<SongPart[]>([])
  useEffect(() => {
    selectedParts.current = parts.filter((p) => selectedIds.some((pId) => pId === p.id))
  }, [selectedIds])

  const [openedDeleteWarning, { open: openDeleteWarning, close: closeDeleteWarning }] =
    useDisclosure(false)

  const [bulkUpdate, { isLoading: bulkUpdateIsLoading }] = useBulkUpdateSongPartsMutation()

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
      <SelectionDrawer
        aria-label={'song-parts-selection-drawer'}
        opened={isClickSelectionActive}
        onClose={clearSelection}
        text={`${selectedIds.length} part${plural(selectedIds)} selected`}
        actionIcons={
          <Tooltip.Group openDelay={200}>
            <Tooltip label={'Add Rehearsals'} openDelay={200}>
              <ActionIcon
                aria-label={'add-rehearsals'}
                variant={'grey-primary'}
                loading={bulkUpdateIsLoading}
                onClick={handleAddRehearsals}
              >
                <IconRefresh size={15} />
              </ActionIcon>
            </Tooltip>
            <Tooltip label={'Delete parts'}>
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

export default SongPartsSelectionDrawer
