import { SearchParamsState } from '../../hooks/useSearchParamsState.ts'
import Filter from '../../types/Filter.ts'

enum SongsSearchParams {
  CurrentPage = 'p',
  ActiveFilters = 'f'
}

interface SongsSearchParamsState {
  currentPage: number
  activeFilters: Map<string, Filter>
}

const songsSearchParamsInitialState: SongsSearchParamsState = {
  currentPage: 1,
  activeFilters: new Map<string, Filter>()
}

const songsSearchParamsSerializer = (state: SongsSearchParamsState) => {
  const params = new URLSearchParams()
  params.set(SongsSearchParams.CurrentPage, state.currentPage.toString())
  if (state.activeFilters.size > 0) {
    params.set(SongsSearchParams.ActiveFilters, JSON.stringify([...state.activeFilters]))
  }
  return params
}

const songsSearchParamsDeserializer = (params: URLSearchParams): SongsSearchParamsState => {
  return {
    currentPage: Number(params.get(SongsSearchParams.CurrentPage) ?? 1),
    activeFilters: new Map<string, Filter>([
      ...(params.get(SongsSearchParams.ActiveFilters)
        ? JSON.parse(params.get(SongsSearchParams.ActiveFilters))
        : [])
    ])
  }
}

const songsSearchParamsState: SearchParamsState<SongsSearchParamsState> = {
  initialState: songsSearchParamsInitialState,
  serialize: songsSearchParamsSerializer,
  deserialize: songsSearchParamsDeserializer
}

export default songsSearchParamsState
