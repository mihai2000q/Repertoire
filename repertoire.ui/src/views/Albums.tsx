import { useState } from 'react'
import { useGetAlbumsQuery } from '../state/albumsApi.ts'
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

function Albums() {
  const [currentPage, setCurrentPage] = useState(1)
  const { data: albums, isLoading, isFetching } = useGetAlbumsQuery({
    pageSize: 20,
    currentPage: currentPage,
    orderBy: ['created_at DESC']
  })

  const showUnknownAlbum = useShowUnknownAlbum()

  const { startCount, endCount, totalPages } = usePaginationInfo(
    albums?.totalCount + (showUnknownAlbum ? 1 : 0),
    20,
    currentPage
  )

  const [openedAddNewAlbumModal, { open: openAddNewAlbumModal, close: closeAddNewAlbumModal }] =
    useDisclosure(false)

  return (
    <Stack h={'100%'} gap={'xs'}>
      <Group gap={4} align={'center'}>
        <Title order={3} fw={800}>
          Albums
        </Title>
        <ActionIcon aria-label={'new-album'} variant={'grey'} size={'lg'} onClick={openAddNewAlbumModal}>
          <IconPlus />
        </ActionIcon>
        <Space flex={1} />
        <ActionIcon aria-label={'order-albums'} variant={'grey'} size={'lg'}>
          <IconArrowsSort size={17} />
        </ActionIcon>
        <ActionIcon aria-label={'filter-albums'} variant={'grey'} size={'lg'}>
          <IconFilterFilled size={17} />
        </ActionIcon>
      </Group>
      {!isFetching && (
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
        {((albums?.totalCount > 0 && currentPage == totalPages) || showUnknownAlbum) && (
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

      <Box style={{ alignSelf: 'center' }} pb={'xs'}>
        {!isFetching ? (
          <Pagination
            data-testid={'albums-pagination'}
            value={currentPage}
            onChange={setCurrentPage}
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
