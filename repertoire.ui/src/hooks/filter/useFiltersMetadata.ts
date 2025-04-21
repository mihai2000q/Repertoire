import { Dispatch, SetStateAction, useEffect, useRef } from 'react'
import Filter, { FilterValue } from '../../types/Filter.ts'

export default function useFiltersMetadata<TMetadata>(
  metadata: TMetadata,
  filters: Map<string, Filter>,
  setFilters: Dispatch<SetStateAction<Map<string, Filter>>>,
  filtersMetadataMap: (metadata: TMetadata) => [string, FilterValue][]
): Map<string, Filter> {
  const initialFilters = useRef(new Map<string, Filter>([...filters]))

  useEffect(() => {
    if (!metadata) return

    filtersMetadataMap(metadata).forEach(([key, value]) =>
      initialFilters.current.set(key, {
        ...filters.get(key),
        value: value
      })
    )
    setFilters(initialFilters.current)
  }, [metadata])

  return initialFilters.current
}
