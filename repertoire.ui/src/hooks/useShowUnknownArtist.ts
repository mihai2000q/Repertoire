import { useGetSongsQuery } from '../state/api/songsApi.ts'
import { useGetAlbumsQuery } from '../state/api/albumsApi.ts'
import useSearchBy from './api/useSearchBy.ts'
import SongProperty from '../types/enums/properties/SongProperty.ts'
import FilterOperator from '../types/enums/FilterOperator.ts'
import AlbumProperty from '../types/enums/properties/AlbumProperty.ts'

export default function useShowUnknownArtist(): boolean {
  const songsSearchBy = useSearchBy([
    { property: SongProperty.ArtistId, operator: FilterOperator.IsNull }
  ])
  const { data: songs } = useGetSongsQuery({
    currentPage: 1,
    pageSize: 1,
    searchBy: songsSearchBy
  })

  const albumsSearchBy = useSearchBy([
    { property: AlbumProperty.ArtistId, operator: FilterOperator.IsNull }
  ])
  const { data: albums } = useGetAlbumsQuery({
    currentPage: 1,
    pageSize: 1,
    searchBy: albumsSearchBy
  })

  return songs?.totalCount > 0 || albums?.totalCount > 0
}
