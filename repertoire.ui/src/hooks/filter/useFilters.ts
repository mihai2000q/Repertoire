import { Dispatch, SetStateAction, useState } from 'react'
import Filter from '../../types/Filter.ts'

export default function useFilters(
  initialFilters?: Filter[]
): [Map<string, Filter>, Dispatch<SetStateAction<Map<string, Filter>>>, number] {
  const [filters, setFilters] = useState(
    new Map<string, Filter>(initialFilters?.map((f) => [f.property + f.operator, f]) ?? [])
  )
  const filtersActiveSize = Array.from(filters.values()).filter((f) => f.isSet !== false).length

  return [filters, setFilters, filtersActiveSize]
}
