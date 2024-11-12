import { Drawer, Image, Stack, Text, Title } from '@mantine/core'
import { useGetAlbumQuery } from '../../state/albumsApi.ts'
import { useAppSelector } from '../../state/store.ts'
import AlbumDrawerLoader from './loader/AlbumDrawerLoader.tsx'
import imagePlaceholder from '../../assets/image-placeholder-1.jpg'
import useTitleBarHeight from '../../hooks/useTitleBarHeight.ts'

interface AlbumDrawerProps {
  opened: boolean
  close: () => void
}

function AlbumDrawer({ opened, close }: AlbumDrawerProps) {
  const titleBarHeight = useTitleBarHeight()

  const albumId = useAppSelector((state) => state.global.albumDrawer.albumId)

  const album = useGetAlbumQuery(albumId, { skip: !albumId })?.data

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
      {!album ? (
        <AlbumDrawerLoader />
      ) : (
        <Stack gap={'xs'}>
          <Image
            src={album.imageUrl}
            fallbackSrc={imagePlaceholder}
            mah={400}
            alt={album.title}
            style={{ alignSelf: 'center' }}
          />

          <Stack px={'md'} gap={4}>
            <Title order={5} fw={600}>
              {album.title}
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

export default AlbumDrawer
