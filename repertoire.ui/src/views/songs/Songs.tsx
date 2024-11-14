import { ReactElement, useState } from 'react'
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
import { useGetSongsQuery } from '../../state/songsApi.ts'
import SongCard from '../../components/songs/SongCard.tsx'
import { IconMusicPlus } from '@tabler/icons-react'
import NewSongCard from '../../components/songs/NewSongCard.tsx'
import { useDisclosure } from '@mantine/hooks'
import AddNewSongModal from '../../components/songs/modal/AddNewSongModal.tsx'
import SongsLoader from '../../components/songs/loader/SongsLoader.tsx'
import usePaginationInfo from '../../hooks/usePaginationInfo.ts'

function Songs(): ReactElement {
  const [currentPage, setCurrentPage] = useState(1)
  const { data: songs, isLoading } = useGetSongsQuery({
    pageSize: 20,
    currentPage: currentPage
  })

  const { startCount, endCount, totalPages } = usePaginationInfo(songs?.totalCount, 20, currentPage)

  const [openedAddNewSongModal, { open: openAddNewSongModal, close: closeAddNewSongModal }] =
    useDisclosure(false)

  return (
    <Stack h={'100%'} gap={'xs'}>
      <AddNewSongModal opened={openedAddNewSongModal} onClose={closeAddNewSongModal} />

      <Title order={3} fw={800}>
        Songs
      </Title>

      <Group>
        <Button
          variant={'gradient'}
          leftSection={<IconMusicPlus size={17} />}
          onClick={openAddNewSongModal}
        >
          New Song
        </Button>
      </Group>
      {!isLoading && (
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
            data-testid={'new-song-card'}
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
        {!isLoading ? (
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
    </Stack>
  )
}

export default Songs
