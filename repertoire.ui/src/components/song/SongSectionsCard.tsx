import { useMoveSongSectionMutation } from '../../state/api/songsApi.ts'
import { ActionIcon, Box, Card, Group, Stack, Text, Tooltip } from '@mantine/core'
import { IconEye, IconEyeOff, IconListNumbers, IconPlus } from '@tabler/icons-react'
import { DragDropContext, Draggable, Droppable } from '@hello-pangea/dnd'
import NewHorizontalCard from '../@ui/card/NewHorizontalCard.tsx'
import AddNewSongSection from './AddNewSongSection.tsx'
import { useDidUpdate, useDisclosure, useListState } from '@mantine/hooks'
import { SongSection as SongSectionModel } from '../../types/models/Song.ts'
import SongSectionCard from './SongSectionCard.tsx'
import { useState } from 'react'
import EditSongSectionsOccurrencesModal from './modal/EditSongSectionsOccurrencesModal.tsx'

interface SongSectionsCardProps {
  sections: SongSectionModel[]
  songId: string
}

function SongSectionsCard({ sections, songId }: SongSectionsCardProps) {
  const [moveSongSection, { isLoading: isMoveLoading }] = useMoveSongSectionMutation()

  const [openedOccurrences, { open: openOccurrences, close: closeOccurrences }] =
    useDisclosure(false)
  const [openedAdd, { open: openAdd, close: closeAdd }] = useDisclosure(false)

  const [internalSections, { reorder, setState }] = useListState<SongSectionModel>(sections)
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
        <Group p={'md'} gap={4}>
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

            <Tooltip label={"Edit Sections' Occurrences"}>
              <ActionIcon
                aria-label={'edit-occurrences'}
                variant={'grey'}
                size={'sm'}
                onClick={openOccurrences}
              >
                <IconListNumbers size={16} />
              </ActionIcon>
            </Tooltip>
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

          <AddNewSongSection songId={songId} opened={openedAdd} onClose={closeAdd} />
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
