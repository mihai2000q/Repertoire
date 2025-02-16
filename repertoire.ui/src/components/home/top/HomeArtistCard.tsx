import Artist from '../../../types/models/Artist.ts'
import { Avatar, Stack, Text } from '@mantine/core'
import artistPlaceholder from '../../../assets/user-placeholder.jpg'
import { useState } from 'react'
import { openArtistDrawer } from '../../../state/slice/globalSlice.ts'
import { useAppDispatch } from '../../../state/store.ts'

interface HomeArtistCardProps {
  artist: Artist
}

function HomeArtistCard({ artist }: HomeArtistCardProps) {
  const dispatch = useAppDispatch()

  const [isImageHovered, setIsImageHovered] = useState(false)

  function handleClick() {
    dispatch(openArtistDrawer(artist.id))
  }

  return (
    <Stack
      aria-label={`artist-card-${artist.name}`}
      align={'center'}
      gap={0}
      style={{ transition: '0.25s', ...(isImageHovered && { transform: 'scale(1.05)' }) }}
      w={150}
    >
      <Avatar
        onMouseEnter={() => setIsImageHovered(true)}
        onMouseLeave={() => setIsImageHovered(false)}
        size={125}
        src={artist.imageUrl ?? artistPlaceholder}
        alt={artist.name}
        onClick={handleClick}
        sx={(theme) => ({
          cursor: 'pointer',
          transition: '0.25s',
          boxShadow: theme.shadows.xl,
          '&:hover': { boxShadow: theme.shadows.xxl }
        })}
      />

      <Stack w={'100%'} pt={'xs'} gap={0} style={{ overflow: 'hidden' }}>
        <Text fw={600} lineClamp={2} ta={'center'}>
          {artist.name}
        </Text>
      </Stack>
    </Stack>
  )
}

export default HomeArtistCard
