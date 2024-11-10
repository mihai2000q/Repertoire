import Artist from '../../types/models/Artist.ts'
import { Avatar, Stack, Text } from '@mantine/core'
import artistPlaceholder from '../../assets/user-placeholder.jpg'
import { useState } from 'react'
import {useNavigate} from "react-router-dom";

interface ArtistCardProps {
  artist: Artist
}

function ArtistCard({ artist }: ArtistCardProps) {
  const navigate = useNavigate()

  const [isAvatarHovered, setIsAvatarHovered] = useState(false)

  function handleClick() {
    navigate(`/artist/${artist.id}`)
  }

  return (
    <Stack
      align={'center'}
      gap={'xs'}
      sx={{
        ...isAvatarHovered && {
          transition: '0.25s',
          transform: 'scale(1.1)',
        }
      }}
    >
      <Avatar
        onMouseEnter={() => setIsAvatarHovered(true)}
        onMouseLeave={() => setIsAvatarHovered(false)}
        src={artist.imageUrl ? artist.imageUrl : artistPlaceholder}
        size={125}
        style={{
          cursor: 'pointer',
          transition: '0.3s',
          boxShadow: `rgba(0, 0, 0, ${isAvatarHovered ? '0.4' : '0.2'}) 0px 10px 36px 0px`,
        }}
        onClick={handleClick}
      />
      <Text fw={600} fz={'lg'}>
        {artist.name}
      </Text>
    </Stack>
  )
}

export default ArtistCard
