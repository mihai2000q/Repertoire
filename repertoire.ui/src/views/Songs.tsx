import { ReactElement, useState } from 'react'
import { Box, Button, Group, Loader, Pagination, Stack, Text, Title } from '@mantine/core'
import { useGetSongsQuery } from '../state/songsApi'
import SongCard from '../components/songs/SongCard'
import { IconMusicPlus } from '@tabler/icons-react'
import NewSongCard from '../components/songs/NewSongCard'
import { useDisclosure } from '@mantine/hooks'
import AddNewSongModal from '../components/songs/modal/AddNewSongModal'

function Songs(): ReactElement {
  const [currentPage, setCurrentPage] = useState(1)
  const { data: songs, isLoading } = useGetSongsQuery({
    pageSize: 20,
    currentPage: currentPage
  })

  const [opened, { open, close }] = useDisclosure(false)

  return (
    <Stack h={'100%'}>
      <AddNewSongModal opened={opened} onClose={close} />

      <Title order={3} fw={800}>
        Songs
      </Title>

      <Group>
        <Button leftSection={<IconMusicPlus size={17} />} onClick={open}>
          New Song
        </Button>
      </Group>

      {songs?.totalCount === 0 && <Text mt={'sm'}>There are no songs yet. Try to add one</Text>}
      <Group flex={1} h={'100%'}>
        {songs?.models.map((song) => <SongCard key={song.id} song={song} />)}
        {songs?.totalCount > 0 && <NewSongCard openModal={open} />}
      </Group>

      <Box style={{ alignSelf: 'center' }}>
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
