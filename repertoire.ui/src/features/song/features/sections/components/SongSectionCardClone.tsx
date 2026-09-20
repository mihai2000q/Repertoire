import {
  Instrument,
  SongSection as SongSectionModel
} from '../../../../../types/models/Song.ts'
import { alpha, Group, Stack, Text } from '@mantine/core'
import { IconChevronDown } from '@tabler/icons-react'
import { DraggableProvided } from '@hello-pangea/dnd'
import { useMemo } from 'react'
import { BandMember } from '../../../../../types/models/Artist.ts'
import SongSectionTypeBadge from './section/SongSectionTypeBadge.tsx'
import SongOutlineConfidenceBar from '../../outline/components/SongOutlineConfidenceBar.tsx'
import SongOutlineProgressBar from '../../outline/components/SongOutlineProgressBar.tsx'
import RehearsalsBadge from '../../outline/components/RehearsalsBadge.tsx'
import BandMembersGroup from '../../outline/components/BandMembersGroup.tsx'
import InstrumentsGroup from './section/InstrumentsGroup.tsx'
import { useSongContext } from '../../../context/SongContext.tsx'

interface SongSectionCardCloneProps {
  section: SongSectionModel
  isDragging: boolean
  maxSectionProgress: number
  isDropAnimating: boolean
  draggableProvided?: DraggableProvided
}

function SongSectionCardClone({
  section,
  isDragging,
  maxSectionProgress,
  isDropAnimating,
  draggableProvided
}: SongSectionCardCloneProps) {
  const { isArtistBand } = useSongContext()

  const [bandMembers, instruments] = useMemo(() => {
    const bandMembers: BandMember[] = []
    const instruments: Instrument[] = []

    section.parts.forEach((part) => {
      part.bandMembers.forEach((bandMember) => {
        if (!bandMembers.some((m) => m.id === bandMember.id)) bandMembers.push(bandMember)
      })
      if (part.instrument && !instruments.some((i) => i.id === part.instrument.id))
        instruments.push(part.instrument)
    })

    return [bandMembers, instruments]
  }, [section.parts])

  return (
    <Stack
      ref={draggableProvided?.innerRef}
      aria-label={`song-section-${section.name}-clone`}
      aria-selected={isDragging}
      gap={0}
      py={'sm'}
      {...draggableProvided?.draggableProps}
      {...draggableProvided?.dragHandleProps}
      style={{
        ...draggableProvided?.draggableProps?.style,
        transition: `${draggableProvided?.draggableProps?.style?.transition ?? ''},
          border-color 0.25s ease,
          border-radius 0.25s ease,
          box-shadow 0.25s ease,
          background-color 0.25s ease`,
        cursor: 'default'
      }}
      sx={(theme) => ({
        borderRadius: 0,
        border: '1px solid transparent',
        boxShadow: theme.shadows.divider,
        ...(isDragging &&
          !isDropAnimating && {
            boxShadow: theme.shadows.xl,
            borderRadius: '16px',
            backgroundColor: alpha(theme.white, 0.33),
            border: `1px solid ${alpha(theme.colors.primary[8], 0.33)}`
          })
      })}
    >
      <Group gap={'sm'} pl={'sm'} pr={'md'}>
        <IconChevronDown color={'gray'} size={16} />

        <Stack gap={'xs'} flex={1}>
          <Group gap={'xs'}>
            <Text fw={600} fz={13} truncate={'end'}>
              {section.name}
            </Text>
            <SongSectionTypeBadge songSectionType={section.songSectionType} />
            <RehearsalsBadge rehearsals={section.rehearsals} />
          </Group>

          <Group gap={'md'}>
            <SongOutlineConfidenceBar w={'8vw'} size={'4px'} confidence={section.confidence} />
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
  )
}

export default SongSectionCardClone
