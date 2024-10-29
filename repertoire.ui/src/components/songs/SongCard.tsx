import Song from '../../types/models/Song'
import demoSong from '../../assets/demoSong.jpg'
import { Card, Group, Image, Text, Tooltip } from '@mantine/core'
import { IconMicrophoneFilled } from '@tabler/icons-react'
import { useAppDispatch } from '../../state/store.ts'
import { setSongId } from '../../state/songsSlice.ts'

interface SongCardProps {
  song: Song
  openDrawer: () => void
}

function SongCard({ song, openDrawer }: SongCardProps) {
  const dispatch = useAppDispatch()

  function handleClick() {
    openDrawer()
    dispatch(setSongId(song.id))
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
        <Image src={demoSong} h={150} alt={song.title} />
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
