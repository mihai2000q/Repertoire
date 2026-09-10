import { useMoveSongSectionMutation } from './state/api/songSectionsApi.ts'
import { Box, Stack } from '@mantine/core'
import { DragDropContext, Draggable, Droppable } from '@hello-pangea/dnd'
import { useDidUpdate, useListState } from '@mantine/hooks'
import { SongSection } from '../../../../../../types/models/Song.ts'
import SongSectionCard from './components/SongSectionCard.tsx'
import { useMemo } from 'react'
import SongSectionsContextMenu from './components/SongSectionsContextMenu.tsx'
import SongSectionsSelectionDrawer from './components/SongSectionsSelectionDrawer.tsx'
import { useClickSelect } from '../../../../../../context/ClickSelectContext.tsx'

interface SongSectionsWidgetProps {
  sections: SongSection[]
  songId: string
  isFetching?: boolean
  showDetails?: boolean
}

function SongSections({ sections, songId, isFetching, showDetails }: SongSectionsWidgetProps) {
  const [moveSongSection, { isLoading: isMoveLoading }] = useMoveSongSectionMutation()

  const [internalSections, { reorder, setState }] = useListState<SongSection>(sections)
  useDidUpdate(() => setState(sections), [sections])

  const [maxSectionRehearsals, maxSectionProgress] = useMemo(() => {
    let rehearsals = 0
    let progress = 0

    sections.forEach((section) => {
      if (section.rehearsals > rehearsals) rehearsals = section.rehearsals
      if (section.progress > progress) progress = section.progress
    })

    return [rehearsals, progress]
  }, [sections])

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
    <Stack gap={0} aria-label={'song-sections'}>
      <SongSectionsContextMenu sections={sections} songId={songId}>
        <span style={{ display: 'contents' }}>
          <DragDropContext onDragEnd={onSectionsDragEnd}>
            <Droppable droppableId="dnd-list" direction="vertical">
              {(provided) => (
                <Box ref={provided.innerRef} {...provided.droppableProps}>
                  {internalSections.map((section, index) => {
                    // eslint-disable-next-line react-hooks/rules-of-hooks
                    const { isClickSelectionActive } = useClickSelect()
                    return (
                      <Draggable
                        key={section.id}
                        index={index}
                        draggableId={section.id}
                        isDragDisabled={isFetching || isMoveLoading || isClickSelectionActive}
                      >
                        {(provided, snapshot) => (
                          <SongSectionCard
                            section={section}
                            songId={songId}
                            isDragging={snapshot.isDragging}
                            showDetails={showDetails}
                            maxSectionProgress={maxSectionProgress}
                            maxSectionRehearsals={maxSectionRehearsals}
                            draggableProvided={provided}
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
      </SongSectionsContextMenu>
      <SongSectionsSelectionDrawer sections={sections} songId={songId} />
    </Stack>
  )
}

export default SongSections
