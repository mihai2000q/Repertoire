import { ActionIcon, Tooltip } from '@mantine/core'
import SelectionDrawer from '../../../../../../../components/drawer/SelectionDrawer.tsx'
import { useDisclosure } from '@mantine/hooks'
import { IconRefresh, IconTrash } from '@tabler/icons-react'
import plural from '../../../../../../../utils/plural.ts'
import { toast } from 'react-toastify'
import { useClickSelect } from '../../../../../../../context/ClickSelectContext.tsx'
import { SongSection } from '../../../../../../../types/models/Song.ts'
import { useMemo } from 'react'
import { useBulkUpdateSongPartsMutation } from '../../parts/state/api/songPartsApi.ts'
import DeleteSongPartsModal from '../../parts/components/modal/DeleteSongPartsModal.tsx'
import useSelectedSectionParts from '../hooks/useSelectedSectionParts.ts'
import DeleteSongSectionsAndPartsModal from './modal/DeleteSongSectionsAndPartsModal.tsx'
import DeleteSongSectionsModal from './modal/DeleteSongSectionsModal.tsx'

interface SongSectionsSelectionDrawerProps {
  sections: SongSection[]
  songId: string
}

function SongSectionsSelectionDrawer({ sections, songId }: SongSectionsSelectionDrawerProps) {
  const { selectedIds, clearSelection, isClickSelectionActive } = useClickSelect()
  const [selectedSections, selectedSectionParts] = useSelectedSectionParts(selectedIds, sections)

  const [openedDeleteWarning, { open: openDeleteWarning, close: closeDeleteWarning }] =
    useDisclosure(false)

  const [bulkUpdate, { isLoading: bulkUpdateIsLoading }] = useBulkUpdateSongPartsMutation()

  const selectionText = useMemo(() => {
    const sectionsLabel = `${selectedSections.length} section${plural(selectedSections)}`
    const partsLabel = `${selectedSectionParts.length} part${plural(selectedSectionParts)}`

    if (selectedSections.length > 0 && selectedSectionParts.length > 0)
      return `${sectionsLabel} and ${partsLabel} selected`
    if (selectedSections.length > 0) return `${sectionsLabel} selected`
    if (selectedSectionParts.length > 0) return `${partsLabel} selected`
    return '0 selected'
  }, [selectedSections, selectedSectionParts])

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
      <SelectionDrawer
        aria-label={'song-sections-selection-drawer'}
        opened={isClickSelectionActive}
        onClose={clearSelection}
        text={selectionText}
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
                <IconRefresh size={15} />
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

      {selectedSections.length > 0 && selectedSectionParts.length > 0 ? (
        <DeleteSongSectionsAndPartsModal
          sections={selectedSections}
          sectionParts={selectedSectionParts}
          songId={songId}
          opened={openedDeleteWarning}
          onClose={closeDeleteWarning}
          onDelete={clearSelection}
        />
      ) : selectedSections.length > 0 ? (
        <DeleteSongSectionsModal
          sections={selectedSections}
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

export default SongSectionsSelectionDrawer
