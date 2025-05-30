import { Divider, Stack } from '@mantine/core'
import { useParams } from 'react-router-dom'
import { useGetPlaylistQuery } from '../state/api/playlistsApi.ts'
import PlaylistLoader from '../components/playlist/PlaylistLoader.tsx'
import PlaylistHeaderCard from '../components/playlist/PlaylistHeaderCard.tsx'
import PlaylistSongsCard from '../components/playlist/PlaylistSongsCard.tsx'
import useDynamicDocumentTitle from '../hooks/useDynamicDocumentTitle.ts'
import { useEffect } from 'react'
import useLocalStorage from '../hooks/useLocalStorage.ts'
import LocalStorageKeys from '../types/enums/LocalStorageKeys.ts'
import playlistSongsOrders from '../data/playlist/playlistSongsOrders.ts'
import useOrderBy from '../hooks/api/useOrderBy.ts'

function Playlist() {
  const params = useParams()
  const setDocumentTitle = useDynamicDocumentTitle()
  const playlistId = params['id'] ?? ''

  const [order, setOrder] = useLocalStorage({
    key: LocalStorageKeys.PlaylistSongsOrder,
    defaultValue: playlistSongsOrders[0]
  })
  const orderBy = useOrderBy([order])

  const {
    data: playlist,
    isLoading,
    isFetching
  } = useGetPlaylistQuery({ id: playlistId, songsOrderBy: orderBy })

  useEffect(() => {
    if (playlist) setDocumentTitle(playlist.title)
  }, [playlist])

  if (isLoading || !playlist) return <PlaylistLoader />

  return (
    <Stack px={'xl'}>
      <PlaylistHeaderCard playlist={playlist} />

      <Divider />

      <PlaylistSongsCard
        playlist={playlist}
        order={order}
        setOrder={setOrder}
        isFetching={isFetching}
      />
    </Stack>
  )
}

export default Playlist
