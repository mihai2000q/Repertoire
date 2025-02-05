import { useGetSongsQuery } from '../state/api/songsApi.ts'
import { useGetAlbumsQuery } from '../state/api/albumsApi.ts'

export default function useShowUnknownArtist(): boolean {
  const { data: songs } = useGetSongsQuery({
    currentPage: 1,
    pageSize: 1,
    searchBy: ['songs.artist_id IS NULL']
  })

  const { data: albums } = useGetAlbumsQuery({
    currentPage: 1,
    pageSize: 1,
    searchBy: ['artist_id IS NULL']
  })

  return songs?.totalCount > 0 || albums?.totalCount > 0
}
