import { SongPart as SongPartModel } from '../../../../../../../types/models/Song.ts'
import {
  ActionIcon,
  alpha,
  Avatar,
  Box,
  Center,
  Collapse,
  Group,
  HoverCard,
  Menu,
  NumberFormatter,
  Progress,
  Stack,
  Text,
  Tooltip
} from '@mantine/core'
import {
  IconCheck,
  IconDots,
  IconEdit,
  IconGripVertical,
  IconLocationPlus,
  IconTrash,
  IconUser
} from '@tabler/icons-react'
import { DraggableProvided } from '@hello-pangea/dnd'
import { useDisclosure, useHover, useMergedRef } from '@mantine/hooks'
import { toast } from 'react-toastify'
import {
  useDeleteSongPartMutation,
  useUpdateSongPartMutation
} from '../state/api/songPartsApi.ts'
import EditSongPartModal from './modal/EditSongPartModal.tsx'
import WarningModal from '../../../../../../../components/modal/WarningModal.tsx'
import { BandMember } from '../../../../../../../types/models/Artist.ts'
import useInstrumentIcon from '../../../../../../../hooks/useInstrumentIcon.tsx'
import { useEffect, useState } from 'react'
import useDoubleMenu from '../../../../../../../hooks/useDoubleMenu.ts'
import { ContextMenu } from '../../../../../../../components/menu/ContextMenu.tsx'
import useClickSelectSelectable from '../../../../../../../hooks/useClickSelectSelectable.ts'

function getRehearsalsMarginLeft(rehearsalsMaxLength: number) {
  return rehearsalsMaxLength > 4
    ? 'xs'
    : rehearsalsMaxLength > 3
      ? 'md'
      : rehearsalsMaxLength > 2
        ? 20
        : rehearsalsMaxLength > 1
          ? 23
          : 27
}

function getRehearsalsWidth(rehearsalsMaxLength: number) {
  return (rehearsalsMaxLength > 2 ? 9 : rehearsalsMaxLength > 1 ? 10 : 12) * rehearsalsMaxLength
}

interface SongPartCardProps {
  part: SongPartModel
  songId: string
  isDragging: boolean
  showDetails: boolean
  maxPartProgress: number
  maxPartRehearsals: number
  draggableProvided?: DraggableProvided
  bandMembers?: BandMember[]
  isArtistBand?: boolean
  showRehearsalsToast?: (name: string) => void
}

function SongPartCard({
  part,
  songId,
  isDragging,
  showDetails,
  maxPartProgress,
  maxPartRehearsals,
  draggableProvided,
  bandMembers,
  isArtistBand,
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

  const [rehearsalsMarginLeft, setRehearsalsMarginLeft] = useState(
    getRehearsalsMarginLeft(maxPartRehearsals.toString().length)
  )
  const [rehearsalsWidth, setRehearsalsWidth] = useState(
    getRehearsalsWidth(maxPartRehearsals.toString().length)
  )
  useEffect(() => {
    const rehearsalsMaxLength = maxPartRehearsals.toString().length
    setRehearsalsMarginLeft(getRehearsalsMarginLeft(rehearsalsMaxLength))
    setRehearsalsWidth(getRehearsalsWidth(rehearsalsMaxLength))
  }, [maxPartRehearsals])

  const [updateSongPartMutation, { isLoading: isUpdateLoading }] = useUpdateSongPartMutation()
  const [deleteSongPartMutation, { isLoading: isDeleteLoading }] = useDeleteSongPartMutation()

  const getInstrumentIcon = useInstrumentIcon()

  const { openedMenu, toggleMenu, openedContextMenu, toggleContextMenu } = useDoubleMenu()

  const [openedEditSongPart, { open: openEditSongPart, close: closeEditSongPart }] =
    useDisclosure(false)
  const [openedDeleteWarning, { open: openDeleteWarning, close: closeDeleteWarning }] =
    useDisclosure(false)

  const isSelected = hovered || openedMenu || openedContextMenu || isDragging || isClickSelected

  async function handleAddRehearsal() {
    await updateSongPartMutation({
      ...part,
      bandMemberId: part.bandMember?.id,
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
      <Menu.Item leftSection={<IconEdit size={14} />} onClick={openEditSongPart}>
        Edit
      </Menu.Item>
      <Menu.Item leftSection={<IconTrash size={14} />} c={'red.5'} onClick={openDeleteWarning}>
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
          sx={(theme) => ({
            cursor: 'default',
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

            {isArtistBand && part.bandMember && (
              <HoverCard openDelay={200} position="top">
                <HoverCard.Target>
                  <Avatar
                    size={25}
                    color={part.bandMember.color}
                    src={part.bandMember.imageUrl}
                    alt={part.bandMember.imageUrl && part.bandMember.name}
                  >
                    <IconUser size={15} />
                  </Avatar>
                </HoverCard.Target>
                <HoverCard.Dropdown>
                  <Group gap={'xs'} maw={200} wrap={'nowrap'}>
                    <Avatar
                      size={60}
                      color={part.bandMember.color}
                      src={part.bandMember.imageUrl}
                      alt={part.bandMember.imageUrl && part.bandMember.name}
                      style={(theme) => ({ boxShadow: theme.shadows.sm })}
                    >
                      <IconUser size={30} />
                    </Avatar>
                    <Stack gap={0}>
                      <Text fw={500} lineClamp={2}>
                        {part.bandMember.name}
                      </Text>
                      {part.bandMember.roles.slice(0, 2).map((role, index) => (
                        <Text key={role.id} c={'dimmed'} fz={'xs'} lineClamp={1} lh={1.05}>
                          {role.name}
                          {index === 1 && part.bandMember.roles.length > 2 && ' ...'}
                        </Text>
                      ))}
                    </Stack>
                  </Group>
                </HoverCard.Dropdown>
              </HoverCard>
            )}

            {part.instrument && (
              <Box aria-label={'instrument-icon'} c={'primary.7'} w={16} h={16}>
                <Tooltip openDelay={200} label={part.instrument?.name} withArrow>
                  {getInstrumentIcon(part.instrument)}
                </Tooltip>
              </Box>
            )}

            <Text flex={1} truncate={'end'}>
              {part.name}
            </Text>

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
                  <IconLocationPlus size={15} />
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
                  >
                    <IconDots size={20} />
                  </ActionIcon>
                </Menu.Target>
                <Menu.Dropdown>{menuDropdown}</Menu.Dropdown>
              </Menu>
            </Group>
          </Group>

          <Collapse expanded={showDetails}>
            <Group
              aria-label={`song-part-details-${part.name}`}
              pt={'md'}
              gap={'lg'}
              pr={'lg'}
            >
              <Tooltip.Floating
                role={'tooltip'}
                label={
                  <>
                    Rehearsals: <NumberFormatter value={part.rehearsals} />
                  </>
                }
              >
                <Text
                  ml={rehearsalsMarginLeft}
                  w={rehearsalsWidth}
                  fz={12}
                  ta={'center'}
                  fw={500}
                  c={'dimmed'}
                  inline
                  data-testid={'rehearsals'}
                >
                  <NumberFormatter value={part.rehearsals} />
                </Text>
              </Tooltip.Floating>

              <Tooltip.Floating role={'tooltip'} label={`Confidence: ${part.confidence}%`}>
                <Progress
                  flex={1}
                  size={'sm'}
                  value={part.confidence}
                  aria-label={'confidence'}
                />
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

      <EditSongPartModal
        opened={openedEditSongPart}
        onClose={closeEditSongPart}
        part={part}
        bandMembers={bandMembers}
      />
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
