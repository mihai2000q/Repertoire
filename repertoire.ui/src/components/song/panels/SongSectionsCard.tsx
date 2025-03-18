import {
  useAddPartialSongRehearsalMutation,
  useAddPerfectSongRehearsalMutation,
  useMoveSongSectionMutation
} from '../../../state/api/songsApi.ts'
import { ActionIcon, alpha, Box, Card, Group, Popover, Stack, Text, Tooltip } from '@mantine/core'
import {
  IconCheck,
  IconChecks,
  IconEye,
  IconEyeOff,
  IconListNumbers,
  IconPlus,
  IconX
} from '@tabler/icons-react'
import { DragDropContext, Draggable, Droppable } from '@hello-pangea/dnd'
import NewHorizontalCard from '../../@ui/card/NewHorizontalCard.tsx'
import AddNewSongSection from '../AddNewSongSection.tsx'
import { useDidUpdate, useDisclosure, useListState } from '@mantine/hooks'
import { SongSection } from '../../../types/models/Song.ts'
import SongSectionCard from '../SongSectionCard.tsx'
import { useState } from 'react'
import EditSongSectionsOccurrencesModal from '../modal/EditSongSectionsOccurrencesModal.tsx'
import { toast } from 'react-toastify'
import { BandMember } from '../../../types/models/Artist.ts'

interface SongSectionsCardProps {
  sections: SongSection[]
  songId: string
  bandMembers?: BandMember[]
  isArtistBand?: boolean
}

function SongSectionsCard({ sections, songId, bandMembers, isArtistBand }: SongSectionsCardProps) {
  const [moveSongSection, { isLoading: isMoveLoading }] = useMoveSongSectionMutation()
  const [addPartialRehearsal, { isLoading: isPartialRehearsalLoading }] =
    useAddPartialSongRehearsalMutation()
  const [addPerfectRehearsal, { isLoading: isPerfectRehearsalLoading }] =
    useAddPerfectSongRehearsalMutation()

  const [openedPartialRehearsalPopover, setOpenedPartialRehearsalPopover] = useState(false)
  const [openedPerfectRehearsalPopover, setOpenedPerfectRehearsalPopover] = useState(false)

  const [openedOccurrences, { open: openOccurrences, close: closeOccurrences }] =
    useDisclosure(false)
  const [openedAdd, { open: openAdd, close: closeAdd }] = useDisclosure(false)

  const [internalSections, { reorder, setState }] = useListState<SongSection>(sections)
  useDidUpdate(() => setState(sections), [sections])

  const maxSectionProgress =
    sections.length > 0
      ? sections.reduce(
          (accumulator, currentValue) => Math.max(accumulator, currentValue.progress),
          sections[0].progress
        )
      : 0

  const [showDetails, setShowDetails] = useState(false)

  function handleShowDetails() {
    setShowDetails(!showDetails)
  }

  async function handleAddPartialRehearsal() {
    await addPartialRehearsal({ id: songId }).unwrap()
    toast.info('Partial rehearsal added!')
    setOpenedPartialRehearsalPopover(false)
  }

  async function handleAddPerfectRehearsal() {
    await addPerfectRehearsal({ id: songId }).unwrap()
    toast.info('Perfect rehearsal added!')
    setOpenedPerfectRehearsalPopover(false)
  }

  function onSectionsDragEnd({ source, destination }) {
    reorder({ from: source.index, to: destination?.index || 0 })

    if (!destination || source.index === destination.index) return

    moveSongSection({
      id: sections[source.index].id,
      overId: sections[destination.index].id,
      songId: songId
    })
  }

  return (
    <Card variant={'panel'} aria-label={'song-sections'} p={0}>
      <Stack gap={0}>
        <Group p={'md'} gap={'xxs'}>
          <Text fw={600} inline>
            Sections
          </Text>

          <Tooltip.Group openDelay={500} closeDelay={100}>
            <Tooltip label={'Add New Section'}>
              <ActionIcon
                aria-label={'add-new-section'}
                variant={'grey'}
                size={'sm'}
                onClick={openedAdd ? closeAdd : openAdd}
              >
                <IconPlus size={16} />
              </ActionIcon>
            </Tooltip>

            <Tooltip label={showDetails ? 'Hide details' : 'Show Details'}>
              <ActionIcon
                aria-label={showDetails ? 'hide-details' : 'show-details'}
                variant={'grey'}
                size={'sm'}
                onClick={handleShowDetails}
              >
                {showDetails ? <IconEyeOff size={16} /> : <IconEye size={16} />}
              </ActionIcon>
            </Tooltip>

            <Tooltip
              label={
                sections.length === 0
                  ? "To edit sections' occurrences you need sections"
                  : "Edit Sections' Occurrences"
              }
            >
              <ActionIcon
                aria-label={'edit-occurrences'}
                variant={'grey'}
                size={'sm'}
                disabled={sections.length === 0}
                onClick={openOccurrences}
              >
                <IconListNumbers size={16} />
              </ActionIcon>
            </Tooltip>

            <Popover
              opened={openedPartialRehearsalPopover}
              onChange={setOpenedPartialRehearsalPopover}
              transitionProps={{ transition: 'fade-up' }}
              position={'top'}
              withArrow
              shadow={'sm'}
              closeOnClickOutside={!isPartialRehearsalLoading}
            >
              <Popover.Target>
                <Tooltip
                  label={
                    sections.length === 0
                      ? 'To add a partial rehearsal you need sections'
                      : 'Add Partial Rehearsal'
                  }
                  disabled={openedPartialRehearsalPopover}
                >
                  <ActionIcon
                    aria-label={'add-partial-rehearsal'}
                    variant={'grey'}
                    size={'sm'}
                    disabled={sections.length === 0}
                    onClick={() =>
                      setOpenedPartialRehearsalPopover(
                        isPartialRehearsalLoading || !openedPartialRehearsalPopover
                      )
                    }
                  >
                    <IconCheck size={16} />
                  </ActionIcon>
                </Tooltip>
              </Popover.Target>

              <Popover.Dropdown>
                <Group gap={'xxs'}>
                  <Text c={'dimmed'} fw={500} fz={'sm'}>
                    Increase sections&#39; rehearsals based on partial occurrences
                  </Text>
                  <Group gap={'xxs'}>
                    <ActionIcon
                      variant={'subtle'}
                      aria-label={'cancel-partial-rehearsal'}
                      disabled={isPartialRehearsalLoading}
                      onClick={() => setOpenedPartialRehearsalPopover(false)}
                      sx={(theme) => ({
                        color: theme.colors.red[4],
                        '&:hover': {
                          color: theme.colors.red[5],
                          backgroundColor: alpha(theme.colors.red[2], 0.35)
                        },
                        '&[data-disabled]': {
                          color: theme.colors.gray[4],
                          backgroundColor: 'transparent'
                        }
                      })}
                    >
                      <IconX size={16} />
                    </ActionIcon>
                    <ActionIcon
                      variant={'subtle'}
                      c={'green'}
                      aria-label={'confirm-partial-rehearsal'}
                      loading={isPartialRehearsalLoading}
                      onClick={handleAddPartialRehearsal}
                      sx={(theme) => ({
                        '&:hover': { backgroundColor: alpha(theme.colors.green[2], 0.35) }
                      })}
                    >
                      <IconCheck size={16} />
                    </ActionIcon>
                  </Group>
                </Group>
              </Popover.Dropdown>
            </Popover>

            <Popover
              opened={openedPerfectRehearsalPopover}
              onChange={setOpenedPerfectRehearsalPopover}
              transitionProps={{ transition: 'fade-up' }}
              position={'top'}
              withArrow
              shadow={'sm'}
              closeOnClickOutside={!isPerfectRehearsalLoading}
            >
              <Popover.Target>
                <Tooltip
                  label={
                    sections.length === 0
                      ? 'To add a perfect rehearsal you need sections'
                      : 'Add Perfect Rehearsal'
                  }
                  disabled={openedPerfectRehearsalPopover}
                >
                  <ActionIcon
                    aria-label={'add-perfect-rehearsal'}
                    variant={'grey'}
                    size={'sm'}
                    disabled={sections.length === 0}
                    onClick={() =>
                      setOpenedPerfectRehearsalPopover(
                        isPerfectRehearsalLoading || !openedPerfectRehearsalPopover
                      )
                    }
                  >
                    <IconChecks size={16} />
                  </ActionIcon>
                </Tooltip>
              </Popover.Target>

              <Popover.Dropdown>
                <Group gap={'xxs'}>
                  <Text c={'dimmed'} fw={500} fz={'sm'}>
                    Increase sections&#39; rehearsals based on occurrences
                  </Text>
                  <Group gap={'xxs'}>
                    <ActionIcon
                      variant={'subtle'}
                      aria-label={'cancel-perfect-rehearsal'}
                      disabled={isPerfectRehearsalLoading}
                      onClick={() => setOpenedPerfectRehearsalPopover(false)}
                      sx={(theme) => ({
                        color: theme.colors.red[4],
                        '&:hover': {
                          color: theme.colors.red[5],
                          backgroundColor: alpha(theme.colors.red[2], 0.35)
                        },
                        '&[data-disabled]': {
                          color: theme.colors.gray[4],
                          backgroundColor: 'transparent'
                        }
                      })}
                    >
                      <IconX size={16} />
                    </ActionIcon>
                    <ActionIcon
                      variant={'subtle'}
                      c={'green'}
                      aria-label={'confirm-perfect-rehearsal'}
                      loading={isPerfectRehearsalLoading}
                      onClick={handleAddPerfectRehearsal}
                      sx={(theme) => ({
                        '&:hover': { backgroundColor: alpha(theme.colors.green[2], 0.35) }
                      })}
                    >
                      <IconCheck size={16} />
                    </ActionIcon>
                  </Group>
                </Group>
              </Popover.Dropdown>
            </Popover>
          </Tooltip.Group>
        </Group>

        <Stack gap={0}>
          <DragDropContext onDragEnd={onSectionsDragEnd}>
            <Droppable droppableId="dnd-list" direction="vertical">
              {(provided) => (
                <Box ref={provided.innerRef} {...provided.droppableProps}>
                  {internalSections.map((section, index) => (
                    <Draggable
                      key={section.id}
                      index={index}
                      draggableId={section.id}
                      isDragDisabled={isMoveLoading}
                    >
                      {(provided, snapshot) => (
                        <SongSectionCard
                          section={section}
                          songId={songId}
                          draggableProvided={provided}
                          isDragging={snapshot.isDragging}
                          showDetails={showDetails}
                          maxSectionProgress={maxSectionProgress}
                          bandMembers={bandMembers}
                          isArtistBand={isArtistBand}
                        />
                      )}
                    </Draggable>
                  ))}
                  {provided.placeholder}
                </Box>
              )}
            </Droppable>
          </DragDropContext>

          {sections.length === 0 && (
            <NewHorizontalCard
              ariaLabel={'add-new-song-section-card'}
              onClick={openedAdd ? closeAdd : openAdd}
            >
              Add New Song Section
            </NewHorizontalCard>
          )}

          <AddNewSongSection
            songId={songId}
            opened={openedAdd}
            onClose={closeAdd}
            bandMembers={bandMembers}
          />
        </Stack>
      </Stack>

      <EditSongSectionsOccurrencesModal
        opened={openedOccurrences}
        onClose={closeOccurrences}
        sections={sections}
        songId={songId}
      />
    </Card>
  )
}

export default SongSectionsCard
