import { SearchParamsState } from '../../hooks/useSearchParamsState.ts'
import Filter from '../../types/Filter.ts'

enum AlbumsSearchParams {
  CurrentPage = 'p',
  ActiveFilters = 'f'
}

interface AlbumsSearchParamsState {
  currentPage: number
  activeFilters: Map<string, Filter>
}

const albumsSearchParamsInitialState: AlbumsSearchParamsState = {
  currentPage: 1,
  activeFilters: new Map<string, Filter>()
}

const albumsSearchParamsSerializer = (state: AlbumsSearchParamsState) => {
  const params = new URLSearchParams()
  params.set(AlbumsSearchParams.CurrentPage, state.currentPage.toString())
  if (state.activeFilters.size > 0) {
    params.set(AlbumsSearchParams.ActiveFilters, JSON.stringify([...state.activeFilters]))
  }
  return params
}

const albumsSearchParamsDeserializer = (params: URLSearchParams): AlbumsSearchParamsState => {
  return {
    currentPage: Number(params.get(AlbumsSearchParams.CurrentPage) ?? 1),
    activeFilters: new Map<string, Filter>([
      ...(params.get(AlbumsSearchParams.ActiveFilters)
        ? JSON.parse(params.get(AlbumsSearchParams.ActiveFilters))
        : [])
    ])
  }
}

const albumsSearchParamsState: SearchParamsState<AlbumsSearchParamsState> = {
  initialState: albumsSearchParamsInitialState,
  serialize: albumsSearchParamsSerializer,
  deserialize: albumsSearchParamsDeserializer
}

export default albumsSearchParamsState
