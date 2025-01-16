import { Divider, Stack } from '@mantine/core'
import { useParams } from 'react-router-dom'
import { useGetPlaylistQuery } from '../state/playlistsApi.ts'
import PlaylistLoader from '../components/playlist/PlaylistLoader.tsx'
import PlaylistHeaderCard from '../components/playlist/PlaylistHeaderCard.tsx'
import PlaylistSongsCard from '../components/playlist/PlaylistSongsCard.tsx'

function Playlist() {
  const params = useParams()
  const playlistId = params['id'] ?? ''

  const { data: playlist, isLoading, isFetching } = useGetPlaylistQuery(playlistId)

  if (isLoading) return <PlaylistLoader />

  return (
    <Stack>
      <PlaylistHeaderCard playlist={playlist} />

      <Divider />

      <PlaylistSongsCard playlist={playlist} isFetching={isFetching} />
    </Stack>
  )
}

export default Playlist
