import { ActionIcon, Group, Tooltip } from '@mantine/core'
import { IconChecks, IconEye, IconEyeOff, IconListNumbers, IconPlus } from '@tabler/icons-react'
import CustomRehearsalButton from './toolbar/CustomRehearsalButton.tsx'
import PopoverConfirmation from '../../../../../components/popover/PopoverConfirmation.tsx'
import SongOutlineSettingsButton from './toolbar/SongOutlineSettingsButton.tsx'
import SongOutlineViewControl from './toolbar/SongOutlineViewControl.tsx'
import { useAddPerfectSongRehearsalMutation } from '../../../../../state/api/songsApi.ts'
import { useMemo, useState } from 'react'
import { useDisclosure } from '@mantine/hooks'
import { toast } from 'react-toastify'
import OutlineView from '../types/enums/OutlineView.ts'
import SongArrangementsModal from '../../arrangements/SongArrangementsModal.tsx'
import { SongPart, SongSection } from '../../../../../types/models/Song.ts'
import { useSongContext } from '../../../context/SongContext.tsx'
import { useSongOutlineContext } from '../context/SongOutlineContext.tsx'

interface SongOutlineToolbarProps {
  toggleAdd: () => void
  parts?: SongPart[]
  sections?: SongSection[]
  scrollIntoView?: () => void
}

function SongOutlineToolbar({
  toggleAdd,
  parts,
  sections,
  scrollIntoView
}: SongOutlineToolbarProps) {
  const { songId, defaultArrangementId: defaultSongArrangementId } = useSongContext()
  const { view: outlineView, areAllDetailsShowing, toggleAllDetails } = useSongOutlineContext()

  const [addPerfectRehearsal, { isLoading: isPerfectRehearsalLoading }] =
    useAddPerfectSongRehearsalMutation()

  const [openedPerfectRehearsalPopover, setOpenedPerfectRehearsalPopover] = useState(false)
  const [openedArrangements, { open: openArrangements, close: closeArrangements }] =
    useDisclosure(false)

  const currentIds = useMemo(
    () =>
      (outlineView === OutlineView.Sections
        ? sections?.map((s) => s.id)
        : parts?.map((p) => p.id)) ?? [],
    [outlineView, sections, parts]
  )

  const sectionParts = sections?.flatMap((s) => s.parts)
  const currentParts = (outlineView === OutlineView.Sections ? sectionParts : parts) ?? []
  function handleShowDetails() {
    toggleAllDetails(currentIds)
    if (!areAllDetailsShowing && scrollIntoView) setTimeout(scrollIntoView, 250)
  }

  async function handleAddPerfectRehearsal() {
    await addPerfectRehearsal({ id: songId }).unwrap()
    toast.info('Perfect rehearsal added!')
    setOpenedPerfectRehearsalPopover(false)
  }

  return (
    <Group aria-label={'outline-toolbar'} gap={'xxs'}>
      <Tooltip.Group openDelay={500} closeDelay={100}>
        <Tooltip label={`Add New ${outlineView === OutlineView.Sections ? 'Section' : 'Part'}`}>
          <ActionIcon
            aria-label={`add-new-${outlineView === OutlineView.Sections ? 'section' : 'part'}`}
            variant={'grey'}
            size={'sm'}
            onClick={toggleAdd}
          >
            <IconPlus size={16} />
          </ActionIcon>
        </Tooltip>

        <Tooltip
          label={
            currentIds.length > 0
              ? `To show details you need ${outlineView === OutlineView.Sections ? 'sections' : 'parts'}`
              : areAllDetailsShowing(currentIds)
                ? 'Hide details'
                : 'Show Details'
          }
        >
          <ActionIcon
            aria-label={areAllDetailsShowing(currentIds) ? 'hide-details' : 'show-details'}
            variant={'grey'}
            size={'sm'}
            disabled={currentIds.length === 0}
            onClick={handleShowDetails}
          >
            {areAllDetailsShowing(currentIds) ? <IconEyeOff size={16} /> : <IconEye size={16} />}
          </ActionIcon>
        </Tooltip>

        <Tooltip label={'Manage Song Arrangements'}>
          <ActionIcon
            aria-label={'manage-song-arrangements'}
            variant={'grey'}
            size={'sm'}
            onClick={openArrangements}
          >
            <IconListNumbers size={16} />
          </ActionIcon>
        </Tooltip>

        <CustomRehearsalButton
          songId={songId}
          defaultSongArrangementId={defaultSongArrangementId}
          partsCount={currentParts.length}
        />

        <PopoverConfirmation
          label={"Increase parts' rehearsals based on occurrences from default arrangement"}
          popoverProps={{
            opened: openedPerfectRehearsalPopover,
            onChange: setOpenedPerfectRehearsalPopover,
            closeOnClickOutside: !isPerfectRehearsalLoading
          }}
          isLoading={isPerfectRehearsalLoading}
          onCancel={() => setOpenedPerfectRehearsalPopover(false)}
          onConfirm={handleAddPerfectRehearsal}
        >
          <Tooltip
            label={
              currentParts.length === 0
                ? 'To add a perfect rehearsal, you need parts'
                : !defaultSongArrangementId
                  ? 'To add a perfect rehearsal, you need a default arrangement'
                  : 'Add Perfect Rehearsal'
            }
            disabled={openedPerfectRehearsalPopover}
          >
            <ActionIcon
              aria-label={'add-perfect-rehearsal'}
              variant={'grey'}
              size={'sm'}
              disabled={currentParts.length === 0 || !defaultSongArrangementId}
              onClick={() =>
                setOpenedPerfectRehearsalPopover(
                  isPerfectRehearsalLoading || !openedPerfectRehearsalPopover
                )
              }
            >
              <IconChecks size={16} />
            </ActionIcon>
          </Tooltip>
        </PopoverConfirmation>

        <SongOutlineSettingsButton parts={currentParts} />

        <SongOutlineViewControl />
      </Tooltip.Group>

      <SongArrangementsModal
        opened={openedArrangements}
        onClose={closeArrangements}
        songId={songId}
        defaultId={defaultSongArrangementId}
      />
    </Group>
  )
}

export default SongOutlineToolbar
