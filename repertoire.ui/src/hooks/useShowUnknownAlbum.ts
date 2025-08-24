import { useGetSongsQuery } from '../state/api/songsApi.ts'
import useSearchBy from './api/useSearchBy.ts'
import SongProperty from '../types/enums/properties/SongProperty.ts'
import FilterOperator from '../types/enums/FilterOperator.ts'

export default function useShowUnknownAlbum(): boolean {
  const searchBy = useSearchBy([
    { property: SongProperty.AlbumId, operator: FilterOperator.IsNull }
  ])
  const { data: songs } = useGetSongsQuery({
    currentPage: 1,
    pageSize: 1,
    searchBy: searchBy
  })

  return songs?.totalCount > 0
}
