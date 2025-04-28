import { useState } from 'react'
import Filter from '../../types/Filter.ts'

interface useSearchParamFiltersProps<TState> {
  initialFilters: Filter[]
  activeFilters: Map<string, Filter>
  setSearchParams: (state: TState | ((state: TState) => TState)) => void
}

type useSearchParamFiltersReturnType = [
  filters: Map<string, Filter>,
  setFilters: (filters: Map<string, Filter>) => void
]

export default function useSearchParamFilters<
  TState extends { activeFilters: Map<string, Filter> }
>({
  initialFilters,
  activeFilters,
  setSearchParams
}: useSearchParamFiltersProps<TState>): useSearchParamFiltersReturnType {
  const [filters, setFilters] = useState(
    new Map<string, Filter>(
      initialFilters.map((f) =>
        activeFilters.get(f.property + f.operator)
          ? [f.property + f.operator, activeFilters.get(f.property + f.operator)]
          : [f.property + f.operator, f]
      )
    )
  )

  const setFiltersAndSearchParams = (newFilters: Map<string, Filter>) => {
    const newActiveFilters = new Map<string, Filter>()
    newFilters.forEach((value, key) => {
      if (value.isSet) {
        newActiveFilters.set(key, value)
      }
    })
    setSearchParams((prevState) => ({
      ...prevState,
      activeFilters: newActiveFilters
    }))
    setFilters(newFilters)
  }
  return [filters, setFiltersAndSearchParams]
}
