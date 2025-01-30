import { SearchParamsState } from '../../hooks/useSearchParamsState.ts'

enum SongsSearchParams {
  CurrentPage = 'p'
}

interface SongsSearchParamsState {
  currentPage: number
}

const songsSearchParamsInitialState: SongsSearchParamsState = {
  currentPage: 1
}

const songsSearchParamsSerializer = (state: SongsSearchParamsState) => {
  const params = new URLSearchParams()
  params.set(SongsSearchParams.CurrentPage, state.currentPage.toString())
  return params
}

const songsSearchParamsDeserializer = (params: URLSearchParams): SongsSearchParamsState => {
  return {
    currentPage: Number(params.get(SongsSearchParams.CurrentPage) ?? 1)
  }
}

const songsSearchParamsState: SearchParamsState<SongsSearchParamsState> = {
  initialState: songsSearchParamsInitialState,
  serialize: songsSearchParamsSerializer,
  deserialize: songsSearchParamsDeserializer
}

export default songsSearchParamsState
