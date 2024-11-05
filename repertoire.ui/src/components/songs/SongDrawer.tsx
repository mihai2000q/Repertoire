import { Drawer, Image, Stack, Text, Title } from '@mantine/core'
import { useGetSongQuery } from '../../state/songsApi.ts'
import { useAppSelector } from '../../state/store.ts'
import SongDrawerLoader from './loader/SongDrawerLoader.tsx'
import demoSong from '../../assets/demoSong.jpg'
import useTitleBarHeight from '../../hooks/useTitleBarHeight.ts'

interface SongDrawerProps {
  opened: boolean
  close: () => void
}

function SongDrawer({ opened, close }: SongDrawerProps) {
  const titleBarHeight = useTitleBarHeight()

  const songId = useAppSelector((state) => state.songs.songId)

  const song = useGetSongQuery(songId, { skip: !songId })?.data

  return (
    <Drawer
      withCloseButton={false}
      opened={opened}
      onClose={close}
      position="right"
      overlayProps={{ backgroundOpacity: 0.1 }}
      shadow="xl"
      radius={'8 0 0 8'}
      styles={{
        overlay: {
          height: `calc(100% - ${titleBarHeight})`,
          marginTop: titleBarHeight
        },
        inner: {
          height: `calc(100% - ${titleBarHeight})`,
          marginTop: titleBarHeight
        },
        body: {
          padding: 0,
          margin: 0
        }
      }}
    >
      {!song ? (
        <SongDrawerLoader />
      ) : (
        <Stack gap={'xs'}>
          <Image
            src={demoSong}
            mah={400}
            alt={song.title}
            style={{ alignSelf: 'center' }}
          />

          <Stack px={'md'} gap={4}>
            <Title order={5} fw={600}>
              {song.title}
            </Title>
            <Text size="sm" c="dimmed">
              With Fjord Tours you can explore more of the
            </Text>
          </Stack>
        </Stack>
      )}
    </Drawer>
  )
}

export default SongDrawer
