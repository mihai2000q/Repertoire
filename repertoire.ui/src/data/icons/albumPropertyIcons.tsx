import AlbumProperty from '../../utils/enums/AlbumProperty.ts'
import { ReactElement } from 'react'
import {
  IconAbc,
  IconCalendarCheck,
  IconCalendarMonth,
  IconCalendarRepeat,
  IconCalendarWeek,
  IconRepeat,
  IconTimeline,
  IconTrendingUp,
  IconUser
} from '@tabler/icons-react'
import CustomIconMusicNoteEighth from '../../components/@ui/icons/CustomIconMusicNoteEighth.tsx'
import SongProperty from '../../utils/enums/SongProperty.ts'

export const albumPropertyIcons = new Map<string, ReactElement>([
  [AlbumProperty.Artist, <IconUser size={'100%'} key={'artist'} />],
  [AlbumProperty.Confidence, <IconTimeline size={'100%'} key={'confidence'} />],
  [AlbumProperty.CreationDate, <IconCalendarWeek size={'100%'} key={'creation-date'} />],
  [AlbumProperty.LastModified, <IconCalendarMonth size={'100%'} key={'last-modified'} />],
  [SongProperty.LastPlayed, <IconCalendarCheck size={'100%'} key={'last-played'} />],
  [AlbumProperty.Progress, <IconTrendingUp size={'100%'} key={'progress'} />],
  [AlbumProperty.Rehearsals, <IconRepeat size={'100%'} key={'rehearsals'} />],
  [AlbumProperty.ReleaseDate, <IconCalendarRepeat size={'100%'} key={'release-date'} />],
  [AlbumProperty.Songs, <CustomIconMusicNoteEighth size={'100%'} key={'songs'} />],
  [AlbumProperty.Title, <IconAbc size={'100%'} key={'title'} />]
])
