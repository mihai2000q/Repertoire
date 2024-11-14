import { useState } from 'react'
import { useGetAlbumsQuery } from '../../state/albumsApi.ts'
import {
  ActionIcon,
  Box,
  Group,
  Loader,
  Pagination,
  Space,
  Stack,
  Text,
  Title
} from '@mantine/core'
import AlbumsLoader from '../../components/albums/loader/AlbumsLoader.tsx'
import AlbumCard from '../../components/albums/AlbumCard.tsx'
import AddNewAlbumModal from '../../components/albums/modal/AddNewAlbumModal.tsx'
import { useDisclosure } from '@mantine/hooks'
import { IconArrowsSort, IconFilterFilled, IconPlus } from '@tabler/icons-react'
import NewAlbumCard from "../../components/albums/NewAlbumCard.tsx";
import usePaginationInfo from "../../hooks/usePaginationInfo.ts";

function Albums() {
  const [currentPage, setCurrentPage] = useState(1)
  const { data: albums, isLoading } = useGetAlbumsQuery({
    pageSize: 20,
    currentPage: currentPage
  })

  const { startCount, endCount, totalPages } = usePaginationInfo(albums?.totalCount, 20, currentPage)

  const [openedAddNewAlbumModal, { open: openAddNewAlbumModal, close: closeAddNewAlbumModal }] =
    useDisclosure(false)

  return (
    <Stack h={'100%'} gap={'xs'}>
      <AddNewAlbumModal opened={openedAddNewAlbumModal} onClose={closeAddNewAlbumModal} />

      <Group gap={4} align={'center'}>
        <Title order={3} fw={800}>
          Albums
        </Title>
        <ActionIcon variant={'grey'} size={'lg'} onClick={openAddNewAlbumModal}>
          <IconPlus />
        </ActionIcon>
        <Space flex={1} />
        <ActionIcon variant={'grey'} size={'lg'}>
          <IconArrowsSort size={17} />
        </ActionIcon>
        <ActionIcon variant={'grey'} size={'lg'}>
          <IconFilterFilled size={17} />
        </ActionIcon>
      </Group>
      {!isLoading && (
        <Text inline mb={'xs'}>
          {startCount} - {endCount} albums out of {albums?.totalCount}
        </Text>
      )}

      {albums?.totalCount === 0 && <Text mt={'sm'}>There are no albums yet. Try to add one</Text>}
      <Group gap={'xl'}>
        {isLoading && <AlbumsLoader />}
        {albums?.models.map((album) => <AlbumCard key={album.id} album={album} />)}
        {albums?.totalCount > 0 && <NewAlbumCard openModal={openAddNewAlbumModal} />}
      </Group>

      <Space flex={1} />

      <Box style={{ alignSelf: 'center' }} pb={'xs'}>
        {!isLoading ? (
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
    </Stack>
  )
}

export default Albums
