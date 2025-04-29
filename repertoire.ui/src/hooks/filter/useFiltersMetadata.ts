import { Dispatch, SetStateAction, useEffect, useRef } from 'react'
import Filter, { FilterValue } from '../../types/Filter.ts'

export default function useFiltersMetadata<TMetadata>(
  initialMetadata: TMetadata,
  metadata: TMetadata,
  filters: Map<string, Filter>,
  setFilters: Dispatch<SetStateAction<Map<string, Filter>>>,
  setInternalFilters: Dispatch<SetStateAction<Map<string, Filter>>>,
  filtersMetadataMap: (metadata: TMetadata) => [string, FilterValue][]
): Map<string, Filter> {
  const initialFilters = useRef(new Map<string, Filter>([...filters]))

  useEffect(() => {
    if (!initialMetadata) return

    const initialFiltersMetadata = new Map<string, FilterValue>([
      ...filtersMetadataMap(initialMetadata)
    ])

    filters.forEach((filter, key) => {
      const newInitialFilter: Filter = initialFiltersMetadata.has(key)
        ? { ...filter, value: initialFiltersMetadata.get(key), isSet: false }
        : { ...filter, isSet: false }

      initialFilters.current.set(key, newInitialFilter)
    })
  }, [initialMetadata])

  useEffect(() => {
    if (!initialMetadata) return

    const newFilters = new Map<string, Filter>([...filters])
    const filtersMetadata = new Map<string, FilterValue>([
      ...filtersMetadataMap(metadata ?? initialMetadata)
    ])

    filters.forEach((filter, key) => {
      if (!filter.isSet) {
        const newFilter: Filter = filtersMetadata.has(key)
          ? { ...filter, value: filtersMetadata.get(key) }
          : { ...filter }
        newFilters.set(key, newFilter)
      }
    })

    setFilters(newFilters)
    setInternalFilters(newFilters)
  }, [initialMetadata, metadata])

  return initialFilters.current
}
