import SongProperty from '../../utils/enums/SongProperty.ts'
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
import CustomIconGuitarHead from '../../components/@ui/icons/CustomIconGuitarHead.tsx'

export const songPropertyIcons = new Map<string, ReactElement>([
  [SongProperty.Title, <IconAbc size={'100%'} key={'title'} />],
  [SongProperty.ReleaseDate, <IconCalendarRepeat size={'100%'} key={'release-date'} />],
  [SongProperty.Artist, <IconUser size={'100%'} key={'artist'} />],
  [SongProperty.GuitarTuning, <CustomIconGuitarHead size={'100%'} key={'guitar-tuning'} />],
  [SongProperty.LastPlayed, <IconCalendarCheck size={'100%'} key={'last-played'} />],
  [SongProperty.Rehearsals, <IconRepeat size={'100%'} key={'rehearsals'} />],
  [SongProperty.Confidence, <IconTimeline size={'100%'} key={'confidence'} />],
  [SongProperty.Progress, <IconTrendingUp size={'100%'} key={'progress'} />],
  [SongProperty.CreationDate, <IconCalendarWeek size={'100%'} key={'creation-date'} />],
  [SongProperty.LastModified, <IconCalendarMonth size={'100%'} key={'last-modified'} />]
])
