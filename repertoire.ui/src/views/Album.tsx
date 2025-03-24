import { Divider, Stack } from '@mantine/core'
import { useParams } from 'react-router-dom'
import { useGetAlbumQuery } from '../state/api/albumsApi.ts'
import AlbumLoader from '../components/album/AlbumLoader.tsx'
import { useGetSongsQuery } from '../state/api/songsApi.ts'
import AlbumHeaderCard from '../components/album/AlbumHeaderCard.tsx'
import AlbumSongsCard from '../components/album/AlbumSongsCard.tsx'
import useDynamicDocumentTitle from '../hooks/useDynamicDocumentTitle.ts'
import { useEffect } from 'react'
import albumSongsOrders from '../data/album/albumSongsOrders.ts'
import { useLocalStorage } from '@mantine/hooks'
import LocalStorageKeys from '../utils/enums/LocalStorageKeys.ts'

function Album() {
  const params = useParams()
  const setDocumentTitle = useDynamicDocumentTitle()
  const albumId = params['id'] ?? ''

  const isUnknownAlbum = albumId === 'unknown'

  const [order, setOrder] = useLocalStorage({
    key: isUnknownAlbum
      ? LocalStorageKeys.UnknownAlbumSongsOrder
      : LocalStorageKeys.AlbumSongsOrder,
    defaultValue: isUnknownAlbum ? albumSongsOrders[1] : albumSongsOrders[0]
  })

  const {
    data: album,
    isLoading,
    isFetching
  } = useGetAlbumQuery(
    {
      id: albumId,
      songsOrderBy: [order.value]
    },
    { skip: isUnknownAlbum }
  )

  useEffect(() => {
    if (isUnknownAlbum) setDocumentTitle('Unknown Album')
    else if (album) setDocumentTitle(album.title)
  }, [album, isUnknownAlbum])

  const {
    data: songs,
    isLoading: isSongsLoading,
    isFetching: isSongsFetching
  } = useGetSongsQuery(
    {
      orderBy: [order.value],
      searchBy: ['album_id IS NULL']
    },
    { skip: !isUnknownAlbum }
  )

  if (isLoading || isSongsLoading || (!album && !isUnknownAlbum) || (!songs && isUnknownAlbum)) return <AlbumLoader />

  return (
    <Stack px={'xl'}>
      <AlbumHeaderCard
        album={album}
        isUnknownAlbum={isUnknownAlbum}
        songsTotalCount={songs?.totalCount}
      />

      <Divider />

      <AlbumSongsCard
        album={album}
        songs={songs?.models}
        isUnknownAlbum={isUnknownAlbum}
        order={order}
        setOrder={setOrder}
        isFetching={isSongsFetching || isFetching}
      />
    </Stack>
  )
}

export default Album
