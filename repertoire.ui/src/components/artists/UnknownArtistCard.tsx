import { Avatar, Center, Stack, Text } from '@mantine/core'
import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { IconQuestionMark } from '@tabler/icons-react'

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
        w={'100%'}
        h={'unset'}
        style={(theme) => ({
          aspectRatio: 1,
          cursor: 'pointer',
          transition: '0.3s',
          boxShadow: isAvatarHovered ? theme.shadows.xxl_hover : theme.shadows.xxl
        })}
        onClick={handleClick}
      >
        <Center c={'gray.6'}>
          <IconQuestionMark
            aria-label={'icon-unknown-artist'}
            size={'100%'}
            strokeWidth={3}
            style={{ padding: '12%' }}
          />
        </Center>
      </Avatar>
      <Text fw={300} ta={'center'} fs={'italic'}>
        Unknown
      </Text>
    </Stack>
  )
}

export default UnknownArtistCard
