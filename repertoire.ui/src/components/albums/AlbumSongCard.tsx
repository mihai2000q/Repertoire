import Song from '../../types/models/Song.ts'
import { alpha, Avatar, Group, Stack, Text } from '@mantine/core'
import songPlaceholder from '../../assets/image-placeholder-1.jpg'
import dayjs from 'dayjs'
import { useAppDispatch } from '../../state/store.ts'
import { openSongDrawer } from '../../state/globalSlice.ts'

interface AlbumSongCardProps {
  song: Song
}

function AlbumSongCard({ song }: AlbumSongCardProps) {
  const dispatch = useAppDispatch()

  function handleClick() {
    dispatch(openSongDrawer(song.id))
  }

  return (
    <Group
      align={'center'}
      wrap={'nowrap'}
      sx={(theme) => ({
        transition: '0.3s',
        '&:hover': {
          boxShadow: theme.shadows.xl,
          backgroundColor: alpha(theme.colors.cyan[0], 0.15)
        }
      })}
      px={'md'}
      py={'xs'}
      onClick={handleClick}
    >
      <Text fw={500} w={35} ta={'center'}>{song.albumTrackNo}</Text>
      <Avatar radius={'8px'} src={song.imageUrl ?? songPlaceholder} />
      <Stack gap={0} style={{ overflow: 'hidden' }}>
        <Group>
          <Text fw={500} truncate={'end'}>
            {song.title}
          </Text>
          {song.album && (
            <Text fz={'sm'} truncate={'end'}>
              {' '}
              - {song.album.title}
            </Text>
          )}
        </Group>
        {song.releaseDate && (
          <Text fz={'xs'} c={'dimmed'}>
            {dayjs(song.releaseDate).format('DD MMM YYYY')}
          </Text>
        )}
      </Stack>
    </Group>
  )
}

export default AlbumSongCard
