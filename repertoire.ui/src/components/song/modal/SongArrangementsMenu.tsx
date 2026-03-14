import { SongArrangement } from '../../../types/models/Song.ts'
import { Dispatch, SetStateAction } from 'react'
import { alpha, Button, Center, Menu, ScrollArea, Stack } from '@mantine/core'
import { IconPlus, IconStarFilled } from '@tabler/icons-react'
import { useDidUpdate, useDisclosure, useListState } from '@mantine/hooks'
import { useMoveSongArrangementMutation } from '../../../state/api/songsApi.ts'
import LoadingOverlayDebounced from '../../@ui/loader/LoadingOverlayDebounced.tsx'
import { DragDropContext, Draggable, Droppable } from '@hello-pangea/dnd'
import fixDraggableInMenu from '../../../utils/fixDraggableInMenu.ts'

const invalidArrangementName = '[Invalid Name]'

interface SongArrangementsMenuProps {
  arrangements: SongArrangement[]
  internalArrangements: Map<string, SongArrangement>
  selectedArrangement: SongArrangement
  setSelectedArrangement: Dispatch<SetStateAction<SongArrangement>>
  songId: string
  openAddPopover: () => void
  isFetching?: boolean
  defaultId?: string
}

function SongArrangementsMenu({
  arrangements,
  internalArrangements,
  selectedArrangement,
  setSelectedArrangement,
  songId,
  openAddPopover,
  isFetching,
  defaultId
}: SongArrangementsMenuProps) {
  const [moveArrangement, { isLoading: isMoveLoading }] = useMoveSongArrangementMutation()

  const [openedMenu, { toggle: toggleMenu, close: closeMenu }] = useDisclosure(false)

  const [internalArrangementsForMovement, { setState, reorder }] =
    useListState<SongArrangement>(arrangements)
  useDidUpdate(() => setState(arrangements), [arrangements])

  function onDragEnd({ source, destination }) {
    reorder({ from: source.index, to: destination?.index || 0 })

    if (!destination || source.index === destination.index) return

    moveArrangement({
      id: arrangements[source.index].id,
      overId: arrangements[destination.index].id,
      songId: songId
    })
  }

  function handleClick(arrangement: SongArrangement) {
    setSelectedArrangement(arrangement)
    closeMenu()
  }

  if (!selectedArrangement) return <></>

  return (
    <Menu opened={openedMenu} onChange={toggleMenu} closeOnItemClick={false}>
      <Menu.Target>
        <Button variant={'subtle'} size={'compact-xs'} styles={{ section: { marginLeft: 4 } }}>
          {internalArrangements.get(selectedArrangement.id)?.name?.trim() === ''
            ? invalidArrangementName
            : (internalArrangements.get(selectedArrangement.id)?.name ?? selectedArrangement.name)}
        </Button>
      </Menu.Target>

      <Menu.Dropdown>
        <ScrollArea.Autosize mah={'40vh'} scrollbars={'y'} scrollbarSize={7}>
          <LoadingOverlayDebounced visible={isFetching || isMoveLoading} timeout={750} />

          <Stack gap={0}>
            <DragDropContext onDragEnd={onDragEnd}>
              <Droppable droppableId="dnd-list" direction="vertical">
                {(provided) => (
                  <Stack ref={provided.innerRef} {...provided.droppableProps} gap={2}>
                    {internalArrangementsForMovement.map((arrangement, index) => (
                      <Draggable
                        key={arrangement.id}
                        index={index}
                        draggableId={arrangement.id}
                        isDragDisabled={isFetching || isMoveLoading}
                      >
                        {(provided, snapshot) => (
                          <Menu.Item
                            component={'div'}
                            ref={provided.innerRef}
                            data-active={selectedArrangement.id === arrangement.id}
                            {...provided.dragHandleProps}
                            {...provided.draggableProps}
                            style={(theme) => ({
                              ...fixDraggableInMenu(provided.draggableProps, snapshot, 47).style,
                              cursor: 'pointer',
                              ...(snapshot.isDragging && {
                                cursor: 'grabbing',
                                borderColor: `${theme.colors.primary[3]}`
                              })
                            })}
                            sx={(theme) => ({
                              transition: '0.2s',
                              border: '1px solid transparent',
                              color: theme.colors.gray[7],
                              fontSize: theme.fontSizes.xs,
                              '.mantine-Menu-itemLabel': { display: 'flex', gap: 8 },
                              '.DefaultIcon': { color: theme.colors.gray[6] },
                              '&:hover': {
                                color: theme.colors.gray[8],
                                '.DefaultIcon': {
                                  color: alpha(theme.colors.gray[7], 0.7)
                                },
                                backgroundColor: alpha(theme.colors.gray[1], 1)
                              },

                              ...(selectedArrangement.id === arrangement.id && {
                                color: theme.colors.primary[4],
                                backgroundColor: alpha(theme.colors.primary[1], 0.4),
                                '.DefaultIcon': { color: theme.colors.primary[4] },
                                // same as above
                                '&:hover': {
                                  color: theme.colors.primary[4],
                                  '.DefaultIcon': { color: theme.colors.primary[4] },
                                  backgroundColor: alpha(theme.colors.primary[1], 0.4)
                                }
                              }),

                              ...(defaultId === arrangement.id && {
                                fontWeight: 500
                              })
                            })}
                            onClick={() => handleClick(arrangement)}
                          >
                            {internalArrangements.get(arrangement.id)?.name?.trim() === ''
                              ? invalidArrangementName
                              : (internalArrangements.get(arrangement.id)?.name ??
                                arrangement.name)}
                            {defaultId === arrangement.id && (
                              <Center className={'DefaultIcon'}>
                                <IconStarFilled size={12} />
                              </Center>
                            )}
                          </Menu.Item>
                        )}
                      </Draggable>
                    ))}
                    {provided.placeholder}
                  </Stack>
                )}
              </Droppable>
            </DragDropContext>
            <Menu.Divider />
            <Menu.Item
              fz={'xs'}
              fw={400}
              leftSection={<IconPlus size={12} />}
              onClick={openAddPopover}
            >
              New Arrangement
            </Menu.Item>
          </Stack>
        </ScrollArea.Autosize>
      </Menu.Dropdown>
    </Menu>
  )
}

export default SongArrangementsMenu
