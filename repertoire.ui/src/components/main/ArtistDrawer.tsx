import { Drawer, Image, Stack, Text, Title } from '@mantine/core'
import { useGetArtistQuery } from '../../state/artistsApi.ts'
import { useAppSelector } from '../../state/store.ts'
import ArtistDrawerLoader from './loader/ArtistDrawerLoader.tsx'
import imagePlaceholder from '../../assets/image-placeholder-1.jpg'
import useTitleBarHeight from '../../hooks/useTitleBarHeight.ts'

interface ArtistDrawerProps {
  opened: boolean
  close: () => void
}

function ArtistDrawer({ opened, close }: ArtistDrawerProps) {
  const titleBarHeight = useTitleBarHeight()

  const artistId = useAppSelector((state) => state.global.artistDrawer.artistId)

  const artist = useGetArtistQuery(artistId, { skip: !artistId })?.data

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
      {!artist ? (
        <ArtistDrawerLoader />
      ) : (
        <Stack gap={'xs'}>
          <Image
            src={artist.imageUrl}
            fallbackSrc={imagePlaceholder}
            mah={400}
            alt={artist.name}
            style={{ alignSelf: 'center' }}
          />

          <Stack px={'md'} gap={4}>
            <Title order={5} fw={600}>
              {artist.name}
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

export default ArtistDrawer
