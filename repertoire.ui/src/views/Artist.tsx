import {
  ActionIcon,
  Avatar,
  Button,
  Card,
  Divider,
  Grid,
  Group,
  LoadingOverlay,
  Menu,
  SimpleGrid,
  Space,
  Stack,
  Text,
  Title
} from '@mantine/core'
import { useNavigate, useParams } from 'react-router-dom'
import {
  useDeleteArtistMutation,
  useGetArtistQuery,
  useRemoveAlbumsFromArtistMutation,
  useRemoveSongsFromArtistMutation
} from '../state/artistsApi.ts'
import artistPlaceholder from '../assets/user-placeholder.jpg'
import unknownPlaceholder from '../assets/unknown-placeholder.png'
import ArtistLoader from '../components/artist/loader/ArtistLoader.tsx'
import { useGetAlbumsQuery } from '../state/albumsApi.ts'
import { useGetSongsQuery } from '../state/songsApi.ts'
import ArtistAlbumsLoader from '../components/artist/loader/ArtistAlbumsLoader.tsx'
import ArtistSongsLoader from '../components/artist/loader/ArtistSongsLoader.tsx'
import ArtistSongCard from '../components/artist/ArtistSongCard.tsx'
import ArtistAlbumCard from '../components/artist/ArtistAlbumCard.tsx'
import { Dispatch, SetStateAction, useState } from 'react'
import Order from '../types/Order.ts'
import artistSongsOrders from '../data/artist/artistSongsOrders.ts'
import {
  IconAlbum,
  IconCaretDownFilled,
  IconCheck,
  IconDisc,
  IconDots,
  IconEdit,
  IconMusicPlus,
  IconPlus,
  IconTrash
} from '@tabler/icons-react'
import AddNewArtistSongModal from '../components/artist/modal/AddNewArtistSongModal.tsx'
import { useDisclosure } from '@mantine/hooks'
import AddExistingArtistSongsModal from '../components/artist/modal/AddExistingArtistSongsModal.tsx'
import AddExistingArtistAlbumsModal from '../components/artist/modal/AddExistingArtistAlbumsModal.tsx'
import AddNewArtistAlbumModal from '../components/artist/modal/AddNewArtistAlbumModal.tsx'
import artistAlbumsOrders from '../data/artist/artistAlbumsOrders.ts'
import NewHorizontalCard from '../components/card/NewHorizontalCard.tsx'
import { toast } from 'react-toastify'
import HeaderPanelCard from '../components/card/HeaderPanelCard.tsx'
import EditArtistHeaderModal from '../components/artist/modal/EditArtistHeaderModal.tsx'

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
  const navigate = useNavigate()

  const params = useParams()
  const artistId = params['id'] ?? ''
  const isUnknownArtist = artistId === 'unknown'

  const [deleteArtistMutation] = useDeleteArtistMutation()

  const { data: artist, isLoading } = useGetArtistQuery(artistId, { skip: isUnknownArtist })

  const [albumsOrder, setAlbumsOrder] = useState<Order>(artistAlbumsOrders[0])
  const [songsOrder, setSongsOrder] = useState<Order>(artistSongsOrders[0])

  const {
    data: albums,
    isLoading: isAlbumsLoading,
    isFetching: isAlbumsFetching
  } = useGetAlbumsQuery({
    orderBy: [albumsOrder.value],
    searchBy: [isUnknownArtist ? 'artist_id IS NULL' : `artist_id = ${artistId}`]
  })
  const {
    data: songs,
    isLoading: isSongsLoading,
    isFetching: isSongsFetching
  } = useGetSongsQuery({
    orderBy: [songsOrder.value],
    searchBy: [isUnknownArtist ? 'artist_id IS NULL' : `artist_id = ${artistId}`]
  })

  const [removeAlbumsFromArtist] = useRemoveAlbumsFromArtistMutation()
  const [removeSongsFromArtist] = useRemoveSongsFromArtistMutation()

  const [openedEditArtistHeader, { open: openEditArtistHeader, close: closeEditArtistHeader }] =
    useDisclosure(false)

  const [openedAddNewAlbum, { open: openAddNewAlbum, close: closeAddNewAlbum }] =
    useDisclosure(false)
  const [openedAddExistingAlbums, { open: openAddExistingAlbums, close: closeAddExistingAlbums }] =
    useDisclosure(false)

  const [openedAddNewSong, { open: openAddNewSong, close: closeAddNewSong }] = useDisclosure(false)
  const [openedAddExistingSongs, { open: openAddExistingSongs, close: closeAddExistingSongs }] =
    useDisclosure(false)

  function handleDelete() {
    deleteArtistMutation(artist.id)
    navigate(`/artists`, { replace: true })
    toast.success(`${artist.name} deleted!`)
  }

  function handleRemoveAlbumsFromArtist(albumIds: string[]) {
    removeAlbumsFromArtist({ albumIds: albumIds, id: artistId })
  }

  function handleRemoveSongsFromArtist(songIds: string[]) {
    removeSongsFromArtist({ songIds: songIds, id: artistId })
  }

  if (isLoading) return <ArtistLoader />

  return (
    <Stack>
      <HeaderPanelCard
        onEditClick={openEditArtistHeader}
        menuDropdown={
          <>
            <Menu.Item leftSection={<IconEdit size={14} />} onClick={openEditArtistHeader}>
              Edit
            </Menu.Item>
            <Menu.Item leftSection={<IconTrash size={14} />} c={'red.5'} onClick={handleDelete}>
              Delete
            </Menu.Item>
          </>
        }
        hideIcons={isUnknownArtist}
      >
        <Group>
          <Avatar
            src={isUnknownArtist ? unknownPlaceholder : (artist?.imageUrl ?? artistPlaceholder)}
            size={125}
            style={(theme) => ({
              boxShadow: theme.shadows.md
            })}
          />
          <Stack
            gap={4}
            style={{ ...(!isUnknownArtist && { alignSelf: 'start', paddingTop: '16px' }) }}
          >
            {!isUnknownArtist && (
              <Text fw={500} inline>
                Artist
              </Text>
            )}
            {isUnknownArtist ? (
              <Title order={3} fw={200} fs={'italic'} mb={2}>
                Unknown
              </Title>
            ) : (
              <Title order={1} fw={700}>
                {artist?.name}
              </Title>
            )}
          </Stack>
        </Group>
      </HeaderPanelCard>

      <Divider />

      <Grid align={'flex-start'}>
        <Grid.Col span={{ sm: 12, md: 6.5 }}>
          <Card variant={'panel'} p={0} h={'100%'} flex={1}>
            {isAlbumsLoading ? (
              <ArtistAlbumsLoader />
            ) : (
              <Stack gap={0}>
                <LoadingOverlay visible={isAlbumsFetching} />

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
                      {!isUnknownArtist && (
                        <Menu.Item
                          leftSection={<IconPlus size={15} />}
                          onClick={openAddExistingAlbums}
                        >
                          Add Existing Albums
                        </Menu.Item>
                      )}
                      <Menu.Item leftSection={<IconDisc size={15} />} onClick={openAddNewAlbum}>
                        Add New Album
                      </Menu.Item>
                    </Menu.Dropdown>
                  </Menu>
                </Group>
                <SimpleGrid
                  cols={{ sm: 1, md: 2 }}
                  spacing={0}
                  verticalSpacing={0}
                  style={{ overflow: 'auto', maxHeight: '55vh' }}
                >
                  {albums.models.map((album) => (
                    <ArtistAlbumCard
                      key={album.id}
                      album={album}
                      handleRemove={() => handleRemoveAlbumsFromArtist([album.id])}
                      isUnknownArtist={isUnknownArtist}
                    />
                  ))}
                  {albums.models.length === albums.totalCount && (
                    <NewHorizontalCard
                      borderRadius={'8px'}
                      onClick={isUnknownArtist ? openAddNewAlbum : openAddExistingAlbums}
                      icon={<IconAlbum size={16} />}
                      p={'10px 9px 6px 9px'}
                    >
                      Add New Albums
                    </NewHorizontalCard>
                  )}
                </SimpleGrid>
              </Stack>
            )}
          </Card>
        </Grid.Col>

        <Grid.Col span={{ sm: 12, md: 5.5 }}>
          <Card variant={'panel'} p={0} h={'100%'} flex={1.05}>
            {isSongsLoading ? (
              <ArtistSongsLoader />
            ) : (
              <Stack gap={0}>
                <LoadingOverlay visible={isSongsFetching} />

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
                      {!isUnknownArtist && (
                        <Menu.Item
                          leftSection={<IconPlus size={15} />}
                          onClick={openAddExistingSongs}
                        >
                          Add Existing Songs
                        </Menu.Item>
                      )}
                      <Menu.Item leftSection={<IconMusicPlus size={15} />} onClick={openAddNewSong}>
                        Add New Song
                      </Menu.Item>
                    </Menu.Dropdown>
                  </Menu>
                </Group>
                <Stack gap={0} style={{ overflow: 'auto', maxHeight: '55vh' }}>
                  {songs.models.map((song) => (
                    <ArtistSongCard
                      key={song.id}
                      song={song}
                      handleRemove={() => handleRemoveSongsFromArtist([song.id])}
                      isUnknownArtist={isUnknownArtist}
                    />
                  ))}
                  {songs.models.length === songs.totalCount && (
                    <NewHorizontalCard
                      onClick={isUnknownArtist ? openAddNewSong : openAddExistingSongs}
                    >
                      Add New Songs
                    </NewHorizontalCard>
                  )}
                </Stack>
              </Stack>
            )}
          </Card>
        </Grid.Col>
      </Grid>

      {!isUnknownArtist && (
        <EditArtistHeaderModal
          artist={artist}
          opened={openedEditArtistHeader}
          onClose={closeEditArtistHeader}
        />
      )}
      <AddNewArtistSongModal
        opened={openedAddNewSong}
        onClose={closeAddNewSong}
        artistId={artist?.id}
      />
      <AddExistingArtistSongsModal
        opened={openedAddExistingSongs}
        onClose={closeAddExistingSongs}
        artistId={artistId}
      />
      <AddNewArtistAlbumModal
        opened={openedAddNewAlbum}
        onClose={closeAddNewAlbum}
        artistId={artist?.id}
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
