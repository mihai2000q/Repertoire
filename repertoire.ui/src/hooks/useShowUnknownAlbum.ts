import { useGetSongsQuery } from '../state/api/songsApi.ts'

export default function useShowUnknownAlbum(): boolean {
  const { data: songs } = useGetSongsQuery({
    currentPage: 1,
    pageSize: 1,
    searchBy: ['album_id IS NULL']
  })

  return songs?.totalCount > 0
}
