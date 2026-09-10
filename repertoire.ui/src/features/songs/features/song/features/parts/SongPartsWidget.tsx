import { useMoveSongPartMutation } from './state/api/songPartsApi.ts'
import { Box, Stack } from '@mantine/core'
import { DragDropContext, Draggable, Droppable } from '@hello-pangea/dnd'
import NewHorizontalCard from '../../../../../../components/card/NewHorizontalCard.tsx'
import AddNewSongPart from './components/AddNewSongPart.tsx'
import { useDidUpdate, useListState } from '@mantine/hooks'
import { SongPart, SongSettings } from '../../../../../../types/models/Song.ts'
import SongPartCard from './components/SongPartCard.tsx'
import { useMemo, useRef } from 'react'
import { toast } from 'react-toastify'
import { BandMember } from '../../../../../../types/models/Artist.ts'
import SongPartsContextMenu from './components/SongPartsContextMenu.tsx'
import SongPartsSelectionDrawer from './components/SongPartsSelectionDrawer.tsx'
import { useClickSelect } from '../../../../../../context/ClickSelectContext.tsx'

interface SongPartsWidgetProps {
  parts: SongPart[]
  settings: SongSettings
  songId: string
  showDetails: boolean
  scrollAddIntoView: () => void
  openedAdd: boolean
  toggleAdd: () => void
  isFetching?: boolean
  bandMembers?: BandMember[]
  isArtistBand?: boolean
}

function SongPartsWidget({
  parts,
  settings,
  songId,
  showDetails,
  scrollAddIntoView,
  openedAdd,
  toggleAdd,
  isFetching,
  bandMembers,
  isArtistBand
}: SongPartsWidgetProps) {
  const [moveSongPart, { isLoading: isMoveLoading }] = useMoveSongPartMutation()

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
                        isDragDisabled={isFetching || isMoveLoading || isClickSelectionActive}
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
        <NewHorizontalCard ariaLabel={'add-new-song-part-card'} onClick={toggleAdd}>
          Add New Song Part
        </NewHorizontalCard>
      )}

      <AddNewSongPart
        songId={songId}
        opened={openedAdd}
        onClose={toggleAdd}
        settings={settings}
        bandMembers={bandMembers}
        scrollIntoView={scrollAddIntoView}
      />
    </Stack>
  )
}

export default SongPartsWidget
