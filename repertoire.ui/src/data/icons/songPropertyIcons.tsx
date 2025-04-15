import SongProperty from '../../utils/enums/SongProperty.ts'
import { ReactElement } from 'react'
import {
  IconAbc,
  IconBombFilled,
  IconCalendarCheck,
  IconCalendarMonth,
  IconCalendarRepeat,
  IconCalendarWeek,
  IconRepeat,
  IconStarFilled,
  IconTimeline,
  IconTrendingUp,
  IconUser
} from '@tabler/icons-react'
import CustomIconGuitarHead from '../../components/@ui/icons/CustomIconGuitarHead.tsx'
import CustomIconAlbumVinyl from '../../components/@ui/icons/CustomIconAlbumVinyl.tsx'
import CustomIconLightningTrio from '../../components/@ui/icons/CustomIconLightningTrio.tsx'

export const songPropertyIcons = new Map<string, ReactElement>([
  [SongProperty.Title, <IconAbc size={'100%'} key={'title'} />],
  [SongProperty.ReleaseDate, <IconCalendarRepeat size={'100%'} key={'release-date'} />],
  [SongProperty.Artist, <IconUser size={'100%'} key={'artist'} />],
  [SongProperty.Album, <CustomIconAlbumVinyl size={'100%'} key={'album'} />],
  [SongProperty.GuitarTuning, <CustomIconGuitarHead size={'100%'} key={'guitar-tuning'} />],
  [SongProperty.Difficulty, <IconStarFilled size={'100%'} key={'difficulty'} />],
  [SongProperty.Solos, <CustomIconLightningTrio size={'100%'} key={'solos'} />],
  [SongProperty.Riffs, <IconBombFilled size={'100%'} key={'riffs'} />],
  [SongProperty.LastPlayed, <IconCalendarCheck size={'100%'} key={'last-played'} />],
  [SongProperty.Rehearsals, <IconRepeat size={'100%'} key={'rehearsals'} />],
  [SongProperty.Confidence, <IconTimeline size={'100%'} key={'confidence'} />],
  [SongProperty.Progress, <IconTrendingUp size={'100%'} key={'progress'} />],
  [SongProperty.CreationDate, <IconCalendarWeek size={'100%'} key={'creation-date'} />],
  [SongProperty.LastModified, <IconCalendarMonth size={'100%'} key={'last-modified'} />]
])
