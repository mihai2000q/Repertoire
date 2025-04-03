import {Avatar, Center, Stack, Text} from '@mantine/core'
import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { IconQuestionMark } from '@tabler/icons-react'

function UnknownArtistCard() {
  const navigate = useNavigate()

  const [isAvatarHovered, setIsAvatarHovered] = useState(false)

  function handleClick() {
    navigate(`/album/unknown`)
  }

  return (
    <Stack
      aria-label={'unknown-album-card'}
      align={'center'}
      gap={0}
      style={{
        alignSelf: 'start',
        transition: '0.3s',
        ...(isAvatarHovered && { transform: 'scale(1.1)' })
      }}
    >
      <Avatar
        onMouseEnter={() => setIsAvatarHovered(true)}
        onMouseLeave={() => setIsAvatarHovered(false)}
        radius={'10%'}
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
            aria-label={'unknown-album'}
            size={'100%'}
            strokeWidth={3}
            style={{ padding: '15%' }}
          />
        </Center>
      </Avatar>

      <Stack pt={'xs'}>
        <Text fw={300} ta={'center'} fs={'italic'}>
          Unknown
        </Text>
      </Stack>
    </Stack>
  )
}

export default UnknownArtistCard
