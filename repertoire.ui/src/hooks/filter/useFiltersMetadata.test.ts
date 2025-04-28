import { renderHook } from '@testing-library/react'
import useFiltersMetadata from './useFiltersMetadata.ts'
import Filter, { FilterValue } from '../../types/Filter.ts'
import FilterOperator from '../../types/enums/FilterOperator.ts'

describe('use Filters Metadata', () => {
  const initialFilters = new Map<string, Filter>([
    ['name=', { property: 'name', operator: FilterOperator.Equal, value: '', isSet: false }],
    ['age>', { property: 'age', operator: FilterOperator.GreaterThan, value: 0, isSet: false }]
  ])

  it('should not update filters when metadata is null', () => {
    const setFilters = vi.fn()
    const setInternalFilters = vi.fn()

    const filtersMetadataMap = vi.fn()
    const { result } = renderHook(() =>
      useFiltersMetadata(null, initialFilters, setFilters, setInternalFilters, filtersMetadataMap)
    )

    expect(filtersMetadataMap).not.toHaveBeenCalled()
    expect(setFilters).not.toHaveBeenCalled()
    expect(setInternalFilters).not.toHaveBeenCalled()
    expect(result.current).toEqual(initialFilters)
  })

  it('should update both filters and internal filters when metadata changes', () => {
    const metadata = { some: 'metadata' }
    const expectedUpdates: [string, FilterValue][] = [
      ['name=', 'John'],
      ['age>', 25]
    ]

    const setFilters = vi.fn()
    const setInternalFilters = vi.fn()
    const filtersMetadataMap = vi.fn().mockReturnValue(expectedUpdates)

    const { result, rerender } = renderHook(
      ({ metadata }) =>
        useFiltersMetadata(
          metadata,
          initialFilters,
          setFilters,
          setInternalFilters,
          filtersMetadataMap
        ),
      { initialProps: { metadata: null } }
    )

    expect(result.current).toEqual(initialFilters)

    // Update metadata
    rerender({ metadata })

    expect(filtersMetadataMap).toHaveBeenCalledWith(metadata)
    expect(setFilters).toHaveBeenCalledOnce()
    expect(setInternalFilters).toHaveBeenCalledOnce()

    const expectedFilters = setFilters.mock.calls[0][0]
    const internalFilters = setInternalFilters.mock.calls[0][0]
    expect(internalFilters).toEqual(expectedFilters)

    expect(expectedFilters.get('name=').value).toBe('John')
    expect(expectedFilters.get('age>').value).toBe(25)
    expect(result.current).toEqual(expect.any(Map))
  })

  it('should update isSet status when updating values to initial filters, and not for those already set', () => {
    const metadata = { some: 'metadata' }
    const existingFilters = new Map(initialFilters)
    existingFilters.set('name=', { ...existingFilters.get('name='), isSet: true, value: 'Michael' })
    const expectedUpdates: [string, FilterValue][] = [['name=', 'John'], ['age>', 12]]

    const setFilters = vi.fn()
    const filtersMetadataMap = vi.fn().mockReturnValue(expectedUpdates)

    const { result } = renderHook(() =>
      useFiltersMetadata(metadata, existingFilters, setFilters, vi.fn(), filtersMetadataMap)
    )

    const updatedFilters = setFilters.mock.calls[0][0]
    expect(updatedFilters.get('name=')).toEqual({
      ...existingFilters.get('name='),
      value: 'Michael',
      isSet: true
    })
    expect(updatedFilters.get('age>')).toEqual({
      ...existingFilters.get('age>'),
      value: 12,
      isSet: false
    })

    expect(result.current.get('name=')).toEqual({
      ...existingFilters.get('name='),
      value: 'John',
      isSet: false
    })
    expect(result.current.get('age>')).toEqual({
      ...existingFilters.get('age>'),
      value: 12,
      isSet: false
    })
  })

  it('should maintain referential equality for initialFilters ref', () => {
    const metadata = { some: 'metadata' }
    const filtersMetadataMap = vi.fn().mockReturnValue([['name=', 'John']])

    const { result, rerender } = renderHook(
      ({ metadata }) =>
        useFiltersMetadata(metadata, initialFilters, vi.fn(), vi.fn(), filtersMetadataMap),
      { initialProps: { metadata: null } }
    )

    const firstRef = result.current
    rerender({ metadata })
    expect(result.current).toBe(firstRef)
  })
})
