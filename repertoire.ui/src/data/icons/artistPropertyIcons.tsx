import ArtistProperty from '../../utils/enums/ArtistProperty.ts'
import { ReactElement } from 'react'
import {
  IconAbc,
  IconCalendarMonth,
  IconCalendarWeek,
  IconRepeat,
  IconTimeline,
  IconTrendingUp,
  IconUsers,
  IconUsersGroup
} from '@tabler/icons-react'
import CustomIconMusicNoteEighth from '../../components/@ui/icons/CustomIconMusicNoteEighth.tsx'
import CustomIconAlbumVinyl from '../../components/@ui/icons/CustomIconAlbumVinyl.tsx'

export const artistPropertyIcons = new Map<string, ReactElement>([
  [ArtistProperty.Albums, <CustomIconAlbumVinyl size={'100%'} key={'albums'} />],
  [ArtistProperty.Band, <IconUsers size={'100%'} key={'band'} />],
  [ArtistProperty.BandMembers, <IconUsersGroup size={'100%'} key={'band-members'} />],
  [ArtistProperty.Confidence, <IconTimeline size={'100%'} key={'confidence'} />],
  [ArtistProperty.CreationDate, <IconCalendarWeek size={'100%'} key={'creation-date'} />],
  [ArtistProperty.LastModified, <IconCalendarMonth size={'100%'} key={'last-modified'} />],
  [ArtistProperty.Name, <IconAbc size={'100%'} key={'name'} />],
  [ArtistProperty.Progress, <IconTrendingUp size={'100%'} key={'progress'} />],
  [ArtistProperty.Rehearsals, <IconRepeat size={'100%'} key={'rehearsals'} />],
  [ArtistProperty.Songs, <CustomIconMusicNoteEighth size={'100%'} key={'songs'} />],
])
