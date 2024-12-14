import { Dispatch, SetStateAction } from 'react'
import ArtistAlbumsLoader from '../loader/ArtistAlbumsLoader.tsx'
import {
  ActionIcon,
  Button,
  Card,
  Group,
  LoadingOverlay,
  Menu,
  SimpleGrid,
  Space,
  Stack,
  Text
} from '@mantine/core'
import artistAlbumsOrders from '../../../data/artist/artistAlbumsOrders.ts'
import {
  IconAlbum,
  IconCaretDownFilled,
  IconCheck,
  IconDisc,
  IconDots,
  IconPlus
} from '@tabler/icons-react'
import ArtistAlbumCard from '../ArtistAlbumCard.tsx'
import NewHorizontalCard from '../../card/NewHorizontalCard.tsx'
import { useDisclosure } from '@mantine/hooks'
import { useRemoveAlbumsFromArtistMutation } from '../../../state/artistsApi.ts'
import Order from '../../../types/Order.ts'
import WithTotalCountResponse from '../../../types/responses/WithTotalCountResponse.ts'
import Album from '../../../types/models/Album.ts'
import AddNewArtistAlbumModal from "../modal/AddNewArtistAlbumModal.tsx";
import AddExistingArtistAlbumsModal from "../modal/AddExistingArtistAlbumsModal.tsx";

interface ArtistAlbumsCardProps {
  albums: WithTotalCountResponse<Album>
  isLoading: boolean
  isFetching: boolean
  isUnknownArtist: boolean
  order: Order
  setOrder: Dispatch<SetStateAction<Order>>
  artistId: string | undefined
}

function ArtistAlbumsCard({
  albums,
  isLoading,
  isFetching,
  isUnknownArtist,
  order,
  setOrder,
  artistId
}: ArtistAlbumsCardProps) {
  const [removeAlbumsFromArtist] = useRemoveAlbumsFromArtistMutation()

  const [openedAddNewAlbum, { open: openAddNewAlbum, close: closeAddNewAlbum }] =
    useDisclosure(false)
  const [openedAddExistingAlbums, { open: openAddExistingAlbums, close: closeAddExistingAlbums }] =
    useDisclosure(false)

  function handleRemoveAlbumsFromArtist(albumIds: string[]) {
    removeAlbumsFromArtist({ albumIds: albumIds, id: artistId })
  }

  return (
    <Card variant={'panel'} p={0} h={'100%'}>
      {isLoading ? (
        <ArtistAlbumsLoader />
      ) : (
        <Stack gap={0}>
          <LoadingOverlay visible={isFetching} />

          <Group px={'md'} py={'xs'} gap={'xs'} align={'center'}>
            <Text fw={600}>Albums</Text>

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
                {artistAlbumsOrders.map((o) => (
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

            <Space flex={1} />

            <Menu position={'bottom-end'}>
              <Menu.Target>
                <ActionIcon size={'md'} variant={'grey'}>
                  <IconDots size={15} />
                </ActionIcon>
              </Menu.Target>
              <Menu.Dropdown>
                {!isUnknownArtist && (
                  <Menu.Item leftSection={<IconPlus size={15} />} onClick={openAddExistingAlbums}>
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
    </Card>
  )
}

export default ArtistAlbumsCard
