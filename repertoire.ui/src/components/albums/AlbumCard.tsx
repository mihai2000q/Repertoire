import Album from '../../types/models/Album.ts'
import { AspectRatio, Image, Stack, Text } from '@mantine/core'
import albumPlaceholder from '../../assets/image-placeholder-1.jpg'
import { useState } from 'react'
import { useAppDispatch } from '../../state/store.ts'
import { openArtistDrawer } from '../../state/globalSlice.ts'
import {useNavigate} from "react-router-dom";

interface AlbumCardProps {
  album: Album
}

function AlbumCard({ album }: AlbumCardProps) {
  const dispatch = useAppDispatch()
  const navigate = useNavigate()

  const [isImageHovered, setIsImageHovered] = useState(false)

  function handleClick() {
    navigate(`/album/${album.id}`)
  }

  function handleArtistClick(artistId: string) {
    dispatch(openArtistDrawer(artistId))
  }

  return (
    <Stack
      align={'center'}
      gap={0}
      style={{ transition: '0.3s', ...(isImageHovered && { transform: 'scale(1.1)' }) }}
      w={150}
    >
      <AspectRatio>
        <Image
          onMouseEnter={() => setIsImageHovered(true)}
          onMouseLeave={() => setIsImageHovered(false)}
          radius={'lg'}
          src={album.imageUrl}
          fallbackSrc={albumPlaceholder}
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
      <Stack w={'100%'} align={'center'} pt={'xs'} gap={0} style={{ overflow: 'hidden' }}>
        <Text fw={600} lineClamp={2} ta={'center'}>
          {album.title}
        </Text>
        {album.artist && (
          <Text
            fw={500}
            ta={'center'}
            c={'dimmed'}
            truncate={'end'}
            onClick={() => handleArtistClick(album.artist.id)}
            sx={{
              cursor: 'pointer',
              '&:hover': {
                textDecoration: 'underline',
              }
            }}
          >
            {album.artist.name}
          </Text>
        )}
      </Stack>
    </Stack>
  )
}

export default AlbumCard
