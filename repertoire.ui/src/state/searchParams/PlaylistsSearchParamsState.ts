import { SearchParamsState } from '../../hooks/useSearchParamsState.ts'

enum PlaylistsSearchParams {
  CurrentPage = 'p'
}

interface PlaylistsSearchParamsState {
  currentPage: number
}

const playlistsSearchParamsInitialState: PlaylistsSearchParamsState = {
  currentPage: 1
}

const playlistsSearchParamsSerializer = (state: PlaylistsSearchParamsState) => {
  const params = new URLSearchParams()
  params.set(PlaylistsSearchParams.CurrentPage, state.currentPage.toString())
  return params
}

const playlistsSearchParamsDeserializer = (params: URLSearchParams): PlaylistsSearchParamsState => {
  return {
    currentPage: Number(params.get(PlaylistsSearchParams.CurrentPage) ?? 1)
  }
}

const playlistsSearchParamsState: SearchParamsState<PlaylistsSearchParamsState> = {
  initialState: playlistsSearchParamsInitialState,
  serialize: playlistsSearchParamsSerializer,
  deserialize: playlistsSearchParamsDeserializer
}

export default playlistsSearchParamsState
