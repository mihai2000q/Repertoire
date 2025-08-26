import { Divider, Stack } from '@mantine/core'
import { useParams } from 'react-router-dom'
import { useGetAlbumQuery } from '../state/api/albumsApi.ts'
import AlbumLoader from '../components/album/AlbumLoader.tsx'
import { useGetSongsQuery } from '../state/api/songsApi.ts'
import AlbumHeader from '../components/album/AlbumHeader.tsx'
import AlbumSongsWidget from '../components/album/AlbumSongsWidget.tsx'
import useDynamicDocumentTitle from '../hooks/useDynamicDocumentTitle.ts'
import { useEffect } from 'react'
import albumSongsOrders from '../data/album/albumSongsOrders.ts'
import LocalStorageKeys from '../types/enums/keys/LocalStorageKeys.ts'
import useOrderBy from '../hooks/api/useOrderBy.ts'
import useLocalStorage from '../hooks/useLocalStorage.ts'
import useSearchBy from '../hooks/api/useSearchBy.ts'
import SongProperty from '../types/enums/properties/SongProperty.ts'
import FilterOperator from '../types/enums/FilterOperator.ts'

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
  const orderBy = useOrderBy([order])

  const {
    data: album,
    isLoading,
    isFetching
  } = useGetAlbumQuery(
    {
      id: albumId,
      songsOrderBy: orderBy
    },
    { skip: isUnknownAlbum }
  )

  useEffect(() => {
    if (isUnknownAlbum) setDocumentTitle('Unknown Album')
    else if (album) setDocumentTitle(album.title)
  }, [album, isUnknownAlbum])

  const songsSearchBy = useSearchBy([
    { property: SongProperty.AlbumId, operator: FilterOperator.IsNull }
  ])
  const {
    data: songs,
    isLoading: isSongsLoading,
    isFetching: isSongsFetching
  } = useGetSongsQuery(
    {
      orderBy: orderBy,
      searchBy: songsSearchBy
    },
    { skip: !isUnknownAlbum }
  )

  if (isLoading || isSongsLoading || (!album && !isUnknownAlbum) || (!songs && isUnknownAlbum))
    return <AlbumLoader />

  return (
    <Stack px={'xl'}>
      <AlbumHeader
        album={album}
        isUnknownAlbum={isUnknownAlbum}
        songsTotalCount={songs?.totalCount}
      />

      <Divider />

      <AlbumSongsWidget
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
