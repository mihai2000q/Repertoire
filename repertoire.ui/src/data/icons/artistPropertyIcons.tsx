import ArtistProperty from '../../utils/enums/ArtistProperty.ts'
import { ReactElement } from 'react'
import {
  IconAbc,
  IconCalendarMonth,
  IconCalendarWeek,
  IconRepeat,
  IconTimeline,
  IconTrendingUp,
  IconUsers
} from '@tabler/icons-react'
import CustomIconMusicNoteEighth from '../../components/@ui/icons/CustomIconMusicNoteEighth.tsx'
import CustomIconAlbumVinyl from '../../components/@ui/icons/CustomIconAlbumVinyl.tsx'

export const artistPropertyIcons = new Map<string, ReactElement>([
  [ArtistProperty.Name, <IconAbc size={'100%'} key={'name'} />],
  [ArtistProperty.Albums, <CustomIconAlbumVinyl size={'100%'} key={'albums'} />],
  [ArtistProperty.Songs, <CustomIconMusicNoteEighth size={'100%'} key={'songs'} />],
  [ArtistProperty.BandMembers, <IconUsers size={'100%'} key={'band-members'} />],
  [ArtistProperty.Rehearsals, <IconRepeat size={'100%'} key={'rehearsals'} />],
  [ArtistProperty.Confidence, <IconTimeline size={'100%'} key={'confidence'} />],
  [ArtistProperty.Progress, <IconTrendingUp size={'100%'} key={'progress'} />],
  [ArtistProperty.CreationDate, <IconCalendarWeek size={'100%'} key={'creation-date'} />],
  [ArtistProperty.LastModified, <IconCalendarMonth size={'100%'} key={'last-modified'} />]
])
