import Album from '../../../types/models/Album.ts'
import { alpha, Avatar, Group, Stack, Text } from '@mantine/core'
import albumPlaceholder from '../../../assets/image-placeholder-1.jpg'
import dayjs from 'dayjs'
import { useAppDispatch } from '../../../state/store.ts'
import { openAlbumDrawer } from '../../../state/globalSlice.ts'

interface ArtistAlbumCardProps {
  album: Album
}

function ArtistAlbumCard({ album }: ArtistAlbumCardProps) {
  const dispatch = useAppDispatch()

  function handleClick() {
    dispatch(openAlbumDrawer(album.id))
  }

  return (
    <Group
      align={'center'}
      wrap={'nowrap'}
      sx={(theme) => ({
        borderRadius: '8px',
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
      <Avatar radius={'8px'} src={album.imageUrl ?? albumPlaceholder} />
      <Stack gap={0} style={{ overflow: 'hidden' }}>
        <Text fw={500} truncate={'end'}>
          {album.title}
        </Text>
        {album.releaseDate && (
          <Text fz={'xs'} c={'dimmed'}>
            {dayjs(album.releaseDate).format('DD MMM YYYY')}
          </Text>
        )}
      </Stack>
    </Group>
  )
}

export default ArtistAlbumCard
