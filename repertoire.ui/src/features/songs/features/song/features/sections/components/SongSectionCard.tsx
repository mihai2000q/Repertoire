import { SongSection as SongSectionModel } from '../../../../../../../types/models/Song.ts'
import {
  ActionIcon,
  alpha,
  Center,
  Collapse,
  Group,
  Menu,
  NumberFormatter,
  Progress,
  Stack,
  Text,
  Tooltip
} from '@mantine/core'
import { IconCheck, IconDots, IconEdit, IconGripVertical, IconTrash } from '@tabler/icons-react'
import { DraggableProvided } from '@hello-pangea/dnd'
import { useDisclosure, useHover, useMergedRef } from '@mantine/hooks'
import { toast } from 'react-toastify'
import { useDeleteSongSectionMutation } from '../state/api/songSectionsApi.ts'
import EditSongSectionModal from './modal/EditSongSectionModal.tsx'
import WarningModal from '../../../../../../../components/modal/WarningModal.tsx'
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

interface SongSectionCardProps {
  section: SongSectionModel
  songId: string
  isDragging: boolean
  showDetails: boolean
  maxSectionProgress: number
  maxSectionRehearsals: number
  draggableProvided?: DraggableProvided
}

function SongSectionCard({
  section,
  songId,
  isDragging,
  showDetails,
  maxSectionProgress,
  maxSectionRehearsals,
  draggableProvided
}: SongSectionCardProps) {
  const { ref: hoverRef, hovered } = useHover()
  const {
    ref: selectableRef,
    isClickSelected,
    isClickSelectionActive,
    isLastInSelection
  } = useClickSelectSelectable(section.id)
  const ref = useMergedRef(hoverRef, draggableProvided?.innerRef, selectableRef)

  const [rehearsalsMarginLeft, setRehearsalsMarginLeft] = useState(
    getRehearsalsMarginLeft(maxSectionRehearsals.toString().length)
  )
  const [rehearsalsWidth, setRehearsalsWidth] = useState(
    getRehearsalsWidth(maxSectionRehearsals.toString().length)
  )
  useEffect(() => {
    const rehearsalsMaxLength = maxSectionRehearsals.toString().length
    setRehearsalsMarginLeft(getRehearsalsMarginLeft(rehearsalsMaxLength))
    setRehearsalsWidth(getRehearsalsWidth(rehearsalsMaxLength))
  }, [maxSectionRehearsals])

  const [deleteSongSectionMutation, { isLoading: isDeleteLoading }] = useDeleteSongSectionMutation()

  const { openedMenu, toggleMenu, openedContextMenu, toggleContextMenu } = useDoubleMenu()

  const [openedEditSongSection, { open: openEditSongSection, close: closeEditSongSection }] =
    useDisclosure(false)
  const [openedDeleteWarning, { open: openDeleteWarning, close: closeDeleteWarning }] =
    useDisclosure(false)

  const isSelected = hovered || openedMenu || openedContextMenu || isDragging || isClickSelected

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
    <ContextMenu
      opened={openedContextMenu}
      onChange={toggleContextMenu}
      disabled={isClickSelectionActive}
    >
      <ContextMenu.Target>
        <Stack
          ref={ref}
          py={'xs'}
          aria-label={`song-section-${section.name}`}
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

            <Text fw={600}>{section.songSectionType.name}</Text>
            <Text flex={1} truncate={'end'}>
              {section.name}
            </Text>

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

          <Collapse expanded={showDetails}>
            <Group
              aria-label={`song-section-details-${section.name}`}
              pt={'md'}
              gap={'lg'}
              pr={'lg'}
            >
              <Tooltip.Floating
                role={'tooltip'}
                label={
                  <>
                    Rehearsals: <NumberFormatter value={section.rehearsals} />
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
      </ContextMenu.Target>

      <ContextMenu.Dropdown>{menuDropdown}</ContextMenu.Dropdown>

      <EditSongSectionModal
        opened={openedEditSongSection}
        onClose={closeEditSongSection}
        section={section}
      />
      <WarningModal
        opened={openedDeleteWarning}
        onClose={closeDeleteWarning}
        title={`Delete Section`}
        description={
          <Group gap={'xxs'}>
            <Text>Are you sure you want to delete</Text>
            <Text fw={600}>{section.name}</Text>
            <Text>?</Text>
          </Group>
        }
        onYes={handleDelete}
        isLoading={isDeleteLoading}
      />
    </ContextMenu>
  )
}

export default SongSectionCard
