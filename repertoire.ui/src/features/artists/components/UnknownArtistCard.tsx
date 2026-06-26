import { Avatar, Center, Stack, Text } from '@mantine/core'
import { useNavigate } from 'react-router-dom'
import { IconQuestionMark } from '@tabler/icons-react'
import { useHover } from '@mantine/hooks'

function UnknownArtistCard() {
  const navigate = useNavigate()
  const { ref, hovered } = useHover()

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
        ...(hovered && {
          transform: 'scale(1.1)'
        })
      }}
    >
      <Avatar
        ref={ref}
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
