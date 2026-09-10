import { ActionIcon, Tooltip } from '@mantine/core'
import SelectionDrawer from '../../../../../../../components/drawer/SelectionDrawer.tsx'
import { useDisclosure } from '@mantine/hooks'
import { IconLocationPlus, IconTrash } from '@tabler/icons-react'
import plural from '../../../../../../../utils/plural.ts'
import DeleteSongSectionsModal from './modal/DeleteSongSectionsModal.tsx'
import { toast } from 'react-toastify'
import { useClickSelect } from '../../../../../../../context/ClickSelectContext.tsx'
import { SongPart, SongSection } from '../../../../../../../types/models/Song.ts'
import { useEffect, useRef, useState } from 'react'
import { useBulkUpdateSongPartsMutation } from '../../parts/state/api/songPartsApi.ts'

interface SongSectionsSelectionDrawerProps {
  sections: SongSection[]
  songId: string
}

function SongSectionsSelectionDrawer({ sections, songId }: SongSectionsSelectionDrawerProps) {
  const { selectedIds, clearSelection, isClickSelectionActive } = useClickSelect()
  const selectedSections = useRef<SongSection[]>([])
  const [selectedSectionParts, setSelectedSectionParts] = useState<SongPart[]>([])
  useEffect(() => {
    selectedSections.current = sections.filter((s) => selectedIds.some((sId) => sId === s.id))
    setSelectedSectionParts(selectedSections.current.flatMap((s) => s.parts))
  }, [selectedIds])

  const [openedDeleteWarning, { open: openDeleteWarning, close: closeDeleteWarning }] =
    useDisclosure(false)

  const [bulkUpdate, { isLoading: bulkUpdateIsLoading }] = useBulkUpdateSongPartsMutation()

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
      <SelectionDrawer
        aria-label={'song-sections-selection-drawer'}
        opened={isClickSelectionActive}
        onClose={clearSelection}
        text={`${selectedIds.length} section${plural(selectedIds)} selected`}
        actionIcons={
          <Tooltip.Group openDelay={200}>
            <Tooltip label={'Add Rehearsals'} openDelay={200}>
              <ActionIcon
                aria-label={'add-rehearsals'}
                variant={'grey-primary'}
                loading={bulkUpdateIsLoading}
                disabled={selectedSectionParts.length === 0}
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
