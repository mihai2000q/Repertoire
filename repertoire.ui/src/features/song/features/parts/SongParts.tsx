import { useMoveSongPartInSongMutation } from './state/api/songPartsApi.ts'
import { Box, Stack } from '@mantine/core'
import { DragDropContext, Draggable, Droppable } from '@hello-pangea/dnd'
import { useDidUpdate, useListState } from '@mantine/hooks'
import { SongPart } from '../../../../types/models/Song.ts'
import SongPartCard from './components/SongPartCard.tsx'
import { useMemo, useRef } from 'react'
import { toast } from 'react-toastify'
import SongPartsContextMenu from './components/SongPartsContextMenu.tsx'
import SongPartsSelectionDrawer from './components/SongPartsSelectionDrawer.tsx'
import { useClickSelect } from '../../../../context/ClickSelectContext.tsx'
import LoadingOverlayDebounced from '../../../../components/loader/LoadingOverlayDebounced.tsx'
import { useSongContext } from '../../context/SongContext.tsx'

interface SongPartsWidgetProps {
  parts: SongPart[]
  isSongFetching?: boolean
  isPartsFetching?: boolean
}

function SongParts({ parts, isSongFetching, isPartsFetching }: SongPartsWidgetProps) {
  const { songId } = useSongContext()

  const [moveSongPartInSong, { isLoading: isMoveLoading }] = useMoveSongPartInSongMutation()

  const [internalParts, { reorder, setState }] = useListState<SongPart>(parts)
  useDidUpdate(() => setState(parts), [parts])

  const [maxPartProgress] = useMemo(() => {
    let progress = 0

    parts.forEach((part) => {
      if (part.progress > progress) progress = part.progress
    })

    return [progress]
  }, [parts])

  const rehearsalsToastId = useRef<number | string>(null)

  function showRehearsalsToast(partName: string) {
    if (rehearsalsToastId.current) toast.dismiss(rehearsalsToastId.current)
    rehearsalsToastId.current = toast.info(`${partName} rehearsals increased by 1!`)
  }

  function onPartsDragEnd({ source, destination }) {
    reorder({ from: source.index, to: destination?.index || 0 })

    if (!destination || source.index === destination.index) return

    moveSongPartInSong({
      id: parts[source.index].id,
      overId: parts[destination.index].id,
      songId: songId
    })
  }

  return (
    <Stack gap={0} aria-label={'song-parts'}>
      <LoadingOverlayDebounced visible={isPartsFetching && !isSongFetching} timeout={750} />
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
                          isSongFetching ||
                          isPartsFetching ||
                          isMoveLoading ||
                          isClickSelectionActive
                        }
                      >
                        {(provided, snapshot) => (
                          <SongPartCard
                            part={part}
                            isDragging={snapshot.isDragging}
                            maxPartProgress={maxPartProgress}
                            draggableProvided={provided}
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
    </Stack>
  )
}

export default SongParts
