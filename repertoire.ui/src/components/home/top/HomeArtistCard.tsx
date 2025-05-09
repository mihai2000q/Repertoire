import Artist from '../../../types/models/Artist.ts'
import { Avatar, Center, Stack, Text } from '@mantine/core'
import { useState } from 'react'
import { openArtistDrawer } from '../../../state/slice/globalSlice.ts'
import { useAppDispatch } from '../../../state/store.ts'
import CustomIconUserAlt from '../../@ui/icons/CustomIconUserAlt.tsx'

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
      w={'max(9vw, 150px)'}
    >
      <Avatar
        onMouseEnter={() => setIsImageHovered(true)}
        onMouseLeave={() => setIsImageHovered(false)}
        size={'max(calc(9vw - 25px), 125px)'}
        src={artist.imageUrl}
        alt={artist.imageUrl && artist.name}
        bg={'gray.0'}
        onClick={handleClick}
        sx={(theme) => ({
          cursor: 'pointer',
          transition: '0.25s',
          boxShadow: theme.shadows.xl,
          '&:hover': { boxShadow: theme.shadows.xxl }
        })}
      >
        <Center c={'gray.7'}>
          <CustomIconUserAlt size={58} />
        </Center>
      </Avatar>

      <Stack w={'100%'} pt={'xs'} gap={0} style={{ overflow: 'hidden' }}>
        <Text fw={600} lineClamp={2} ta={'center'}>
          {artist.name}
        </Text>
      </Stack>
    </Stack>
  )
}

export default HomeArtistCard
