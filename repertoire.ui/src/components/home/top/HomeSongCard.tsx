import Song from '../../../types/models/Song.ts'
import { AspectRatio, Image, Stack, Text } from '@mantine/core'
import songPlaceholder from '../../../assets/image-placeholder-1.jpg'
import { useState } from 'react'
import { openArtistDrawer, openSongDrawer } from '../../../state/slice/globalSlice.ts'
import { useAppDispatch } from '../../../state/store.ts'

interface HomeSongCardProps {
  song: Song
}

function HomeSongCard({ song }: HomeSongCardProps) {
  const dispatch = useAppDispatch()

  const [isImageHovered, setIsImageHovered] = useState(false)

  function handleClick() {
    dispatch(openSongDrawer(song.id))
  }

  function handleArtistClick() {
    dispatch(openArtistDrawer(song.artist.id))
  }

  return (
    <Stack
      aria-label={`song-card-${song.title}`}
      align={'center'}
      gap={0}
      style={{ transition: '0.25s', ...(isImageHovered && { transform: 'scale(1.05)' }) }}
      w={150}
    >
      <AspectRatio>
        <Image
          onMouseEnter={() => setIsImageHovered(true)}
          onMouseLeave={() => setIsImageHovered(false)}
          radius={'lg'}
          src={song.imageUrl ?? song.album?.imageUrl}
          fallbackSrc={songPlaceholder}
          alt={song.title}
          onClick={handleClick}
          sx={(theme) => ({
            cursor: 'pointer',
            transition: '0.25s',
            boxShadow: theme.shadows.xl,
            '&:hover': { boxShadow: theme.shadows.xxl }
          })}
        />
      </AspectRatio>

      <Stack w={'100%'} pt={'xs'} gap={0} style={{ overflow: 'hidden' }}>
        <Text fw={600} lineClamp={2} ta={'center'}>
          {song.title}
        </Text>
        {song.artist && (
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
            {song.artist.name}
          </Text>
        )}
      </Stack>
    </Stack>
  )
}

export default HomeSongCard
