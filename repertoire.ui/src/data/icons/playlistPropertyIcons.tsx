import PlaylistProperty from '../../utils/enums/PlaylistProperty.ts'
import { ReactElement } from 'react'
import { IconAbc, IconCalendarMonth, IconCalendarWeek } from '@tabler/icons-react'
import CustomIconMusicNoteEighth from '../../components/@ui/icons/CustomIconMusicNoteEighth.tsx'

export const playlistPropertyIcons = new Map<string, ReactElement>([
  [PlaylistProperty.Title, <IconAbc size={'100%'} key={'title'} />],
  [PlaylistProperty.Songs, <CustomIconMusicNoteEighth size={'100%'} key={'songs'} />],
  [PlaylistProperty.CreationDate, <IconCalendarWeek size={'100%'} key={'creation-date'} />],
  [PlaylistProperty.LastModified, <IconCalendarMonth size={'100%'} key={'last-modified'} />]
])
