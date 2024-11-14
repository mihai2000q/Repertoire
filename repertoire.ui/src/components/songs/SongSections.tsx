import { useMoveSongSectionMutation } from '../../state/songsApi.ts'
import { ActionIcon, alpha, Card, Group, Stack, Text } from '@mantine/core'
import { IconDots, IconGripVertical, IconPlus } from '@tabler/icons-react'
import { DragDropContext, Draggable, Droppable } from '@hello-pangea/dnd'
import NewHorizontalCard from '../card/NewHorizontalCard.tsx'
import AddNewSongSection from './AddNewSongSection.tsx'
import { useDisclosure } from '@mantine/hooks'
import { SongSection } from '../../types/models/Song.ts'

interface SongSectionsProps {
  sections: SongSection[]
  songId: string
}

function SongSections({ sections, songId }: SongSectionsProps) {
  const [moveSongSectionMutation, { isLoading: isMoveLoading }] = useMoveSongSectionMutation()

  const [openedAddSongSection, { open: openAddSongSection, close: closeAddSongSection }] =
    useDisclosure(false)

  function onSectionsDragEnd({ destination, source }) {
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
                  {sections.map((section, index) => (
                    <Draggable key={section.id} index={index} draggableId={section.id} isDragDisabled={isMoveLoading}>
                      {(provided, snapshot) => (
                        <Group
                          key={section.id}
                          align={'center'}
                          gap={'xs'}
                          py={'xs'}
                          px={'md'}
                          sx={(theme) => ({
                            transition: '0.25s',
                            borderRadius: snapshot.isDragging ? '16px' : '0px',
                            border: snapshot.isDragging ? `1px solid ${alpha(theme.colors.cyan[9], 0.33)}` : '1px solid transparent',

                            '&:hover': {
                              boxShadow: theme.shadows.xl,
                              backgroundColor: alpha(theme.colors.cyan[0], 0.15)
                            }
                          })}
                          ref={provided.innerRef}
                          {...provided.draggableProps}
                        >
                          <ActionIcon variant={'subtle'} size={'lg'} {...provided.dragHandleProps}>
                            <IconGripVertical size={20} />
                          </ActionIcon>

                          <Text inline fw={600}>{section.songSectionType.name}</Text>
                          <Text flex={1} inline>{section.name}</Text>

                          <ActionIcon variant={'subtle'} size={'lg'}>
                            <IconDots size={20} />
                          </ActionIcon>
                        </Group>
                      )}
                    </Draggable>
                  ))}
                  {provided.placeholder}
                </Stack>
              )}
            </Droppable>
          </DragDropContext>
          {sections.length === 0 && (
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
