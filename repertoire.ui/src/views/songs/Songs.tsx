import { ReactElement, useState } from 'react'
import { Box, Button, Group, Loader, Pagination, Space, Stack, Text, Title } from '@mantine/core'
import { useGetSongsQuery } from '../../state/songsApi.ts'
import SongCard from '../../components/songs/SongCard.tsx'
import { IconMusicPlus } from '@tabler/icons-react'
import NewSongCard from '../../components/songs/NewSongCard.tsx'
import { useDisclosure } from '@mantine/hooks'
import AddNewSongModal from '../../components/songs/modal/AddNewSongModal.tsx'
import SongsLoader from '../../components/songs/loader/SongsLoader.tsx'

function Songs(): ReactElement {
  const [currentPage, setCurrentPage] = useState(1)
  const { data: songs, isLoading } = useGetSongsQuery({
    pageSize: 20,
    currentPage: currentPage
  })

  const [openedAddNewSongModal, { open: openAddNewSongModal, close: closeAddNewSongModal }] =
    useDisclosure(false)

  return (
    <Stack h={'100%'}>
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

      {songs?.totalCount === 0 && <Text mt={'sm'}>There are no songs yet. Try to add one</Text>}
      <Group>
        {isLoading && <SongsLoader />}
        {songs?.models.map((song) => <SongCard key={song.id} song={song} />)}
        {songs?.totalCount > 0 && <NewSongCard openModal={openAddNewSongModal} />}
      </Group>

      <Space flex={1} />

      <Box style={{ alignSelf: 'center' }} pb={'xs'}>
        {!isLoading ? (
          <Pagination
            data-testid={'songs-pagination'}
            value={currentPage}
            onChange={setCurrentPage}
            total={songs?.totalCount > 0 ? songs?.totalCount / songs?.models.length : 0}
          />
        ) : (
          <Loader size={25} />
        )}
      </Box>
    </Stack>
  )
}

export default Songs
