import { alpha, Box, Group, Text } from '@mantine/core'
import { IconMusicPlus } from '@tabler/icons-react'

interface NewArtistSongCardProps {
  openModal: () => void
}

function NewArtistSongCard({ openModal }: NewArtistSongCardProps) {
  return (
    <Group
      align={'center'}
      px={'md'}
      py={'sm'}
      style={{ cursor: 'pointer' }}
      sx={(theme) => ({
        transition: '0.3s',
        color: theme.colors.gray[6],
        '&:hover': {
          boxShadow: theme.shadows.xl,
          color: theme.colors.gray[7],
          backgroundColor: alpha(theme.colors.cyan[0], 0.15)
        }
      })}
      onClick={openModal}
    >
      <Box bd={'1px dashed gray'} p={'11px 9px 7px 9px'} style={{ borderRadius: '8px' }}>
        <IconMusicPlus size={18} />
      </Box>
      <Text fw={500} c={'inherit'} truncate={'end'}>
        Add New Songs
      </Text>
    </Group>
  )
}

export default NewArtistSongCard
