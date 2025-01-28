import { ReactElement, useState } from 'react'
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
import { useGetSongsQuery } from '../state/songsApi.ts'
import SongCard from '../components/songs/SongCard.tsx'
import { IconArrowsSort, IconFilterFilled, IconMusicPlus, IconPlus } from '@tabler/icons-react'
import { useDisclosure } from '@mantine/hooks'
import AddNewSongModal from '../components/songs/modal/AddNewSongModal.tsx'
import SongsLoader from '../components/songs/SongsLoader.tsx'
import usePaginationInfo from '../hooks/usePaginationInfo.ts'
import useFixedDocumentTitle from '../hooks/useFixedDocumentTitle.ts'

function Songs(): ReactElement {
  useFixedDocumentTitle('Songs')

  const pageSize = 40
  const [currentPage, setCurrentPage] = useState(1)
  const {
    data: songs,
    isLoading,
    isFetching
  } = useGetSongsQuery({
    pageSize: pageSize,
    currentPage: currentPage,
    orderBy: ['created_at DESC']
  })

  const { startCount, endCount, totalPages } = usePaginationInfo(
    songs?.totalCount,
    pageSize,
    currentPage
  )

  const [openedAddNewSongModal, { open: openAddNewSongModal, close: closeAddNewSongModal }] =
    useDisclosure(false)

  return (
    <Stack h={'100%'} gap={'xs'}>
      <Group gap={4}>
        <Title order={3} fw={800}>
          Songs
        </Title>
        <ActionIcon
          aria-label={'new-song'}
          variant={'grey'}
          size={'lg'}
          onClick={openAddNewSongModal}
        >
          <IconPlus />
        </ActionIcon>
        <Space flex={1} />
        <ActionIcon aria-label={'order-songs'} variant={'grey'} size={'lg'}>
          <IconArrowsSort size={17} />
        </ActionIcon>
        <ActionIcon aria-label={'filter-songs'} variant={'grey'} size={'lg'}>
          <IconFilterFilled size={17} />
        </ActionIcon>
      </Group>
      {!isFetching && (
        <Text inline mb={'xs'}>
          {startCount} - {endCount} songs out of {songs?.totalCount}
        </Text>
      )}

      {songs?.totalCount === 0 && <Text mt={'sm'}>There are no songs yet. Try to add one</Text>}
      <Group gap={'lg'} align={'stretch'}>
        {isLoading && <SongsLoader />}
        {songs?.models.map((song) => <SongCard key={song.id} song={song} />)}
        {songs?.totalCount > 0 && currentPage == totalPages && (
          <Card
            aria-label={'new-song-card'}
            variant={'add-new'}
            radius={'lg'}
            mih={175}
            w={175}
            onClick={openAddNewSongModal}
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
            data-testid={'songs-pagination'}
            value={currentPage}
            onChange={setCurrentPage}
            total={totalPages}
          />
        ) : (
          <Loader size={25} />
        )}
      </Box>

      <AddNewSongModal opened={openedAddNewSongModal} onClose={closeAddNewSongModal} />
    </Stack>
  )
}

export default Songs
