import Song from '../../../types/models/Song.ts'
import { Avatar, Center, Stack, Text } from '@mantine/core'
import { useState } from 'react'
import { openArtistDrawer, openSongDrawer } from '../../../state/slice/globalSlice.ts'
import { useAppDispatch } from '../../../state/store.ts'
import CustomIconMusicNote from '../../@ui/icons/CustomIconMusicNote.tsx'

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
      w={'max(10vw, 150px)'}
    >
      <Avatar
        onMouseEnter={() => setIsImageHovered(true)}
        onMouseLeave={() => setIsImageHovered(false)}
        radius={'10%'}
        w={'100%'}
        h={'unset'}
        src={song.imageUrl ?? song.album?.imageUrl}
        alt={(song.imageUrl ?? song.album?.imageUrl) && song.title}
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
          <CustomIconMusicNote aria-label={`default-icon-${song.title}`} size={50} />
        </Center>
      </Avatar>

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
