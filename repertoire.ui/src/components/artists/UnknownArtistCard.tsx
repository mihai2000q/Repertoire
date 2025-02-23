import { Avatar, Stack, Text } from '@mantine/core'
import unknownPlaceholder from '../../assets/unknown-placeholder.png'
import { useState } from 'react'
import { useNavigate } from 'react-router-dom'

function UnknownArtistCard() {
  const navigate = useNavigate()

  const [isAvatarHovered, setIsAvatarHovered] = useState(false)

  function handleClick() {
    navigate(`/artist/unknown`)
  }

  return (
    <Stack
      aria-label={'unknown-artist-card'}
      align={'center'}
      gap={'xs'}
      style={{
        transition: '0.25s',
        ...(isAvatarHovered && {
          transform: 'scale(1.1)'
        })
      }}
    >
      <Avatar
        onMouseEnter={() => setIsAvatarHovered(true)}
        onMouseLeave={() => setIsAvatarHovered(false)}
        src={unknownPlaceholder}
        alt={'unknown-artist'}
        w={'100%'}
        h={'unset'}
        style={(theme) => ({
          aspectRatio: 1,
          cursor: 'pointer',
          transition: '0.3s',
          boxShadow: isAvatarHovered ? theme.shadows.xxl_hover : theme.shadows.xxl
        })}
        onClick={handleClick}
      />
      <Text fw={300} ta={'center'} fs={'italic'}>
        Unknown
      </Text>
    </Stack>
  )
}

export default UnknownArtistCard
