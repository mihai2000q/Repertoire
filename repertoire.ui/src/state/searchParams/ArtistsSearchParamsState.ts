import { SearchParamsState } from '../../hooks/useSearchParamsState.ts'
import Filter from '../../types/Filter.ts'

enum ArtistsSearchParams {
  CurrentPage = 'p',
  ActiveFilters = 'f'
}

interface ArtistsSearchParamsState {
  currentPage: number
  activeFilters: Map<string, Filter>
}

const artistsSearchParamsInitialState: ArtistsSearchParamsState = {
  currentPage: 1,
  activeFilters: new Map<string, Filter>()
}

const artistsSearchParamsSerializer = (state: ArtistsSearchParamsState) => {
  const params = new URLSearchParams()
  params.set(ArtistsSearchParams.CurrentPage, state.currentPage.toString())
  if (state.activeFilters.size > 0) {
    params.set(ArtistsSearchParams.ActiveFilters, JSON.stringify([...state.activeFilters]))
  }
  return params
}

const artistsSearchParamsDeserializer = (params: URLSearchParams): ArtistsSearchParamsState => {
  return {
    currentPage: Number(params.get(ArtistsSearchParams.CurrentPage) ?? 1),
    activeFilters: new Map<string, Filter>([
      ...(params.get(ArtistsSearchParams.ActiveFilters)
        ? JSON.parse(params.get(ArtistsSearchParams.ActiveFilters))
        : [])
    ])
  }
}

const artistsSearchParamsState: SearchParamsState<ArtistsSearchParamsState> = {
  initialState: artistsSearchParamsInitialState,
  serialize: artistsSearchParamsSerializer,
  deserialize: artistsSearchParamsDeserializer
}

export default artistsSearchParamsState
