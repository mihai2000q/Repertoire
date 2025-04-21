import { Dispatch, SetStateAction, useState } from 'react'
import Filter from '../../types/Filter.ts'

export default function useFilters(
  initialFilters?: Filter[]
): [Map<string, Filter>, number, Dispatch<SetStateAction<Map<string, Filter>>>] {
  const [filters, setFilters] = useState(
    new Map<string, Filter>(initialFilters?.map((f) => [f.property + f.operator, f]) ?? [])
  )
  const filtersActiveSize = Array.from(filters.values()).filter((f) => f.isSet).length

  return [filters, filtersActiveSize, setFilters]
}
