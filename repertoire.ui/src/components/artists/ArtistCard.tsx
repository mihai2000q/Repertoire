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
      style={{
        transition: '0.25s',
        ...isAvatarHovered && {
          transform: 'scale(1.1)',
        }
      }}
    >
      <Avatar
        onMouseEnter={() => setIsAvatarHovered(true)}
        onMouseLeave={() => setIsAvatarHovered(false)}
        src={artist.imageUrl ?? artistPlaceholder}
        size={125}
        style={(theme) => ({
          cursor: 'pointer',
          transition: '0.3s',
          boxShadow: isAvatarHovered ? theme.shadows.xxl_hover : theme.shadows.xxl,
        })}
        onClick={handleClick}
      />
      <Text fw={600} ta={'center'} lineClamp={2}>
        {artist.name}
      </Text>
    </Stack>
  )
}

export default ArtistCard
