import { ReactElement, useState } from 'react'
import { Box, Group, Loader, Pagination, Space, Stack, Title } from '@mantine/core'
import { useGetSongsQuery } from '../state/api'
import SongCard from '../components/SongCard'

function SongsView(): ReactElement {
  const [currentPage, setCurrentPage] = useState(1)
  const { data: songs, isLoading } = useGetSongsQuery({
    pageSize: 20,
    currentPage: currentPage
  })

  return (
    <Stack h={'100%'}>
      <Title order={3}>Songs</Title>
      <Group>{songs?.models.map((song) => <SongCard key={song.id} song={song} />)}</Group>
      <Space flex={1} />
      <Box style={{ alignSelf: 'center' }}>
        {!isLoading ? (
          <Pagination
            value={currentPage}
            onChange={setCurrentPage}
            total={songs?.totalCount / songs?.models.length}
          />
        ) : (
          <Loader size={25} />
        )}
      </Box>
    </Stack>
  )
}

export default SongsView
