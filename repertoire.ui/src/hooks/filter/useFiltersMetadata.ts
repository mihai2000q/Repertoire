import { Dispatch, SetStateAction, useEffect, useRef } from 'react'
import Filter, { FilterValue } from '../../types/Filter.ts'

export default function useFiltersMetadata<TMetadata>(
  metadata: TMetadata,
  filters: Map<string, Filter>,
  setFilters: Dispatch<SetStateAction<Map<string, Filter>>>,
  setInternalFilters: Dispatch<SetStateAction<Map<string, Filter>>>,
  filtersMetadataMap: (metadata: TMetadata) => [string, FilterValue][]
): Map<string, Filter> {
  const initialFilters = useRef(new Map<string, Filter>([...filters]))

  useEffect(() => {
    if (!metadata) return

    const newFilters = new Map<string, Filter>([...filters])
    filtersMetadataMap(metadata).forEach(([key, value]) => {
      const curr = filters.get(key)
      const newFilter: Filter = {
        ...curr,
        value: value
      }

      initialFilters.current.set(key, { ...newFilter, isSet: false })
      if (!curr.isSet) newFilters.set(key, newFilter)
    })
    setFilters(newFilters)
    setInternalFilters(newFilters)
  }, [metadata])

  return initialFilters.current
}
