import { SongPart as SongPartModel } from '../../../../../../../types/models/Song.ts'
import {
  ActionIcon,
  alpha,
  Center,
  Collapse,
  Group,
  Menu,
  NumberFormatter,
  Progress,
  Space,
  Stack,
  Text,
  Tooltip
} from '@mantine/core'
import {
  IconCheck,
  IconDots,
  IconEdit,
  IconGripVertical,
  IconRefresh,
  IconTrash
} from '@tabler/icons-react'
import { DraggableProvided } from '@hello-pangea/dnd'
import { useDisclosure, useHover, useMergedRef } from '@mantine/hooks'
import { toast } from 'react-toastify'
import { useDeleteSongPartMutation, useUpdateSongPartMutation } from '../state/api/songPartsApi.ts'
import EditSongPartModal from './modal/EditSongPartModal.tsx'
import WarningModal from '../../../../../../../components/modal/WarningModal.tsx'
import useInstrumentIcon from '../../../../../../../hooks/useInstrumentIcon.tsx'
import { MouseEvent, useEffect, useState } from 'react'
import useDoubleMenu from '../../../../../../../hooks/useDoubleMenu.ts'
import { ContextMenu } from '../../../../../../../components/menu/ContextMenu.tsx'
import useClickSelectSelectable from '../../../../../../../hooks/useClickSelectSelectable.ts'
import { useSongContext } from '../../../context/SongContext.tsx'
import RehearsalsBadge from '../../outline/components/RehearsalsBadge.tsx'
import BandMembersGroup from '../../outline/components/BandMembersGroup.tsx'
import { useSongOutlineContext } from '../../outline/context/SongOutlineContext.tsx'

interface SongPartCardProps {
  part: SongPartModel
  isDragging: boolean
  maxPartProgress: number
  draggableProvided?: DraggableProvided
  showRehearsalsToast?: (name: string) => void
}

function SongPartCard({
  part,
  isDragging,
  maxPartProgress,
  draggableProvided,
  showRehearsalsToast
}: SongPartCardProps) {
  const { ref: hoverRef, hovered } = useHover()
  const {
    ref: selectableRef,
    isClickSelected,
    isClickSelectionActive,
    isLastInSelection
  } = useClickSelectSelectable(part.id)
  const ref = useMergedRef(hoverRef, draggableProvided?.innerRef, selectableRef)

  const { songId, isArtistBand } = useSongContext()
  const { showDetails } = useSongOutlineContext()

  const [updateSongPartMutation, { isLoading: isUpdateLoading }] = useUpdateSongPartMutation()
  const [deleteSongPartMutation, { isLoading: isDeleteLoading }] = useDeleteSongPartMutation()

  const getInstrumentIcon = useInstrumentIcon()

  const { openedMenu, toggleMenu, openedContextMenu, toggleContextMenu } = useDoubleMenu()
  const [openedDetails, setOpenedDetails] = useState(false)
  useEffect(() => setOpenedDetails(showDetails), [showDetails])

  const [openedEditSongPart, { open: openEditSongPart, close: closeEditSongPart }] =
    useDisclosure(false)
  const [openedDeleteWarning, { open: openDeleteWarning, close: closeDeleteWarning }] =
    useDisclosure(false)

  const isSelected = hovered || openedMenu || openedContextMenu || isDragging || isClickSelected

  function handleClick(e: MouseEvent) {
    if (e.ctrlKey || e.shiftKey) return
    e.stopPropagation()
    setOpenedDetails(!openedDetails)
  }

  async function handleAddRehearsal(e: MouseEvent) {
    e.stopPropagation()

    await updateSongPartMutation({
      ...part,
      instrumentId: part.instrument?.id,
      rehearsals: part.rehearsals + 1
    }).unwrap()
    showRehearsalsToast?.(part.name)
  }

  async function handleDelete() {
    await deleteSongPartMutation({ id: part.id, songId: songId }).unwrap()
    toast.success(`${part.name} deleted!`)
  }

  const menuDropdown = (
    <>
      <Menu.Item
        leftSection={<IconEdit size={14} />}
        onClick={(e) => {
          e.stopPropagation()
          openEditSongPart()
        }}
      >
        Edit
      </Menu.Item>
      <Menu.Item
        leftSection={<IconTrash size={14} />}
        c={'red.5'}
        onClick={(e) => {
          e.stopPropagation()
          openDeleteWarning()
        }}
      >
        Delete
      </Menu.Item>
    </>
  )

  return (
    <ContextMenu
      opened={openedContextMenu}
      onChange={toggleContextMenu}
      disabled={isClickSelectionActive}
    >
      <ContextMenu.Target>
        <Stack
          ref={ref}
          py={'xs'}
          aria-label={`song-part-${part.name}`}
          aria-selected={isSelected}
          gap={0}
          onClick={handleClick}
          sx={(theme) => ({
            cursor: 'pointer',
            transition: '0.25s',
            borderRadius: 0,
            border: '1px solid transparent',
            ...(isSelected && {
              boxShadow: theme.shadows.md,
              backgroundColor: alpha(theme.colors.primary[0], 0.15)
            }),

            ...(isClickSelected && {
              boxShadow: 'none',
              backgroundColor: alpha(theme.colors.primary[0], 0.15),
              ...(hovered && {
                boxShadow: theme.shadows.xs,
                backgroundColor: alpha(theme.colors.primary[0], 0.35)
              }),
              ...(isLastInSelection && {
                boxShadow: theme.shadows.lg
              })
            }),

            ...(isDragging && {
              boxShadow: theme.shadows.xl,
              borderRadius: '16px',
              backgroundColor: alpha(theme.white, 0.33),
              border: `1px solid ${alpha(theme.colors.primary[9], 0.33)}`
            })
          })}
          {...draggableProvided?.draggableProps}
        >
          <Group gap={'xs'} px={'md'}>
            {!isClickSelected ? (
              <ActionIcon
                aria-label={'drag-handle'}
                variant={'subtle'}
                size={'lg'}
                {...draggableProvided?.dragHandleProps}
                disabled={isClickSelectionActive}
                sx={{ '&[data-disabled="true"]': { backgroundColor: 'transparent' } }}
              >
                <IconGripVertical size={20} />
              </ActionIcon>
            ) : (
              <Center
                m={6.5}
                data-testid={'selected-checkmark'}
                w={21}
                h={21}
                style={(theme) => ({
                  borderRadius: '100%',
                  backgroundColor: alpha(theme.colors.green[2], 0.95)
                })}
              >
                <IconCheck color={'white'} size={'75%'} />
              </Center>
            )}

            {isArtistBand && part.bandMembers.length > 0 && (
              <BandMembersGroup bandMembers={part.bandMembers} />
            )}
            {part.instrument && (
              <Tooltip openDelay={200} label={part.instrument?.name} withArrow>
                <Center aria-label={'instrument-icon'} c={'primary.7'} w={16} h={16}>
                  {getInstrumentIcon(part.instrument)}
                </Center>
              </Tooltip>
            )}
            <Text fw={500} truncate={'end'}>
              {part.name}
            </Text>
            <RehearsalsBadge rehearsals={part.rehearsals} />

            <Space flex={1} />

            <Group gap={2}>
              <Tooltip label={'Add Rehearsal'} openDelay={200} disabled={isClickSelectionActive}>
                <ActionIcon
                  variant={'subtle'}
                  size={'md'}
                  loading={isUpdateLoading}
                  aria-label={'add-rehearsal'}
                  disabled={isClickSelectionActive}
                  sx={{ '&[data-disabled="true"]': { backgroundColor: 'transparent' } }}
                  onClick={handleAddRehearsal}
                >
                  <IconRefresh size={15} />
                </ActionIcon>
              </Tooltip>

              <Menu opened={openedMenu} onChange={toggleMenu}>
                <Menu.Target>
                  <ActionIcon
                    variant={'subtle'}
                    size={'md'}
                    aria-label={'more-menu'}
                    disabled={isClickSelectionActive}
                    sx={{ '&[data-disabled="true"]': { backgroundColor: 'transparent' } }}
                    onClick={(e) => e.stopPropagation()}
                  >
                    <IconDots size={20} />
                  </ActionIcon>
                </Menu.Target>
                <Menu.Dropdown>{menuDropdown}</Menu.Dropdown>
              </Menu>
            </Group>
          </Group>

          <Collapse expanded={openedDetails}>
            <Group
              aria-label={`song-part-details-${part.name}`}
              pt={'md'}
              pb={'xxs'}
              pl={'60px'}
              pr={'40px'}
            >
              <Tooltip.Floating role={'tooltip'} label={`Confidence: ${part.confidence}%`}>
                <Progress flex={1} size={'sm'} value={part.confidence} aria-label={'confidence'} />
              </Tooltip.Floating>

              <Tooltip.Floating
                role={'tooltip'}
                label={
                  <>
                    Progress: <NumberFormatter value={part.progress} />
                  </>
                }
              >
                <Progress
                  flex={1}
                  size={'sm'}
                  aria-label={'progress'}
                  value={part.progress === 0 ? 0 : (part.progress / maxPartProgress) * 100}
                  color={'green'}
                />
              </Tooltip.Floating>
            </Group>
          </Collapse>
        </Stack>
      </ContextMenu.Target>

      <ContextMenu.Dropdown>{menuDropdown}</ContextMenu.Dropdown>

      <EditSongPartModal opened={openedEditSongPart} onClose={closeEditSongPart} part={part} />
      <WarningModal
        opened={openedDeleteWarning}
        onClose={closeDeleteWarning}
        title={`Delete Part`}
        description={
          <Group gap={'xxs'}>
            <Text>Are you sure you want to delete</Text>
            <Text fw={600}>{part.name}</Text>
            <Text>?</Text>
          </Group>
        }
        onYes={handleDelete}
        isLoading={isDeleteLoading}
      />
    </ContextMenu>
  )
}

export default SongPartCard
