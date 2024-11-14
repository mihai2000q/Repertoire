import { alpha, Card, Center } from '@mantine/core'
import { IconMusicPlus } from '@tabler/icons-react'

interface NewAlbumCardProps {
  openModal: () => void
}

function NewAlbumCard({ openModal }: NewAlbumCardProps) {
  return (
    <Card
      w={150}
      h={150}
      radius={'lg'}
      onClick={openModal}
      sx={(theme) => ({
        cursor: 'pointer',
        alignSelf: 'start',
        transition: '0.3s',
        boxShadow: theme.shadows.xxl,
        color: theme.colors.cyan[7],
        '&:hover': {
          boxShadow: theme.shadows.xxl_hover,
          color: theme.colors.cyan[8],
          backgroundColor: alpha(theme.colors.cyan[0], 0.2),
          transform: 'scale(1.1)'
        }
      })}
    >
      <Center h={'100%'}>
        <IconMusicPlus size={40} />
      </Center>
    </Card>
  )
}

export default NewAlbumCard
