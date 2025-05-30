import PlaylistProperty from '../../types/enums/PlaylistProperty.ts'
import { ReactElement } from 'react'
import { IconAbc, IconCalendarMonth, IconCalendarWeek } from '@tabler/icons-react'
import CustomIconMusicNoteEighth from '../../components/@ui/icons/CustomIconMusicNoteEighth.tsx'

export const playlistPropertyIcons = new Map<string, ReactElement>([
  [PlaylistProperty.CreationDate, <IconCalendarWeek size={'100%'} key={'creation-date'} />],
  [PlaylistProperty.LastModified, <IconCalendarMonth size={'100%'} key={'last-modified'} />],
  [PlaylistProperty.Songs, <CustomIconMusicNoteEighth size={'100%'} key={'songs'} />],
  [PlaylistProperty.Title, <IconAbc size={'100%'} key={'title'} />],
])
