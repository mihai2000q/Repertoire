import { SearchParamsState } from '../../hooks/useSearchParamsState.ts'

enum ArtistsSearchParams {
  CurrentPage = 'p'
}

interface ArtistsSearchParamsState {
  currentPage: number
}

const artistsSearchParamsInitialState: ArtistsSearchParamsState = {
  currentPage: 1
}

const artistsSearchParamsSerializer = (state: ArtistsSearchParamsState) => {
  const params = new URLSearchParams()
  params.set(ArtistsSearchParams.CurrentPage, state.currentPage.toString())
  return params
}

const artistsSearchParamsDeserializer = (params: URLSearchParams): ArtistsSearchParamsState => {
  return {
    currentPage: Number(params.get(ArtistsSearchParams.CurrentPage) ?? 1)
  }
}

const artistsSearchParamsState: SearchParamsState<ArtistsSearchParamsState> = {
  initialState: artistsSearchParamsInitialState,
  serialize: artistsSearchParamsSerializer,
  deserialize: artistsSearchParamsDeserializer
}

export default artistsSearchParamsState
