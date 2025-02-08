import { BandMember } from '../../../types/models/Artist.ts'
import {
  ActionIcon,
  Avatar,
  Card,
  Group,
  Menu,
  ScrollArea,
  Space,
  Stack,
  Text
} from '@mantine/core'
import { IconChevronLeft, IconChevronRight, IconDots, IconUserPlus } from '@tabler/icons-react'
import { useDidUpdate, useDisclosure, useListState } from '@mantine/hooks'
import AddNewBandMemberModal from '../modal/AddNewBandMemberModal.tsx'
import BandMemberCard from '../BandMemberCard.tsx'
import { useRef, useState } from 'react'
import { useMoveBandMemberMutation } from '../../../state/api/artistsApi.ts'
import { DragDropContext, Draggable, Droppable } from '@hello-pangea/dnd'

interface BandMembersCardProps {
  bandMembers: BandMember[]
  artistId: string
}

function BandMembersCard({ bandMembers, artistId }: BandMembersCardProps) {
  const [moveBandMember, { isLoading: isMoveLoading }] = useMoveBandMemberMutation()

  const viewportRef = useRef<HTMLDivElement>(null)

  const [disableBack, setDisableBack] = useState(false)
  const [disableForward, setDisableForward] = useState(false)
  useDidUpdate(() => {
    setDisableBack(viewportRef.current?.scrollLeft === 0)
    setDisableForward(viewportRef.current?.scrollWidth === viewportRef.current?.clientWidth)
  }, [viewportRef.current])

  const [internalMembers, { reorder, setState }] = useListState<BandMember>(bandMembers)
  useDidUpdate(() => setState(bandMembers), [bandMembers])

  const [openedAddNewBandMember, { open: openAddNewBandMember, close: closeAddNewBandMember }] =
    useDisclosure(false)

  const handleMembersNav = (direction: 'left' | 'right') => {
    if (!viewportRef.current) return
    viewportRef.current.scrollBy({ left: direction === 'left' ? -100 : 100, behavior: 'smooth' })
  }

  const handleOnScroll = () => {
    const viewport = viewportRef.current
    setDisableBack(viewport?.scrollLeft <= 0)
    setDisableForward(viewport?.scrollWidth <= viewport?.clientWidth + viewport?.scrollLeft)
  }

  function onMembersDragEnd({ source, destination }) {
    reorder({ from: source.index, to: destination?.index || 0 })

    if (!destination || source.index === destination.index) return

    moveBandMember({
      id: bandMembers[source.index].id,
      overId: bandMembers[destination.index].id,
      artistId: artistId
    })
  }

  return (
    <Card variant={'panel'} aria-label={'band-members-card'} p={0} mih={140}>
      <Stack gap={0}>
        <Group px={'md'} pt={'xs'} gap={'xs'}>
          <Text fw={600}>Band Members</Text>

          {bandMembers.length > 0 && (
            <Group gap={4}>
              <ActionIcon
                aria-label={'back-button'}
                size={'sm'}
                variant={'grey'}
                radius={'50%'}
                disabled={disableBack}
                style={{ cursor: disableBack ? 'default' : 'pointer' }}
                onClick={() => handleMembersNav('left')}
              >
                <IconChevronLeft size={14} />
              </ActionIcon>

              <ActionIcon
                aria-label={'forward-button'}
                size={'sm'}
                variant={'grey'}
                radius={'50%'}
                disabled={disableForward}
                style={{ cursor: disableForward ? 'default' : 'pointer' }}
                onClick={() => handleMembersNav('right')}
              >
                <IconChevronRight size={14} />
              </ActionIcon>
            </Group>
          )}

          <Space flex={1} />

          <Menu position={'bottom-end'}>
            <Menu.Target>
              <ActionIcon size={'md'} variant={'grey'} aria-label={'band-members-more-menu'}>
                <IconDots size={15} />
              </ActionIcon>
            </Menu.Target>

            <Menu.Dropdown>
              <Menu.Item leftSection={<IconUserPlus size={15} />} onClick={openAddNewBandMember}>
                Add New Member
              </Menu.Item>
            </Menu.Dropdown>
          </Menu>
        </Group>

        <ScrollArea
          scrollbars={'x'}
          scrollbarSize={5}
          viewportRef={viewportRef}
          viewportProps={{ onScroll: handleOnScroll }}
        >
          <Group gap={'xs'} wrap={'nowrap'} align={'start'} px={'lg'} pb={'lg'} pt={'xs'}>
            <DragDropContext onDragEnd={onMembersDragEnd}>
              <Droppable droppableId="dnd-list" direction={'horizontal'}>
                {(provided) => (
                  <Group
                    gap={'xs'}
                    wrap={'nowrap'}
                    align={'start'}
                    ref={provided.innerRef}
                    {...provided.droppableProps}
                  >
                    {internalMembers.map((bandMember, index) => (
                      <Draggable
                        key={bandMember.id}
                        index={index}
                        draggableId={bandMember.id}
                        isDragDisabled={isMoveLoading}
                      >
                        {(provided) => (
                          <BandMemberCard
                            key={bandMember.id}
                            bandMember={bandMember}
                            artistId={artistId}
                            draggableProvided={provided}
                          />
                        )}
                      </Draggable>
                    ))}
                    {provided.placeholder}
                  </Group>
                )}
              </Droppable>
            </DragDropContext>

            {bandMembers.length === 0 && (
              <Stack
                aria-label={`add-new-band-member-card`}
                align={'center'}
                w={75}
                gap={4}
                sx={{
                  cursor: 'pointer',
                  transition: '0.3s',
                  '&:hover': { transform: 'scale(1.1)' }
                }}
                onClick={openAddNewBandMember}
              >
                <Avatar
                  size={'lg'}
                  sx={(theme) => ({
                    boxShadow: theme.shadows.md,
                    '&:hover': { boxShadow: theme.shadows.lg }
                  })}
                >
                  <IconUserPlus size={25} />
                </Avatar>
                <Text c={'dimmed'} ta={'center'} fw={500} lh={'xs'}>
                  Add New Member
                </Text>
              </Stack>
            )}
          </Group>
        </ScrollArea>
      </Stack>

      <AddNewBandMemberModal
        opened={openedAddNewBandMember}
        onClose={closeAddNewBandMember}
        artistId={artistId}
      />
    </Card>
  )
}

export default BandMembersCard
