import { SearchParamsState } from '../../hooks/useSearchParamsState.ts'
import Filter from '../../types/Filter.ts'

enum PlaylistsSearchParams {
  CurrentPage = 'p',
  ActiveFilters = 'f'
}

interface PlaylistsSearchParamsState {
  currentPage: number
  activeFilters: Map<string, Filter>
}

const playlistsSearchParamsInitialState: PlaylistsSearchParamsState = {
  currentPage: 1,
  activeFilters: new Map<string, Filter>()
}

const playlistsSearchParamsSerializer = (state: PlaylistsSearchParamsState) => {
  const params = new URLSearchParams()
  params.set(PlaylistsSearchParams.CurrentPage, state.currentPage.toString())
  if (state.activeFilters.size > 0) {
    params.set(PlaylistsSearchParams.ActiveFilters, JSON.stringify([...state.activeFilters]))
  }
  return params
}

const playlistsSearchParamsDeserializer = (params: URLSearchParams): PlaylistsSearchParamsState => {
  return {
    currentPage: Number(params.get(PlaylistsSearchParams.CurrentPage) ?? 1),
    activeFilters: new Map<string, Filter>([
      ...(params.get(PlaylistsSearchParams.ActiveFilters)
        ? JSON.parse(params.get(PlaylistsSearchParams.ActiveFilters))
        : [])
    ])
  }
}

const playlistsSearchParamsState: SearchParamsState<PlaylistsSearchParamsState> = {
  initialState: playlistsSearchParamsInitialState,
  serialize: playlistsSearchParamsSerializer,
  deserialize: playlistsSearchParamsDeserializer
}

export default playlistsSearchParamsState
