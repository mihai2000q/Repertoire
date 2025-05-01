import { renderHook } from '@testing-library/react'
import Filter from '../../types/Filter.ts'
import FilterOperator from '../../types/enums/FilterOperator.ts'
import useSearchBy from './useSearchBy.ts'

describe('use Search By', () => {
  it('should return empty array when no filters are set', () => {
    const filters = new Map<string, Filter>([
      ['name=', { property: 'name', operator: FilterOperator.Equal, isSet: false }],
      ['age>', { property: 'age', operator: FilterOperator.GreaterThan, isSet: false }]
    ])

    const { result } = renderHook(() => useSearchBy(filters))
    expect(result.current).toStrictEqual([])
  })

  it('should construct filters correctly for set filters', () => {
    const filters = new Map<string, Filter>([
      ['name=', { property: 'name', operator: FilterOperator.Equal, value: 'John', isSet: true }],
      ['age>', { property: 'age', operator: FilterOperator.GreaterThan, value: 18, isSet: true }],
      [
        'age<=',
        { property: 'age', operator: FilterOperator.LessThanOrEqual, value: 65 } // if not set, then don't ignore
      ],
      [
        'email<',
        { property: 'email', operator: FilterOperator.LessThan, value: 'example', isSet: false }
      ]
    ])

    const { result } = renderHook(() => useSearchBy(filters))
    expect(result.current).toStrictEqual(['name = John', 'age > 18', 'age <= 65'])
  })

  it('should handle IsNull and IsNotNull operators without values', () => {
    const filters = new Map<string, Filter>([
      ['deletedISNULL', { property: 'deleted', operator: FilterOperator.IsNull, isSet: true }],
      ['activeISNOTNULL', { property: 'active', operator: FilterOperator.IsNotNull, isSet: true }]
    ])

    const { result } = renderHook(() => useSearchBy(filters))
    expect(result.current).toStrictEqual(['deleted IS NULL', 'active IS NOT NULL'])
  })

  it('should update filters when the input changes', () => {
    const initialFilters = new Map<string, Filter>([
      [
        'name=',
        { property: 'name', operator: FilterOperator.Equal, value: 'Michael', isSet: true }
      ],
      [
        'status=',
        { property: 'name', operator: FilterOperator.Equal, value: 'not active', isSet: false }
      ]
    ])

    const { result, rerender } = renderHook(({ filters }) => useSearchBy(filters), {
      initialProps: { filters: initialFilters }
    })

    expect(result.current).toStrictEqual(['name = Michael'])

    const updatedFilters = new Map<string, Filter>([
      ['name=', { property: 'name', operator: FilterOperator.Equal, value: 'John', isSet: true }],
      [
        'status=',
        { property: 'status', operator: FilterOperator.Equal, value: 'active', isSet: true }
      ]
    ])

    rerender({ filters: updatedFilters })

    expect(result.current).toStrictEqual(['name = John', 'status = active'])
  })

  it('should handle empty filters map', () => {
    const filters = new Map<string, Filter>()
    const { result } = renderHook(() => useSearchBy(filters))
    expect(result.current).toEqual([])
  })
})
