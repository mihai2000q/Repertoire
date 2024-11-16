import Song from '../../types/models/Song.ts'
import { alpha, Avatar, Group, Space, Stack, Text } from '@mantine/core'
import songPlaceholder from '../../assets/image-placeholder-1.jpg'
import dayjs from 'dayjs'
import { useAppDispatch } from '../../state/store.ts'
import { openAlbumDrawer, openSongDrawer } from '../../state/globalSlice.ts'
import { MouseEvent } from 'react'

interface ArtistSongCardProps {
  song: Song
}

function ArtistSongCard({ song }: ArtistSongCardProps) {
  const dispatch = useAppDispatch()

  function handleClick() {
    dispatch(openSongDrawer(song.id))
  }

  function handleAlbumClick(e: MouseEvent) {
    e.stopPropagation()
    dispatch(openAlbumDrawer(song.album.id))
  }

  return (
    <Group
      align={'center'}
      wrap={'nowrap'}
      sx={(theme) => ({
        cursor: 'default',
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
      <Avatar radius={'8px'} src={song.imageUrl ?? songPlaceholder} />
      <Stack gap={0} style={{ overflow: 'hidden' }}>
        <Group gap={4}>
          <Text fw={500} truncate={'end'}>
            {song.title}
          </Text>
          {song.album && (
            <>
              <Text fz={'sm'}>-</Text>
              <Text
                fz={'sm'}
                c={'dimmed'}
                truncate={'end'}
                sx={{ '&:hover': { textDecoration: 'underline' } }}
                onClick={handleAlbumClick}
              >
                {song.album.title}
              </Text>
            </>
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
