import Song from '../../types/models/Song.ts'
import { alpha, Avatar, Group, Stack, Text } from '@mantine/core'
import songPlaceholder from '../../assets/image-placeholder-1.jpg'
import { useAppDispatch } from '../../state/store.ts'
import { openSongDrawer } from '../../state/globalSlice.ts'

interface AlbumSongCardProps {
  song: Song
  isUnknownAlbum: boolean
}

function AlbumSongCard({ song, isUnknownAlbum }: AlbumSongCardProps) {
  const dispatch = useAppDispatch()

  function handleClick() {
    dispatch(openSongDrawer(song.id))
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
      {!isUnknownAlbum && (
        <Text fw={500} w={35} ta={'center'}>
          {song.albumTrackNo}
        </Text>
      )}
      <Avatar radius={'8px'} src={song.imageUrl ?? songPlaceholder} />
      <Stack style={{ overflow: 'hidden' }}>
        <Text fw={500} truncate={'end'}>
          {song.title}
        </Text>
      </Stack>
    </Group>
  )
}

export default AlbumSongCard
