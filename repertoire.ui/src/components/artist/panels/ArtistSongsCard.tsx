import WithTotalCountResponse from '../../../types/responses/WithTotalCountResponse.ts'
import Order from '../../../types/Order.ts'
import { Dispatch, SetStateAction } from 'react'
import Song from '../../../types/models/Song.ts'
import ArtistSongsLoader from '../loader/ArtistSongsLoader.tsx'
import {
  ActionIcon,
  Button,
  Card,
  Group,
  LoadingOverlay,
  Menu,
  Space,
  Stack,
  Text
} from '@mantine/core'
import artistSongsOrders from '../../../data/artist/artistSongsOrders.ts'
import {
  IconCaretDownFilled,
  IconCheck,
  IconDots,
  IconMusicPlus,
  IconPlus
} from '@tabler/icons-react'
import ArtistSongCard from '../ArtistSongCard.tsx'
import NewHorizontalCard from '../../card/NewHorizontalCard.tsx'
import AddNewArtistSongModal from '../modal/AddNewArtistSongModal.tsx'
import AddExistingArtistSongsModal from '../modal/AddExistingArtistSongsModal.tsx'
import { useDisclosure } from '@mantine/hooks'
import { useRemoveSongsFromArtistMutation } from '../../../state/artistsApi.ts'

interface ArtistSongsCardProps {
  songs: WithTotalCountResponse<Song>
  isLoading: boolean
  isFetching: boolean
  isUnknownArtist: boolean
  order: Order
  setOrder: Dispatch<SetStateAction<Order>>
  artistId: string | undefined
}

function ArtistSongsCard({
  songs,
  isLoading,
  isFetching,
  isUnknownArtist,
  order,
  setOrder,
  artistId
}: ArtistSongsCardProps) {
  const [removeSongsFromArtist] = useRemoveSongsFromArtistMutation()

  const [openedAddNewSong, { open: openAddNewSong, close: closeAddNewSong }] = useDisclosure(false)
  const [openedAddExistingSongs, { open: openAddExistingSongs, close: closeAddExistingSongs }] =
    useDisclosure(false)

  function handleRemoveSongsFromArtist(songIds: string[]) {
    removeSongsFromArtist({ songIds: songIds, id: artistId })
  }

  return (
    <Card variant={'panel'} p={0} h={'100%'}>
      {isLoading ? (
        <ArtistSongsLoader />
      ) : (
        <Stack gap={0}>
          <LoadingOverlay visible={isFetching} />

          <Group px={'md'} py={'xs'} gap={'xs'} align={'center'}>
            <Text fw={600}>Songs</Text>

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
                {artistSongsOrders.map((o) => (
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
                  <Menu.Item leftSection={<IconPlus size={15} />} onClick={openAddExistingSongs}>
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
              <NewHorizontalCard onClick={isUnknownArtist ? openAddNewSong : openAddExistingSongs}>
                Add New Songs
              </NewHorizontalCard>
            )}
          </Stack>
        </Stack>
      )}

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
    </Card>
  )
}

export default ArtistSongsCard
