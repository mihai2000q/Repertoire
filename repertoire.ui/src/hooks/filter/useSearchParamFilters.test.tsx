import { act, renderHook } from '@testing-library/react'
import useSearchParamFilters from './useSearchParamFilters.ts'
import Filter from '../../types/Filter.ts'
import FilterOperator from '../../types/enums/FilterOperator.ts'

describe('use Search Param Filters', () => {
  const initialFilters: Filter[] = [
    { property: 'name', operator: FilterOperator.Equal, value: '', isSet: false },
    { property: 'age', operator: FilterOperator.GreaterThan, value: 0, isSet: false }
  ]

  const activeFilters = new Map<string, Filter>([
    ['name=', { property: 'name', operator: FilterOperator.Equal, value: 'John', isSet: true }]
  ])

  it('should handle initial state when no active filters exist', () => {
    const { result } = renderHook(() =>
      useSearchParamFilters({
        initialFilters,
        activeFilters: new Map(),
        setSearchParams: vi.fn()
      })
    )

    const [filters] = result.current
    expect(filters.get('name=')).toStrictEqual(initialFilters[0])
    expect(filters.get('age>')).toStrictEqual(initialFilters[1])
  })

  it('should initialize with active filters taking precedence', () => {
    const setSearchParams = vi.fn()

    const { result } = renderHook(() =>
      useSearchParamFilters({ initialFilters, activeFilters, setSearchParams })
    )

    const [filters] = result.current
    expect(filters.get('name=')).toStrictEqual(activeFilters.get('name='))
    expect(filters.get('age>')).toStrictEqual(initialFilters[1])
  })

  it('should re-initialize state when active filters change', () => {
    const { result, rerender } = renderHook(
      ({ activeFilters }) =>
        useSearchParamFilters({
          initialFilters,
          activeFilters: activeFilters,
          setSearchParams: vi.fn()
        }),
      { initialProps: { activeFilters: new Map() } }
    )

    const [filters] = result.current
    expect(filters.get('name=')).toStrictEqual(initialFilters[0])
    expect(filters.get('age>')).toStrictEqual(initialFilters[1])

    rerender({ activeFilters: activeFilters })

    const [rerenderedFilters] = result.current
    expect(rerenderedFilters.get('name=')).toStrictEqual(activeFilters.get('name='))
    expect(rerenderedFilters.get('age>')).toStrictEqual(initialFilters[1])
  })

  it('should update filters and search params when setting new filters', () => {
    const setSearchParams = vi.fn()

    const { result } = renderHook(() =>
      useSearchParamFilters({ initialFilters, activeFilters, setSearchParams })
    )
    const [, setFilters] = result.current

    const newFilters = new Map<string, Filter>([
      ['name=', { ...activeFilters.get('name='), value: 'Alice' }],
      ['age>', { property: 'age', operator: FilterOperator.GreaterThan, value: 25, isSet: true }]
    ])

    act(() => setFilters(newFilters))

    // Should update local state
    const [updatedFilters] = result.current
    expect(updatedFilters.get('name=').value).toBe('Alice')
    expect(updatedFilters.get('age>').value).toBe(25)

    // Should update search params with only active filters
    expect(setSearchParams).toHaveBeenCalledOnce()
    const searchParamsUpdate = setSearchParams.mock.calls[0][0]
    const updatedSearchParams = searchParamsUpdate({})
    expect(updatedSearchParams.activeFilters.size).toBe(2)
    expect(updatedSearchParams.activeFilters.get('name=').value).toBe('Alice')
    expect(updatedSearchParams.activeFilters.get('age>').value).toBe(25)
  })

  it('should update current page on search params when there is one', () => {
    const setSearchParams = vi.fn()

    const { result } = renderHook(() =>
      useSearchParamFilters({ initialFilters, activeFilters, setSearchParams })
    )
    const [, setFilters] = result.current

    const newFilters = new Map<string, Filter>([
      ['name=', { ...activeFilters.get('name='), value: 'Alice' }],
      ['age>', { property: 'age', operator: FilterOperator.GreaterThan, value: 25, isSet: true }]
    ])

    act(() => setFilters(newFilters))

    expect(setSearchParams).toHaveBeenCalledOnce()
    const searchParamsUpdate = setSearchParams.mock.calls[0][0]
    const updatedSearchParams = searchParamsUpdate({})
    expect(updatedSearchParams.activeFilters.size).toBe(2)
    expect(updatedSearchParams.activeFilters.get('name=').value).toBe('Alice')
    expect(updatedSearchParams.activeFilters.get('age>').value).toBe(25)
  })

  it('should update only filters when setting new filters and search params is disabled', () => {
    const setSearchParams = vi.fn()

    const { result } = renderHook(() =>
      useSearchParamFilters({ initialFilters, activeFilters, setSearchParams })
    )
    const [, setFilters] = result.current

    const newFilters = new Map<string, Filter>([
      ['name=', { ...activeFilters.get('name='), value: 'Alice' }],
      ['age>', { property: 'age', operator: FilterOperator.GreaterThan, value: 25, isSet: true }]
    ])

    act(() => setFilters(newFilters, false))

    const [updatedFilters] = result.current
    expect(updatedFilters.get('name=').value).toBe('Alice')
    expect(updatedFilters.get('age>').value).toBe(25)

    expect(setSearchParams).not.toHaveBeenCalled()
  })

  it('should only include filters isSet true in search params', () => {
    const setSearchParams = vi.fn()

    const { result } = renderHook(() =>
      useSearchParamFilters({ initialFilters, activeFilters, setSearchParams })
    )
    const [, setFilters] = result.current

    const newFilters = new Map<string, Filter>([
      ['name=', { property: 'name', operator: FilterOperator.Equal, value: 'Alice', isSet: false }],
      ['age>', { property: 'age', operator: FilterOperator.GreaterThan, value: 25, isSet: true }]
    ])

    act(() => setFilters(newFilters))

    const searchParamsUpdate = setSearchParams.mock.calls[0][0]
    const updatedSearchParams = searchParamsUpdate({ activeFilters: new Map() })
    expect(updatedSearchParams.activeFilters.size).toBe(1)
    expect(updatedSearchParams.activeFilters.has('age>')).toBeTruthy()
    expect(updatedSearchParams.activeFilters.has('name=')).toBeFalsy()
  })
})
