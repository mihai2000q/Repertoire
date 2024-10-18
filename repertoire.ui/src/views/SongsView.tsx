import { ReactElement } from 'react'
import {Group, Pagination, Stack, Title} from "@mantine/core";
import {useGetSongsQuery} from "../state/api";
import SongCard from "../components/SongCard";

function SongsView(): ReactElement {

  const data = useGetSongsQuery({})?.data

  return (
    <Stack h={'100%'}>
      <Title order={3}>Songs</Title>
      <Group>
        {data?.songs.map(song => (
          <SongCard key={song.id} song={song} />
        ))}
      </Group>
      <Pagination total={data?.totalCount / data?.songs.length} />
    </Stack>
  )
}

export default SongsView
