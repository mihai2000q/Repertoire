import { SongSection as SongSectionModel } from '../../types/models/Song.ts'
import {
  ActionIcon,
  alpha,
  Avatar,
  Box,
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
  IconDots,
  IconEdit,
  IconGripVertical,
  IconLocationPlus,
  IconTrash,
  IconUser
} from '@tabler/icons-react'
import { DraggableProvided } from '@hello-pangea/dnd'
import { useDisclosure } from '@mantine/hooks'
import { toast } from 'react-toastify'
import {
  useDeleteSongSectionMutation,
  useUpdateSongSectionMutation
} from '../../state/api/songsApi.ts'
import EditSongSectionModal from './modal/EditSongSectionModal.tsx'
import WarningModal from '../@ui/modal/WarningModal.tsx'
import useContextMenu from '../../hooks/useContextMenu.ts'
import { BandMember } from '../../types/models/Artist.ts'
import useInstrumentIcon from '../../hooks/useInstrumentIcon.tsx'

interface SongSectionCardProps {
  section: SongSectionModel
  songId: string
  isDragging: boolean
  showDetails: boolean
  maxSectionProgress: number
  draggableProvided?: DraggableProvided
  bandMembers?: BandMember[]
}

function SongSectionCard({
  section,
  songId,
  isDragging,
  showDetails,
  maxSectionProgress,
  draggableProvided,
  bandMembers
}: SongSectionCardProps) {
  const [updateSongSectionMutation, { isLoading: isUpdateLoading }] = useUpdateSongSectionMutation()
  const [deleteSongSectionMutation, { isLoading: isDeleteLoading }] = useDeleteSongSectionMutation()

  const getInstrumentIcon = useInstrumentIcon()

  const [openedMenu, menuDropdownProps, { openMenu, closeMenu }] = useContextMenu()

  const [openedEditSongSection, { open: openEditSongSection, close: closeEditSongSection }] =
    useDisclosure(false)
  const [openedDeleteWarning, { open: openDeleteWarning, close: closeDeleteWarning }] =
    useDisclosure(false)

  function handleAddRehearsal() {
    updateSongSectionMutation({
      ...section,
      typeId: section.songSectionType.id,
      rehearsals: section.rehearsals + 1
    })
  }

  async function handleDelete() {
    await deleteSongSectionMutation({ id: section.id, songId: songId }).unwrap()
    toast.success(`${section.name} deleted!`)
  }

  const menuDropdown = (
    <>
      <Menu.Item leftSection={<IconEdit size={14} />} onClick={openEditSongSection}>
        Edit
      </Menu.Item>
      <Menu.Item leftSection={<IconTrash size={14} />} c={'red.5'} onClick={openDeleteWarning}>
        Delete
      </Menu.Item>
    </>
  )

  return (
    <Menu shadow={'lg'} opened={openedMenu} onClose={closeMenu}>
      <Menu.Target>
        <Stack
          py={'xs'}
          px={'md'}
          aria-label={`song-section-${section.name}`}
          sx={(theme) => ({
            cursor: 'default',
            transition: '0.25s',
            borderRadius: isDragging ? '16px' : '0px',
            border: isDragging
              ? `1px solid ${alpha(theme.colors.primary[9], 0.33)}`
              : '1px solid transparent',

            '&:hover': {
              boxShadow: theme.shadows.xl,
              backgroundColor: alpha(theme.colors.primary[0], 0.15)
            }
          })}
          ref={draggableProvided?.innerRef}
          {...draggableProvided?.draggableProps}
          onContextMenu={openMenu}
        >
          <Group gap={'xs'}>
            <ActionIcon
              aria-label={'drag-handle'}
              variant={'subtle'}
              size={'lg'}
              {...draggableProvided?.dragHandleProps}
            >
              <IconGripVertical size={20} />
            </ActionIcon>

            {section.bandMember && (
              <HoverCard withArrow={true} openDelay={200} position="top" shadow={'md'}>
                <HoverCard.Target>
                  <Avatar
                    size={25}
                    color={section.bandMember.color}
                    src={section.bandMember.imageUrl ?? section.bandMember.color}
                    alt={section.bandMember.name}
                  >
                    <IconUser size={15} />
                  </Avatar>
                </HoverCard.Target>
                <HoverCard.Dropdown>
                  <Group gap={'xs'} maw={200} wrap={'nowrap'}>
                    <Avatar
                      size={60}
                      color={section.bandMember.color}
                      src={section.bandMember.imageUrl ?? section.bandMember.color}
                      alt={section.bandMember.name}
                      style={(theme) => ({ boxShadow: theme.shadows.sm })}
                    >
                      <IconUser size={30} />
                    </Avatar>
                    <Stack gap={0}>
                      <Text fw={500} lineClamp={2}>
                        {section.bandMember.name}
                      </Text>
                      {section.bandMember.roles.slice(0, 2).map((role, index) => (
                        <Text key={role.id} c={'dimmed'} fz={'xs'} lineClamp={1} lh={1.05}>
                          {role.name}
                          {index === 1 && section.bandMember.roles.length > 2 && ' ...'}
                        </Text>
                      ))}
                    </Stack>
                  </Group>
                </HoverCard.Dropdown>
              </HoverCard>
            )}

            {section.instrument && (
              <Box aria-label={'instrument-icon'} c={'primary.7'} w={16} h={16}>
                <Tooltip openDelay={200} label={section.instrument?.name} withArrow>
                  {getInstrumentIcon(section.instrument)}
                </Tooltip>
              </Box>
            )}

            <Text inline fw={600}>
              {section.songSectionType.name}
            </Text>
            <Text flex={1} inline truncate={'end'}>
              {section.name}
            </Text>

            <Group gap={2}>
              <Tooltip label={'Add Rehearsal'} openDelay={200}>
                <ActionIcon
                  variant={'subtle'}
                  size={'md'}
                  disabled={isUpdateLoading}
                  aria-label={'add-rehearsal'}
                  onClick={handleAddRehearsal}
                >
                  <IconLocationPlus size={15} />
                </ActionIcon>
              </Tooltip>

              <Menu>
                <Menu.Target>
                  <ActionIcon variant={'subtle'} size={'lg'} aria-label={'more-menu'}>
                    <IconDots size={20} />
                  </ActionIcon>
                </Menu.Target>
                <Menu.Dropdown>{menuDropdown}</Menu.Dropdown>
              </Menu>
            </Group>
          </Group>

          <Collapse in={showDetails}>
            <Group aria-label={`song-section-details-${section.name}`} gap={'xl'} px={'md'}>
              <Tooltip.Floating
                role={'tooltip'}
                label={
                  <>
                    Rehearsals: <NumberFormatter value={section.rehearsals} />
                  </>
                }
              >
                <Text fw={500} c={'dimmed'} fz={'md'} inline>
                  <NumberFormatter value={section.rehearsals} />
                </Text>
              </Tooltip.Floating>

              <Tooltip.Floating role={'tooltip'} label={`Confidence: ${section.confidence}%`}>
                <Progress
                  flex={1}
                  size={'sm'}
                  value={section.confidence}
                  aria-label={'confidence'}
                />
              </Tooltip.Floating>

              <Tooltip.Floating
                role={'tooltip'}
                label={
                  <>
                    Progress: <NumberFormatter value={section.progress} />
                  </>
                }
              >
                <Progress
                  flex={1}
                  size={'sm'}
                  aria-label={'progress'}
                  value={section.progress === 0 ? 0 : (section.progress / maxSectionProgress) * 100}
                  color={'green'}
                />
              </Tooltip.Floating>
            </Group>
          </Collapse>
        </Stack>
      </Menu.Target>

      <Menu.Dropdown {...menuDropdownProps}>{menuDropdown}</Menu.Dropdown>

      <EditSongSectionModal
        opened={openedEditSongSection}
        onClose={closeEditSongSection}
        section={section}
        bandMembers={bandMembers}
      />
      <WarningModal
        opened={openedDeleteWarning}
        onClose={closeDeleteWarning}
        title={`Delete Section`}
        description={
          <Group gap={4}>
            <Text>Are you sure you want to delete</Text>
            <Text fw={600}>{section.name}</Text>
            <Text>?</Text>
          </Group>
        }
        onYes={handleDelete}
        isLoading={isDeleteLoading}
      />
    </Menu>
  )
}

export default SongSectionCard
