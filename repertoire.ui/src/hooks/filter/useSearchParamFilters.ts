import { useState } from 'react'
import Filter from '../../types/Filter.ts'
import { useDidUpdate } from '@mantine/hooks'

interface useSearchParamFiltersProps<TState> {
  initialFilters: Filter[]
  activeFilters: Map<string, Filter>
  setSearchParams: (state: TState | ((state: TState) => TState)) => void
}

type useSearchParamFiltersReturnType = [
  filters: Map<string, Filter>,
  setFilters: (filters: Map<string, Filter>, withSearchParams?: boolean) => void
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
  useDidUpdate(
    () =>
      setFilters(
        new Map<string, Filter>(
          Array.from(filters.values()).map((f) =>
            activeFilters.get(f.property + f.operator)
              ? [f.property + f.operator, activeFilters.get(f.property + f.operator)]
              : [f.property + f.operator, { ...f, isSet: false }]
          )
        )
      ),
    [JSON.stringify([...activeFilters])]
  )

  const setFiltersAndSearchParams = (newFilters: Map<string, Filter>, withSearchParams = true) => {
    setFilters(newFilters)
    if (!withSearchParams) return
    const newActiveFilters = new Map<string, Filter>()
    newFilters.forEach((value, key) => {
      if (value.isSet) {
        newActiveFilters.set(key, value)
      }
    })
    setSearchParams((prevState) => ({
      ...prevState,
      activeFilters: newActiveFilters,
      ...('currentPage' in prevState ? { currentPage: 1 } : {})
    }))
  }
  return [filters, setFiltersAndSearchParams]
}
