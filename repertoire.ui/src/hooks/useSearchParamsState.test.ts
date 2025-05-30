import useSearchParamsState, { SearchParamsState } from './useSearchParamsState.ts'
import { act } from '@testing-library/react'
import { routerRenderHook } from '../test-utils.tsx'
import { afterEach } from 'vitest'

describe('use Search Params State', () => {
  interface TestSearchParamsState {
    testValue: number
  }

  const testSearchParamsInitialState: TestSearchParamsState = {
    testValue: 1
  }

  const testSearchParamsSerializer = (state: TestSearchParamsState) => {
    const params = new URLSearchParams()
    params.set('testValue', state.testValue.toString())
    return params
  }

  const testSearchParamsDeserializer = (params: URLSearchParams): TestSearchParamsState => {
    return {
      testValue: Number(params.get('testValue') ?? testSearchParamsInitialState.testValue)
    }
  }

  const testProps: SearchParamsState<TestSearchParamsState> = {
    initialState: testSearchParamsInitialState,
    serialize: testSearchParamsSerializer,
    deserialize: testSearchParamsDeserializer
  }

  afterEach(() => window.location.search = '')

  it('should update state with value', () => {
    const { result } = routerRenderHook(() => useSearchParamsState(testProps))
    const [state, setState] = result.current

    expect(state).toStrictEqual(testSearchParamsInitialState)

    const newValue = 23
    act(() => setState({ testValue: newValue }))

    expect(result.current[0].testValue).toBe(newValue)
    expect(window.location.search).toBe(`?testValue=${newValue}`)
  })

  it('should update state with function', () => {
    const { result } = routerRenderHook(() => useSearchParamsState(testProps))
    const [state, setState] = result.current

    expect(state).toStrictEqual(testSearchParamsInitialState)

    const increment = 5
    act(() => setState((prev) => ({ testValue: prev.testValue + increment })))

    const expectedValue = testSearchParamsInitialState.testValue + increment
    expect(result.current[0].testValue).toBe(expectedValue)
    expect(window.location.search).toBe(`?testValue=${expectedValue}`)
  })
})
