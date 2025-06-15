import WithTotalCountResponse from '../../../types/responses/WithTotalCountResponse.ts'
import Order from '../../../types/Order.ts'
import { Dispatch, SetStateAction } from 'react'
import Song from '../../../types/models/Song.ts'
import ArtistSongsLoader from '../loader/ArtistSongsLoader.tsx'
import {
  ActionIcon,
  Card,
  Group,
  LoadingOverlay,
  Menu,
  ScrollArea,
  Space,
  Stack,
  Text
} from '@mantine/core'
import artistSongsOrders from '../../../data/artist/artistSongsOrders.ts'
import { IconDots, IconMusicPlus, IconPlus } from '@tabler/icons-react'
import ArtistSongCard from '../ArtistSongCard.tsx'
import NewHorizontalCard from '../../@ui/card/NewHorizontalCard.tsx'
import AddNewArtistSongModal from '../modal/AddNewArtistSongModal.tsx'
import AddExistingArtistSongsModal from '../modal/AddExistingArtistSongsModal.tsx'
import { useDisclosure } from '@mantine/hooks'
import CompactOrderButton from '../../@ui/button/CompactOrderButton.tsx'

interface ArtistSongsCardProps {
  songs: WithTotalCountResponse<Song>
  isLoading: boolean
  isUnknownArtist: boolean
  order: Order
  setOrder: Dispatch<SetStateAction<Order>>
  artistId: string | undefined
  isFetching?: boolean
}

function ArtistSongsCard({
  songs,
  isLoading,
  isUnknownArtist,
  order,
  setOrder,
  artistId,
  isFetching
}: ArtistSongsCardProps) {
  const [openedAddNewSong, { open: openAddNewSong, close: closeAddNewSong }] = useDisclosure(false)
  const [openedAddExistingSongs, { open: openAddExistingSongs, close: closeAddExistingSongs }] =
    useDisclosure(false)

  if (isLoading || !songs) return <ArtistSongsLoader />

  return (
    <Card aria-label={'songs-card'} variant={'panel'} p={0} mah={'100%'}>
      <Stack gap={0} mah={'100%'}>
        <LoadingOverlay visible={isFetching} />

        <Group px={'md'} py={'xs'} gap={'xs'}>
          <Text fw={600}>Songs</Text>

          <CompactOrderButton
            availableOrders={artistSongsOrders}
            order={order}
            setOrder={setOrder}
          />

          <Space flex={1} />

          <Menu position={'bottom-end'}>
            <Menu.Target>
              <ActionIcon size={'md'} variant={'grey'} aria-label={'songs-more-menu'}>
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

        <ScrollArea.Autosize
          scrollbars={'y'}
          scrollbarSize={7}
          mah={'100%'}
          styles={{
            viewport: {
              '> div': {
                width: 0,
                minWidth: '100%'
              }
            }
          }}
        >
          <Stack gap={0} style={{ overflow: 'hidden' }}>
            {songs.models.map((song) => (
              <ArtistSongCard
                key={song.id}
                song={song}
                artistId={artistId}
                isUnknownArtist={isUnknownArtist}
                order={order}
              />
            ))}
            {songs.models.length === songs.totalCount && (
              <NewHorizontalCard
                ariaLabel={'new-songs-card'}
                onClick={isUnknownArtist ? openAddNewSong : openAddExistingSongs}
              >
                Add New Songs
              </NewHorizontalCard>
            )}
          </Stack>
        </ScrollArea.Autosize>
      </Stack>

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
