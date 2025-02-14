import { useGetPlaylistsQuery } from '../state/api/playlistsApi.ts'
import usePaginationInfo from '../hooks/usePaginationInfo.ts'
import { useDisclosure } from '@mantine/hooks'
import {
  ActionIcon,
  Box,
  Card,
  Center,
  Group,
  Loader,
  Pagination,
  Space,
  Stack,
  Text,
  Title
} from '@mantine/core'
import AddNewPlaylistModal from '../components/playlists/modal/AddNewPlaylistModal.tsx'
import { IconArrowsSort, IconFilterFilled, IconMusicPlus, IconPlus } from '@tabler/icons-react'
import PlaylistsLoader from '../components/playlists/PlaylistsLoader.tsx'
import PlaylistCard from '../components/playlists/PlaylistCard.tsx'
import useFixedDocumentTitle from '../hooks/useFixedDocumentTitle.ts'
import useSearchParamsState from '../hooks/useSearchParamsState.ts'
import songsSearchParamsState from '../state/searchParams/SongsSearchParamsState.ts'

function Playlists() {
  useFixedDocumentTitle('Playlists')
  const [searchParams, setSearchParams] = useSearchParamsState(songsSearchParamsState)
  const { currentPage } = searchParams

  const pageSize = 40
  const {
    data: playlists,
    isLoading,
    isFetching
  } = useGetPlaylistsQuery({
    pageSize: pageSize,
    currentPage: currentPage,
    orderBy: ['created_at DESC']
  })

  const { startCount, endCount, totalPages } = usePaginationInfo(
    playlists?.totalCount,
    pageSize,
    currentPage
  )

  const [
    openedAddNewPlaylistModal,
    { open: openAddNewPlaylistModal, close: closeAddNewPlaylistModal }
  ] = useDisclosure(false)

  const handleCurrentPageChange = (p: number) => {
    setSearchParams({ ...searchParams, currentPage: p })
  }

  return (
    <Stack h={'100%'} gap={'xs'} px={'xl'}>
      <Group gap={'xxs'}>
        <Title order={3} fw={800} fz={'max(2.5vw, 32px)'}>
          Playlists
        </Title>
        <ActionIcon
          aria-label={'new-playlist'}
          variant={'grey'}
          size={'lg'}
          onClick={openAddNewPlaylistModal}
        >
          <IconPlus />
        </ActionIcon>
        <Space flex={1} />
        <ActionIcon aria-label={'order-playlists'} variant={'grey'} size={'lg'} disabled={isLoading}>
          <IconArrowsSort size={17} />
        </ActionIcon>
        <ActionIcon aria-label={'filter-playlists'} variant={'grey'} size={'lg'} disabled={isLoading}>
          <IconFilterFilled size={17} />
        </ActionIcon>
      </Group>
      {!isLoading && (
        <Text inline mb={'xs'}>
          {startCount} - {endCount} playlists out of {playlists?.totalCount}
        </Text>
      )}

      {playlists?.totalCount === 0 && (
        <Text mt={'sm'}>There are no playlists yet. Try to add one</Text>
      )}
      <Group gap={'lg'} align={'stretch'}>
        {isLoading && <PlaylistsLoader />}
        {playlists?.models.map((playlist) => (
          <PlaylistCard key={playlist.id} playlist={playlist} />
        ))}
        {playlists?.totalCount > 0 && currentPage == totalPages && (
          <Card
            aria-label={'new-playlist-card'}
            variant={'add-new'}
            radius={'lg'}
            w={150}
            h={150}
            onClick={openAddNewPlaylistModal}
          >
            <Center h={'100%'}>
              <IconMusicPlus size={40} />
            </Center>
          </Card>
        )}
      </Group>

      <Space flex={1} />

      <Box style={{ alignSelf: 'center' }} pb={'md'}>
        {!isFetching ? (
          <Pagination
            data-testid={'playlists-pagination'}
            value={currentPage}
            onChange={handleCurrentPageChange}
            total={totalPages}
          />
        ) : (
          <Loader size={25} />
        )}
      </Box>

      <AddNewPlaylistModal opened={openedAddNewPlaylistModal} onClose={closeAddNewPlaylistModal} />
    </Stack>
  )
}

export default Playlists
