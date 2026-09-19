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
import LoadingOverlayDebounced from '../../../../../../components/loader/LoadingOverlayDebounced.tsx'
import { useSongContext } from '../../context/SongContext.tsx'

interface SongSectionsWidgetProps {
  sections: SongSection[]
  isSongFetching?: boolean
  isSectionsFetching?: boolean
  scrollIntoView?: () => void
}

function SongSections({
  sections,
  isSectionsFetching,
  isSongFetching,
  scrollIntoView
}: SongSectionsWidgetProps) {
  const { songId } = useSongContext()

  const [moveSongSection, { isLoading: isMoveLoading }] = useMoveSongSectionMutation()

  const [internalSections, { reorder, setState }] = useListState<SongSection>(sections)
  useDidUpdate(() => setState(sections), [sections])

  const [maxSectionProgress] = useMemo(() => {
    let progress = 0

    sections.forEach((section) => {
      if (section.progress > progress) progress = section.progress
    })

    return [progress]
  }, [sections])

  function onSectionsDragEnd({ source, destination }) {
    if (!destination || source.index === destination.index) return

    reorder({ from: source.index, to: destination.index })
    moveSongSection({
      id: sections[source.index].id,
      overId: sections[destination.index].id,
      songId: songId
    })
  }

  return (
    <Stack gap={0} aria-label={'song-sections'}>
      <LoadingOverlayDebounced visible={isSectionsFetching && !isSongFetching} timeout={750} />
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
                        isDragDisabled={
                          isSongFetching ||
                          isSectionsFetching ||
                          isMoveLoading ||
                          isClickSelectionActive
                        }
                      >
                        {(provided, snapshot) => (
                          <SongSectionCard
                            section={section}
                            maxSectionProgress={maxSectionProgress}
                            isDragging={snapshot.isDragging}
                            draggableProvided={provided}
                            scrollIntoView={scrollIntoView}
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
