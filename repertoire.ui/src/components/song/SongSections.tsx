import { useMoveSongSectionMutation } from '../../state/songsApi.ts'
import { ActionIcon, Card, Group, Stack, Text, Tooltip } from '@mantine/core'
import { IconEye, IconEyeOff, IconPlus } from '@tabler/icons-react'
import { DragDropContext, Draggable, Droppable } from '@hello-pangea/dnd'
import NewHorizontalCard from '../@ui/card/NewHorizontalCard.tsx'
import AddNewSongSection from './AddNewSongSection.tsx'
import { useDidUpdate, useDisclosure, useListState } from '@mantine/hooks'
import { SongSection as SongSectionModel } from '../../types/models/Song.ts'
import SongSection from './SongSection.tsx'
import { useState } from 'react'

interface SongSectionsProps {
  sections: SongSectionModel[]
  songId: string
}

function SongSections({ sections, songId }: SongSectionsProps) {
  const [moveSongSectionMutation, { isLoading: isMoveLoading }] = useMoveSongSectionMutation()

  const [openedAddSongSection, { open: openAddSongSection, close: closeAddSongSection }] =
    useDisclosure(false)

  const [internalSections, { reorder, setState }] = useListState<SongSectionModel>(sections)
  useDidUpdate(() => {
    setState(sections)
  }, [sections])

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

    if (source.index === destination.index || !destination) return

    moveSongSectionMutation({
      id: sections[source.index].id,
      overId: sections[destination.index].id,
      songId: songId
    })
  }

  return (
    <Card variant={'panel'} aria-label={'song-sections'} p={0}>
      <Stack gap={0}>
        <Group align={'center'} p={'md'} gap={4}>
          <Text fw={600} inline>
            Sections
          </Text>

          <Tooltip.Group openDelay={500} closeDelay={100}>
            <Tooltip label={'Add New Section'}>
              <ActionIcon
                aria-label={'add-new-section'}
                variant={'grey'}
                size={'sm'}
                onClick={openedAddSongSection ? closeAddSongSection : openAddSongSection}
              >
                <IconPlus size={17} />
              </ActionIcon>
            </Tooltip>

            <Tooltip label={showDetails ? 'Hide details' : 'Show Details'}>
              <ActionIcon
                aria-label={showDetails ? 'hide-details' : 'show-details'}
                variant={'grey'}
                size={'sm'}
                onClick={handleShowDetails}
              >
                {showDetails ? <IconEyeOff size={17} /> : <IconEye size={17} />}
              </ActionIcon>
            </Tooltip>
          </Tooltip.Group>
        </Group>

        <Stack gap={0}>
          <DragDropContext onDragEnd={onSectionsDragEnd}>
            <Droppable droppableId="dnd-list" direction="vertical">
              {(provided) => (
                <Stack gap={0} ref={provided.innerRef} {...provided.droppableProps}>
                  {internalSections.map((section, index) => (
                    <Draggable
                      key={section.id}
                      index={index}
                      draggableId={section.id}
                      isDragDisabled={isMoveLoading}
                    >
                      {(provided, snapshot) => (
                        <SongSection
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
                </Stack>
              )}
            </Droppable>
          </DragDropContext>
          {internalSections.length === 0 && (
            <NewHorizontalCard
              ariaLabel={'add-new-song-section-card'}
              onClick={openedAddSongSection ? closeAddSongSection : openAddSongSection}
            >
              Add New Song Section
            </NewHorizontalCard>
          )}
          <AddNewSongSection
            songId={songId}
            opened={openedAddSongSection}
            onClose={closeAddSongSection}
          />
        </Stack>
      </Stack>
    </Card>
  )
}

export default SongSections
