import AlbumProperty from '../../utils/enums/AlbumProperty.ts'
import { ReactElement } from 'react'
import {
  IconAbc,
  IconCalendarMonth,
  IconCalendarRepeat,
  IconCalendarWeek,
  IconRepeat,
  IconTimeline,
  IconTrendingUp,
  IconUser
} from '@tabler/icons-react'
import CustomIconMusicNoteEighth from '../../components/@ui/icons/CustomIconMusicNoteEighth.tsx'

export const albumPropertyIcons = new Map<string, ReactElement>([
  [AlbumProperty.Title, <IconAbc size={'100%'} key={'title'} />],
  [AlbumProperty.ReleaseDate, <IconCalendarRepeat size={'100%'} key={'release-date'} />],
  [AlbumProperty.Artist, <IconUser size={'100%'} key={'artist'} />],
  [AlbumProperty.Songs, <CustomIconMusicNoteEighth size={'100%'} key={'songs'} />],
  [AlbumProperty.Rehearsals, <IconRepeat size={'100%'} key={'rehearsals'} />],
  [AlbumProperty.Confidence, <IconTimeline size={'100%'} key={'confidence'} />],
  [AlbumProperty.Progress, <IconTrendingUp size={'100%'} key={'progress'} />],
  [AlbumProperty.CreationDate, <IconCalendarWeek size={'100%'} key={'creation-date'} />],
  [AlbumProperty.LastModified, <IconCalendarMonth size={'100%'} key={'last-modified'} />]
])
