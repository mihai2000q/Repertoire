import SongProperty from '../../types/enums/properties/SongProperty.ts'
import { ReactElement } from 'react'
import {
  IconAbc,
  IconBombFilled,
  IconCalendarCheck,
  IconCalendarMonth,
  IconCalendarRepeat,
  IconCalendarWeek,
  IconList,
  IconRepeat,
  IconStarFilled,
  IconTimeline,
  IconTrendingUp,
  IconUser
} from '@tabler/icons-react'
import CustomIconGuitarHead from '../../components/@ui/icons/CustomIconGuitarHead.tsx'
import CustomIconAlbumVinyl from '../../components/@ui/icons/CustomIconAlbumVinyl.tsx'
import CustomIconLightningTrio from '../../components/@ui/icons/CustomIconLightningTrio.tsx'
import CustomIconMetronome from '../../components/@ui/icons/CustomIconMetronome.tsx'

export const songPropertyIcons = new Map<string, ReactElement>([
  [SongProperty.AlbumTitle, <CustomIconAlbumVinyl size={'100%'} key={'album'} />],
  [SongProperty.ArtistName, <IconUser size={'100%'} key={'artist'} />],
  [SongProperty.BPM, <CustomIconMetronome size={'100%'} key={'bpm'} />],
  [SongProperty.CreationDate, <IconCalendarWeek size={'100%'} key={'creation-date'} />],
  [SongProperty.Confidence, <IconTimeline size={'100%'} key={'confidence'} />],
  [SongProperty.Difficulty, <IconStarFilled size={'100%'} key={'difficulty'} />],
  [SongProperty.GuitarTuningOrder, <CustomIconGuitarHead size={'100%'} key={'guitar-tuning'} />],
  [SongProperty.LastModified, <IconCalendarMonth size={'100%'} key={'last-modified'} />],
  [SongProperty.LastPlayed, <IconCalendarCheck size={'100%'} key={'last-played'} />],
  [SongProperty.Progress, <IconTrendingUp size={'100%'} key={'progress'} />],
  [SongProperty.ReleaseDate, <IconCalendarRepeat size={'100%'} key={'release-date'} />],
  [SongProperty.Riffs, <IconBombFilled size={'100%'} key={'riffs'} />],
  [SongProperty.Rehearsals, <IconRepeat size={'100%'} key={'rehearsals'} />],
  [SongProperty.Sections, <IconList size={'100%'} key={'riffs'} />],
  [SongProperty.Solos, <CustomIconLightningTrio size={'100%'} key={'solos'} />],
  [SongProperty.Title, <IconAbc size={'100%'} key={'title'} />]
])
