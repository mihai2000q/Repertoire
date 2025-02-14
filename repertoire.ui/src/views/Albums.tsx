import { useGetAlbumsQuery } from '../state/api/albumsApi.ts'
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
import AlbumsLoader from '../components/albums/AlbumsLoader.tsx'
import AlbumCard from '../components/albums/AlbumCard.tsx'
import AddNewAlbumModal from '../components/albums/modal/AddNewAlbumModal.tsx'
import { useDisclosure } from '@mantine/hooks'
import { IconArrowsSort, IconFilterFilled, IconMusicPlus, IconPlus } from '@tabler/icons-react'
import usePaginationInfo from '../hooks/usePaginationInfo.ts'
import useShowUnknownAlbum from '../hooks/useShowUnknownAlbum.ts'
import UnknownAlbumCard from '../components/albums/UnknownAlbumCard.tsx'
import useFixedDocumentTitle from '../hooks/useFixedDocumentTitle.ts'
import useSearchParamsState from '../hooks/useSearchParamsState.ts'
import albumsSearchParamsState from '../state/searchParams/AlbumsSearchParamsState.ts'

function Albums() {
  useFixedDocumentTitle('Albums')
  const [searchParams, setSearchParams] = useSearchParamsState(albumsSearchParamsState)
  const { currentPage } = searchParams

  const pageSize = 40
  const {
    data: albums,
    isLoading,
    isFetching
  } = useGetAlbumsQuery({
    pageSize: pageSize,
    currentPage: currentPage,
    orderBy: ['created_at DESC']
  })

  const showUnknownAlbum = useShowUnknownAlbum()

  const { startCount, endCount, totalPages } = usePaginationInfo(
    albums?.totalCount + (showUnknownAlbum ? 1 : 0),
    pageSize,
    currentPage
  )

  const [openedAddNewAlbumModal, { open: openAddNewAlbumModal, close: closeAddNewAlbumModal }] =
    useDisclosure(false)

  const handleCurrentPageChange = (p: number) => {
    setSearchParams({ ...searchParams, currentPage: p })
  }

  return (
    <Stack h={'100%'} gap={'xs'} px={'xl'}>
      <Group gap={'xxs'}>
        <Title order={3} fw={800} fz={'max(2.5vw, 32px)'}>
          Albums
        </Title>
        <ActionIcon
          aria-label={'new-album'}
          variant={'grey'}
          size={'lg'}
          onClick={openAddNewAlbumModal}
        >
          <IconPlus />
        </ActionIcon>
        <Space flex={1} />
        <ActionIcon aria-label={'order-albums'} variant={'grey'} size={'lg'} disabled={isLoading}>
          <IconArrowsSort size={17} />
        </ActionIcon>
        <ActionIcon aria-label={'filter-albums'} variant={'grey'} size={'lg'} disabled={isLoading}>
          <IconFilterFilled size={17} />
        </ActionIcon>
      </Group>
      {!isLoading && (
        <Text inline mb={'xs'}>
          {startCount} - {endCount} albums out of {albums?.totalCount + (showUnknownAlbum ? 1 : 0)}
        </Text>
      )}

      {albums?.totalCount === 0 && !showUnknownAlbum && (
        <Text mt={'sm'}>There are no albums yet. Try to add one</Text>
      )}
      <Group gap={'xl'} align={'start'}>
        {isLoading && <AlbumsLoader />}
        {albums?.models.map((album) => <AlbumCard key={album.id} album={album} />)}
        {showUnknownAlbum && currentPage == totalPages && <UnknownAlbumCard />}
        {((albums?.totalCount > 0 && currentPage == totalPages) ||
          (albums?.totalCount === 0 && showUnknownAlbum)) && (
          <Card
            aria-label={'new-album-card'}
            w={150}
            h={150}
            radius={'lg'}
            onClick={openAddNewAlbumModal}
            variant={'add-new'}
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
            data-testid={'albums-pagination'}
            value={currentPage}
            onChange={handleCurrentPageChange}
            total={totalPages}
          />
        ) : (
          <Loader size={25} />
        )}
      </Box>

      <AddNewAlbumModal opened={openedAddNewAlbumModal} onClose={closeAddNewAlbumModal} />
    </Stack>
  )
}

export default Albums
