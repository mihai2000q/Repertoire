import { AspectRatio, Image, Stack, Text } from '@mantine/core'
import unknownPlaceholder from '../../assets/unknown-placeholder.png'
import { useState } from 'react'
import { useNavigate } from 'react-router-dom'

function UnknownArtistCard() {
  const navigate = useNavigate()

  const [isImageHovered, setIsImageHovered] = useState(false)

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
        ...(isImageHovered && { transform: 'scale(1.1)' })
      }}
    >
      <AspectRatio>
        <Image
          onMouseEnter={() => setIsImageHovered(true)}
          onMouseLeave={() => setIsImageHovered(false)}
          radius={'lg'}
          src={unknownPlaceholder}
          onClick={handleClick}
          alt={'unknown-album'}
          sx={(theme) => ({
            cursor: 'pointer',
            transition: '0.3s',
            boxShadow: theme.shadows.xxl,
            '&:hover': { boxShadow: theme.shadows.xxl_hover }
          })}
        />
      </AspectRatio>

      <Stack pt={'xs'}>
        <Text fw={300} ta={'center'} fs={'italic'}>
          Unknown
        </Text>
      </Stack>
    </Stack>
  )
}

export default UnknownArtistCard
