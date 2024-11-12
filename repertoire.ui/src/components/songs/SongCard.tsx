import Song from '../../types/models/Song'
import imagePlaceholder from '../../assets/image-placeholder-1.jpg'
import { Card, Group, Image, Text, Tooltip } from '@mantine/core'
import { IconMicrophoneFilled } from '@tabler/icons-react'
import { useAppDispatch } from '../../state/store.ts'
import { openSongDrawer } from '../../state/globalSlice.ts'

interface SongCardProps {
  song: Song
}

function SongCard({ song }: SongCardProps) {
  const dispatch = useAppDispatch()

  function handleClick() {
    dispatch(openSongDrawer(song.id))
  }

  return (
    <Card
      data-testid={`song-card-${song.id}`}
      p="sm"
      shadow="md"
      h={253}
      w={175}
      onClick={handleClick}
    >
      <Card.Section>
        <Image src={song.imageUrl} fallbackSrc={imagePlaceholder} h={150} alt={song.title} />
      </Card.Section>

      <Group justify="space-between" mt="sm" mb="xs">
        <Text fw={500} lineClamp={2}>
          {song.title}
        </Text>
      </Group>

      <Text size="sm" c="dimmed" mb="xs">
        With Fjord Tours you can explore more of the
      </Text>

      <Group c={'cyan.9'}>
        {song.isRecorded && (
          <Tooltip label={'This song is recorded'} openDelay={200} position="bottom">
            <IconMicrophoneFilled size={18} />
          </Tooltip>
        )}
      </Group>
    </Card>
  )
}

export default SongCard
