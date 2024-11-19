import { useMoveSongSectionMutation } from '../../state/songsApi.ts'
import { ActionIcon, Card, Group, Stack, Text } from '@mantine/core'
import { IconPlus } from '@tabler/icons-react'
import { DragDropContext, Draggable, Droppable } from '@hello-pangea/dnd'
import NewHorizontalCard from '../card/NewHorizontalCard.tsx'
import AddNewSongSection from './AddNewSongSection.tsx'
import { useDisclosure, useListState } from '@mantine/hooks'
import { SongSection as SongSectionType } from '../../types/models/Song.ts'
import SongSection from './SongSection.tsx'

interface SongSectionsProps {
  sections: SongSectionType[]
  songId: string
}

function SongSections({ sections, songId }: SongSectionsProps) {
  const [moveSongSectionMutation, { isLoading: isMoveLoading }] = useMoveSongSectionMutation()

  const [openedAddSongSection, { open: openAddSongSection, close: closeAddSongSection }] =
    useDisclosure(false)

  const [internalSections, { reorder }] = useListState<SongSectionType>(sections)

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
    <Card variant={'panel'} p={0}>
      <Stack gap={0}>
        <Group align={'center'} p={'md'} gap={4}>
          <Text fw={600} inline>
            Sections
          </Text>
          <ActionIcon
            variant={'grey'}
            size={'sm'}
            onClick={openedAddSongSection ? closeAddSongSection : openAddSongSection}
          >
            <IconPlus size={17} />
          </ActionIcon>
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
                          draggableProvided={provided}
                          isDragging={snapshot.isDragging}
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
