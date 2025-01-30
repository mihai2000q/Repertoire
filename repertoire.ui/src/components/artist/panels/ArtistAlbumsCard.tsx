import { Dispatch, SetStateAction } from 'react'
import ArtistAlbumsLoader from '../loader/ArtistAlbumsLoader.tsx'
import {
  ActionIcon,
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
import { IconAlbum, IconDisc, IconDots, IconPlus } from '@tabler/icons-react'
import ArtistAlbumCard from '../ArtistAlbumCard.tsx'
import NewHorizontalCard from '../../@ui/card/NewHorizontalCard.tsx'
import { useDisclosure } from '@mantine/hooks'
import { useRemoveAlbumsFromArtistMutation } from '../../../state/api/artistsApi.ts'
import Order from '../../../types/Order.ts'
import WithTotalCountResponse from '../../../types/responses/WithTotalCountResponse.ts'
import Album from '../../../types/models/Album.ts'
import AddNewArtistAlbumModal from '../modal/AddNewArtistAlbumModal.tsx'
import AddExistingArtistAlbumsModal from '../modal/AddExistingArtistAlbumsModal.tsx'
import CompactOrderButton from '../../@ui/button/CompactOrderButton.tsx'

interface ArtistAlbumsCardProps {
  albums: WithTotalCountResponse<Album>
  isLoading: boolean
  isUnknownArtist: boolean
  order: Order
  setOrder: Dispatch<SetStateAction<Order>>
  artistId: string | undefined
  isFetching?: boolean
}

function ArtistAlbumsCard({
  albums,
  isLoading,
  isUnknownArtist,
  order,
  setOrder,
  artistId,
  isFetching
}: ArtistAlbumsCardProps) {
  const [removeAlbumsFromArtist] = useRemoveAlbumsFromArtistMutation()

  const [openedAddNewAlbum, { open: openAddNewAlbum, close: closeAddNewAlbum }] =
    useDisclosure(false)
  const [openedAddExistingAlbums, { open: openAddExistingAlbums, close: closeAddExistingAlbums }] =
    useDisclosure(false)

  function handleRemoveAlbumsFromArtist(albumIds: string[]) {
    removeAlbumsFromArtist({ albumIds: albumIds, id: artistId })
  }

  if (isLoading) return <ArtistAlbumsLoader />

  return (
    <Card variant={'panel'} aria-label={'albums-card'} p={0} h={'100%'} mb={'lg'}>
      <Stack gap={0}>
        <LoadingOverlay visible={isFetching} />

        <Group px={'md'} py={'xs'} gap={'xs'}>
          <Text fw={600}>Albums</Text>

          <CompactOrderButton
            availableOrders={artistAlbumsOrders}
            order={order}
            setOrder={setOrder}
          />

          <Space flex={1} />

          <Menu position={'bottom-end'}>
            <Menu.Target>
              <ActionIcon size={'md'} variant={'grey'} aria-label={'albums-more-menu'}>
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

        <SimpleGrid cols={{ sm: 1, md: 2, xl: 3 }} spacing={0} verticalSpacing={0}>
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
              ariaLabel={'new-albums-card'}
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
