import { alpha, Box, Group, Text } from '@mantine/core'
import { IconMusicPlus } from '@tabler/icons-react'

interface NewArtistAlbumCardProps {
  openModal: () => void
}

function NewArtistAlbumCard({ openModal }: NewArtistAlbumCardProps) {
  return (
    <div>
      <Group
        align={'center'}
        wrap={'nowrap'}
        px={'md'}
        py={'xs'}
        sx={(theme) => ({
          borderRadius: '8px',
          cursor: 'pointer',
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
          Add New Albums
        </Text>
      </Group>
    </div>
  )
}

export default NewArtistAlbumCard
