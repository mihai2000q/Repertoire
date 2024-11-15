import {
  ActionIcon,
  Avatar,
  Button,
  Card,
  Divider,
  Group,
  Menu,
  SimpleGrid,
  Space,
  Stack,
  Text,
  Title
} from '@mantine/core'
import { useParams } from 'react-router-dom'
import { useGetArtistQuery } from '../../state/artistsApi.ts'
import artistPlaceholder from '../../assets/user-placeholder.jpg'
import ArtistLoader from '../../components/artists/loader/ArtistLoader.tsx'
import { useGetAlbumsQuery } from '../../state/albumsApi.ts'
import { useGetSongsQuery } from '../../state/songsApi.ts'
import ArtistAlbumsLoader from '../../components/artists/loader/ArtistAlbumsLoader.tsx'
import ArtistSongsLoader from '../../components/artists/loader/ArtistSongsLoader.tsx'
import ArtistSongCard from '../../components/artists/card/ArtistSongCard.tsx'
import ArtistAlbumCard from '../../components/artists/card/ArtistAlbumCard.tsx'
import { Dispatch, SetStateAction, useState } from 'react'
import Order from '../../types/Order.ts'
import artistSongsOrders from '../../data/artist/artistSongsOrders.ts'
import {
  IconCaretDownFilled,
  IconCheck,
  IconDisc,
  IconDots,
  IconMusicPlus,
  IconPlus
} from '@tabler/icons-react'
import AddNewArtistSongModal from '../../components/artists/modal/AddNewArtistSongModal.tsx'
import { useDisclosure } from '@mantine/hooks'
import AddExistingArtistSongsModal from '../../components/artists/modal/AddExistingArtistSongsModal.tsx'
import AddExistingArtistAlbumsModal from '../../components/artists/modal/AddExistingArtistAlbumsModal.tsx'
import AddNewArtistAlbumModal from '../../components/artists/modal/AddNewArtistAlbumModal.tsx'
import artistAlbumsOrders from '../../data/artist/artistAlbumsOrders.ts'
import NewHorizontalCard from '../../components/card/NewHorizontalCard.tsx'

const SortButton = ({
  order,
  setOrder,
  orders
}: {
  order: Order
  setOrder: Dispatch<SetStateAction<Order>>
  orders: Order[]
}) => (
  <Menu shadow={'sm'}>
    <Menu.Target>
      <Button
        variant={'subtle'}
        size={'compact-xs'}
        rightSection={<IconCaretDownFilled size={11} />}
        styles={{ section: { marginLeft: 4 } }}
      >
        {order.label}
      </Button>
    </Menu.Target>

    <Menu.Dropdown>
      {orders.map((o) => (
        <Menu.Item
          key={o.value}
          leftSection={order === o && <IconCheck size={12} />}
          onClick={() => setOrder(o)}
        >
          {o.label}
        </Menu.Item>
      ))}
    </Menu.Dropdown>
  </Menu>
)

function Artist() {
  const params = useParams()
  const artistId = params['id'] ?? ''

  const { data: artist, isLoading } = useGetArtistQuery(artistId)

  const [albumsOrder, setAlbumsOrder] = useState<Order>(artistAlbumsOrders[0])
  const [songsOrder, setSongsOrder] = useState<Order>(artistSongsOrders[0])

  const { data: albums, isLoading: isAlbumsLoading } = useGetAlbumsQuery({
    orderBy: [albumsOrder.value]
  })
  const { data: songs, isLoading: isSongsLoading } = useGetSongsQuery({
    orderBy: [songsOrder.value]
  })

  const [openedAddNewAlbum, { open: openAddNewAlbum, close: closeAddNewAlbum }] =
    useDisclosure(false)
  const [openedAddExistingAlbums, { open: openAddExistingAlbums, close: closeAddExistingAlbums }] =
    useDisclosure(false)

  const [openedAddNewSong, { open: openAddNewSong, close: closeAddNewSong }] = useDisclosure(false)
  const [openedAddExistingSongs, { open: openAddExistingSongs, close: closeAddExistingSongs }] =
    useDisclosure(false)

  if (isLoading) return <ArtistLoader />

  return (
    <Stack>
      <Group>
        <Avatar
          src={artist.imageUrl ?? artistPlaceholder}
          size={125}
          sx={(theme) => ({
            boxShadow: theme.shadows.sm
          })}
        />
        <Title order={3} fw={700}>
          {artist.name}
        </Title>
      </Group>

      <Divider />

      <Group align={'start'}>
        <Card variant={'panel'} p={0} h={'100%'} flex={1}>
          {isAlbumsLoading ? (
            <ArtistAlbumsLoader />
          ) : (
            <Stack gap={0}>
              <Group px={'md'} py={'xs'} gap={'xs'} align={'center'}>
                <Text fw={500}>Albums</Text>
                <SortButton
                  order={albumsOrder}
                  setOrder={setAlbumsOrder}
                  orders={artistAlbumsOrders}
                />
                <Space flex={1} />
                <Menu position={'bottom-end'}>
                  <Menu.Target>
                    <ActionIcon size={'md'} variant={'grey'}>
                      <IconDots size={15} />
                    </ActionIcon>
                  </Menu.Target>
                  <Menu.Dropdown>
                    <Menu.Item leftSection={<IconPlus size={15} />} onClick={openAddExistingAlbums}>
                      Add Existing Albums
                    </Menu.Item>
                    <Menu.Item leftSection={<IconDisc size={15} />} onClick={openAddNewAlbum}>
                      Add New Album
                    </Menu.Item>
                  </Menu.Dropdown>
                </Menu>
              </Group>
              <SimpleGrid
                cols={{ sm: 1, md: 2, xl: 3 }}
                spacing={0}
                verticalSpacing={0}
                style={{ overflow: 'auto', maxHeight: '55vh' }}
              >
                {albums.models.map((album) => (
                  <ArtistAlbumCard key={album.id} album={album} />
                ))}
                {albums.models.length === albums.totalCount && (
                  <NewHorizontalCard onClick={openAddExistingAlbums} borderRadius={'8px'}>
                    Add New Albums
                  </NewHorizontalCard>
                )}
              </SimpleGrid>
            </Stack>
          )}
        </Card>

        <Card variant={'panel'} p={0} h={'100%'} flex={1.05}>
          {isSongsLoading ? (
            <ArtistSongsLoader />
          ) : (
            <Stack gap={0}>
              <Group px={'md'} py={'xs'} gap={'xs'} align={'center'}>
                <Text fw={600}>Songs</Text>
                <SortButton
                  order={songsOrder}
                  setOrder={setSongsOrder}
                  orders={artistSongsOrders}
                />
                <Space flex={1} />
                <Menu position={'bottom-end'}>
                  <Menu.Target>
                    <ActionIcon size={'md'} variant={'grey'}>
                      <IconDots size={15} />
                    </ActionIcon>
                  </Menu.Target>
                  <Menu.Dropdown>
                    <Menu.Item leftSection={<IconPlus size={15} />} onClick={openAddExistingSongs}>
                      Add Existing Songs
                    </Menu.Item>
                    <Menu.Item leftSection={<IconMusicPlus size={15} />} onClick={openAddNewSong}>
                      Add New Song
                    </Menu.Item>
                  </Menu.Dropdown>
                </Menu>
              </Group>
              <Stack gap={0} style={{ overflow: 'auto', maxHeight: '55vh' }}>
                {songs.models.map((song) => (
                  <ArtistSongCard key={song.id} song={song} />
                ))}
                {songs.models.length === songs.totalCount && (
                  <NewHorizontalCard onClick={openAddExistingSongs}>
                    Add New Songs
                  </NewHorizontalCard>
                )}
              </Stack>
            </Stack>
          )}
        </Card>
      </Group>

      <AddNewArtistSongModal
        opened={openedAddNewSong}
        onClose={closeAddNewSong}
        artistId={artistId}
      />
      <AddExistingArtistSongsModal
        opened={openedAddExistingSongs}
        onClose={closeAddExistingSongs}
        artistId={artistId}
      />
      <AddNewArtistAlbumModal
        opened={openedAddNewAlbum}
        onClose={closeAddNewAlbum}
        artistId={artistId}
      />
      <AddExistingArtistAlbumsModal
        opened={openedAddExistingAlbums}
        onClose={closeAddExistingAlbums}
        artistId={artistId}
      />
    </Stack>
  )
}

export default Artist
