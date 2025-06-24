import { Divider, Stack } from '@mantine/core'
import { useParams } from 'react-router-dom'
import { useGetPlaylistQuery } from '../state/api/playlistsApi.ts'
import PlaylistLoader from '../components/playlist/loader/PlaylistLoader.tsx'
import PlaylistHeaderCard from '../components/playlist/PlaylistHeaderCard.tsx'
import PlaylistSongsCard from '../components/playlist/PlaylistSongsCard.tsx'
import useDynamicDocumentTitle from '../hooks/useDynamicDocumentTitle.ts'
import { useEffect } from 'react'

function Playlist() {
  const params = useParams()
  const setDocumentTitle = useDynamicDocumentTitle()
  const playlistId = params['id'] ?? ''

  const { data: playlist, isLoading } = useGetPlaylistQuery({ id: playlistId })

  useEffect(() => {
    if (playlist) setDocumentTitle(playlist.title)
  }, [playlist])

  if (isLoading || !playlist) return <PlaylistLoader />

  return (
    <Stack px={'xl'}>
      <PlaylistHeaderCard playlist={playlist} />

      <Divider />

      <PlaylistSongsCard playlistId={playlistId} />
    </Stack>
  )
}

export default Playlist
