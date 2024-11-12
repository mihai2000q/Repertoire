import Song from '../../../types/models/Song.ts'
import { alpha, Avatar, Group, Space, Stack, Text } from '@mantine/core'
import songPlaceholder from '../../../assets/image-placeholder-1.jpg'
import dayjs from 'dayjs'
import { useAppDispatch } from '../../../state/store.ts'
import { openSongDrawer } from '../../../state/globalSlice.ts'

interface ArtistSongCardProps {
  song: Song
}

function ArtistSongCard({ song }: ArtistSongCardProps) {
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
      <Avatar radius={'8px'} src={song.imageUrl ? song.imageUrl : songPlaceholder} />
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
      <Space flex={1} />
    </Group>
  )
}

export default ArtistSongCard
