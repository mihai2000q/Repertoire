import { useSearchParams } from 'react-router-dom'

export interface SearchParamsState<TState> {
  initialState: TState
  serialize: (state: TState) => URLSearchParams
  deserialize: (params: URLSearchParams) => TState
}

export default function useSearchParamsState<TState>({
  initialState,
  serialize,
  deserialize
}: SearchParamsState<TState>): [TState, (state: TState) => void] {
  const [searchParams, setSearchParams] = useSearchParams(serialize(initialState))

  const state = deserialize(searchParams)
  const setState = (state: TState) => {
    setSearchParams(serialize(state))
  }

  return [state, setState]
}
