import { renderHook } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import useFiltersMetadata from './useFiltersMetadata.ts'
import Filter, { FilterValue } from '../../types/Filter.ts'
import FilterOperator from '../../types/enums/FilterOperator.ts'

describe('useFiltersMetadata', () => {
  const mockFiltersMetadataMap = vi.fn()
  const initialFilters = new Map<string, Filter>([
    ['name=', { property: 'name', operator: FilterOperator.Equal, value: '', isSet: false }],
    ['age>', { property: 'age', operator: FilterOperator.GreaterThan, value: 0, isSet: false }]
  ])

  beforeEach(() => mockFiltersMetadataMap.mockClear())

  it('should not update filters when metadata is null', () => {
    const setFilters = vi.fn()
    const { result } = renderHook(() =>
      useFiltersMetadata(null, initialFilters, setFilters, mockFiltersMetadataMap)
    )

    expect(mockFiltersMetadataMap).not.toHaveBeenCalled()
    expect(setFilters).not.toHaveBeenCalled()
    expect(result.current).toEqual(initialFilters)
  })

  it('should update filters when metadata changes', () => {
    const metadata = { some: 'metadata' }
    const expectedUpdates: [string, FilterValue][] = [
      ['name=', 'John'],
      ['age>', 25]
    ]

    mockFiltersMetadataMap.mockReturnValue(expectedUpdates)
    const setFilters = vi.fn()

    const { result, rerender } = renderHook(
      ({ metadata }) =>
        useFiltersMetadata(metadata, initialFilters, setFilters, mockFiltersMetadataMap),
      { initialProps: { metadata: null } }
    )

    expect(result.current).toEqual(initialFilters)

    // update metadata
    rerender({ metadata })

    expect(mockFiltersMetadataMap).toHaveBeenCalledWith(metadata)
    expect(setFilters).toHaveBeenCalledTimes(1)

    const updatedFilters = setFilters.mock.calls[0][0]
    expect(updatedFilters.get('name=').value).toBe('John')
    expect(updatedFilters.get('age>').value).toBe(25)
    expect(result.current).toEqual(updatedFilters)
  })

  it('should preserve existing filter properties when updating values', () => {
    const metadata = { some: 'metadata' }
    const expectedUpdates: [string, FilterValue][] = [['name=', 'John']]
    mockFiltersMetadataMap.mockReturnValue(expectedUpdates)
    const setFilters = vi.fn()

    renderHook(() =>
      useFiltersMetadata(metadata, initialFilters, setFilters, mockFiltersMetadataMap)
    )

    const updatedFilters = setFilters.mock.calls[0][0]
    const updatedNameFilter = updatedFilters.get('name=')
    expect(updatedNameFilter).toEqual({
      ...initialFilters.get('name='),
      value: 'John'
    })
  })

  it('should maintain referential equality for initialFilters ref between renders', () => {
    const metadata = { some: 'metadata' }
    mockFiltersMetadataMap.mockReturnValue([['name', 'John']])
    const setFilters = vi.fn()

    const { result, rerender } = renderHook(
      ({ metadata }) =>
        useFiltersMetadata(metadata, initialFilters, setFilters, mockFiltersMetadataMap),
      { initialProps: { metadata: null } }
    )

    const firstRenderRef = result.current
    rerender({ metadata })
    const secondRenderRef = result.current

    expect(firstRenderRef).toBe(secondRenderRef)
  })
})
