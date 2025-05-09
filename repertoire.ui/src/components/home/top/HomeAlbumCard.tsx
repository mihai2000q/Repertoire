import Album from '../../../types/models/Album.ts'
import { Avatar, Center, Stack, Text } from '@mantine/core'
import { useState } from 'react'
import { openAlbumDrawer, openArtistDrawer } from '../../../state/slice/globalSlice.ts'
import { useAppDispatch } from '../../../state/store.ts'
import CustomIconAlbumVinyl from '../../@ui/icons/CustomIconAlbumVinyl.tsx'

interface HomeAlbumCardProps {
  album: Album
}

function HomeAlbumCard({ album }: HomeAlbumCardProps) {
  const dispatch = useAppDispatch()

  const [isImageHovered, setIsImageHovered] = useState(false)

  function handleClick() {
    dispatch(openAlbumDrawer(album.id))
  }

  function handleArtistClick() {
    dispatch(openArtistDrawer(album.artist.id))
  }

  return (
    <Stack
      aria-label={`album-card-${album.title}`}
      align={'center'}
      gap={0}
      style={{ transition: '0.25s', ...(isImageHovered && { transform: 'scale(1.05)' }) }}
      w={'max(10vw, 150px)'}
    >
      <Avatar
        onMouseEnter={() => setIsImageHovered(true)}
        onMouseLeave={() => setIsImageHovered(false)}
        radius={'10%'}
        w={'100%'}
        h={'unset'}
        src={album.imageUrl}
        alt={album.imageUrl && album.title}
        bg={'gray.5'}
        onClick={handleClick}
        sx={(theme) => ({
          aspectRatio: 1,
          cursor: 'pointer',
          transition: '0.25s',
          boxShadow: theme.shadows.xl,
          '&:hover': { boxShadow: theme.shadows.xxl }
        })}
      >
        <Center c={'white'}>
          <CustomIconAlbumVinyl aria-label={`default-icon-${album.title}`} size={40} />
        </Center>
      </Avatar>

      <Stack w={'100%'} pt={'xs'} gap={0} style={{ overflow: 'hidden' }}>
        <Text fw={600} lineClamp={2} ta={'center'}>
          {album.title}
        </Text>
        {album.artist && (
          <Text
            fw={500}
            ta={'center'}
            c={'dimmed'}
            truncate={'end'}
            onClick={handleArtistClick}
            sx={{
              cursor: 'pointer',
              '&:hover': { textDecoration: 'underline' }
            }}
          >
            {album.artist.name}
          </Text>
        )}
      </Stack>
    </Stack>
  )
}

export default HomeAlbumCard
