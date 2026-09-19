import { SongPart as SongPartModel } from '../../../../../../../types/models/Song.ts'
import {
  ActionIcon,
  alpha,
  Center,
  Collapse,
  Group,
  Space,
  Stack,
  Text,
  Tooltip
} from '@mantine/core'
import { IconEdit, IconRefresh, IconTrash } from '@tabler/icons-react'
import { DraggableProvided } from '@hello-pangea/dnd'
import { useDisclosure, useHover, useMergedRef } from '@mantine/hooks'
import { toast } from 'react-toastify'
import WarningModal from '../../../../../../../components/modal/WarningModal.tsx'
import useInstrumentIcon from '../../../../../../../hooks/useInstrumentIcon.tsx'
import { ContextMenu } from '../../../../../../../components/menu/ContextMenu.tsx'
import useClickSelectSelectable from '../../../../../../../hooks/useClickSelectSelectable.ts'
import {
  useDeleteSongPartMutation,
  useUpdateSongPartMutation
} from '../../parts/state/api/songPartsApi.ts'
import EditSongPartModal from '../../parts/components/modal/EditSongPartModal.tsx'
import BandMemberAvatar from '../../../../../../../components/avatar/BandMemberAvatar.tsx'
import RehearsalsBadge from '../../../../../../../components/badge/RehearsalsBadge.tsx'
import SongOutlineProgressBar from '../../../../../../../components/bar/SongOutlineProgressBar.tsx'
import SongOutlineConfidenceBar from '../../../../../../../components/bar/SongOutlineConfidenceBar.tsx'
import { MouseEvent, useState } from 'react'
import { useSongContext } from '../../../context/SongContext.tsx'

interface SongSectionPartCardProps {
  part: SongPartModel
  sectionId: string
  isDragging: boolean
  maxPartProgress: number
  draggableProvided?: DraggableProvided
  showRehearsalsToast?: (name: string) => void
}

function SongSectionPartCard({
  part,
  sectionId,
  isDragging,
  maxPartProgress,
  draggableProvided,
  showRehearsalsToast
}: SongSectionPartCardProps) {
  const { ref: hoverRef, hovered } = useHover()
  const {
    ref: selectableRef,
    isClickSelected,
    isClickSelectionActive,
    isLastInSelection
  } = useClickSelectSelectable('part-' + part.id + ':' + sectionId)
  const ref = useMergedRef(hoverRef, draggableProvided?.innerRef, selectableRef)

  const { songId, isArtistBand } = useSongContext()

  const [updateSongPartMutation, { isLoading: isUpdateLoading }] = useUpdateSongPartMutation()
  const [deleteSongPartMutation, { isLoading: isDeleteLoading }] = useDeleteSongPartMutation()

  const getInstrumentIcon = useInstrumentIcon()

  const [openedContextMenu, { toggle: toggleContextMenu }] = useDisclosure(false)
  const [openedDetails, setOpenedDetails] = useState(false)

  const [openedEditSongPart, { open: openEditSongPart, close: closeEditSongPart }] =
    useDisclosure(false)
  const [openedDeleteWarning, { open: openDeleteWarning, close: closeDeleteWarning }] =
    useDisclosure(false)

  const isSelected = hovered || openedContextMenu || isDragging || isClickSelected

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

  return (
    <ContextMenu
      opened={openedContextMenu}
      onChange={toggleContextMenu}
      disabled={isClickSelectionActive}
    >
      <ContextMenu.Target>
        <Stack
          ref={ref}
          py={'6px'}
          aria-label={`song-section-part-${part.name}`}
          aria-selected={isSelected}
          gap={0}
          px={'45px'}
          onClick={handleClick}
          sx={(theme) => ({
            cursor: 'pointer',
            transition: '0.25s',
            borderRadius: 0,
            border: '1px solid transparent',
            ...(isSelected && {
              boxShadow: theme.shadows.md,
              backgroundColor: alpha(theme.colors.gray[1], 0.35)
            }),

            ...(isClickSelected && {
              boxShadow: 'none',
              backgroundColor: alpha(theme.colors.gray[1], 0.35),
              ...(hovered && {
                boxShadow: theme.shadows.xs,
                backgroundColor: alpha(theme.colors.gray[1], 0.5)
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
          <Group gap={'xs'}>
            {/*It is expected that the part will only have a band member*/}
            {isArtistBand && part.bandMembers.length > 0 && (
              <BandMemberAvatar size={23} iconSize={14} bandMember={part.bandMembers[0]} />
            )}
            {part.instrument && (
              <Tooltip openDelay={200} label={part.instrument?.name} withArrow>
                <Center aria-label={'instrument-icon'} c={'primary.7'} w={15} h={15}>
                  {getInstrumentIcon(part.instrument)}
                </Center>
              </Tooltip>
            )}

            <Text fz={'sm'} fw={500} truncate={'end'}>
              {part.name}
            </Text>
            <RehearsalsBadge rehearsals={part.rehearsals} />

            <Space flex={1} />

            <Group
              gap={2}
              style={{ transition: '0.25s', opacity: hovered || openedContextMenu ? 1 : 0 }}
            >
              <Tooltip label={'Add Rehearsal'} openDelay={200} disabled={isClickSelectionActive}>
                <ActionIcon
                  variant={'grey'}
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
            </Group>
          </Group>

          <Collapse expanded={openedDetails}>
            <Group
              aria-label={`song-section-part-details-${part.name}`}
              pb={'4px'}
              pt={'8px'}
              gap={'lg'}
            >
              <SongOutlineConfidenceBar w={'33%'} confidence={part.confidence} />
              <SongOutlineProgressBar
                w={'33%'}
                progress={part.progress}
                maxProgress={maxPartProgress}
              />
            </Group>
          </Collapse>
        </Stack>
      </ContextMenu.Target>

      <ContextMenu.Dropdown>
        <ContextMenu.Label>Part</ContextMenu.Label>
        <ContextMenu.Item leftSection={<IconEdit size={14} />} onClick={openEditSongPart}>
          Edit
        </ContextMenu.Item>
        <ContextMenu.Item
          leftSection={<IconTrash size={14} />}
          c={'red.5'}
          onClick={openDeleteWarning}
        >
          Delete
        </ContextMenu.Item>
      </ContextMenu.Dropdown>

      <EditSongPartModal
        opened={openedEditSongPart}
        onClose={closeEditSongPart}
        part={part}
        sectionId={sectionId}
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

export default SongSectionPartCard
