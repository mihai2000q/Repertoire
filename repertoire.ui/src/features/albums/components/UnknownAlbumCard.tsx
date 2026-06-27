import { Avatar, Center, Stack, Text } from '@mantine/core'
import { useNavigate } from 'react-router-dom'
import { IconQuestionMark } from '@tabler/icons-react'
import { useHover } from '@mantine/hooks'

function UnknownArtistCard() {
  const navigate = useNavigate()
  const { ref, hovered } = useHover()

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
        ...(hovered && { transform: 'scale(1.1)' })
      }}
    >
      <Avatar
        ref={ref}
        radius={'10%'}
        w={'100%'}
        h={'unset'}
        style={(theme) => ({
          aspectRatio: 1,
          cursor: 'pointer',
          transition: '0.3s',
          boxShadow: hovered ? theme.shadows.xxl_hover : theme.shadows.xxl
        })}
        onClick={handleClick}
      >
        <Center c={'gray.6'}>
          <IconQuestionMark
            aria-label={'icon-unknown-album'}
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
