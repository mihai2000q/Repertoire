import { Card, Center } from '@mantine/core'
import { IconMusicPlus } from '@tabler/icons-react'

interface NewSongCardProps {
  openModal: () => void
}

function NewSongCard({ openModal }: NewSongCardProps) {
  return (
    <Card data-testid={'new-song-card'} padding={0} shadow="md" h={253} w={175} onClick={openModal}>
      <Center c={'cyan.7'} bg={'gray.1'} h={'100%'}>
        <IconMusicPlus size={35} />
      </Center>
    </Card>
  )
}

export default NewSongCard
