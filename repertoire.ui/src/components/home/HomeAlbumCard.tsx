import Album from '../../types/models/Album.ts'
import { AspectRatio, Image, Stack, Text } from '@mantine/core'
import albumPlaceholder from '../../assets/image-placeholder-1.jpg'
import { useState } from 'react'
import { openAlbumDrawer, openArtistDrawer } from '../../state/slice/globalSlice.ts'
import { useAppDispatch } from '../../state/store.ts'

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
      style={{ transition: '0.3s', ...(isImageHovered && { transform: 'scale(1.05)' }) }}
      w={150}
    >
      <AspectRatio>
        <Image
          onMouseEnter={() => setIsImageHovered(true)}
          onMouseLeave={() => setIsImageHovered(false)}
          radius={'lg'}
          src={album.imageUrl}
          fallbackSrc={albumPlaceholder}
          alt={album.title}
          onClick={handleClick}
          sx={(theme) => ({
            cursor: 'pointer',
            transition: '0.3s',
            boxShadow: theme.shadows.xxl,
            '&:hover': {
              boxShadow: theme.shadows.xxl_hover
            }
          })}
        />
      </AspectRatio>

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
