import { useEffect, useState } from 'react'
import Filter from '../../types/Filter.ts'
import FilterOperator from '../../types/enums/FilterOperator.ts'

export default function useSearchBy(filters: Map<string, Filter>): string[] {
  const [searchBy, setSearchBy] = useState(constructFilters(filters))
  useEffect(() => setSearchBy(constructFilters(filters)), [filters])
  return searchBy
}

function constructFilters(filters: Map<string, Filter>) {
  const newFilters: string[] = []
  filters.forEach((filter) => {
    if (filter.isSet === false) return
    newFilters.push(createFilter(filter))
  })
  return newFilters
}

function createFilter(filter: Filter): string {
  if (filter.operator === FilterOperator.IsNull || filter.operator === FilterOperator.IsNotNull)
    return `${filter.property} ${filter.operator}`
  return `${filter.property} ${filter.operator} ${filter.value}`
}
