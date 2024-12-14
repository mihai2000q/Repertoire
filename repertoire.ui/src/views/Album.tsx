import { Divider, Stack } from '@mantine/core'
import { useParams } from 'react-router-dom'
import { useGetAlbumQuery } from '../state/albumsApi.ts'
import AlbumLoader from '../components/album/AlbumLoader.tsx'
import { useGetSongsQuery } from '../state/songsApi.ts'
import AlbumHeaderCard from '../components/album/AlbumHeaderCard.tsx'
import AlbumSongsCard from '../components/album/AlbumSongsCard.tsx'

function Album() {
  const params = useParams()
  const albumId = params['id'] ?? ''

  const isUnknownAlbum = albumId === 'unknown'

  const { data: album, isLoading, isFetching } = useGetAlbumQuery(albumId, { skip: isUnknownAlbum })

  const {
    data: songs,
    isLoading: isSongsLoading,
    isFetching: isSongsFetching
  } = useGetSongsQuery(
    {
      orderBy: ['title'],
      searchBy: ['album_id IS NULL']
    },
    { skip: !isUnknownAlbum }
  )

  if (isLoading || isSongsLoading) return <AlbumLoader />

  return (
    <Stack>
      <AlbumHeaderCard
        album={album}
        isUnknownAlbum={isUnknownAlbum}
        songsTotalCount={songs?.totalCount}
      />

      <Divider />

      <AlbumSongsCard
        album={album}
        songs={songs?.models}
        isFetching={isSongsFetching || isFetching}
        isUnknownAlbum={isUnknownAlbum}
      />
    </Stack>
  )
}

export default Album
