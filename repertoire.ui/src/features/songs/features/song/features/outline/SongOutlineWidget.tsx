import { useAddPerfectSongRehearsalMutation } from '../../../../../../state/api/songsApi.ts'
import {
  ActionIcon,
  Card,
  Center,
  Group,
  ScrollArea,
  SegmentedControl,
  Stack,
  Text,
  Tooltip
} from '@mantine/core'
import {
  IconChecks,
  IconEye,
  IconEyeOff,
  IconList,
  IconListNumbers,
  IconListTree,
  IconPlus
} from '@tabler/icons-react'
import { useDisclosure } from '@mantine/hooks'
import { SongPart, SongSection, SongSettings } from '../../../../../../types/models/Song.ts'
import { useEffect, useRef, useState } from 'react'
import SongArrangementsModal from '../arrangements/SongArrangementsModal.tsx'
import { toast } from 'react-toastify'
import { BandMember } from '../../../../../../types/models/Artist.ts'
import PopoverConfirmation from '../../../../../../components/popover/PopoverConfirmation.tsx'
import SongOutlineSettingsButton from './components/SongOutlineSettingsButton.tsx'
import LoadingOverlayDebounced from '../../../../../../components/loader/LoadingOverlayDebounced.tsx'
import { useMain } from '../../../../../../context/MainContext.tsx'
import { ClickSelectProvider } from '../../../../../../context/ClickSelectContext.tsx'
import CustomRehearsalButton from './components/CustomRehearsalButton.tsx'
import SongParts from '../parts/SongParts.tsx'
import NewHorizontalCard from '../../../../../../components/card/NewHorizontalCard.tsx'
import AddNewSongPart from '../parts/components/AddNewSongPart.tsx'
import AddNewSongSection from '../sections/components/AddNewSongSection.tsx'
import SongSections from '../sections/SongSections.tsx'

enum OutlineView {
  Sections,
  Parts
}

interface SongOutlineWidgetProps {
  sections: SongSection[]
  parts: SongPart[]
  songId: string
  settings: SongSettings
  defaultSongArrangementId?: string
  isFetching?: boolean
  bandMembers?: BandMember[]
  isArtistBand?: boolean
}

function SongOutlineWidget({
  sections,
  parts,
  settings,
  songId,
  defaultSongArrangementId,
  isFetching,
  bandMembers,
  isArtistBand
}: SongOutlineWidgetProps) {
  const [addPerfectRehearsal, { isLoading: isPerfectRehearsalLoading }] =
    useAddPerfectSongRehearsalMutation()

  const [showDetails, setShowDetails] = useState(false)
  const [openedPerfectRehearsalPopover, setOpenedPerfectRehearsalPopover] = useState(false)
  const [openedArrangements, { open: openArrangements, close: closeArrangements }] =
    useDisclosure(false)
  const [openedAdd, { toggle: toggleAdd }] = useDisclosure(false)
  const [outlineView, setOutlineView] = useState(OutlineView.Sections)
  const part_or_section = outlineView === OutlineView.Sections ? 'section' : 'part'
  const partOrSection = outlineView === OutlineView.Sections ? 'Section' : 'Part'

  useEffect(() => setShowDetails(false), [songId])

  const ref = useRef<HTMLDivElement>(null)
  const scrollableRef = useRef<HTMLDivElement>(null)
  const { mainScroll } = useMain()

  const scrollAddIntoView = () => {
    scrollableRef.current.scrollTo({ top: scrollableRef.current.scrollHeight, behavior: 'smooth' })
    mainScroll.ref.current?.scrollTo({
      top: mainScroll.ref.current.scrollHeight,
      behavior: 'smooth'
    })
  }

  function handleShowDetails() {
    setShowDetails(!showDetails)
    if (!showDetails) setTimeout(() => ref.current?.scrollIntoView({ behavior: 'smooth' }), 250)
  }

  async function handleAddPerfectRehearsal() {
    await addPerfectRehearsal({ id: songId }).unwrap()
    toast.info('Perfect rehearsal added!')
    setOpenedPerfectRehearsalPopover(false)
  }

  return (
    <ClickSelectProvider data={parts}>
      <Card ref={ref} variant={'widget'} aria-label={'outline-widget'} p={0}>
        <Stack gap={0}>
          <LoadingOverlayDebounced visible={isFetching} timeout={750} />

          <Group px={'md'} pt={'md'} pb={'sm'} gap={'xxs'}>
            <Text fw={600} inline>
              Outline
            </Text>

            <Tooltip.Group openDelay={500} closeDelay={100}>
              <Tooltip label={`Add New ${partOrSection}`}>
                <ActionIcon
                  aria-label={`add-new-${part_or_section}`}
                  variant={'grey'}
                  size={'sm'}
                  onClick={toggleAdd}
                >
                  <IconPlus size={16} />
                </ActionIcon>
              </Tooltip>

              <Tooltip
                label={
                  parts.length > 0
                    ? showDetails
                      ? 'Hide details'
                      : 'Show Details'
                    : `To show details you need ${part_or_section}s`
                }
              >
                <ActionIcon
                  aria-label={showDetails ? 'hide-details' : 'show-details'}
                  variant={'grey'}
                  size={'sm'}
                  disabled={parts.length === 0}
                  onClick={handleShowDetails}
                >
                  {showDetails ? <IconEyeOff size={16} /> : <IconEye size={16} />}
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
                partsCount={parts.length}
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
                    parts.length === 0
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
                    disabled={parts.length === 0 || !defaultSongArrangementId}
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

              <SongOutlineSettingsButton
                settings={settings}
                parts={parts}
                songId={songId}
                bandMembers={bandMembers}
              />

              <SegmentedControl<OutlineView>
                value={outlineView}
                onChange={setOutlineView}
                size={'xs'}
                color={'gray'}
                radius={'12px'}
                data={[
                  {
                    value: OutlineView.Sections,
                    label: (
                      <Tooltip label={'Sections View'}>
                        <Center>
                          <IconListTree size={16} aria-label={'sections-view'} />
                        </Center>
                      </Tooltip>
                    )
                  },
                  {
                    value: OutlineView.Parts,
                    label: (
                      <Tooltip label={'Parts View'}>
                        <Center>
                          <IconList size={16} aria-label={'parts-view'} />
                        </Center>
                      </Tooltip>
                    )
                  }
                ]}
                styles={{
                  label: {
                    width: '24px',
                    padding: '3px'
                  }
                }}
              />
            </Tooltip.Group>
          </Group>

          <ScrollArea.Autosize
            viewportRef={scrollableRef}
            scrollbars={'y'}
            scrollbarSize={7}
            mah={(showDetails ? 2 : 1) * 383.35}
            style={{ transition: 'max-height 0.25s' }}
          >
            <Stack gap={0}>
              {outlineView === OutlineView.Sections && (
                <SongSections
                  sections={sections}
                  songId={songId}
                  showDetails={showDetails}
                  isFetching={isFetching}
                />
              )}
              {outlineView === OutlineView.Parts && (
                <SongParts
                  parts={parts}
                  songId={songId}
                  showDetails={showDetails}
                  isFetching={isFetching}
                  bandMembers={bandMembers}
                  isArtistBand={isArtistBand}
                />
              )}

              {outlineView === OutlineView.Sections && sections.length === 0 && (
                <NewHorizontalCard ariaLabel={'add-new-song-section-card'} onClick={toggleAdd}>
                  Add New Song Section
                </NewHorizontalCard>
              )}
              {outlineView === OutlineView.Parts && parts.length === 0 && (
                <NewHorizontalCard ariaLabel={'add-new-song-part-card'} onClick={toggleAdd}>
                  Add New Song Part
                </NewHorizontalCard>
              )}

              {outlineView === OutlineView.Sections && (
                <AddNewSongSection
                  songId={songId}
                  opened={openedAdd}
                  onClose={toggleAdd}
                  scrollIntoView={scrollAddIntoView}
                />
              )}
              {outlineView === OutlineView.Parts && (
                <AddNewSongPart
                  songId={songId}
                  opened={openedAdd}
                  onClose={toggleAdd}
                  settings={settings}
                  bandMembers={bandMembers}
                  scrollIntoView={scrollAddIntoView}
                />
              )}
            </Stack>
          </ScrollArea.Autosize>
        </Stack>

        <SongArrangementsModal
          opened={openedArrangements}
          onClose={closeArrangements}
          songId={songId}
          defaultId={defaultSongArrangementId}
        />
      </Card>
    </ClickSelectProvider>
  )
}

export default SongOutlineWidget
