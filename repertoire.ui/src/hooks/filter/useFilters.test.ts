import { act, renderHook } from '@testing-library/react'
import useFilters from './useFilters.ts'
import Filter from '../../types/Filter.ts'
import FilterOperator from '../../types/enums/FilterOperator.ts'

describe('use Filters', () => {
  it('should initialize with provided filters', () => {
    const initialFilters: Filter[] = [
      { property: 'name', operator: FilterOperator.Equal, value: '', isSet: false },
      { property: 'age', operator: FilterOperator.GreaterThan, value: 0, isSet: true }
    ]

    const { result } = renderHook(() => useFilters(initialFilters))

    const [filters, activeCount, _] = result.current
    expect(filters.size).toBe(2)
    expect(filters.get('name=')).toEqual(initialFilters[0])
    expect(filters.get('age>')).toEqual(initialFilters[1])
    expect(activeCount).toBe(1)
  })

  it('should correctly count active filters', () => {
    const initialFilters: Filter[] = [
      { property: 'name', operator: FilterOperator.Equal, value: 'John', isSet: true },
      { property: 'age', operator: FilterOperator.GreaterThan, value: 18, isSet: true },
      { property: 'email', operator: FilterOperator.LessThan, value: '', isSet: false }
    ]

    const { result } = renderHook(() => useFilters(initialFilters))

    const [_, activeCount] = result.current
    expect(activeCount).toBe(2)
  })

  it('should update filters and active count when setFilters is called', () => {
    const { result } = renderHook(() => useFilters())
    const [_, __, setFilters] = result.current

    act(() => {
      const newFilters = new Map<string, Filter>()
      newFilters.set('name=', {
        property: 'name',
        operator: FilterOperator.Equal,
        value: 'Alice',
        isSet: true
      })
      setFilters(newFilters)
    })

    const [filters, activeCount] = result.current
    expect(filters.size).toBe(1)
    expect(activeCount).toBe(1)
    expect(filters.get('name=')).toEqual({
      property: 'name',
      operator: FilterOperator.Equal,
      value: 'Alice',
      isSet: true
    })
  })

  it('should handle duplicate filter keys by overwriting', () => {
    const initialFilters: Filter[] = [
      { property: 'name', operator: FilterOperator.Equal, value: 'John', isSet: false }
    ]

    const { result } = renderHook(() => useFilters(initialFilters))
    const [_, __, setFilters] = result.current

    act(() => {
      const newFilters = new Map<string, Filter>()
      newFilters.set('name=', {
        property: 'name',
        operator: FilterOperator.Equal,
        value: 'Alice',
        isSet: true
      })
      setFilters(newFilters)
    })

    const [filters] = result.current
    expect(filters.size).toBe(1)
    expect(filters.get('name=')?.value).toBe('Alice')
    expect(filters.get('name=')?.isSet).toBe(true)
  })

  it('should initialize with empty map when no initial filters provided', () => {
    const { result } = renderHook(() => useFilters())

    const [filters, activeCount, _] = result.current
    expect(filters.size).toBe(0)
    expect(activeCount).toBe(0)
  })

  it('should handle empty initial filters array', () => {
    const { result } = renderHook(() => useFilters([]))

    const [filters, activeCount] = result.current
    expect(filters.size).toBe(0)
    expect(activeCount).toBe(0)
  })
})
