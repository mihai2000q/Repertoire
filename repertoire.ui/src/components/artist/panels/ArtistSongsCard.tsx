import WithTotalCountResponse from '../../../types/responses/WithTotalCountResponse.ts'
import Order from '../../../types/Order.ts'
import { Dispatch, memo, SetStateAction, useEffect, useRef } from 'react'
import Song from '../../../types/models/Song.ts'
import ArtistSongsLoader from '../loader/ArtistSongsLoader.tsx'
import {
  ActionIcon,
  Box,
  Card,
  Group,
  Loader,
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
import { useDisclosure, useIntersection } from '@mantine/hooks'
import CompactOrderButton from '../../@ui/button/CompactOrderButton.tsx'
import LoadingOverlayDebounced from '../../@ui/loader/LoadingOverlayDebounced.tsx'

interface ArtistSongsCardProps {
  songs: WithTotalCountResponse<Song>
  isUnknownArtist: boolean
  order: Order
  setOrder: Dispatch<SetStateAction<Order>>
  artistId: string | undefined
  isLoading?: boolean
  isFetching?: boolean
  isFetchingNextPage?: boolean
  fetchNextPage?: () => void
}

function ArtistSongsCard({
  songs,
  isLoading,
  isUnknownArtist,
  order,
  setOrder,
  artistId,
  isFetching,
  isFetchingNextPage,
  fetchNextPage
}: ArtistSongsCardProps) {
  const [openedAddNewSong, { open: openAddNewSong, close: closeAddNewSong }] = useDisclosure(false)
  const [openedAddExistingSongs, { open: openAddExistingSongs, close: closeAddExistingSongs }] =
    useDisclosure(false)

  const scrollRef = useRef()
  const { ref: lastSongRef, entry } = useIntersection({
    root: scrollRef.current,
    threshold: 0.1
  })
  useEffect(() => {
    if (entry?.isIntersecting === true) fetchNextPage()
  }, [entry?.isIntersecting])

  if (isLoading || !songs) return <ArtistSongsLoader />

  return (
    <Card aria-label={'songs-card'} variant={'panel'} p={0} mah={'100%'}>
      <Stack gap={0} mah={'100%'}>
        <LoadingOverlayDebounced visible={isFetching && !isFetchingNextPage} />

        <Group px={'md'} py={'xs'} gap={'xs'}>
          <Text fw={600}>Songs</Text>

          <CompactOrderButton
            availableOrders={artistSongsOrders}
            order={order}
            setOrder={setOrder}
          />

          <Space flex={1} />

          <Menu>
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
          viewportRef={scrollRef}
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
            <Songs
              songs={songs}
              artistId={artistId}
              isUnknownArtist={isUnknownArtist}
              order={order}
            />

            <Stack gap={0} align={'center'}>
              <Box ref={lastSongRef} w={1} h={1} />
              {isFetchingNextPage && <Loader size={30} mt={'xs'} mb={'md'} />}
            </Stack>

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

const Songs = memo(
  ({
    songs,
    artistId,
    isUnknownArtist,
    order
  }: {
    songs: WithTotalCountResponse<Song>
    artistId: string
    isUnknownArtist: boolean
    order: Order
  }) => {
    return (
      <div>
        {songs.models.map((song) => (
          <ArtistSongCard
            key={song.id}
            song={song}
            artistId={artistId}
            isUnknownArtist={isUnknownArtist}
            order={order}
          />
        ))}
      </div>
    )
  },
  (prevProps, nextProps) => {
    return (
      JSON.stringify(prevProps.songs) === JSON.stringify(nextProps.songs) &&
      JSON.stringify(prevProps.order) === JSON.stringify(nextProps.order) &&
      prevProps.artistId === nextProps.artistId &&
      prevProps.isUnknownArtist === nextProps.isUnknownArtist
    )
  }
)

Songs.displayName = 'Songs'

export default ArtistSongsCard
