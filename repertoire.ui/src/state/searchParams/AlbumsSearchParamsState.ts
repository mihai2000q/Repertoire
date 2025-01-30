import { SearchParamsState } from '../../hooks/useSearchParamsState.ts'

enum AlbumsSearchParams {
  CurrentPage = 'p'
}

interface AlbumsSearchParamsState {
  currentPage: number
}

const albumsSearchParamsInitialState: AlbumsSearchParamsState = {
  currentPage: 1
}

const albumsSearchParamsSerializer = (state: AlbumsSearchParamsState) => {
  const params = new URLSearchParams()
  params.set(AlbumsSearchParams.CurrentPage, state.currentPage.toString())
  return params
}

const albumsSearchParamsDeserializer = (params: URLSearchParams): AlbumsSearchParamsState => {
  return {
    currentPage: Number(params.get(AlbumsSearchParams.CurrentPage) ?? 1)
  }
}

const albumsSearchParamsState: SearchParamsState<AlbumsSearchParamsState> = {
  initialState: albumsSearchParamsInitialState,
  serialize: albumsSearchParamsSerializer,
  deserialize: albumsSearchParamsDeserializer
}

export default albumsSearchParamsState
