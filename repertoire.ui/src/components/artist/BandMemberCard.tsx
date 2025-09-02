import { Avatar, Group, Stack, Text } from '@mantine/core'
import { MouseEvent } from 'react'
import { useDisclosure, useHover } from '@mantine/hooks'
import { IconPencil, IconTrash, IconUser } from '@tabler/icons-react'
import WarningModal from '../@ui/modal/WarningModal.tsx'
import { BandMember } from '../../types/models/Artist.ts'
import { useDeleteBandMemberMutation } from '../../state/api/artistsApi.ts'
import { toast } from 'react-toastify'
import BandMemberDetailsModal from './modal/BandMemberDetailsModal.tsx'
import EditBandMemberModal from './modal/EditBandMemberModal.tsx'
import { DraggableProvided } from '@hello-pangea/dnd'
import { ContextMenu } from '../@ui/menu/ContextMenu.tsx'

interface BandMemberCardProps {
  bandMember: BandMember
  artistId: string
  draggableProvided?: DraggableProvided
}

function BandMemberCard({ bandMember, artistId, draggableProvided }: BandMemberCardProps) {
  const [deleteBandMember, { isLoading: isDeleteLoading }] = useDeleteBandMemberMutation()
  const { ref, hovered } = useHover()

  const [openedMenu, { toggle: toggleMenu }] = useDisclosure(false)

  const isSelected = openedMenu || hovered

  const [openedDetails, { open: openDetails, close: closeDetails }] = useDisclosure(false)
  const [openedEdit, { open: openEdit, close: closeEdit }] = useDisclosure(false)
  const [openedDeleteWarning, { open: openDeleteWarning, close: closeDeleteWarning }] =
    useDisclosure(false)

  function handleClick() {
    openDetails()
  }

  function handleOpenEdit(e: MouseEvent) {
    e.stopPropagation()
    openEdit()
  }

  function handleOpenDeleteWarning(e: MouseEvent) {
    e.stopPropagation()
    openDeleteWarning()
  }

  async function handleDelete() {
    await deleteBandMember({ id: bandMember.id, artistId: artistId }).unwrap()
    toast.success(`${bandMember.name} deleted!`)
  }

  return (
    <Stack
      ref={draggableProvided?.innerRef}
      w={75}
      align={'center'}
      aria-label={`band-member-card-${bandMember.name}`}
      gap={'xxs'}
      sx={{ transition: '0.3s', ...(isSelected && { transform: 'scale(1.1)' }) }}
      {...draggableProvided?.draggableProps}
    >
      <ContextMenu opened={openedMenu} onChange={toggleMenu}>
        <ContextMenu.Target>
          <Avatar
            ref={ref}
            variant={'light'}
            size={'lg'}
            color={bandMember.color}
            src={bandMember.imageUrl}
            alt={bandMember.imageUrl && bandMember.name}
            sx={(theme) => ({
              transition: '0.3s',
              boxShadow: theme.shadows.md,
              '&:hover': { boxShadow: theme.shadows.lg }
            })}
            onClick={handleClick}
            {...draggableProvided?.dragHandleProps}
            style={{ cursor: 'pointer' }} // to override the drag handle
          >
            <IconUser aria-label={`icon-${bandMember.name}`} size={25} />
          </Avatar>
        </ContextMenu.Target>

        <ContextMenu.Dropdown>
          <ContextMenu.Item leftSection={<IconPencil size={14} />} onClick={handleOpenEdit}>
            Edit
          </ContextMenu.Item>
          <ContextMenu.Item
            leftSection={<IconTrash size={14} />}
            c={'red.5'}
            onClick={handleOpenDeleteWarning}
          >
            Delete
          </ContextMenu.Item>
        </ContextMenu.Dropdown>
      </ContextMenu>

      <Text ta={'center'} fw={500} lh={'xs'} lineClamp={2}>
        {bandMember.name}
      </Text>

      <EditBandMemberModal opened={openedEdit} onClose={closeEdit} bandMember={bandMember} />
      <BandMemberDetailsModal
        opened={openedDetails}
        onClose={closeDetails}
        bandMember={bandMember}
      />
      <WarningModal
        opened={openedDeleteWarning}
        onClose={closeDeleteWarning}
        title={`Delete Band Member`}
        description={
          <Group gap={'xxs'}>
            <Text>Are you sure you want to delete</Text>
            <Text fw={600}>{bandMember.name}</Text>
            <Text>?</Text>
          </Group>
        }
        onYes={handleDelete}
        isLoading={isDeleteLoading}
      />
    </Stack>
  )
}

export default BandMemberCard
