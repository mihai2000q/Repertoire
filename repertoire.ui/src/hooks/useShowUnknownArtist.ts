import { useGetSongsQuery } from '../state/songsApi.ts'
import { useGetAlbumsQuery } from '../state/albumsApi.ts'

export default function useShowUnknownArtist(): boolean {
  const { data: songs } = useGetSongsQuery({
    currentPage: 1,
    pageSize: 1,
    searchBy: ['artist_id IS NULL']
  })

  const { data: albums } = useGetAlbumsQuery({
    currentPage: 1,
    pageSize: 1,
    searchBy: ['artist_id IS NULL']
  })

  return songs?.totalCount > 0 || albums?.totalCount > 0
}
