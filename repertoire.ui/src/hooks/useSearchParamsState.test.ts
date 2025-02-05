import useSearchParamsState, { SearchParamsState } from './useSearchParamsState.ts'
import { act } from '@testing-library/react'
import { routerRenderHook } from '../test-utils.tsx'

describe('use Search Params State', () => {
  it('should render hook', () => {
    // setup
    const initialTestValue = 1

    interface TestSearchParamsState {
      testValue: number
    }

    const testSearchParamsInitialState: TestSearchParamsState = {
      testValue: initialTestValue
    }

    const testSearchParamsSerializer = (state: TestSearchParamsState) => {
      const params = new URLSearchParams()
      params.set('testValue', state.testValue.toString())
      return params
    }

    const testSearchParamsDeserializer = (params: URLSearchParams): TestSearchParamsState => {
      return {
        testValue: Number(params.get('testValue') ?? 1)
      }
    }

    const songsSearchParamsState: SearchParamsState<TestSearchParamsState> = {
      initialState: testSearchParamsInitialState,
      serialize: testSearchParamsSerializer,
      deserialize: testSearchParamsDeserializer
    }

    const { result } = routerRenderHook(() => useSearchParamsState(songsSearchParamsState))
    const [state, setState] = result.current

    expect(state).toEqual(testSearchParamsInitialState)

    const newTestValue = 23

    act(() => setState({ testValue: newTestValue }))

    expect(result.current[0].testValue).toBe(newTestValue)
    expect(window.location.search).toBe('?testValue=' + newTestValue)
  })
})
