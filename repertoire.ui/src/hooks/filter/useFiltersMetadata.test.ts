import { renderHook } from '@testing-library/react'
import useFiltersMetadata from './useFiltersMetadata.ts'
import Filter, { FilterValue } from '../../types/Filter.ts'
import FilterOperator from '../../types/enums/FilterOperator.ts'

describe('use Filters Metadata', () => {
  const filters = new Map<string, Filter>([
    ['name=', { property: 'name', operator: FilterOperator.Equal, value: '', isSet: false }],
    ['age>', { property: 'age', operator: FilterOperator.GreaterThan, value: 0, isSet: false }],
    ['age<', { property: 'age', operator: FilterOperator.LessThanOrEqual, value: 80, isSet: true }],
    ['allergies_null', { property: 'allergies', operator: FilterOperator.IsNull, isSet: true }]
  ])

  it('should not update filters when both metadata are null', () => {
    const setFilters = vi.fn()
    const setInternalFilters = vi.fn()

    const filtersMetadataMap = vi.fn()
    const { result } = renderHook(() =>
      useFiltersMetadata(null, null, filters, setFilters, setInternalFilters, filtersMetadataMap)
    )

    expect(filtersMetadataMap).not.toHaveBeenCalled()
    expect(setFilters).not.toHaveBeenCalled()
    expect(setInternalFilters).not.toHaveBeenCalled()
    expect(result.current).toEqual(filters)
  })

  it('should update initial filters when initial metadata changes', () => {
    const initialMetadata = { some: 'initial-metadata' }

    const metadataMap: [string, FilterValue][] = [['name=', 'John']]
    const filtersMetadataMap = vi.fn().mockReturnValue(metadataMap)

    const { result } = renderHook(() =>
      useFiltersMetadata(initialMetadata, null, filters, vi.fn(), vi.fn(), filtersMetadataMap)
    )

    expect(filtersMetadataMap).toHaveBeenCalledOnce()
    expect(filtersMetadataMap).toHaveBeenCalledWith(initialMetadata)

    // initial filters
    expect(result.current.size).toBe(4)
    expect(result.current.get('name=')).toStrictEqual({
      ...filters.get('name='),
      value: 'John'
    })
    expect(result.current.get('age>')).toStrictEqual({
      ...filters.get('age>')
    })
    expect(result.current.get('age<')).toStrictEqual({
      ...filters.get('age<'),
      isSet: false
    })
    expect(result.current.get('allergies_null')).toStrictEqual({
      ...filters.get('allergies_null'),
      isSet: false
    })
  })

  it('should update filters and internal filters, when metadata changes', () => {
    const metadata = { some: 'metadata' }

    const setFilters = vi.fn()
    const setInternalFilters = vi.fn()
    const metadataMap: [string, FilterValue][] = [['name=', 'John']]
    const filtersMetadataMap = vi.fn().mockReturnValue(metadataMap)

    renderHook(() =>
      useFiltersMetadata(
        null,
        metadata,
        filters,
        setFilters,
        setInternalFilters,
        filtersMetadataMap
      )
    )

    expect(filtersMetadataMap).toHaveBeenCalledOnce()
    expect(filtersMetadataMap).toHaveBeenCalledWith(metadata)

    // filters
    expect(setFilters).toHaveBeenCalledOnce()
    expect(setInternalFilters).toHaveBeenCalledOnce()

    const expectedFilters = setFilters.mock.calls[0][0]
    const expectedWithSearchParams = setFilters.mock.calls[0][1]
    const internalFilters = setInternalFilters.mock.calls[0][0]
    expect(internalFilters).toEqual(expectedFilters)

    expect(expectedFilters.size).toBe(4)
    expect(expectedFilters.get('name=')).toStrictEqual({
      ...filters.get('name='),
      value: 'John'
    })
    expect(expectedFilters.get('age>')).toStrictEqual(filters.get('age>'))
    expect(expectedFilters.get('age<')).toStrictEqual(filters.get('age<'))
    expect(expectedFilters.get('allergies_null')).toStrictEqual(filters.get('allergies_null'))
    expect(expectedWithSearchParams).toBeFalsy()
  })

  it('should update initial filters, filters and internal filters, when metadata and initial metadata change', () => {
    const initialMetadata = { some: 'initial-metadata' }
    const metadata = { some: 'metadata' }

    const setFilters = vi.fn()
    const setInternalFilters = vi.fn()
    const initialMetadataMap: [string, FilterValue][] = [['name=', 'J']]
    const metadataMap: [string, FilterValue][] = [['name=', 'John']]
    const filtersMetadataMap = vi.fn((args) => {
      if (args === initialMetadata) return initialMetadataMap
      if (args === metadata) return metadataMap
    })

    const { result } = renderHook(() =>
      useFiltersMetadata(
        initialMetadata,
        metadata,
        filters,
        setFilters,
        setInternalFilters,
        filtersMetadataMap
      )
    )

    expect(filtersMetadataMap).toHaveBeenCalledTimes(2)
    expect(filtersMetadataMap).toHaveBeenCalledWith(initialMetadata)
    expect(filtersMetadataMap).toHaveBeenCalledWith(metadata)

    // initial filters
    expect(result.current.size).toBe(4)
    expect(result.current.get('name=')).toStrictEqual({
      ...filters.get('name='),
      value: 'J'
    })

    // filters
    expect(setFilters).toHaveBeenCalledOnce()
    expect(setInternalFilters).toHaveBeenCalledOnce()

    const expectedFilters = setFilters.mock.calls[0][0]
    const internalFilters = setInternalFilters.mock.calls[0][0]
    expect(internalFilters).toEqual(expectedFilters)

    expect(expectedFilters.size).toBe(4)
    expect(expectedFilters.get('name=')).toStrictEqual({
      ...filters.get('name='),
      value: 'John'
    })
  })

  it('should maintain referential equality for initial filters', () => {
    const initialMetadata = { some: 'initial-metadata' }
    const metadata = { some: 'metadata' }
    const filtersMetadataMap = vi.fn().mockReturnValue([['name=', 'John']])

    const { result, rerender } = renderHook(
      ({ initialMetadata, metadata }) =>
        useFiltersMetadata(
          initialMetadata,
          metadata,
          filters,
          vi.fn(),
          vi.fn(),
          filtersMetadataMap
        ),
      { initialProps: { initialMetadata: null, metadata: null } }
    )

    const firstRef = result.current
    rerender({ initialMetadata, metadata })
    expect(result.current).toBe(firstRef)
  })
})
