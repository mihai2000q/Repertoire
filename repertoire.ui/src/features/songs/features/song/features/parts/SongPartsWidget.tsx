import { useAddPerfectSongRehearsalMutation } from '../../../../../../state/api/songsApi.ts'
import { useMoveSongPartMutation } from './state/api/songPartsApi.ts'
import { ActionIcon, Box, Card, Group, ScrollArea, Stack, Text, Tooltip } from '@mantine/core'
import { IconChecks, IconEye, IconEyeOff, IconListNumbers, IconPlus } from '@tabler/icons-react'
import { DragDropContext, Draggable, Droppable } from '@hello-pangea/dnd'
import NewHorizontalCard from '../../../../../../components/card/NewHorizontalCard.tsx'
import AddNewSongPart from './components/AddNewSongPart.tsx'
import { useDidUpdate, useDisclosure, useListState } from '@mantine/hooks'
import { SongPart, SongSettings } from '../../../../../../types/models/Song.ts'
import SongPartCard from './components/SongPartCard.tsx'
import { useEffect, useMemo, useRef, useState } from 'react'
import SongArrangementsModal from '../arrangements/SongArrangementsModal.tsx'
import { toast } from 'react-toastify'
import { BandMember } from '../../../../../../types/models/Artist.ts'
import PopoverConfirmation from '../../../../../../components/popover/PopoverConfirmation.tsx'
import SongPartsSettingsButton from './components/toolbar/SongPartsSettingsButton.tsx'
import LoadingOverlayDebounced from '../../../../../../components/loader/LoadingOverlayDebounced.tsx'
import { useMain } from '../../../../../../context/MainContext.tsx'
import SongPartsContextMenu from './components/SongPartsContextMenu.tsx'
import SongPartsSelectionDrawer from './components/SongPartsSelectionDrawer.tsx'
import {
  ClickSelectProvider,
  useClickSelect
} from '../../../../../../context/ClickSelectContext.tsx'
import CustomRehearsalButton from './components/toolbar/CustomRehearsalButton.tsx'

interface SongPartsWidgetProps {
  parts: SongPart[]
  settings: SongSettings
  songId: string
  defaultSongArrangementId?: string
  isFetching?: boolean
  bandMembers?: BandMember[]
  isArtistBand?: boolean
}

function SongPartsWidget({
  parts,
  settings,
  songId,
  defaultSongArrangementId,
  isFetching,
  bandMembers,
  isArtistBand
}: SongPartsWidgetProps) {
  const [moveSongPart, { isLoading: isMoveLoading }] = useMoveSongPartMutation()
  const [addPerfectRehearsal, { isLoading: isPerfectRehearsalLoading }] =
    useAddPerfectSongRehearsalMutation()

  const [showDetails, setShowDetails] = useState(false)
  const [openedPerfectRehearsalPopover, setOpenedPerfectRehearsalPopover] = useState(false)

  const [openedArrangements, { open: openArrangements, close: closeArrangements }] =
    useDisclosure(false)
  const [openedAdd, { open: openAdd, close: closeAdd }] = useDisclosure(false)

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

  const [internalParts, { reorder, setState }] = useListState<SongPart>(parts)
  useDidUpdate(() => setState(parts), [parts])

  const [maxPartRehearsals, maxPartProgress] = useMemo(() => {
    let rehearsals = 0
    let progress = 0

    parts.forEach((part) => {
      if (part.rehearsals > rehearsals) rehearsals = part.rehearsals
      if (part.progress > progress) progress = part.progress
    })

    return [rehearsals, progress]
  }, [parts])

  const rehearsalsToastId = useRef<number | string>(null)

  function showRehearsalsToast(partName: string) {
    if (rehearsalsToastId.current) toast.dismiss(rehearsalsToastId.current)
    rehearsalsToastId.current = toast.info(`${partName} rehearsals' have been increased by 1!`)
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

  function onPartsDragEnd({ source, destination }) {
    reorder({ from: source.index, to: destination?.index || 0 })

    if (!destination || source.index === destination.index) return

    moveSongPart({
      id: parts[source.index].id,
      overId: parts[destination.index].id,
      songId: songId
    })
  }

  return (
    <ClickSelectProvider data={parts}>
      <Card ref={ref} variant={'widget'} aria-label={'parts-widget'} p={0}>
        <Stack gap={0}>
          <LoadingOverlayDebounced visible={isFetching || isMoveLoading} timeout={750} />

          <Group px={'md'} pt={'md'} pb={'sm'} gap={'xxs'}>
            <Text fw={600} inline>
              Parts
            </Text>

            <Tooltip.Group openDelay={500} closeDelay={100}>
              <Tooltip label={'Add New Part'}>
                <ActionIcon
                  aria-label={'add-new-part'}
                  variant={'grey'}
                  size={'sm'}
                  onClick={openedAdd ? closeAdd : openAdd}
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
                    : 'To show details you need parts'
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
                label={
                  "Increase parts' rehearsals based on occurrences from default arrangement"
                }
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

              <SongPartsSettingsButton
                settings={settings}
                parts={parts}
                songId={songId}
                bandMembers={bandMembers}
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
              <SongPartsContextMenu parts={parts} songId={songId}>
                <span style={{ display: 'contents' }}>
                  <DragDropContext onDragEnd={onPartsDragEnd}>
                    <Droppable droppableId="dnd-list" direction="vertical">
                      {(provided) => (
                        <Box ref={provided.innerRef} {...provided.droppableProps}>
                          {internalParts.map((part, index) => {
                            // eslint-disable-next-line react-hooks/rules-of-hooks
                            const { isClickSelectionActive } = useClickSelect()
                            return (
                              <Draggable
                                key={part.id}
                                index={index}
                                draggableId={part.id}
                                isDragDisabled={
                                  isFetching || isMoveLoading || isClickSelectionActive
                                }
                              >
                                {(provided, snapshot) => (
                                  <SongPartCard
                                    part={part}
                                    songId={songId}
                                    isDragging={snapshot.isDragging}
                                    showDetails={showDetails}
                                    maxPartProgress={maxPartProgress}
                                    maxPartRehearsals={maxPartRehearsals}
                                    draggableProvided={provided}
                                    bandMembers={bandMembers}
                                    isArtistBand={isArtistBand}
                                    showRehearsalsToast={showRehearsalsToast}
                                  />
                                )}
                              </Draggable>
                            )
                          })}
                          {provided.placeholder}
                        </Box>
                      )}
                    </Droppable>
                  </DragDropContext>
                </span>
              </SongPartsContextMenu>
              <SongPartsSelectionDrawer parts={parts} songId={songId} />

              {parts.length === 0 && (
                <NewHorizontalCard
                  ariaLabel={'add-new-song-part-card'}
                  onClick={openedAdd ? closeAdd : openAdd}
                >
                  Add New Song Part
                </NewHorizontalCard>
              )}

              <AddNewSongPart
                songId={songId}
                opened={openedAdd}
                onClose={closeAdd}
                settings={settings}
                bandMembers={bandMembers}
                scrollIntoView={scrollAddIntoView}
              />
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

export default SongPartsWidget
