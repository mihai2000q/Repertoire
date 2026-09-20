import {
  Instrument,
  SongSection as SongSectionModel
} from '../../../../../../../types/models/Song.ts'
import { alpha, Box, Center, Collapse, Group, Stack, Text } from '@mantine/core'
import { IconCheck, IconChevronDown, IconEdit, IconRefresh, IconTrash } from '@tabler/icons-react'
import { DraggableProvided } from '@hello-pangea/dnd'
import { useDisclosure, useHover, useMergedRef } from '@mantine/hooks'
import EditSongSectionModal from './modal/EditSongSectionModal.tsx'
import { MouseEvent, useEffect, useMemo, useRef, useState } from 'react'
import { ContextMenu } from '../../../../../../../components/menu/ContextMenu.tsx'
import useClickSelectSelectable from '../../../../../../../hooks/useClickSelectSelectable.ts'
import SongSectionPartCard from './section/SongSectionPartCard.tsx'
import { BandMember } from '../../../../../../../types/models/Artist.ts'
import SongSectionTypeBadge from './section/SongSectionTypeBadge.tsx'
import SongOutlineConfidenceBar from '../../outline/components/SongOutlineConfidenceBar.tsx'
import SongOutlineProgressBar from '../../outline/components/SongOutlineProgressBar.tsx'
import RehearsalsBadge from '../../outline/components/RehearsalsBadge.tsx'
import BandMembersGroup from '../../outline/components/BandMembersGroup.tsx'
import InstrumentsGroup from './section/InstrumentsGroup.tsx'
import { useSongContext } from '../../../context/SongContext.tsx'
import DeleteSongSectionModal from './modal/DeleteSongSectionModal.tsx'
import { toast } from 'react-toastify'
import { useBulkUpdateSongPartsMutation } from '../../parts/state/api/songPartsApi.ts'
import plural from '../../../../../../../utils/plural.ts'
import MenuItemConfirmation from '../../../../../../../components/menu/item/MenuItemConfirmation.tsx'
import { useSongOutlineContext } from '../../outline/context/SongOutlineContext.tsx'
import AddNewSongSectionPart from './section/AddNewSongSectionPart.tsx'

interface SongSectionCardProps {
  section: SongSectionModel
  maxSectionProgress: number
  isDragging: boolean
  isDragStarting?: boolean
  draggableProvided?: DraggableProvided
  scrollIntoView?: () => void
}

function SongSectionCard({
  section,
  maxSectionProgress,
  isDragging,
  isDragStarting,
  draggableProvided,
  scrollIntoView
}: SongSectionCardProps) {
  const {
    ref: selectableRef,
    isClickSelected,
    isClickSelectionActive,
    isLastInSelection
  } = useClickSelectSelectable('section-' + section.id)
  const ref = useMergedRef(draggableProvided?.innerRef, selectableRef)

  const { ref: hoverRef, hovered } = useHover()
  const sectionRef = useMergedRef(hoverRef)

  const { songId, isArtistBand } = useSongContext()
  const { showDetails } = useSongOutlineContext()

  const [updateSongParts, { isLoading: isUpdateSongPartsLoading }] =
    useBulkUpdateSongPartsMutation()

  const [openedContextMenu, { toggle: toggleContextMenu }] = useDisclosure(false)
  const [openedDetails, setOpenedDetails] = useState(false)
  useEffect(() => setOpenedDetails(showDetails), [showDetails])

  const [openedEdit, { open: openEdit, close: closeEdit }] = useDisclosure(false)
  const [openedDeleteWarning, { open: openDeleteWarning, close: closeDeleteWarning }] =
    useDisclosure(false)

  const isSelected = hovered || openedContextMenu || isDragging || isClickSelected

  const rehearsalsToastId = useRef<number | string>(null)

  function showRehearsalsToast(partName: string) {
    if (rehearsalsToastId.current) toast.dismiss(rehearsalsToastId.current)
    rehearsalsToastId.current = toast.info(`${partName} rehearsals increased by 1!`)
  }

  const [maxPartProgress, bandMembers, instruments] = useMemo(() => {
    let progress = 0
    const bandMembers: BandMember[] = []
    const instruments: Instrument[] = []

    section.parts.forEach((part) => {
      if (part.progress > progress) progress = section.progress
      part.bandMembers.forEach((bandMember) => {
        if (!bandMembers.some((m) => m.id === bandMember.id)) bandMembers.push(bandMember)
      })
      if (part.instrument && !instruments.some((i) => i.id === part.instrument.id))
        instruments.push(part.instrument)
    })

    return [progress, bandMembers, instruments]
  }, [section.parts])

  function handleClick(e: MouseEvent) {
    if (e.ctrlKey || e.shiftKey) return
    e.stopPropagation()
    setOpenedDetails(!openedDetails)
  }

  async function handleAddRehearsal() {
    await updateSongParts({
      songId: songId,
      requests: section.parts.map((part) => ({
        id: part.id,
        confidence: part.confidence,
        rehearsals: part.rehearsals + 1
      }))
    }).unwrap()
    toast.success(
      `Rehearsals added to ${section.name}'s ${section.parts.length} ` +
        `part${plural(section.parts)}!`
    )
  }

  return (
    <Stack
      ref={ref}
      gap={0}
      {...draggableProvided?.draggableProps}
      {...draggableProvided?.dragHandleProps}
      style={{
        ...draggableProvided?.draggableProps?.style,
        cursor: 'default'
      }}
      sx={(theme) => ({
        transition: '0.25s',
        borderRadius: 0,
        border: '1px solid transparent',
        boxShadow: theme.shadows.divider,
        ...(isDragging && {
          boxShadow: theme.shadows.xl,
          borderRadius: '16px',
          backgroundColor: alpha(theme.white, 0.33),
          border: `1px solid ${alpha(theme.colors.primary[8], 0.33)}`
        })
      })}
    >
      <ContextMenu
        opened={openedContextMenu}
        onChange={toggleContextMenu}
        disabled={isClickSelectionActive}
      >
        <ContextMenu.Target>
          <Stack
            ref={sectionRef}
            aria-label={`song-section-${section.name}`}
            aria-selected={isSelected}
            gap={0}
            py={'sm'}
            onClick={handleClick}
            sx={(theme) => ({
              cursor: 'pointer',
              transition: '0.25s',
              borderRadius: 0,
              ...(isSelected && {
                boxShadow: theme.shadows.md,
                backgroundColor: theme.colors.gray[0]
              }),

              ...(isClickSelected && {
                boxShadow: 'none',
                backgroundColor: theme.colors.gray[0],
                ...(hovered && {
                  boxShadow: theme.shadows.xs,
                  backgroundColor: alpha(theme.colors.gray[1], 0.5)
                }),
                ...(isLastInSelection && {
                  boxShadow: theme.shadows.lg
                })
              }),

              ...(isDragging && {
                transition: '0s',
                boxShadow: 'none',
                borderRadius: '16px'
              })
            })}
          >
            <Group gap={'sm'} pl={'sm'} pr={'md'}>
              {!isClickSelected ? (
                <IconChevronDown
                  color={'gray'}
                  size={16}
                  style={{
                    transition: 'transform 200ms ease',
                    transform: openedDetails ? 'rotate(180deg)' : 'rotate(0deg)'
                  }}
                />
              ) : (
                <Center
                  data-testid={'selected-checkmark'}
                  w={16}
                  h={16}
                  style={(theme) => ({
                    borderRadius: '100%',
                    backgroundColor: alpha(theme.colors.green[2], 0.95)
                  })}
                >
                  <IconCheck color={'white'} size={'75%'} />
                </Center>
              )}

              <Stack gap={'xs'} flex={1}>
                <Group gap={'xs'}>
                  <Text fw={600} fz={13} truncate={'end'}>
                    {section.name}
                  </Text>
                  <SongSectionTypeBadge songSectionType={section.songSectionType} />
                  <RehearsalsBadge rehearsals={section.rehearsals} />
                </Group>

                <Group gap={'md'}>
                  <SongOutlineConfidenceBar
                    w={'8vw'}
                    size={'4px'}
                    confidence={section.confidence}
                  />
                  <SongOutlineProgressBar
                    w={'8vw'}
                    size={'4px'}
                    progress={section.progress}
                    maxProgress={maxSectionProgress}
                  />
                </Group>
              </Stack>

              <Group gap={'xxs'}>
                {instruments.length > 0 && (
                  <InstrumentsGroup aria-label={'instruments'} instruments={instruments} />
                )}
                {isArtistBand && bandMembers.length > 0 && (
                  <BandMembersGroup aria-label={'band-members'} bandMembers={bandMembers} />
                )}
              </Group>
            </Group>
          </Stack>
        </ContextMenu.Target>

        <ContextMenu.Dropdown>
          <ContextMenu.Label>Section</ContextMenu.Label>
          <ContextMenu.Item leftSection={<IconEdit size={14} />} onClick={openEdit}>
            Edit
          </ContextMenu.Item>
          <MenuItemConfirmation
            isLoading={isUpdateSongPartsLoading}
            onConfirm={handleAddRehearsal}
            leftSection={<IconRefresh size={14} />}
            disabled={section.parts.length === 0}
          >
            Add Rehearsal
          </MenuItemConfirmation>
          <ContextMenu.Divider />

          <ContextMenu.Item
            leftSection={<IconTrash size={14} />}
            c={'red.5'}
            onClick={openDeleteWarning}
          >
            Delete
          </ContextMenu.Item>
        </ContextMenu.Dropdown>

        <EditSongSectionModal opened={openedEdit} onClose={closeEdit} section={section} />
        <DeleteSongSectionModal
          opened={openedDeleteWarning}
          onClose={closeDeleteWarning}
          section={section}
          songId={songId}
        />
      </ContextMenu>

      {!isDragStarting && (
        <Collapse expanded={openedDetails} onTransitionEnd={scrollIntoView}>
          <Box pos={'relative'}>
            <Box
              pos={'absolute'}
              top={0}
              left={35}
              bottom={0}
              my={'4px'}
              style={(theme) => ({
                borderLeft: `2px solid ${theme.colors.gray[1]}`
              })}
            />
            <Stack gap={0} pl={35}>
              {section.parts.map((part) => (
                <SongSectionPartCard
                  key={part.id}
                  part={part}
                  sectionId={section.id}
                  maxPartProgress={maxPartProgress}
                  isDragging={false}
                  showRehearsalsToast={showRehearsalsToast}
                />
              ))}
              <AddNewSongSectionPart section={section} />
            </Stack>
          </Box>
        </Collapse>
      )}
    </Stack>
  )
}

export default SongSectionCard
