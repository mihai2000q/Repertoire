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
}: SearchParamsState<TState>): [TState, (state: TState | ((state: TState) => TState)) => void] {
  const [searchParams, setSearchParams] = useSearchParams(serialize(initialState))

  const state = deserialize(searchParams)
  const setState = (newState: TState | ((state: TState) => TState)) => {
    if (isFunction(newState)) {
      // eslint-disable-next-line @typescript-eslint/ban-ts-comment
      // @ts-ignore
      setSearchParams(serialize(newState(state)))
    } else {
      // eslint-disable-next-line @typescript-eslint/ban-ts-comment
      // @ts-ignore
      setSearchParams(serialize(newState))
    }
  }

  return [state, setState]
}

function isFunction(functionToCheck: unknown) {
  return functionToCheck && {}.toString.call(functionToCheck) === '[object Function]'
}
